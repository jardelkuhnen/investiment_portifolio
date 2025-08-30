package main

import (
	"fmt"

	"github.com/jardelkuhnen/investiment_portifolio/database"
)

func main() {
	actualPortifolio := database.GetActualPortfolio()

	// Simulate a new aporte de R$4000,00
	contribution := 4000.0
	fmt.Printf("New contribution: R$%.2f\n", contribution)
	fmt.Println("--------------------------------")

	suggestions := actualPortifolio.RebalanceSuggestion(contribution)

	// Group allocation by class
	classAllocations, classNames := actualPortifolio.GroupAllocationByClass(suggestions, actualPortifolio)

	fmt.Println("--------------------------------")

	fmt.Println("Allocation By Class:")
	for classID, amount := range classAllocations {
		fmt.Printf("- %s: R$%.2f\n", classNames[classID], amount)
	}

	fmt.Println("--------------------------------")

	fmt.Println("Suggested allocation to maintain target balance:")
	for _, s := range suggestions {
		fmt.Printf("- Invest R$%.2f in %s\n", s.Amount, s.AssetName)
	}
}
