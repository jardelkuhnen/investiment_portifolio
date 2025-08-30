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
		// First, calculate the total amount needed to bring all classes to their targets (sum of ToInvest)
		var totalToInvest float64
		for _, b := range balances {
			totalToInvest += b.ToInvest
		}
		// If totalToInvest is zero (all classes at or above target), just skip
		if totalToInvest == 0 {
			break
		}
		// Calculate the percentage of allocation for this class based on its ToInvest
		classAlloc := alloc
		if totalToInvest > 0 {
			classAlloc = (bal.ToInvest / totalToInvest) * remaining
			if classAlloc > alloc {
				classAlloc = alloc // Don't allocate more than this class needs
			}
		}
		// Distribute classAlloc equally among assets in the class (could be improved to use scores)
		equalAlloc := classAlloc / float64(len(assets))
		for _, asset := range assets {
			suggestions = append(suggestions, Suggestion{
				AssetID:   asset.ID,
				AssetName: asset.Name,
				Amount:    equalAlloc,
			})
		}
		remaining -= classAlloc

	}

	return suggestions
}

func (p *Portfolio) GroupAllocationByClass(suggestions []Suggestion, actualPortifolio Portfolio) (map[int]float64, map[int]string) {
	// Group allocation by class
	classAllocations := make(map[int]float64)
	classNames := make(map[int]string)
	for _, s := range suggestions {
		// Find the class of this asset
		for _, asset := range actualPortifolio.Assets {
			if asset.ID == s.AssetID {
				classAllocations[asset.ClassID] += s.Amount
				// Save class name for output
				if _, ok := classNames[asset.ClassID]; !ok {
					for _, class := range actualPortifolio.AssetClasses {
						if class.ID == asset.ClassID {
							classNames[asset.ClassID] = class.Name
							break
						}
					}
				}
				break
			}
		}
	}

	return classAllocations, classNames
}
