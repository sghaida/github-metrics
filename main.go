package main

import (
	"flag"
	"fmt"
	"github.com/sghaida/github-metrics/src"
	"golang.org/x/net/context"
	"log"
	"time"
)

var (
	outputFilePath string
	//from           = &time.Time{}
	//to             = &time.Time{}
	from string
	to   string
)

func init() {
	flag.StringVar(&outputFilePath, "out", "/tmp", "define where to dump the generate excel file")
	flag.StringVar(&from, "from", BeginningOfMonth(time.Now()).Format("2006-01-02"), "from date")
	flag.StringVar(&to, "to", EndOfMonth(time.Now()).Format("2006-01-02"), "to date")
	flag.Parse()
}

func main() {

	config := src.GetConfig()
	ac := src.NewClient(config.Token)
	client := ac.Create(context.Background())

	// build teams hierarchy
	contributors := make(map[string]src.SquadMember)
	for squad, teams := range config.Teams {
		for _, team := range teams {
			for _, ic := range team.Members {
				contributors[ic] = src.SquadMember{
					LoginName: ic,
					SquadName: squad,
					Team:      team.Name,
				}
			}
		}
	}

	from, err := time.Parse("2006-01-02", from)
	if err != nil {
		log.Fatalf("error parsing from date: %s", err.Error())
	}
	to, err := time.Parse("2006-01-02", to)
	if err != nil {
		log.Fatalf("error parsing to date: %s", err.Error())
	}

	prChan := make(chan src.RepoPrs, 10)
	prProcessor := src.NewPRProcessor(client, config, contributors, prChan)
	prProcessor.GetPrs(from, to)

	excel := src.ExcelOps{}
	err = excel.NewExcelFile()
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
	path := fmt.Sprintf("%s/metrics.xlsx", outputFilePath)
	if err := excel.WriteFile(path); err != nil {
		fmt.Println(err)
	}
}

func BeginningOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 0, -date.Day()+1)
}

func EndOfMonth(date time.Time) time.Time {
	return date.AddDate(0, 1, -date.Day())
}
