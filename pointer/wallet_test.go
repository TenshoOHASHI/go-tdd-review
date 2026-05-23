package pointer_test

import (
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

func (w *Wallet) Withdraw(amount Bitcoin) Bitcoin {
	w.balance -= amount

	return w.balance

}
func TestWallet(t *testing.T) {

	assertBalance := func(t *testing.T, wallet Wallet, want Bitcoin) {
		// 残高を取得
		got := wallet.Balance()

		if got != want {
			t.Errorf("got %s want %s", got, want)

		}
	}

	t.Run("Wallet", func(t *testing.T) {

		// コレクションを初期化
		wallet := Wallet{}
		// ポインターレシーバーなので、自動で&walletのアドレスを渡している
		err := wallet.Deposit(Bitcoin(10)) // 自身の新しい型で型変換

		if err != nil {
			t.Errorf("deposit failed: %v", err)
		}
		// fmtが最終的に呼ばれるため、%sにBitcoinのString()が呼ばれます。
		// got.String()
		//
		assertBalance(t, wallet, Bitcoin(10))

	})

	t.Run("Withdraw", func(t *testing.T) {
		// 預金に20追加し、初期化
		wallet := Wallet{balance: Bitcoin(20)}
		// 預金から引き出す
		wallet.Withdraw(Bitcoin(10))

		assertBalance(t, wallet, Bitcoin(10))

	})
}
