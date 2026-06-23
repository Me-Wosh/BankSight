package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"math"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func drawPieChart(transactions transactions) error {
	const (
		totalKey        = "total"
		regularFontSize = 14
	)

	timePeriod := transactions.timePeriod
	spendings := transactions.spendings
	incomes := transactions.incomes
	previousBalance := transactions.previousBalance
	closingBalance := transactions.closingBalance
	categorizedBalance := transactions.categorizedBalance
	difference := transactions.difference
	leftOver := incomes - math.Abs(spendings)

	spendingsIncomesData := []opts.PieData{
		{Name: "Spendings", Value: fmt.Sprintf("%.2f", spendings*-1), Tooltip: &opts.Tooltip{Show: opts.Bool(false)}},
		{Name: "Left over", Value: fmt.Sprintf("%.2f", leftOver), Tooltip: &opts.Tooltip{Show: opts.Bool(false)}},
	}

	var categorizedSpendingsData []opts.PieData
	categories := maps.Keys(categorizedBalance)

	for category := range categories {
		shops := maps.Keys(categorizedBalance[category])
		var tooltip strings.Builder

		for shop := range shops {
			if shop == totalKey {
				continue
			}

			tooltip.WriteString(fmt.Sprintf("%s: %.2f zł<br/>", shop, categorizedBalance[category][shop]))
		}

		tooltip.WriteString(fmt.Sprintf("%s: %.2f zł", totalKey, categorizedBalance[category][totalKey]))

		categorizedSpendingsData = append(categorizedSpendingsData, opts.PieData{
			Name:  category,
			Value: fmt.Sprintf("%.2f", categorizedBalance[category][totalKey]*-1),
			Tooltip: &opts.Tooltip{
				Formatter: types.FuncStr("<b>{b}</b><br/>" + tooltip.String()),
			},
		})
	}

	var sign string

	if difference > 0 {
		sign = "+"
	}

	pieChart := charts.NewPie()

	pieChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Total transactions",
			TitleStyle: &opts.TextStyle{
				FontSize: 20,
			},
			Subtitle: fmt.Sprintf(
				"Time period: %s\n\nPrevious balance: %.2f zł\nClosing balance: %.2f zł\nDifference: %s%.2f zł",
				timePeriod,
				previousBalance,
				closingBalance,
				sign,
				difference,
			),
			SubtitleStyle: &opts.TextStyle{
				FontSize:   15,
				LineHeight: 20,
			},
		}),
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{Bottom: "0", TextStyle: &opts.TextStyle{FontSize: regularFontSize}}),
	)

	pieChart.AddSeries(
		"",
		spendingsIncomesData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"0", "35%"}}),
		charts.WithLabelOpts(opts.Label{Position: "inside", Formatter: "{b}\n{d}%\n{c} zł", FontSize: regularFontSize}),
	)

	pieChart.AddSeries(
		"",
		categorizedSpendingsData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"50%", "70%"}}),
		charts.WithLabelOpts(opts.Label{Formatter: "{b}: {c} zł", FontSize: regularFontSize}),
	)

	page := components.NewPage()
	page.SetPageTitle(appName)
	page.SetLayout(components.PageFullLayout)
	page.AddCustomizedHeaders(`
		<style>
			body {
				margin: 0;
			}

			.container {
				align-items: center;
				box-sizing: border-box;
				display: flex;
				justify-content: center;
				padding: 20px;
			}
		</style>
	`)
	page.AddCharts(pieChart)

	file, err := createTrackedTempFile(tempFilePatternStatement)

	if err != nil {
		log.Println("Error while creating temporary HTML file: ", err)
		return err
	}

	if err := page.Render(io.MultiWriter(file)); err != nil {
		file.Close()
		removeTrackedTempFile(file.Name())
		log.Println("Error while rendering temporary HTML file: ", err)
		return err
	}

	if err := openFileWithDefaultApp(file.Name()); err != nil {
		removeTrackedTempFile(file.Name())
		log.Println("ERROR: could not open temporary HTML file: ", err)
		return err
	}

	return nil
}
