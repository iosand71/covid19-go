# CLI tool per dati sintetici covid19

Esercizio di sviluppo con Go

## Installazione

- `go build`
- `go install`

## Cross compilazione

É molto semplice produrre binari per altri os/architetture: é sufficiente impostare
due variabili di ambiente, ad ex:

- **MacOS**: `env GOOS=darwin GOARCH=amd64 go build`
- **Linux**: `env GOOS=linux GOARCH=amd64 go build`
- **Windows**: `env GOOS=windows GOARCH=amd64 go build`

## Riga di comando

```
Usage: covid19 [-r [-a]] [-d]

Daily stats for covid19 in Italy.

Options:
  -v, --version      Show the version and exit
  -r, --region       Specify a region
  -a, --availables   Print available regions
  -d, --date         Date in yyyy-mm-dd format (default 0001-01-01)
```

