package main

import (
	"os"

	"github.com/gocarina/gocsv"
)

func loadFromCSV(path string) ([]*Order, error) {
	orders := make([]*Order, count)

	if csvFile, err := os.OpenFile(path, os.O_RDONLY, os.ModePerm); err != nil {
		return nil, err
	} else if err := gocsv.UnmarshalFile(csvFile, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}

func saveToCSV(path string, orders []*Order) error {
	return nil
}
