package main

import (
	"fmt"
	"github.com/sghaida/github-metrics/src"
	"golang.org/x/net/context"
	"time"
)

func main() {
	config := src.GetConfig()
	ac := src.NewClient(config.Token)
	client := ac.Create(context.Background())

	// Specify the date range
	from := time.Date(2023, time.May, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2023, time.May, 31, 23, 59, 59, 0, time.UTC)

	prChan := make(chan src.RepoPrs, 10)
	prProcessor := src.NewPRProcessor(client, config, prChan)

	prProcessor.GetPrs(from, to)

	excel := src.ExcelOps{}
	err := excel.NewExcelFile()
	if err != nil {
		fmt.Println(err)
	}

	for repoData := range prChan {
		err := excel.AppendData(repoData)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if err := excel.WriteFile("metrics.xlsx"); err != nil {
		fmt.Println(err)
	}
}
