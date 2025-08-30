package entities

import (
	"fmt"
	"sort"
)

// Portfolio holds the current state of the investment portfolio
type (
	Portfolio struct {
		AssetClasses []AssetClass
		Assets       []Asset
	}
)

// Suggestion represents a recommendation for new investment
type Suggestion struct {
	AssetID   int
	AssetName string
	Amount    float64
}

// getClassByID returns the AssetClass by its ID
func (p *Portfolio) getClassByID(classID int) *AssetClass {
	for i := range p.AssetClasses {
		if p.AssetClasses[i].ID == classID {
			return &p.AssetClasses[i]
		}
	}
	return nil
}

// getAssetsByClass returns all assets belonging to a given class
func (p *Portfolio) getAssetsByClass(classID int) []Asset {
	var assets []Asset
	for _, a := range p.Assets {
		if a.ClassID == classID {
			assets = append(assets, a)
		}
	}
	return assets
}

// totalPortfolioValue calculates the total value of the portfolio
func (p *Portfolio) totalPortfolioValue() float64 {
	total := 0.0
	for _, a := range p.Assets {
		total += a.Quantity * a.UnitPrice
	}
	return total
}

// classCurrentValue calculates the current value of a class in the portfolio
func (p *Portfolio) classCurrentValue(classID int) float64 {
	total := 0.0
	for _, a := range p.Assets {
		if a.ClassID == classID {
			total += a.Quantity * a.UnitPrice
		}
	}
	return total
}

// rebalanceSuggestion calculates how to allocate a new contribution to maintain target allocation
func (p *Portfolio) RebalanceSuggestion(contribution float64) []Suggestion {
	actualTotalValue := p.totalPortfolioValue()
	fmt.Printf("Actual Total portfolio value: R$%.2f\n", actualTotalValue)

	totalValue := actualTotalValue + contribution

	fmt.Printf("Next totalportfolio value: R$%.2f\n", totalValue)

	// Calculate current and target values for each class
	type classBalance struct {
		Class      AssetClass
		CurrentVal float64
		TargetVal  float64
		ToInvest   float64
	}
	var balances []classBalance
	for _, class := range p.AssetClasses {
		current := p.classCurrentValue(class.ID)
		target := totalValue * class.TargetPct / 100.0
		toInvest := target - current
		if toInvest < 0 {
			toInvest = 0 // Do not suggest selling, only new investments
		}
		balances = append(balances, classBalance{
			Class:      class,
			CurrentVal: current,
			TargetVal:  target,
			ToInvest:   toInvest,
		})
	}

	// Sort classes by how much they are under target (descending)
	sort.Slice(balances, func(i, j int) bool {
		return balances[i].ToInvest > balances[j].ToInvest
	})

	remaining := contribution
	var suggestions []Suggestion

	for _, bal := range balances {
		if remaining <= 0 {
			break
		}
		alloc := bal.ToInvest
		if alloc > remaining {
			alloc = remaining
		}
		if alloc <= 0 {
			continue
		}

		// Within the class, allocate to the best scored asset(s)
		assets := p.getAssetsByClass(bal.Class.ID)
		if len(assets) == 0 {
			continue
		}
		// Sort assets by score descending
		sort.Slice(assets, func(i, j int) bool {
			return assets[i].Score > assets[j].Score
		})

		// Distribute allocation equally among all assets in the class to balance their quantities
		// TODO evolucionate the code to distribute the allocation based on the assets scores
		equalAlloc := alloc / float64(len(assets))
		for _, asset := range assets {
			suggestions = append(suggestions, Suggestion{
				AssetID:   asset.ID,
				AssetName: asset.Name,
				Amount:    equalAlloc,
			})
		}
		remaining -= alloc

	}

	return suggestions
}
