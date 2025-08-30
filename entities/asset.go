package entities

// AssetClass represents a class of assets (e.g., Fixed Income, Stocks, Crypto)
type AssetClass struct {
	ID        int
	Name      string
	TargetPct float64 // Target percentage in the portfolio (0-100)
}

// Asset represents an individual investment asset
type Asset struct {
	ID        int
	Name      string
	ClassID   int
	Score     int     // 1-10, 10 is best
	Quantity  float64 // Current quantity held
	UnitPrice float64 // Current price per unit
}
