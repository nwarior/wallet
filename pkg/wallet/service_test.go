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
