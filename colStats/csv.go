package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

// defines a generic statistical function
type statsFunc func(data []float64) float64

// column: starting from 1, as it's more natural for users to understand
func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)

	// adjusting for 0 based index
	column--

	allData, err := cr.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("cannot read data from file: %w", err)
	}

	var data []float64
	for i, row := range allData {
		if i == 0 {
			continue
		}
		if len(row) <= column {
			return nil,
				fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}
		data = append(data, v)
	}

	return data, nil
}
