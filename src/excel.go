package src

import (
	"fmt"
	"github.com/sghaida/fpv2/src/collections/list"
	"github.com/xuri/excelize/v2"
	"math"
	"sort"
	"time"
)

//type summaries struct {
//	contributions int
//	avgTime       float32
//	mdnTime       float32
//}

type ExcelOps struct {
	f         *excelize.File
	Comments  [][]interface{}
	Prs       [][]interface{}
	summaries [][]interface{}
}

func (e *ExcelOps) NewExcelFile() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheets.
	_, err := f.NewSheet("prs")
	if err != nil {
		return err
	}
	_, err = f.NewSheet("comments")
	if err != nil {
		return err
	}
	_, err = f.NewSheet("summaries")
	if err != nil {
		return err
	}
	// remove the default one
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		return err
	}

	prsData := [][]interface{}{
		{"Repo", "PR Owner", "status", "Created At", "Updated At", "Merged At", "UntilMerged"},
	}
	commentsData := [][]interface{}{
		{"Repo", "PR Owner", "Comment Owner", "Created At", "Updated At"},
	}

	prsSummaryData := [][]interface{}{
		{"Repo", "contributions", "AVG Waiting time", "MDN Waiting time"},
	}

	e.f = f
	e.Prs = prsData
	e.Comments = commentsData
	e.summaries = prsSummaryData

	return nil
}

func (e *ExcelOps) WriteFile(filepath string) error {
	if err := e.f.SaveAs(filepath); err != nil {
		return err
	}
	return nil
}

func (e *ExcelOps) AppendData(prs RepoPrs) error {

	var durations []float64

	for _, pr := range prs.Prs {

		createdAt := pr.CreatedAt.GetTime()
		creationDate := createdAt.Format("01/02/2006 15:04:05")

		updatedAt := pr.UpdatedAt.GetTime()
		updateDate := updatedAt.Format("01/02/2006 15:04:05")

		// cater for open prs
		if pr.MergedAt.Before(*createdAt) {

			hours := time.Since(*createdAt).Hours()
			hours = math.Round(hours)
			durations = append(durations, hours)

			e.Prs = append(e.Prs, []interface{}{
				prs.Repo, pr.OwnerName, "Open", creationDate, updateDate, nil, hours,
			})

		} else {

			hours := pr.MergedAt.Sub(*createdAt).Hours()
			hours = math.Round(hours)
			durations = append(durations, hours)

			mergedAt := pr.MergedAt.GetTime()
			MergeDate := mergedAt.Format("01/02/2006 15:04:05")

			e.Prs = append(e.Prs, []interface{}{
				prs.Repo, pr.OwnerName, "Close", creationDate, updateDate, MergeDate, hours,
			})
		}

		// add comments
		for _, comment := range pr.CommentInfo {

			createdAt := comment.CreatedAt.GetTime()
			creationDate := createdAt.Format("01/02/2006 15:04:05")

			updatedAt := comment.UpdatedAt.GetTime()
			updateDate := updatedAt.Format("01/02/2006 15:04:05")

			e.Comments = append(e.Comments, []interface{}{
				prs.Repo, pr.OwnerName, comment.OwnerName, creationDate, updateDate,
			})
		}

	}
	avg, mdn := e.calculateSummaries(durations)
	e.summaries = append(e.summaries, []interface{}{
		prs.Repo, len(durations), math.Round(avg), mdn,
	})

	return e.setSheetsData()
}

func (e *ExcelOps) calculateSummaries(durations []float64) (avg float64, mdn float64) {
	// calculate the summaries
	sort.Slice(durations, func(i, j int) bool { return i > j })
	// calculate median
	if len(durations) == 0 {
		mdn = 0
	} else if len(durations) == 1 {
		mdn = durations[0]
	} else {
		mod := len(durations) % 2
		if mod == 0 {
			center := len(durations) / 2
			mdn = durations[center]
		}
		if mod != 0 {
			center := (len(durations) + 1) / 2
			mdn = durations[center]
		}
	}

	// calculate the average
	sum := list.FoldLeft(durations, 0.0, func(acc float64, value float64) float64 {
		return acc + value
	})
	avg = sum / float64(len(durations))
	return
}

func (e *ExcelOps) setSheetsData() error {
	//  prs
	for idx, row := range e.Prs {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		err = e.f.SetSheetRow("prs", cell, &row)
		if err != nil {
			return err
		}
	}
	// comments
	for idx, row := range e.Comments {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		err = e.f.SetSheetRow("comments", cell, &row)
		if err != nil {
			return err
		}
	}

	// summaries
	for idx, row := range e.summaries {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			return err
		}
		err = e.f.SetSheetRow("summaries", cell, &row)
		if err != nil {
			return err
		}
	}
	return nil
}
