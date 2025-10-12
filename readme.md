# BankSight

![BankSight Demo](https://github.com/user-attachments/assets/4b4f63cb-43a9-4472-a54f-dc142b3b1542 "BankSight Demo")

## Motivation

Since my main personal bank never adds any useful features and I was shocked that they even added icons next to the shops in the transactions menu, I had to program this basic functionality myself. Also I always wanted to learn Go.

## Prerequisites

* latest [golang](https://go.dev/dl/) version (≥ 1.24.5)
* latest [pdftotext](https://poppler.freedesktop.org) version (≥ 25.07.0) (if using PKO PDF statements)

## Supported banks and file formats

- **PKO BP**: PDF
- **ING**: MT940 (.sta)

## Usage

Download your bank statement, then run:

`go run . -f "path_to_file" -b "your_bank"` or `go run . --file "path_to_file" --bank "your_bank"`

## Flags

```
  -h help
  -b string
        Alias for -bank
  -bank string
        (Required) Name of the bank that generated the statement
  -d    Alias for -debug
  -debug
        (Optional) Enable debugging info
  -f string
        Alias for -file
  -file string
        (Required) Path to the file containing bank statement lines
  -o    Alias for -open-default-app (default true)
  -open-default-app
        (Optional) Open output file with the default application (default true)
```
