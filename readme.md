# BankSight

![BankSight Demo](https://github.com/user-attachments/assets/4b4f63cb-43a9-4472-a54f-dc142b3b1542 "BankSight Demo")

## Motivation

Since my main personal bank never adds any useful features and I was shocked that they even added icons next to the shops in the transactions menu, I had to program this basic functionality myself. Also I always wanted to learn Go.

## Prerequisites

* latest [golang](https://go.dev/dl/) version (≥ 1.24.5)
* latest [pdftotext](https://poppler.freedesktop.org) version (≥ 25.07.0)

## Supported banks

- PKO BP

## Usage

Download your bank statement as a PDF file, then run:

`go run . -f "path_to_file"` or `go run . --file "path_to_file"`

## Flags

```
  -h    Help
  -d    Alias for -debug
  -debug
        (Optional) Enable debugging info
  -f string
        Alias for -file
  -file string
        (Required) Path to the file containing bank statement lines
```
