# BankSight

![BankSight Demo](https://github.com/user-attachments/assets/7eb1e5a9-913d-47c7-8bc2-53bc1385744f "BankSight Demo")

![Generated Insights Demo](https://github.com/user-attachments/assets/6cdf3f69-6570-4698-9eb9-466f786c61f6 "Generated Insights Demo")


## Motivation

My main personal bank rarely adds useful features (I was genuinely surprised when they managed to put store icons next to transactions) so I decided to build the account balance tracking functionality myself. Plus, it gave me a great excuse to learn Go.

## Prerequisites

* latest [golang](https://go.dev/dl/) version (≥ 1.24.5) (if building yourself)
* latest [pdftotext](https://poppler.freedesktop.org) version (≥ 25.07.0) (if using PKO PDF statements)

## Supported banks and file formats

- **PKO BP**: PDF
- **ING**: MT940 (.sta)
- **mBank**: CSV

## Usage

1. Download the app from Releases tab or, if you want to build it yourself, run `go run .` inside root folder
2. Press appropriate radio button to select your bank
3. Upload your bank statement either by pressing "Upload (or drop)" button or by dragging and dropping the statement file onto the app window
4. Press "Get insights" button
5. The browser tab should now open with transactions analysis
