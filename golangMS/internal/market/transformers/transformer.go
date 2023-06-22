package transformers

import (
	"github.com/luizeduu/imersao13/go/internal/dto"
	"github.com/luizeduu/imersao13/go/internal/market/entities"
)

func TransformInput(input dto.TradeInput) *entities.Order {
	asset := entities.NewAsset(input.AssetID, input.AssetID, 1000)
	investor := entities.NewInvestor(input.InverstorID)
	order := entities.NewOrder(input.OrderId, investor, asset, input.Shares, input.Price, input.OrderType)

	if input.CurrentShares > 0 {
		assetPosition := entities.NewInvestorAssetPosition(input.AssetID, input.CurrentShares)
		investor.AddAssetPosition(assetPosition)
	}

	return order
}

func TransformOutput(order *entities.Order) *dto.OrderOutput {
	output := &dto.OrderOutput{
		OrderId:     order.ID,
		InverstorID: order.Investor.ID,
		AssetID:     order.Asset.ID,
		OrderType:   order.OrderType,
		Status:      order.Status,
		Partial:     order.PendingShares,
		Shares:      order.Shares,
	}

	var transactionsOutput []*dto.TransactionOutput

	for _, value := range order.Transactions {
		transactionOutput := &dto.TransactionOutput{
			TransactionID: value.ID,
			BuyerID:       value.BuyingOrder.ID,
			SellerID:      value.SellingOrder.ID,
			AssetID:       value.SellingOrder.Asset.ID,
			Price:         value.Price,
			Shares:        value.SellingOrder.Shares,
		}

		transactionsOutput = append(transactionsOutput, transactionOutput)
	}

	output.TransactionOutput = transactionsOutput
	return output
}
