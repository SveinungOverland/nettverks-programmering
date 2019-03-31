package main

import (
	"fmt"
	"oving5/account"
	"oving5/config"
	"sync"
)

func main() {
	defer config.DB.Close()

	account.CreateAccount("123456789", "Helene Larsen", 4324)
	account.CreateAccount("987654321", "Jørgen Sandvik", 3588)
	account.CreateAccount("543216789", "Sondre Henriksen", 8773)
	account.CreateAccount("123459876", "Hanne Løkke", 8476)
	account.CreateAccount("341256789", "Trond Ødegaard", 7546)
	account.CreateAccount("982143657", "Herman Berge", 9760)

	accounts := account.FindWhere("balance >= ?", 5000)
	fmt.Println(accounts)

	richest := accounts[0]

	fmt.Printf("This is the richest of the found accounts %+v\n", richest)
	fmt.Println("Withdrawing 5000 from account")

	richest.Withdraw(5000)
	fmt.Printf("This is the result %+v\n", richest)

	moneyBeforeLockup := richest.Balance

	// Use 10 threads to withdraw 100 from the richest account [THREAD-UNSAFE]
	fmt.Println("Withdrawing 100 from the account across ten threads")
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			richest.Withdraw(10)
			wg.Done()
		}()
	}
	wg.Wait()
	account.Sync(&richest)

	fmt.Printf("%f - 100 should be %f, but is %f\n", moneyBeforeLockup, moneyBeforeLockup-100, richest.Balance)
	// Use 10 threads to withdraw 100 from the richest account [THREAD-SAFE]
	moneyBeforeLockup = richest.Balance
	fmt.Println("Withdrawing 100 from the account across ten threads safely")
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			richest.WithdrawSafely(10)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Printf("%f - 100 should be %f, and is %f\n", moneyBeforeLockup, moneyBeforeLockup-100, richest.Balance)

	fmt.Println("Closing DB Connection")
}
