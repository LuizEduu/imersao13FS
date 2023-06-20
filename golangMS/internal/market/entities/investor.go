package entities

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition //com o ponteiro, qualquer valor que seja alterado no slice vai ser refletido e alterado na memória
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (investor *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	investor.AssetPosition = append(investor.AssetPosition, assetPosition)
}

func (investor *Investor) UpdateAssetPosition(assetId string, shares int) {
	assetPosition := investor.GetAssetPosition(assetId)

	if assetPosition != nil {
		assetPosition.Shares += shares
	} else {
		investor.AddNewAssetPosition(assetId, shares)
	}
}

func (investor *Investor) GetAssetPosition(assetId string) *InvestorAssetPosition {
	for _, assetPosition := range investor.AssetPosition {
		if assetPosition.ID == assetId {
			return assetPosition
		}

	}

	return nil
}

func (investor *Investor) AddNewAssetPosition(assetId string, shares int) {
	investor.AssetPosition = append(investor.AssetPosition, NewInvestorAssetPosition(assetId, shares))
}
