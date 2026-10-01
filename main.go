package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/Me-Wosh/BankSight/internal/file"
	"github.com/Me-Wosh/BankSight/internal/input"
	"github.com/Me-Wosh/BankSight/internal/ui"
)

func main() {
	a := app.New()

	originalStdout := os.Stdout
	originalStderr := os.Stderr
	reader, writer, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	shopCategories, err := input.ParseShopCategories(resourceShopcategoriesYaml.StaticContent)

	if err != nil {
		panic(err)
	}

	w := ui.CreateMainWindow(a, reader, writer, shopCategories)
	w.ShowAndRun()

	file.CleanupTempFiles()
	writer.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	log.SetOutput(originalStderr)
}
