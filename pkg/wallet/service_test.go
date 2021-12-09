package wallet

import (
	"fmt"
	"reflect"
	"testing"
	"github.com/google/uuid"
	"github.com/nwarior/wallet/pkg/types"
)

type testService struct {
	*Service
}

func newTestService() *testService {
	return &testService{Service: &Service{}}
}

type testAccount struct {
	phone types.Phone
	balance types.Money
	payments []struct {
		amount types.Money
		category types.PaymentCategory
	}
}

var defaultTestAccount = testAccount {
	phone: "+992000000001",
	balance: 10_000_00,
	payments: []struct {
		amount	types.Money
		category types.PaymentCategory
	}{
		{amount: 1_000_00, category: "auto"},
	},
}

func (s *testService) addAccount(data testAccount)(*types.Account, []*types.Payment, error) {
	// регистрируем там пользователя
	account, err := s.RegisterAccount(data.phone)
	if err != nil {
		return nil, nil, fmt.Errorf("can't register account, error = %v", err)
	}

	// пополняем его счёт
	err = s.Deposit(account.ID, data.balance)
	if err != nil {
		return nil, nil, fmt.Errorf("can't deposit accoun, error = %v", err)
	}

	// выполняем платеже
	// можем создать слайс сразу нужной длины, поскольку знаем размер
	payments := make([]*types.Payment, len(data.payments))
	for i, payment := range data.payments {
		// тогда здесь работаем просто через index, а не через append
		payments[i], err = s.Pay(account.ID, payment.amount, payment.category)
		if err != nil {
			return nil, nil, fmt.Errorf("can't make payment, error = %v", err)
		}
	}

	return account, payments, nil
}


/*
func TestService_RegisterAccount_success(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		t.Errorf("invalid result, expected: nil, actual: %v", err)
	}
}
*/

func TestService_RegisterAccount_alreadyRegistered(t *testing.T) {
	svc := &Service{}
	_, err1 := svc.RegisterAccount("+992000000001")
	_, err1 = svc.RegisterAccount("+992000000001")
	
	if err1 != ErrPhoneRegistered {
		t.Errorf("invalid result, expected: %v, actual: %v", ErrPhoneRegistered, err1)
	}
}

/*
func TestService_FindAccountByID_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	got, err := s.FindAccountByID(account.ID)
	if err != nil {
		t.Error("FindAccountByID(): can't find account")
		return
	}

	if !reflect.DeepEqual(account, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", account, got)
	}
}

func TestService_FindAccountByID_fail(t *testing.T) {
	// создаём сервис
	s := newTestService()
	account, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	got, err := s.FindAccountByID(514152)
	if err == nil {
		t.Error("result must be ErrAccountNotFound")
		return
	}

	if reflect.DeepEqual(account, got) {
		t.Errorf("invalid result, expected: %v, actual: %v", account, got)
	}
}
*/
func TestService_FindPaymentByID_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти платёж
	payment := payments[0]
	got, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("FindPaymentByID(): error = %v", err)
		return
	}

	// сравниваем платежи
	if !reflect.DeepEqual(payment, got) {
		t.Errorf("FindPaymentByID(): wrong payment returned = %v", err)
		return
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, _, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	// пробуем найти несуществующий платёж
	_, err = s.FindPaymentByID(uuid.New().String())
	if err == nil {
		t.Error("FindPaymentByID(): must return error, returned nil")
		return
	}

	if err != ErrPaymentNotFound {
		t.Errorf("FindPaymentByID(): must return ErrPaymentNotFound, returned = %v", err)
		return
	}
}
/*
func TestService_FindPaymentByID_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.Deposit(1, 1000)

	payment, _ := svc.Pay(1, 300, "ресторан")

	payment2, _ := svc.FindPaymentByID(payment.ID)

	if payment.ID != payment2.ID {
		t.Errorf("invalid result, expected: %v, actual: %v", payment.ID, payment2 )
	}
}

func TestService_FindPaymentByID_fail(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.Deposit(1, 1000)

	_, err1 := svc.Pay(1, 300, "ресторан")

	_, err := svc.FindPaymentByID("12342143")

	if err1 == err {
		t.Errorf("invalid result, expected: %v, actual: %v", err1, err )
	}
}


func TestService_Reject_success(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.Deposit(1, 1000)

	payment, _ := svc.Pay(1, 200, "ресторан")

	err := svc.Reject(payment.ID)

	if err != nil {
		t.Errorf("invalid result, expected: nil, actual: %v", err)
	}
}


func TestService_Reject_fail(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.Deposit(1, 1000)
	svc.Pay(1, 200, "ресторан")
	err := svc.Reject("12333132")

	if err != ErrPaymentNotFound {
		t.Errorf("invalid result, expected: %v, actual: %v", ErrPaymentNotFound, err)
	}
}
*/

