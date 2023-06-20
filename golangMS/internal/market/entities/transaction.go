package entities

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID           string
	SellingOrder *Order
	BuyingOrder  *Order
	Shares       int
	Price        float64
	Total        float64
	DateTime     time.Time
}

func NewTransaction(sellingOrder *Order, buyingOrder *Order, shares int, price float64) *Transaction {
	total := float64(shares) * price

	return &Transaction{
		ID:           uuid.New().String(),
		SellingOrder: sellingOrder,
		BuyingOrder:  buyingOrder,
		Shares:       shares,
		Total:        total,
		DateTime:     time.Now(),
	}
}

func (transaction *Transaction) CalculateTotal(shares int, price float64) {
	transaction.Total = float64(shares) * price

}

func (transaction *Transaction) AddSellingOrderPendingShares(minShares int) {
	transaction.SellingOrder.PendingShares += minShares
}

func (transaction *Transaction) AddBuyingOrderPendingShares(minShares int) {
	transaction.BuyingOrder.PendingShares += minShares
}

func (transaction *Transaction) CloseBuyOrderTransaction() {
	if transaction.BuyingOrder.PendingShares == 0 {
		transaction.BuyingOrder.Status = "CLOSED"
	}
}

func (transaction *Transaction) CloseSellingOrderTransaction() {
	if transaction.SellingOrder.PendingShares == 0 {
		transaction.SellingOrder.Status = "CLOSED"
	}
}
