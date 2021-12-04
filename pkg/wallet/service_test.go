package wallet

import (
	"testing"
)

func TestService_RegisterAccount_success(t *testing.T) {
	svc := &Service{}
	_, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		t.Errorf("invalid result, expected: nil, actual: %v", err)
	}
}

func TestService_RegisterAccount_alreadyRegistered(t *testing.T) {
	svc := &Service{}
	_, err1 := svc.RegisterAccount("+992000000001")
	_, err1 = svc.RegisterAccount("+992000000001")
	
	if err1 != ErrPhoneRegistered {
		t.Errorf("invalid result, expected: %v, actual: %v", ErrPhoneRegistered, err1)
	}
}

func TestService_FindAccountByID_nil(t *testing.T) {
	svc := &Service{}
	svc.RegisterAccount("+992000000001")
	svc.RegisterAccount("+992000000002")
	svc.RegisterAccount("+992000000003")
	svc.RegisterAccount("+992000000004")

	_, err := svc.FindAccountByID(1)

	if err != nil {
		t.Errorf("invalid result, expected: nil, actual: %v", err)
	}
}

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

/*
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
