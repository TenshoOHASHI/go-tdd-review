package pointer_test

import (
	"errors"
	"fmt"
	"testing"
)

// Type MyName Original Type
type Bitcoin int

// 既存のタイプの上にドメイン固有の機能を追加する場合に役立つ
type Stringer interface {
	String() string
}

type Wallet struct {
	// 同じパッケージ内でしかアクセスできない
	balance Bitcoin
}

// 独自の型にメソッドを追加できる（カスタマイズ）
func (b Bitcoin) String() string {
	return fmt.Sprintf("%d BTC", b)
}

// ポインターレシーバーを通じで、内部のフィールにアクセスする
// 通常のレシーバの場合、wは呼び出し元のコピーのため、呼び出しものwallet.balanceは変わらない（&wallet.balance内部のアドレスと異なる）
func (w *Wallet) Deposit(amount Bitcoin) error {
	if amount < 0 {
		return fmt.Errorf("Deposit is zero")
	}
	// w は *Wallet なので、Go が自動で参照を外し、w.balance と書ける。
	// 実際の意味は (*w).balance += amount に近い。
	w.balance += amount
	return nil
}

func (w *Wallet) Balance() Bitcoin {
	return w.balance
}

var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {
	if amount > w.balance {
		return ErrInsufficientFunds
	}
	w.balance -= amount

	return nil

}

func assertError(t *testing.T, got error, want error) {
	t.Helper()
	// エラーを想定していだけど、エラーでなければ、強制終了
	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if !errors.Is(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func assertBalance(t *testing.T, wallet Wallet, want Bitcoin) {
	t.Helper()

	// 残高を取得
	got := wallet.Balance()

	if got != want {
		t.Errorf("got %s want %s", got, want)

	}
}

func TestWallet(t *testing.T) {

	t.Run("Deposit", func(t *testing.T) {
		wallet := Wallet{}
		err := wallet.Deposit(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("Withdraw with funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))
		assertNoError(t, err)
	})

	t.Run("Withdraw insufficient funds", func(t *testing.T) {
		wallet := Wallet{Bitcoin(20)}
		err := wallet.Withdraw(Bitcoin(100))

		assertBalance(t, wallet, Bitcoin(20))
		assertError(t, err, ErrInsufficientFunds)
	})
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Fatal("got an error but didn't want one")
	}
}
