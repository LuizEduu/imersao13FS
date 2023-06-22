package entities

import (
	"container/heap"
	"sync"
)

type Book struct {
	Orders              []*Order
	Transactions        []*Transaction
	OrdersChannel       chan *Order
	OrdersChannelOutput chan *Order
	WaitGroup           *sync.WaitGroup
}

func NewBook(orderChannel chan *Order, orderChannelOutput chan *Order, waitGroup *sync.WaitGroup) *Book {
	return &Book{
		Orders:              []*Order{},
		Transactions:        []*Transaction{},
		OrdersChannel:       orderChannel,
		OrdersChannelOutput: orderChannelOutput,
		WaitGroup:           waitGroup,
	}
}

func (book *Book) Trade() {
	buyOrders := make(map[string]*OrderQueue)
	sellOrders := make(map[string]*OrderQueue)

	for order := range book.OrdersChannel {

		asset := order.Asset.ID

		if buyOrders[asset] == nil {
			buyOrders[asset] = NewOrderQueue()
			heap.Init(buyOrders[asset])
		}

		if sellOrders[asset] == nil {
			sellOrders[asset] = NewOrderQueue()
			heap.Init(sellOrders[asset])
		}

		if order.OrderType == "BUY" {
			buyOrders[asset].Push(order)

			if sellOrders[asset].Len() > 0 && sellOrders[asset].Orders[0].Price <= order.Price {
				sellOrder := sellOrders[asset].Pop().(*Order)

				if sellOrder.PendingShares > 0 {
					transaction := NewTransaction(sellOrder, order, order.Shares, sellOrder.Price)

					book.AddTransaction(transaction, book.WaitGroup)
					sellOrder.Transactions = append(sellOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)
					book.OrdersChannelOutput <- sellOrder
					book.OrdersChannelOutput <- order

					if sellOrder.PendingShares > 0 {
						sellOrders[asset].Push(sellOrder)
					}
				}
			}
		} else if order.OrderType == "SELL" {
			sellOrders[asset].Push(order)

			if buyOrders[asset].Len() > 0 && buyOrders[asset].Orders[0].Price >= order.Price {
				buyOrder := buyOrders[asset].Pop().(*Order)

				if buyOrder.PendingShares > 0 {
					transaction := NewTransaction(order, buyOrder, order.Shares, buyOrder.Price)
					book.AddTransaction(transaction, book.WaitGroup)
					buyOrder.Transactions = append(buyOrder.Transactions, transaction)
					order.Transactions = append(order.Transactions, transaction)

					book.OrdersChannelOutput <- buyOrder
					book.OrdersChannelOutput <- order

					if buyOrder.PendingShares > 0 {
						buyOrders[asset].Push(buyOrder)
					}
				}
			}
		}
	}
}

func (book *Book) AddTransaction(transaction *Transaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	sellingShares := transaction.SellingOrder.PendingShares
	buyingShares := transaction.BuyingOrder.PendingShares

	minShares := sellingShares
	if buyingShares < sellingShares {
		minShares = buyingShares
	}

	transaction.SellingOrder.Investor.UpdateAssetPosition(transaction.SellingOrder.Asset.ID, -minShares)
	transaction.AddSellingOrderPendingShares(-minShares)
	transaction.BuyingOrder.Investor.UpdateAssetPosition(transaction.BuyingOrder.Asset.ID, minShares)
	transaction.AddBuyingOrderPendingShares(-minShares)

	transaction.CalculateTotal(transaction.Shares, transaction.BuyingOrder.Price)

	transaction.CloseBuyOrderTransaction()
	transaction.CloseSellingOrderTransaction()

	book.Transactions = append(book.Transactions, transaction)
}
