package main

import (
	"log"
	"os"

	"fyne.io/fyne/v2/app"
)

func main() {
	a := app.NewWithID("com.mewosh.banksight")

	originalStdout := os.Stdout
	originalStderr := os.Stderr
	reader, writer, err := os.Pipe()

	if err != nil {
		panic(err)
	}

	w := createMainWindow(a, reader, writer)
	w.ShowAndRun()

	cleanupTempFiles()
	writer.Close()
	os.Stdout = originalStdout
	os.Stderr = originalStderr
	log.SetOutput(originalStderr)
}
