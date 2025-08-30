package database

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/jardelkuhnen/investiment_portifolio/entities"
)

func readAssetClassesFromCSV(path string) ([]entities.AssetClass, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var classes []entities.AssetClass
	for i, rec := range records {
		if i == 0 {
			// skip header
			continue
		}
		if len(rec) < 3 {
			continue
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			continue
		}
		targetPct, err := strconv.ParseFloat(rec[2], 64)
		if err != nil {
			continue
		}
		classes = append(classes, entities.AssetClass{
			ID:        id,
			Name:      rec[1],
			TargetPct: targetPct,
		})
	}
	return classes, nil
}

func readAssetsFromCSV(path string) ([]entities.Asset, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var assets []entities.Asset
	for i, rec := range records {
		if i == 0 {
			// skip header
			continue
		}
		if len(rec) < 6 {
			continue
		}
		id, err := strconv.Atoi(rec[0])
		if err != nil {
			continue
		}
		classID, err := strconv.Atoi(rec[1])
		if err != nil {
			continue
		}
		score, err := strconv.Atoi(rec[3])
		if err != nil {
			continue
		}
		quantity, err := strconv.ParseFloat(rec[4], 64)
		if err != nil {
			continue
		}
		unitPrice, err := strconv.ParseFloat(rec[5], 64)
		if err != nil {
			continue
		}
		assets = append(assets, entities.Asset{
			ID:        id,
			ClassID:   classID,
			Name:      rec[2],
			Score:     score,
			Quantity:  quantity,
			UnitPrice: unitPrice,
		})
	}
	return assets, nil
}

func GetActualPortfolio() entities.Portfolio {
	classes, err := readAssetClassesFromCSV(os.Getenv("CLASSES_PATH"))
	if err != nil {
		panic("Failed to read asset classes: " + err.Error())
	}
	assets, err := readAssetsFromCSV(os.Getenv("ASSETS_PATH"))
	if err != nil {
		panic("Failed to read assets: " + err.Error())
	}

	return entities.Portfolio{
		AssetClasses: classes,
		Assets:       assets,
	}
}
