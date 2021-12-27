package wallet

import (
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/nwarior/wallet/pkg/types"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAmountMustBePositive = errors.New("amount must be greater than zero")
var ErrAccountNotFound = errors.New("account not found")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}

	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)

	return account, nil
}

func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount <= 0 {
		return ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return ErrAccountNotFound
	}

	// зачисление средств пока не рассматриваем как платёж
	account.Balance += amount
	return nil
}

func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}

	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}

	if account == nil {
		return nil, ErrAccountNotFound
	}

	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}

	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: accountID,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			return acc, nil
		}
	}

	return nil, ErrAccountNotFound
}

func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}

	return nil, ErrPaymentNotFound
}

func (s *Service) Reject(paymentID string) error {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return err
	}

	account, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		return err
	}

	payment.Status = types.PaymentStatusFail
	account.Balance += payment.Amount
	return nil
}

func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	pay, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	payment, err := s.Pay(pay.AccountID, pay.Amount, pay.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        uuid.New().String(),
		AccountID: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}

func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var favorite *types.Favorite
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			favorite = fav
			break
		}
	}
	if favorite == nil {
		return nil, ErrFavoriteNotFound
	}

	payment, err := s.Pay(favorite.AccountID, favorite.Amount, favorite.Category)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for _, acc := range s.accounts {
		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.ID), 10)))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(acc.Phone))
		if err != nil {
			if err != nil {
				return err
			}
		}
		_, err = file.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.Balance), 10)))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte("|"))
		if err != nil {
			return err
		}

	}

	return nil
}

func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return err
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Print(err)
		}
	}()

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			content = append(content, buf[:read]...)
			break
		}

		if err != nil {
			log.Print(err)
			return err
		}

		content = append(content, buf[:read]...)
	}

	data := string(content)
	accounts := strings.Split(data, "|")
	if len(accounts) > 0 {
		accounts = accounts[:len(accounts)-1]
	}

	for _, acc := range accounts {
		fields := strings.Split(acc, ";")  // добавляем к каждому полю ";" fields - это поля
		id, err := strconv.Atoi(fields[0]) // id это первое поле поэтому fields[0]
		if err != nil {
			log.Print(err)
			return err
		}
		phone := fields[1] // здесь и так string поэтому проверка на ошибок не будет

		balance, err := strconv.Atoi(fields[2])
		if err != nil {
			log.Print(err)
			return err
		}
		// добавление всех полей в целую структуру.
		account := &types.Account{
			ID:      int64(id),
			Phone:   types.Phone(phone),
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, account)
	}

	return nil
}

func (s *Service) Export(dir string) error {
	// file or accounts
	file, err := os.Create(dir + "accounts.dump")
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for _, acc := range s.accounts {
		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.ID), 10)))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(acc.Phone))
		if err != nil {
			if err != nil {
				return err
			}
		}
		_, err = file.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte(strconv.FormatInt(int64(acc.Balance), 10)))
		if err != nil {
			return err
		}

		_, err = file.Write([]byte("\n"))
		if err != nil {
			return err
		}

	}

	filePayment, err := os.Create(dir + "payments.dump")
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if cerr := filePayment.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for _, pay := range s.payments {
		_, err = filePayment.Write([]byte(pay.ID))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(strconv.FormatInt(int64(pay.AccountID), 10)))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(strconv.FormatInt(int64(pay.Amount), 10)))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(string(pay.Category)))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte(string(pay.Status)))
		if err != nil {
			return err
		}

		_, err = filePayment.Write([]byte("\n"))
		if err != nil {
			return err
		}

	}

	fileFavorite, err := os.Create(dir + "favorites.dump")
	if err != nil {
		log.Print(err)
	}
	defer func() {
		if cerr := fileFavorite.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()

	for _, fav := range s.favorites {
		_, err = fileFavorite.Write([]byte(fav.ID))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(strconv.FormatInt(int64(fav.AccountID), 10)))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(fav.Name))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(strconv.FormatInt(int64(fav.Amount), 10)))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(";"))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte(fav.Category))
		if err != nil {
			return err
		}

		_, err = fileFavorite.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
