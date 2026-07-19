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

func drawPieChart(transactions transactions, isDarkTheme bool) error {
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
	dataColors := []string{
		// Original palette
		"#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de",
		"#3ba272", "#fc8452", "#9a60b4", "#ea7ccc",

		"#11c8dc", "#734cb7", "#e14f22", "#10a246", "#e33971",
	}
	i := 0

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
		colorIndex := i % len(dataColors) // Wraps around the color palette if there are more categories than colors

		categorizedSpendingsData = append(categorizedSpendingsData, opts.PieData{
			Name:  category,
			Value: fmt.Sprintf("%.2f", categorizedBalance[category][totalKey]*-1),
			Tooltip: &opts.Tooltip{
				Formatter: types.FuncStr("<b>{b}</b><br/>" + tooltip.String()),
			},
			ItemStyle: &opts.ItemStyle{
				Color: dataColors[colorIndex],
			},
		})

		i++
	}

	var sign string

	if difference > 0 {
		sign = "+"
	}

	pieChart := charts.NewPie()

	theme := "white"

	if isDarkTheme {
		theme = "dark"
	}

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
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%", Theme: theme}),
		charts.WithLegendOpts(opts.Legend{Bottom: "0", TextStyle: &opts.TextStyle{FontSize: regularFontSize}}),
	)

	spendingsToIncomesRatio := (spendings * -1) / incomes * 100

	if spendingsToIncomesRatio > 100 {
		pieChart.AddSeries(
			"",
			spendingsIncomesData,
			charts.WithPieChartOpts(opts.PieChart{Radius: []string{"0", "35%"}}),
			charts.WithLabelOpts(opts.Label{
				Position:  "inside",
				Formatter: fmt.Sprintf("{b}\n%.1f%%\n{c} zł", spendingsToIncomesRatio),
				FontSize:  regularFontSize,
			}),
		)
	} else {
		pieChart.AddSeries(
			"",
			spendingsIncomesData,
			charts.WithPieChartOpts(opts.PieChart{Radius: []string{"0", "35%"}}),
			charts.WithLabelOpts(opts.Label{
				Position:  "inside",
				Formatter: "{b}\n{d}%\n{c} zł",
				FontSize:  regularFontSize,
			}),
		)
	}

	pieChart.AddSeries(
		"",
		categorizedSpendingsData,
		charts.WithPieChartOpts(opts.PieChart{Radius: []string{"50%", "70%"}}),
		charts.WithLabelOpts(opts.Label{Formatter: "{b}: {c} zł", FontSize: regularFontSize}),
	)

	pieChart.AddJSFuncs(`
		window.addEventListener('resize', (event) => {
			const chartContainer = document.querySelector('.item');
			const chartId = chartContainer.getAttribute('_echarts_instance_');
			const chartInstance = window.echarts.getInstanceById(chartId);
			if (chartInstance) {
				chartInstance.resize();
			}
		});
	`)

	backgroundColor := "rgb(255, 255, 255)"

	if isDarkTheme {
		backgroundColor = "rgb(16, 12, 42)"
	}

	page := components.NewPage()
	page.SetPageTitle(appName)
	page.SetLayout(components.PageFullLayout)
	page.AddCustomizedHeaders(`
		<link rel="icon" type="image/x-icon" href="data:image/x-icon;base64,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"/>
	`)
	page.AddCustomizedHeaders(fmt.Sprintf(`
		<style>
			body {
				margin: 0;
				background-color: %s;
			}

			.container {
				align-items: center;
				box-sizing: border-box;
				display: flex;
				justify-content: center;
				padding: 50px;
			}
		</style>
	`, backgroundColor))
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
