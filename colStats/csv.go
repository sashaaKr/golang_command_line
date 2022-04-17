package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
)

type statsFunc func(data []float64) float64

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

func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true
	column--

	// allData, err := cr.ReadAll()
	// if err != nil {
	// 	return nil, fmt.Errorf("could not read csv: %w", err)
	// }

	var data []float64

	for i := 0; ; i++ {
		row, err := cr.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, fmt.Errorf("Cannot read dat from file: %w", err)
		}

		if i == 0 {
			continue
		}

		if len(row) <= column {
			return nil, fmt.Errorf("%w: file has only %d columns", ErrInvalidColumn, len(row))
		}

		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, row[column])
		}

		data = append(data, v)
	}


	// for i, row := range allData {
	// 	if i == 0 {
	// 		continue
	// 	}

	// 	if len(row) <= column {
	// 		return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
	// 	}

	// 	v, err := strconv.ParseFloat(row[column], 64)
	// 	if err != nil {
	// 		return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
	// 	}

	// 	data = append(data, v)
	// }

	return data, nil
}