# Go Interfaces --- Complete Notes

> A practical reference for understanding and using interfaces in Go,
> with examples from a simple banking application.

## Table of contents

1.  [What is an interface?](#1-what-is-an-interface)
2.  [Why interfaces are useful](#2-why-interfaces-are-useful)
3.  [Defining and implementing an
    interface](#3-defining-and-implementing-an-interface)
4.  [Implicit interface
    implementation](#4-implicit-interface-implementation)
5.  [Passing interfaces to
    functions](#5-passing-interfaces-to-functions)
6.  [Polymorphism](#6-polymorphism)
7.  [Pointer receivers and interface
    satisfaction](#7-pointer-receivers-and-interface-satisfaction)
8.  [One type can implement multiple
    interfaces](#8-one-type-can-implement-multiple-interfaces)
9.  [Interface composition](#9-interface-composition)
10. [The built-in error interface](#10-the-built-in-error-interface)
11. [Creating custom errors](#11-creating-custom-errors)
12. [errors.Is and errors.As](#12-errorsis-and-errorsas)
13. [The empty interface and any](#13-the-empty-interface-and-any)
14. [Type assertions and type
    switches](#14-type-assertions-and-type-switches)
15. [Best practices](#15-best-practices)
16. [Common mistakes](#16-common-mistakes)
17. [Complete banking application
    example](#17-complete-banking-application-example)
18. [Quick revision checklist](#18-quick-revision-checklist)

------------------------------------------------------------------------

## 1. What is an interface?

An **interface** is a type that defines a set of method signatures. It
describes **what a value can do**, rather than the data it stores or
exactly how it performs an operation.

For example:

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}
```

This interface says:

> Any type that has a method named `Pay`, accepting a `float64` and
> returning an `error`, satisfies `PaymentProcessor`.

An interface does not normally contain the implementation of the method.
Concrete types provide the implementation.

### Struct vs. method vs. interface

  -----------------------------------------------------------------------
  Concept                 Main purpose            Banking example
  ----------------------- ----------------------- -----------------------
  Struct                  Groups data into a type `BankAccount` stores
                                                  holder and balance

  Method                  Defines behaviour       `Pay()` performs a
                          associated with a type  payment

  Interface               Describes required      `PaymentProcessor`
                          behaviour               requires `Pay()`
  -----------------------------------------------------------------------

Example struct:

``` go
type BankAccount struct {
    AccountHolder string
    Balance       float64
}
```

The struct stores the account data. An interface can describe the
behaviour that other code expects from it.

------------------------------------------------------------------------

## 2. Why interfaces are useful

Interfaces let different concrete types be used through a common
behaviour.

Imagine an application that supports payments from:

-   Bank accounts
-   M-Pesa wallets
-   Credit cards
-   Cryptocurrency wallets

Each payment method stores different data and may process payments
differently. However, the rest of the application may only need to know
that each one can make a payment.

An interface provides that common contract.

Benefits include:

-   **Flexibility:** different types can be used by the same function.
-   **Loose coupling:** code can depend on behaviour rather than a
    specific implementation.
-   **Testability:** test doubles can implement the same interface as
    production dependencies.
-   **Extensibility:** a new payment type can be added without rewriting
    the function that uses the interface.
-   **Clear responsibilities:** small interfaces make it clear what a
    caller needs.

------------------------------------------------------------------------

## 3. Defining and implementing an interface

Define an interface:

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}
```

Create a concrete type:

``` go
type BankAccount struct {
    AccountHolder string
    Balance       float64
}
```

Give the type the required method:

``` go
func (account *BankAccount) Pay(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("payment amount must be greater than zero")
    }

    if amount > account.Balance {
        return fmt.Errorf("insufficient funds")
    }

    account.Balance -= amount
    return nil
}
```

Because `*BankAccount` has the required `Pay(float64) error` method, it
satisfies `PaymentProcessor`.

The method signature must match. A method such as
`Pay(amount int) error` would **not** satisfy this interface because
`int` is not `float64`.

------------------------------------------------------------------------

## 4. Implicit interface implementation

Go uses **implicit interface implementation**.

Unlike some languages, Go does not require a type to write something
like:

``` text
BankAccount implements PaymentProcessor
```

There is no `implements` keyword.

The compiler checks whether the type has all the methods required by the
interface.

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}

type Mpesa struct {
    PhoneNumber string
    Balance     float64
}

func (m *Mpesa) Pay(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("payment amount must be greater than zero")
    }

    if amount > m.Balance {
        return fmt.Errorf("insufficient M-Pesa balance")
    }

    m.Balance -= amount
    return nil
}
```

`*Mpesa` satisfies `PaymentProcessor` because it has the required
method.

A useful compile-time check is:

``` go
var _ PaymentProcessor = (*Mpesa)(nil)
```

This declaration asks the compiler to verify that `*Mpesa` implements
`PaymentProcessor`. It is a compile-time check; it does not create a
usable `Mpesa` instance.

------------------------------------------------------------------------

## 5. Passing interfaces to functions

A function can accept an interface instead of a specific concrete type:

``` go
func ProcessPayment(processor PaymentProcessor, amount float64) error {
    return processor.Pay(amount)
}
```

Now it can accept any value whose type satisfies `PaymentProcessor`.

``` go
account := &BankAccount{
    AccountHolder: "Alex",
    Balance:       1000,
}

mpesa := &Mpesa{
    PhoneNumber: "0712345678",
    Balance:     2000,
}

if err := ProcessPayment(account, 300); err != nil {
    fmt.Println("Bank payment failed:", err)
}

if err := ProcessPayment(mpesa, 500); err != nil {
    fmt.Println("M-Pesa payment failed:", err)
}
```

There is no need to write a separate `ProcessBankPayment` and
`ProcessMpesaPayment` function if both types satisfy the same interface.

The function depends on the behaviour it needs: `Pay`.

------------------------------------------------------------------------

## 6. Polymorphism

**Polymorphism** means that different concrete types can be used through
a common interface while each type provides its own behaviour.

For example:

``` go
func ProcessPayment(processor PaymentProcessor, amount float64) error {
    return processor.Pay(amount)
}
```

When called with a `*BankAccount`, Go invokes `BankAccount.Pay`.

When called with a `*Mpesa`, Go invokes `Mpesa.Pay`.

The function remains the same; the concrete value determines which
method implementation runs.

Conceptually:

``` text
                 ProcessPayment
                       |
                PaymentProcessor
                       |
             +---------+---------+
             |                   |
             v                   v
        *BankAccount           *Mpesa
             |                   |
          Pay()                 Pay()
             |                   |
       Bank-specific        M-Pesa-specific
       implementation      implementation
```

------------------------------------------------------------------------

## 7. Pointer receivers and interface satisfaction

This is an important Go rule.

Consider:

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}

type BankAccount struct {
    Balance float64
}

func (account *BankAccount) Pay(amount float64) error {
    account.Balance -= amount
    return nil
}
```

`Pay` has a **pointer receiver**: `*BankAccount`.

Therefore:

-   `*BankAccount` satisfies `PaymentProcessor`.
-   `BankAccount` does not satisfy `PaymentProcessor` in this example.

Use a pointer when passing the account:

``` go
account := BankAccount{Balance: 1000}

ProcessPayment(&account, 100) // Correct
// ProcessPayment(account, 100) // Does not compile
```

### Method-set rule

For a defined type `T`:

-   The method set of `T` includes methods declared with receiver `T`.
-   The method set of `*T` includes methods declared with receiver `T`
    and `*T`.

This is why a pointer to a type has at least the methods available to
the value type, plus pointer-receiver methods.

Go may automatically take the address of an addressable variable when
calling a pointer-receiver method directly:

``` go
account.Pay(100)
```

This is effectively like:

``` go
(&account).Pay(100)
```

However, that convenience does **not** change interface satisfaction. If
the interface requires a pointer-receiver method, pass `&account` when
the interface expects it.

------------------------------------------------------------------------

## 8. One type can implement multiple interfaces

A concrete type can satisfy several interfaces at the same time.

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}

type RefundProcessor interface {
    Refund(amount float64) error
}

type AccountReporter interface {
    GetBalance() float64
}
```

A `BankAccount` can implement all three:

``` go
type BankAccount struct {
    Balance float64
}

func (a *BankAccount) Pay(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    if amount > a.Balance {
        return fmt.Errorf("insufficient funds")
    }

    a.Balance -= amount
    return nil
}

func (a *BankAccount) Refund(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("refund amount must be positive")
    }

    a.Balance += amount
    return nil
}

func (a *BankAccount) GetBalance() float64 {
    return a.Balance
}
```

Different functions can accept the interface they need:

``` go
func ProcessPayment(p PaymentProcessor, amount float64) error {
    return p.Pay(amount)
}

func ProcessRefund(r RefundProcessor, amount float64) error {
    return r.Refund(amount)
}

func DisplayBalance(r AccountReporter) {
    fmt.Println("Balance:", r.GetBalance())
}
```

The same `*BankAccount` can be passed to all three functions.

**Design principle:** a function should accept the smallest interface
that provides the methods it actually needs.

------------------------------------------------------------------------

## 9. Interface composition

Interface composition means combining smaller interfaces into a larger
interface by embedding them.

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}

type RefundProcessor interface {
    Refund(amount float64) error
}

type TransactionProcessor interface {
    PaymentProcessor
    RefundProcessor
}
```

`TransactionProcessor` requires both methods:

``` go
Pay(amount float64) error
Refund(amount float64) error
```

A type must implement both methods to satisfy `TransactionProcessor`.

Example:

``` go
type BankAccount struct {
    Balance float64
}

func (a *BankAccount) Pay(amount float64) error {
    if amount <= 0 || amount > a.Balance {
        return fmt.Errorf("invalid payment or insufficient funds")
    }

    a.Balance -= amount
    return nil
}

func (a *BankAccount) Refund(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("refund amount must be positive")
    }

    a.Balance += amount
    return nil
}
```

Now a function can require both behaviours:

``` go
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
```

### Why compose interfaces?

-   Reuse smaller behaviour contracts.
-   Build larger interfaces only when a caller needs multiple
    behaviours.
-   Avoid forcing every type to implement methods it does not support.

For example, a payment-only M-Pesa integration may satisfy
`PaymentProcessor` without satisfying `TransactionProcessor` if it does
not support refunds.

------------------------------------------------------------------------

## 10. The built-in `error` interface

Go's built-in `error` type is an interface. Its definition is
essentially:

``` go
type error interface {
    Error() string
}
```

Any type with an `Error() string` method satisfies the `error`
interface.

Functions commonly return an error as their last return value:

``` go
func Withdraw(account *BankAccount, amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("withdrawal amount must be greater than zero")
    }

    if amount > account.Balance {
        return fmt.Errorf("insufficient funds")
    }

    account.Balance -= amount
    return nil
}
```

Call the function and check the error:

``` go
err := Withdraw(account, 500)

if err != nil {
    fmt.Println("Withdrawal failed:", err)
    return
}

fmt.Println("Withdrawal successful")
```

Important:

-   A non-`nil` error means the operation reported a problem.
-   `nil` means no error was returned.
-   Returning an error does not automatically stop execution. The caller
    decides how to handle it.

A common Go pattern is:

``` go
result, err := SomeFunction()
if err != nil {
    // Handle the error.
    return
}

// Continue using result.
```

For simple errors, use the standard library's `errors.New` or
`fmt.Errorf`.

``` go
return errors.New("insufficient funds")
```

``` go
return fmt.Errorf("account %s has insufficient funds", accountNumber)
```

------------------------------------------------------------------------

## 11. Creating custom errors

A custom error is useful when the caller needs structured information,
not just an error message.

Define a struct:

``` go
type InsufficientFundsError struct {
    Requested float64
    Available float64
}
```

Implement `Error() string`:

``` go
func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf(
        "insufficient funds: requested %.2f, available %.2f",
        e.Requested,
        e.Available,
    )
}
```

Because it has the required `Error() string` method,
`*InsufficientFundsError` satisfies the built-in `error` interface.

Return it from a function:

``` go
func Withdraw(account *BankAccount, amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("withdrawal amount must be positive")
    }

    if amount > account.Balance {
        return &InsufficientFundsError{
            Requested: amount,
            Available: account.Balance,
        }
    }

    account.Balance -= amount
    return nil
}
```

The function's return type is `error`, but it can return a concrete
`*InsufficientFundsError` because that type implements `error`.

------------------------------------------------------------------------

## 12. `errors.Is` and `errors.As`

Go provides two important functions for inspecting errors.

### `errors.Is`

Use `errors.Is` to check whether an error matches a particular error
value, including through a wrapping chain.

``` go
var ErrInsufficientFunds = errors.New("insufficient funds")

func Withdraw(account *BankAccount, amount float64) error {
    if amount > account.Balance {
        return fmt.Errorf(
            "withdrawal rejected: %w",
            ErrInsufficientFunds,
        )
    }

    account.Balance -= amount
    return nil
}
```

Check it:

``` go
err := Withdraw(account, 2000)

if errors.Is(err, ErrInsufficientFunds) {
    fmt.Println("The account has insufficient funds.")
}
```

The `%w` verb wraps an error so that functions such as `errors.Is` can
inspect it.

### `errors.As`

Use `errors.As` when you want to check for a particular error type and
access its fields.

``` go
err := Withdraw(account, 2000)

var fundsErr *InsufficientFundsError

if errors.As(err, &fundsErr) {
    fmt.Println("Requested:", fundsErr.Requested)
    fmt.Println("Available:", fundsErr.Available)
}
```

Quick comparison:

  -----------------------------------------------------------------------
  Function                            What it checks
  ----------------------------------- -----------------------------------
  `errors.Is`                         Whether an error matches a
                                      particular error value

  `errors.As`                         Whether an error can be treated as
                                      a particular error type
  -----------------------------------------------------------------------

Both functions can inspect wrapped errors.

------------------------------------------------------------------------

## 13. The empty interface and `any`

An empty interface has no methods:

``` go
interface{}
```

Because it requires no methods, every Go value satisfies it.

Modern Go also provides the alias:

``` go
any
```

So these are equivalent:

``` go
func PrintValue(value interface{}) {
    fmt.Println(value)
}
```

``` go
func PrintValue(value any) {
    fmt.Println(value)
}
```

An `any` value can hold values of different types:

``` go
var value any

value = 100
value = "Hello"
value = true
value = BankAccount{Balance: 500}
```

Use `any` when a function genuinely needs to accept values of different,
unknown types. Do not use it automatically: a specific interface gives
the compiler more information and makes code easier to understand.

**Important:** an interface value is not simply the concrete value
itself. Conceptually, it carries a dynamic concrete type and a dynamic
value. A `nil` interface and an interface holding a typed `nil` pointer
are not the same:

``` go
var p *BankAccount = nil
var processor PaymentProcessor = p

fmt.Println(processor == nil) // false
```

The interface is non-`nil` because it contains the dynamic type
`*BankAccount`, even though the pointer value is nil. Calling a method
through it may panic if that method does not handle a nil receiver.

------------------------------------------------------------------------

## 14. Type assertions and type switches

Sometimes a function receives an interface but needs to inspect the
concrete value inside it.

### Type assertion

A type assertion checks whether an interface value holds a particular
concrete type:

``` go
var processor PaymentProcessor = &BankAccount{
    Balance: 1000,
}

account, ok := processor.(*BankAccount)

if ok {
    fmt.Println("Bank balance:", account.Balance)
}
```

The two-result form is safer:

``` go
value, ok := someInterface.(SomeType)
```

-   `value` contains the asserted value if successful.
-   `ok` is `true` if the assertion succeeds.
-   If it fails, `ok` is `false` and `value` is the zero value of
    `SomeType`.

The one-result form panics if the assertion fails:

``` go
account := processor.(*BankAccount)
```

Use it only when failure is impossible or intentionally treated as a
programming error.

### Type switch

A type switch handles several possible concrete types:

``` go
func DescribePaymentMethod(p PaymentProcessor) {
    switch value := p.(type) {
    case *BankAccount:
        fmt.Println("Bank account balance:", value.Balance)

    case *Mpesa:
        fmt.Println("M-Pesa number:", value.PhoneNumber)

    default:
        fmt.Println("Unknown payment method")
    }
}
```

Use type assertions and switches when concrete-type-specific behaviour
is genuinely needed. If you only need the interface's methods, call
those methods directly instead.

------------------------------------------------------------------------

## 15. Best practices

1.  **Keep interfaces small.** Prefer a few focused methods over a large
    interface with unrelated responsibilities.

2.  **Define interfaces around the needs of the caller.** The code using
    a dependency should describe the behaviour it requires.

3.  **Accept interfaces, return concrete types** as a useful default
    guideline. It is not an absolute rule, but it often keeps APIs
    simple.

4.  **Do not create interfaces without a reason.** An interface is
    useful when it enables substitution, decoupling, testing, or a clear
    abstraction.

5.  **Use pointer receivers when methods need to modify the original
    value** or when copying the value is undesirable.

6.  **Be consistent with receiver choices** across methods of a type,
    especially when the type is intended to satisfy interfaces.

7.  **Check errors explicitly.** Do not ignore an error unless you have
    a deliberate reason.

8.  **Wrap errors with context** using `%w` when callers may need to
    inspect the underlying error.

9.  **Use `errors.Is` for error values and `errors.As` for error
    types.**

10. **Avoid unnecessary `any`.** Prefer a specific interface when the
    required behaviour is known.

11. **Use compile-time interface checks** when they make implementation
    requirements clearer:

    ``` go
    var _ PaymentProcessor = (*BankAccount)(nil)
    ```

12. **Avoid oversized interfaces.** A type that can make payments should
    not be forced to support refunds if it cannot perform them.

------------------------------------------------------------------------

## 16. Common mistakes

### Mistake 1: Method signature does not match

``` go
type PaymentProcessor interface {
    Pay(amount float64) error
}

func (a *BankAccount) Pay(amount int) error {
    return nil
}
```

This does not satisfy the interface. The parameter types must match
exactly.

### Mistake 2: Passing a value when only its pointer implements the interface

If `Pay` has a pointer receiver, `*BankAccount` implements the
interface, not `BankAccount`.

``` go
ProcessPayment(&account, 100) // Correct for pointer receiver
```

### Mistake 3: Forgetting to check an error

``` go
Withdraw(account, 500) // Error is ignored.
```

Prefer:

``` go
if err := Withdraw(account, 500); err != nil {
    fmt.Println("Withdrawal failed:", err)
    return
}
```

### Mistake 4: Assuming `nil` always means a nil interface

An interface holding a typed nil pointer is itself non-`nil`. Be careful
when returning pointers as interface values.

### Mistake 5: Making one interface responsible for everything

A large interface makes it harder for types to satisfy it and makes
functions depend on more behaviour than they need. Split
responsibilities into smaller interfaces and compose them when
necessary.

------------------------------------------------------------------------

## 17. Complete banking application example

This program brings together:

-   Multiple concrete types (`BankAccount` and `CryptoWallet`).
-   Implicit interface implementation.
-   Multiple interfaces.
-   Interface composition.
-   Functions that accept interfaces.
-   The built-in `error` interface.
-   A custom error type.
-   `errors.As` for inspecting a custom error.

Save this as `main.go`:

``` go
package main

import (
    "errors"
    "fmt"
)

// 1. Interfaces

type PaymentProcessor interface {
    Pay(amount float64) error
}

type RefundProcessor interface {
    Refund(amount float64) error
}

type AccountReporter interface {
    GetBalance() float64
}

// Interface composition: requires both Pay and Refund.
type TransactionProcessor interface {
    PaymentProcessor
    RefundProcessor
}

// 2. Custom error type

type InsufficientFundsError struct {
    Requested float64
    Available float64
}

func (e *InsufficientFundsError) Error() string {
    return fmt.Sprintf(
        "insufficient funds: requested %.2f, available %.2f",
        e.Requested,
        e.Available,
    )
}

// 3. Bank account

type BankAccount struct {
    AccountHolder string
    Balance       float64
}

func (a *BankAccount) Pay(amount float64) error {
    if amount <= 0 {
        return errors.New("payment must be greater than zero")
    }

    if amount > a.Balance {
        return &InsufficientFundsError{
            Requested: amount,
            Available: a.Balance,
        }
    }

    a.Balance -= amount
    fmt.Println("Bank account payment successful!")
    return nil
}

func (a *BankAccount) Refund(amount float64) error {
    if amount <= 0 {
        return errors.New("refund must be greater than zero")
    }

    a.Balance += amount
    fmt.Println("Bank account refund successful!")
    return nil
}

func (a *BankAccount) GetBalance() float64 {
    return a.Balance
}

// 4. Crypto wallet

type CryptoWallet struct {
    WalletAddress string
    Balance       float64
}

func (w *CryptoWallet) Pay(amount float64) error {
    if amount <= 0 {
        return errors.New("payment must be greater than zero")
    }

    if amount > w.Balance {
        return &InsufficientFundsError{
            Requested: amount,
            Available: w.Balance,
        }
    }

    w.Balance -= amount
    fmt.Println("Crypto wallet payment successful!")
    return nil
}

func (w *CryptoWallet) Refund(amount float64) error {
    if amount <= 0 {
        return errors.New("refund must be greater than zero")
    }

    w.Balance += amount
    fmt.Println("Crypto wallet refund successful!")
    return nil
}

func (w *CryptoWallet) GetBalance() float64 {
    return w.Balance
}

// 5. Functions that accept interfaces

func ProcessPayment(p PaymentProcessor, amount float64) error {
    return p.Pay(amount)
}

func ProcessRefund(r RefundProcessor, amount float64) error {
    return r.Refund(amount)
}

func DisplayBalance(r AccountReporter) {
    fmt.Printf("Current balance: %.2f\n", r.GetBalance())
}

func ProcessTransaction(
    p TransactionProcessor,
    payment float64,
    refund float64,
) error {
    if err := p.Pay(payment); err != nil {
        return err
    }

    if err := p.Refund(refund); err != nil {
        return err
    }

    return nil
}

// 6. Main function

func main() {
    account := BankAccount{
        AccountHolder: "Alex",
        Balance:       5000,
    }

    wallet := CryptoWallet{
        WalletAddress: "0x123456789",
        Balance:       2.5,
    }

    fmt.Println("=== BANK ACCOUNT ===")

    if err := ProcessPayment(&account, 1000); err != nil {
        fmt.Println("Payment failed:", err)
    }

    if err := ProcessRefund(&account, 200); err != nil {
        fmt.Println("Refund failed:", err)
    }

    DisplayBalance(&account)

    fmt.Println("\n=== COMPOSED TRANSACTION ===")

    if err := ProcessTransaction(&account, 500, 100); err != nil {
        fmt.Println("Transaction failed:", err)
    }

    DisplayBalance(&account)

    fmt.Println("\n=== CRYPTO WALLET ===")

    if err := ProcessTransaction(&wallet, 0.5, 0.2); err != nil {
        fmt.Println("Transaction failed:", err)
    }

    DisplayBalance(&wallet)

    fmt.Println("\n=== CUSTOM ERROR HANDLING ===")

    err := ProcessPayment(&account, 10000)

    if err != nil {
        fmt.Println("Payment failed:", err)

        var fundsErr *InsufficientFundsError

        if errors.As(err, &fundsErr) {
            fmt.Printf("Requested: %.2f\n", fundsErr.Requested)
            fmt.Printf("Available: %.2f\n", fundsErr.Available)
        }
    }
}
```

Run it:

``` bash
go run main.go
```

The example is for learning. A production financial application would
also need stronger validation, currency-safe monetary representation
(for example, integer minor units or a decimal type), authentication,
authorization, audit logging, and database transactions for atomic
financial operations.

------------------------------------------------------------------------

## 18. Quick revision checklist

Use this checklist to test your understanding.

-   [ ] I can explain the difference between a struct, a method, and an
    interface.
-   [ ] I understand that interfaces describe behaviour through method
    signatures.
-   [ ] I understand implicit interface implementation in Go.
-   [ ] I can explain why method signatures must match exactly.
-   [ ] I understand the difference between `T` and `*T` method sets.
-   [ ] I can pass different concrete types to one function through an
    interface.
-   [ ] I understand polymorphism.
-   [ ] I know that one type can satisfy multiple interfaces.
-   [ ] I can compose interfaces by embedding smaller interfaces.
-   [ ] I understand that `error` is a built-in interface.
-   [ ] I can return `nil` when an operation succeeds and an error when
    it fails.
-   [ ] I can create a custom error type by implementing
    `Error() string`.
-   [ ] I know when to use `errors.Is` and `errors.As`.
-   [ ] I understand `interface{}` and its alias `any`.
-   [ ] I can use type assertions and type switches when necessary.

## Final takeaway

The most important idea to remember is:

> **An interface describes what a type can do, not what the type is.**

A `BankAccount`, `Mpesa`, and `CryptoWallet` are different types with
different data and internal logic. But if they all implement
`Pay(amount float64) error`, they can all be used wherever a
`PaymentProcessor` is expected.

That is what makes Go interfaces so useful for building flexible,
testable, and maintainable applications.
