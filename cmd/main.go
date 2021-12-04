package main

import (
	"github.com/nwarior/wallet/pkg/wallet"
	"fmt"
)

func main() {
	svc := &wallet.Service{}
	svc.RegisterAccount("+992000000002")
	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = svc.Deposit(account.ID, 10)
	err = svc.Deposit(account.ID, 10)
	err = svc.Deposit(account.ID, 10)
	if err != nil {
		switch err {
		case wallet.ErrAmountMustBePositive:
			fmt.Println("Сумма должна быть положительной")
		case wallet.ErrAccountNotFound:
			fmt.Println("Аккаунт пользователя не найден")
		}
		return
	}

	fmt.Println(account.Balance) // 10
	fmt.Println(account)
}
