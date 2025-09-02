package main

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

func drawPieChart(transactions transactions, subtitle, filePath string) {
	const (
		totalKey        = "total"
		regularFontSize = 14
	)

	spendings := transactions.spendings
	incomes := transactions.incomes
	previousBalance := transactions.previousBalance
	closingBalance := transactions.closingBalance
	categoriesBalance := transactions.categoriesBalance

	spendingsIncomesData := []opts.PieData{
		{Name: "Spendings", Value: fmt.Sprintf("%.2f", spendings*-1), Tooltip: &opts.Tooltip{Show: opts.Bool(false)}},
		{Name: "Incomes", Value: fmt.Sprintf("%.2f", incomes), Tooltip: &opts.Tooltip{Show: opts.Bool(false)}},
	}

	var categorizedSpendingsData []opts.PieData
	categories := maps.Keys(categoriesBalance)

	for category := range categories {
		shops := maps.Keys(categoriesBalance[category])
		var tooltip strings.Builder

		for shop := range shops {
			if shop == totalKey {
				continue
			}

			tooltip.WriteString(fmt.Sprintf("%s: %.2f zł<br/>", shop, categoriesBalance[category][shop]))
		}

		tooltip.WriteString(fmt.Sprintf("%s: %.2f zł", totalKey, categoriesBalance[category][totalKey]))

		categorizedSpendingsData = append(categorizedSpendingsData, opts.PieData{
			Name:  category,
			Value: fmt.Sprintf("%.2f", categoriesBalance[category][totalKey]*-1),
			Tooltip: &opts.Tooltip{
				Formatter: types.FuncStr("<b>{b}</b><br/>" + tooltip.String()),
			},
		})
	}

	pieChart := charts.NewPie()

	pieChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Total transactions",
			TitleStyle: &opts.TextStyle{
				FontSize: 20,
			},
			Subtitle: fmt.Sprintf(
				"Time period: %s\n\nPrevious balance: %.2f zł\nClosing balance: %.2f zł",
				subtitle,
				previousBalance,
				closingBalance,
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
	page.AddCustomizedCSSAssets(cssFilePath)
	page.AddCharts(pieChart)

	file, err := os.Create(strings.Replace(filePath, ".txt", ".html", 1))

	if err != nil {
		panic(fmt.Sprintf("Error while creating HTML file: %v", err))
	}

	if err := page.Render(io.MultiWriter(file)); err != nil {
		panic(fmt.Sprintf("Error while rendering HTML file: %v", err))
	}

	if err := openFileWithDefaultApp(file.Name()); err != nil {
		fmt.Println("WARNING: could not open HTML file: ", err)
	}
}
