package entities

type OrderQueue struct {
	Orders []*Order
}

func (orderQueue *OrderQueue) Less(i, j int) bool {
	return orderQueue.Orders[i].Price < orderQueue.Orders[j].Price
}

func (orderQueue *OrderQueue) Swap(i, j int) {
	orderQueue.Orders[i], orderQueue.Orders[j] = orderQueue.Orders[j], orderQueue.Orders[i]
}

func (orderQueue *OrderQueue) Len() int {
	return len(orderQueue.Orders)
}

func (orderQueue *OrderQueue) Push(order any) {
	orderQueue.Orders = append(orderQueue.Orders, order.(*Order))
}

func (orderQueue *OrderQueue) Pop() any {
	oldValue := orderQueue.Orders
	orderQuantity := len(oldValue)
	item := oldValue[orderQuantity-1]
	orderQueue.Orders = oldValue[0 : orderQuantity-1]
	return item

}

func NewOrderQueue() *OrderQueue {
	return &OrderQueue{}
}
