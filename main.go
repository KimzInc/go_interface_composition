package main

import (
	"errors"
	"fmt"
)

// ==================================================
// 1. INTERFACES
// ==================================================

// A type implementing this interface must have Pay().
type PaymentProcessor interface {
	Pay(amount float64) error
}

// A type implementing this interface must have Refund().
type RefundProcessor interface {
	Refund(amount float64) error
}

// A type implementing this interface must have GetBalance().
type AccountReporter interface {
	GetBalance() float64
}

// INTERFACE COMPOSITION:
// A type implementing TransactionProcessor must implement
// both PaymentProcessor and RefundProcessor.
type TransactionProcessor interface {
	PaymentProcessor
	RefundProcessor
}

// ==================================================
// 2. CUSTOM ERROR TYPE
// ==================================================

type InsufficientFundsError struct {
	Requested float64
	Available float64
}

// This method makes *InsufficientFundsError satisfy
// Go's built-in error interface.
func (e *InsufficientFundsError) Error() string {
	return fmt.Sprintf(
		"insufficient funds: requested %.2f, available %.2f",
		e.Requested,
		e.Available,
	)
}

// ==================================================
// 3. BANK ACCOUNT STRUCT
// ==================================================

type BankAccount struct {
	AccountHolder string
	Balance       float64
}

// BankAccount implements PaymentProcessor.
func (account *BankAccount) Pay(amount float64) error {
	if amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	if amount > account.Balance {
		return &InsufficientFundsError{
			Requested: amount,
			Available: account.Balance,
		}
	}

	account.Balance -= amount
	fmt.Println("Bank account payment successful!")

	return nil
}

// BankAccount implements RefundProcessor.
func (account *BankAccount) Refund(amount float64) error {
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	account.Balance += amount
	fmt.Println("Bank account refund successful!")

	return nil
}

// BankAccount implements AccountReporter.
func (account *BankAccount) GetBalance() float64 {
	return account.Balance
}

// ==================================================
// 4. CRYPTO WALLET STRUCT
// ==================================================

type CryptoWallet struct {
	WalletAddress string
	Balance       float64
}

// CryptoWallet implements PaymentProcessor.
func (wallet *CryptoWallet) Pay(amount float64) error {
	if amount <= 0 {
		return errors.New("payment amount must be greater than zero")
	}

	if amount > wallet.Balance {
		return &InsufficientFundsError{
			Requested: amount,
			Available: wallet.Balance,
		}
	}

	wallet.Balance -= amount
	fmt.Println("Crypto wallet payment successful!")

	return nil
}

// CryptoWallet implements RefundProcessor.
func (wallet *CryptoWallet) Refund(amount float64) error {
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	wallet.Balance += amount
	fmt.Println("Crypto wallet refund successful!")

	return nil
}

// CryptoWallet implements AccountReporter.
func (wallet *CryptoWallet) GetBalance() float64 {
	return wallet.Balance
}

// ==================================================
// 5. FUNCTIONS THAT ACCEPT INTERFACES
// ==================================================

// Accepts any type that implements PaymentProcessor.
func ProcessPayment(processor PaymentProcessor, amount float64) error {
	return processor.Pay(amount)
}

// Accepts any type that implements RefundProcessor.
func ProcessRefund(processor RefundProcessor, amount float64) error {
	return processor.Refund(amount)
}

// Accepts any type that implements AccountReporter.
func DisplayBalance(reporter AccountReporter) {
	fmt.Printf("Current balance: %.2f\n", reporter.GetBalance())
}

// Accepts any type that implements TransactionProcessor.
// Remember: TransactionProcessor requires BOTH Pay() and Refund().
func ProcessTransaction(
	processor TransactionProcessor,
	payment float64,
	refund float64,
) error {

	if err := processor.Pay(payment); err != nil {
		return err
	}

	if err := processor.Refund(refund); err != nil {
		return err
	}

	return nil
}

// ==================================================
// 6. MAIN FUNCTION
// ==================================================

func main() {

	// ----------------------------------------------
	// Create a bank account
	// ----------------------------------------------

	account := BankAccount{
		AccountHolder: "Kim",
		Balance:       5000,
	}

	fmt.Println("========== BANK ACCOUNT ==========")

	// BankAccount is passed to PaymentProcessor.
	if err := ProcessPayment(&account, 1000); err != nil {
		fmt.Println("Payment failed:", err)
	}

	// BankAccount is passed to RefundProcessor.
	if err := ProcessRefund(&account, 200); err != nil {
		fmt.Println("Refund failed:", err)
	}

	// BankAccount is passed to AccountReporter.
	DisplayBalance(&account)

	// ----------------------------------------------
	// Interface composition
	// ----------------------------------------------

	fmt.Println("\n========== COMBINED TRANSACTION ==========")

	// BankAccount satisfies TransactionProcessor
	// because it implements BOTH Pay() and Refund().
	if err := ProcessTransaction(&account, 500, 100); err != nil {
		fmt.Println("Transaction failed:", err)
	}

	DisplayBalance(&account)

	// ----------------------------------------------
	// Create a crypto wallet
	// ----------------------------------------------

	wallet := CryptoWallet{
		WalletAddress: "0x123456789",
		Balance:       2.5,
	}

	fmt.Println("\n========== CRYPTO WALLET ==========")

	// CryptoWallet also satisfies TransactionProcessor.
	if err := ProcessTransaction(&wallet, 0.5, 0.2); err != nil {
		fmt.Println("Transaction failed:", err)
	}

	DisplayBalance(&wallet)

	// ----------------------------------------------
	// Demonstrate custom error handling
	// ----------------------------------------------

	fmt.Println("\n========== ERROR HANDLING ==========")

	err := ProcessPayment(&account, 10000)

	if err != nil {
		fmt.Println("Payment failed:", err)

		// Check whether the returned error is specifically
		// an InsufficientFundsError.
		var fundsError *InsufficientFundsError

		if errors.As(err, &fundsError) {
			fmt.Println("Requested amount:", fundsError.Requested)
			fmt.Println("Available balance:", fundsError.Available)
		}
	}

	fmt.Println("\n========== FINAL BALANCES ==========")

	fmt.Print("Bank account: ")
	DisplayBalance(&account)

	fmt.Print("Crypto wallet: ")
	DisplayBalance(&wallet)
}
