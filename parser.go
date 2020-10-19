package main

import (
	"encoding/csv"
	"io"

	"github.com/shopspring/decimal"
)

type Product struct {
	Name  string
	Price decimal.Decimal
}

func ParseCSV(data io.ReadCloser) (products []Product, err error) {
	products = make([]Product, 0)
	defer data.Close()
	reader := csv.NewReader(data)
	reader.Comma = ';'
	for {
		values, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return
		}

		if len(values) < 2 {
			continue
		}

		price, err := decimal.NewFromString(values[1])
		if err != nil {
			continue
		}

		products = append(products, Product{
			Name:  values[0],
			Price: price,
		})

	}

}
