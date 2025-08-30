package main

import (
	"fmt"

	"github.com/jardelkuhnen/investiment_portifolio/database"
)

func main() {
	actualPortifolio := database.GetActualPortfolio()

	// Simulate a new aporte de R$4000,00
	contribution := 40000.0
	fmt.Printf("New contribution: R$%.2f\n", contribution)
	fmt.Println("--------------------------------")

	suggestions := actualPortifolio.RebalanceSuggestion(contribution)

	fmt.Println("--------------------------------")
	fmt.Println("Suggested allocation to maintain target balance:")
	for _, s := range suggestions {
		fmt.Printf("- Invest R$%.2f in %s\n", s.Amount, s.AssetName)
	}
}
