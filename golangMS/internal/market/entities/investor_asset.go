package entities

type InvestorAssetPosition struct {
	ID     string
	Shares int
}

func NewInvestorAssetPosition(id string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		ID:     id,
		Shares: shares,
	}
}
