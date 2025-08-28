### Prerequisites

* latest [golang](https://go.dev/dl/) version (≥ 1.24.5)
* latest [pdftotext](https://poppler.freedesktop.org) version (≥ 25.07.0)

### Usage

`go run . -f "path_to_file"` or `go run . --file "path_to_file"`

### Flags

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
