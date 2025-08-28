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
)

func drawPieChart(spendings, incomes float64, categoriesBalance map[string]float64, subtitle, fileName string) {
	spendingsIncomesData := []opts.PieData{
		{Name: "Spendings", Value: fmt.Sprintf("%.2f", spendings*-1)},
		{Name: "Incomes", Value: fmt.Sprintf("%.2f", incomes)},
	}

	var categorizedSpendingsData []opts.PieData
	categories := maps.Keys(categoriesBalance)
	for category := range categories {
		categorizedSpendingsData = append(categorizedSpendingsData, opts.PieData{
			Name:  category,
			Value: fmt.Sprintf("%.2f", categoriesBalance[category]*-1),
		})
	}

	pieChart := charts.NewPie()

	pieChart.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Total transactions", Subtitle: "Time period: " + subtitle}),
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{Bottom: "0"}),
	)

	pieChart.AddSeries(
		"",
		spendingsIncomesData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"0", "35%"}}),
		charts.WithLabelOpts(opts.Label{Position: "inside", Formatter: "{b}:\n{c} zł"}),
	)

	pieChart.AddSeries(
		"",
		categorizedSpendingsData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"50%", "70%"}}),
		charts.WithLabelOpts(opts.Label{Formatter: "{b}: {c} zł"}),
	)

	page := components.NewPage()
	page.SetPageTitle(appName)
	page.SetLayout(components.PageFullLayout)
	page.AddCustomizedCSSAssets(cssFilePath)
	page.AddCharts(pieChart)

	file, err := os.Create(strings.Replace(fileName, ".pdf", ".html", 1))

	if err != nil {
		panic(fmt.Sprintf("Error while creating HTML file: %v", err))
	}

	if err := page.Render(io.MultiWriter(file)); err != nil {
		panic(fmt.Sprintf("Error while rendering HTML file: %v", err))
	}
}