func TestService_Reject_success(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	// пробуем отменить платёж
	payment := payments[0]
	err = s.Reject(payment.ID)
	if err != nil {
		t.Errorf("Reject(): error = %v", err)
		return
	}

	savedPayment, err := s.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status != types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v,", savedPayment)
		return
	}

	savedAccount, err := s.FindAccountByID(payment.AccountID)
	if err != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance != defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}

/*
func TestService_Reject_fail(t *testing.T) {
	// создаём сервис
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}
	
	// пробуем отменить платёж
	payment := payments[0]
	err = s.Reject(uuid.New().String())
	if err == nil {
		t.Error("Reject(): must be error")
		return
	}

	savedPayment, err1 := s.FindPaymentByID(payment.ID)
	if err1 != nil {
		t.Errorf("Reject(): can't find payment by id, error = %v", err)
		return
	}
	if savedPayment.Status == types.PaymentStatusFail {
		t.Errorf("Reject(): status didn't changed, payment = %v,", savedPayment)
		return
	}

	savedAccount, err2 := s.FindAccountByID(payment.AccountID)
	if err2 != nil {
		t.Errorf("Reject(): can't find account by id, error = %v", err)
		return
	}
	if savedAccount.Balance == defaultTestAccount.balance {
		t.Errorf("Reject(): balance didn't changed, account = %v", savedAccount)
		return
	}
}
*/
func TestService_Repeat_success(t *testing.T) {
	svc := Service{}
 	svc.RegisterAccount("+9920000001")

 	account, err := svc.FindAccountByID(1)
 	if err != nil {
  		t.Errorf("\ngot > %v \nwant > nil", err)
 	}

 	err = svc.Deposit(account.ID, 1000_00)
 	if err != nil {
	 	t.Errorf("\ngot > %v \nwant > nil", err)
	}

 	payment, err := svc.Pay(account.ID, 100_00, "auto")
 	if err != nil {
  		t.Errorf("\ngot > %v \nwant > nil", err)
 	}

 	pay, err := svc.FindPaymentByID(payment.ID)
 	if err != nil {
 		t.Errorf("\ngot > %v \nwant > nil", err)
 	}

 	pay, err = svc.Repeat(pay.ID)
 	if err != nil {
  		t.Errorf("Repeat(): Error(): can't pay for an account(%v): %v", pay.ID, err)
 	}
}

func TestService_FavoritePayment_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	result, err := s.FavoritePayment(payment.ID, "kiki")
	if err != nil {
		t.Error(err)
		return
	}

	expected := &types.Favorite{
		ID: payment.ID,
		AccountID: payment.AccountID,
		Name: "kiki",
		Amount: payment.Amount,
		Category: payment.Category,
	}

	if !reflect.DeepEqual(expected, result) {
		t.Errorf("invalid error, expected: %v, actual: %v", expected, result)
	}
}

func TestService_FavoritePayment_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	result, err := s.FavoritePayment(payment.ID, "kiki")
	if err != nil {
		t.Error(err)
		return
	}

	expected := &types.Favorite{
		ID: payment.ID,
		AccountID: payment.AccountID,
		Name: "santo",
		Amount: payment.Amount,
		Category: payment.Category,
	}

	if reflect.DeepEqual(expected, result) {
		t.Errorf("invalid error, expected: %v, actual: %v", expected, result)
	}
}

func TestService_PayFromFavorite_success(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	result, err := s.FavoritePayment(payment.ID, "kiki")
	if err != nil {
		t.Error(err)
		return
	}

	fromFavor, err := s.PayFromFavorite(result.ID)
	if err != nil  {
		t.Error(err)
		return
	}

	expected := &types.Payment{
		ID:	fromFavor.ID,
		AccountID: fromFavor.AccountID,
		Amount: fromFavor.Amount, 
		Category: fromFavor.Category,
		Status: fromFavor.Status,
	}

	if !reflect.DeepEqual(expected, fromFavor) {
		t.Errorf("invalid error, expected: %v, actual: %v", expected, fromFavor)
	}
}

func TestService_PayFromFavorite_fail(t *testing.T) {
	s := newTestService()
	_, payments, err := s.addAccount(defaultTestAccount)
	if err != nil {
		t.Error(err)
		return
	}

	payment := payments[0]
	result, err := s.FavoritePayment(payment.ID, "kiki")
	if err != nil {
		t.Error(err)
		return
	}

	fromFavor, err := s.PayFromFavorite(result.ID)
	if err != nil  {
		t.Error(err)
		return
	}

	expected := &types.Payment{
		ID:	fromFavor.ID,
		AccountID: fromFavor.AccountID,
		Amount: fromFavor.Amount, 
		Category: "santo",
		Status: fromFavor.Status,
	}

	if reflect.DeepEqual(expected, fromFavor) {
		t.Errorf("invalid error, expected: %v, actual: %v", expected, fromFavor)
	}
}