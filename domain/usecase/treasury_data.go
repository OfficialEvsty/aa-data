package usecase

type TreasuryItem struct {
	ItemID   uint64 `json:"item_id"`
	Name     string `json:"name"`
	IconURL  string `json:"icon_url"`
	Quantity uint64 `json:"quantity"`
}

type TreasuryItemList struct {
	Items                 map[uint64]*TreasuryItem
	TreasuryItemsQuantity uint64
	FierceEssenceQuantity uint64
	TrophyEssenceQuantity uint64
}
