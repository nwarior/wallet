package main

import (
	"github.com/nwarior/wallet/pkg/wallet"
)

func main() {
	s := &wallet.Service{}

	s.RegisterAccount("+992000000001")
	s.RegisterAccount("+992000000002")
	s.RegisterAccount("+992000000003")
	s.RegisterAccount("+992000000004")

	//s.ExportToFile("data/accounts.txt")

	s.ImportFromFile("data/accounts1.txt")
}
