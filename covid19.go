package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"text/template"
	"time"

	cli "github.com/jawher/mow.cli"
	dataframe "github.com/rocketlaunchr/dataframe-go"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Config holder for cli configs
type Config struct {
	Region string
}

// ItalyDataURL CSV data at national level
const ItalyDataURL = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv"

// RegionsDataURL CSV data at regional level
const RegionsDataURL = "https://github.com/pcm-dpc/COVID-19/raw/master/dati-regioni/dpc-covid19-ita-regioni.csv"

// ProvincesDataURL CSV data at provincial level
const ProvincesDataURL = "https://github.com/pcm-dpc/COVID-19/raw/master/dati-province/dpc-covid19-ita-province.csv"

// GlobalDataURL CSV data at international level
const GlobalDataURL = "https://data.humdata.org/hxlproxy/api/data-preview.csv?url=https%3A%2F%2Fraw.githubusercontent.com%2FCSSEGISandData%2FCOVID-19%2Fmaster%2Fcsse_covid_19_data%2Fcsse_covid_19_time_series%2Ftime_series_covid19_confirmed_global.csv&filename=time_series_covid19_confirmed_global.csv"

const timeFormat = "2006-01-02T15:04:05"

var cols = [...]string{"date", "ricoverati_con_sintomi", "totale_casi", "totale_positivi", "nuovi_positivi",
	"variazione_totale_positivi", "deceduti", "nuovi_decessi", "terapia_intensiva", "totale_ospedalizzati", "dimessi_guariti", "isolamento_domiciliare",
	"casi_da_sospetto_diagnostico", "casi_da_screening", "tamponi", "casi_testati"}

func main() {
	// here main
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard)

	var (
		app = cli.App("covid19", `
Daily stats for covid19 in Italy.

Available regions: Abruzzo, Basilicata, Calabria, Campania, Emilia-
  Romagna, 'Friuli Venezia Giulia', Lazio, Liguria, Lombardia, Marche,
  Molise, 'P.A. Bolzano', 'P.A. Trento', Piemonte, Puglia, Sardegna,
  Sicilia, Toscana, Umbria, "Valle d'Aosta", Veneto.`)
		cfg Config
	)

	app.Version("v version", "covid19 0.0.1")
	//TODO provincial data
	//TODO international data
	app.Spec = "[-r]"

	app.StringOptPtr(&cfg.Region, "r region", "", "Specify a region")

	app.Action = func() {
		mainAction(&cfg)
	}

	app.Run(os.Args)
}

func mainAction(cfg *Config) {
	fmt.Println("Covid19 - dati sintetici")

	var csv string
	if cfg.Region != "" {
		csv = getData(RegionsDataURL)
		fmt.Printf("regione selezionata: %v\n", cfg.Region)
	} else {
		csv = getData(ItalyDataURL)
	}

	df := loadDataFrame(csv)
	if cfg.Region != "" {
		filterByRegion(df, cfg.Region)
	}

	printSummary(df)
	printPercentages(df)
}

func filterByRegion(df *dataframe.DataFrame, region string) {
	var ctx = context.Background()
	filterFn := dataframe.FilterDataFrameFn(
		func(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
			if vals["denominazione_regione"] != region {
				return dataframe.DROP, nil
			}
			return dataframe.KEEP, nil
		})
	dataframe.Filter(ctx, df, filterFn, dataframe.FilterOptions{InPlace: true})
}

func printSummary(df *dataframe.DataFrame) {
	nrows := df.NRows()
	lastRow := df.Row(nrows-1, false, dataframe.SeriesName)
	secondLastRow := df.Row(nrows-2, false, dataframe.SeriesName)
	last2days := map[string]interface{}{
		"today":     lastRow,
		"yesterday": secondLastRow,
	}

	tmplStr := `
Aggiornamento: {{LocaleTimeFmt .today.data}} 		Oggi 			Ieri 	   	  Differenza (%)
-------------------------------------------------------------------------------------------------------

Totale casi: 		{{LocaleIntFmt .today.totale_casi}}	{{LocaleIntFmt .yesterday.totale_casi}}	{{sub .today.totale_casi .yesterday.totale_casi}}
Nuovi casi: 		{{LocaleIntFmt .today.nuovi_positivi}} 	{{LocaleIntFmt .yesterday.nuovi_positivi}}	{{sub .today.nuovi_positivi .yesterday.nuovi_positivi}}
Totale positivi: 	{{LocaleIntFmt .today.totale_positivi}} 	{{LocaleIntFmt .yesterday.totale_positivi}}	{{sub .today.totale_positivi .yesterday.totale_positivi}}
Variazione positivi: 	{{ LocaleIntFmt .today.variazione_totale_positivi }} 	{{ LocaleIntFmt .yesterday.variazione_totale_positivi }}	{{sub .today.variazione_totale_positivi .yesterday.variazione_totale_positivi}}
Totale decessi: 	{{LocaleIntFmt .today.deceduti}} 	{{LocaleIntFmt .yesterday.deceduti}}	{{sub .today.deceduti .yesterday.deceduti}}
Terapia intensiva: 	{{LocaleIntFmt .today.terapia_intensiva}} 	{{LocaleIntFmt .yesterday.terapia_intensiva}}	{{sub .today.terapia_intensiva .yesterday.terapia_intensiva}}
Ingressi in intensiva: 	{{LocaleIntFmt .today.ingressi_terapia_intensiva}} 	{{LocaleIntFmt .yesterday.ingressi_terapia_intensiva}}	{{sub .today.ingressi_terapia_intensiva .yesterday.ingressi_terapia_intensiva}}
Ospedalizzati: 		{{LocaleIntFmt .today.totale_ospedalizzati}} 	{{LocaleIntFmt .yesterday.totale_ospedalizzati}}	{{sub .today.totale_ospedalizzati .yesterday.totale_ospedalizzati}}
Dimessi: 		{{LocaleIntFmt .today.dimessi_guariti}} 	{{LocaleIntFmt .yesterday.dimessi_guariti}}	{{sub .today.dimessi_guariti .yesterday.dimessi_guariti}}
Totale tamponi: 	{{LocaleIntFmt .today.tamponi}} 	{{LocaleIntFmt .yesterday.tamponi}}	{{sub .today.tamponi .yesterday.tamponi}}
Totale testati: 	{{LocaleIntFmt .today.casi_testati}} 	{{LocaleIntFmt .yesterday.casi_testati}}	{{sub .today.casi_testati .yesterday.casi_testati}}
`
	p := message.NewPrinter(language.Italian)
	tmpl := template.Must(
		template.New("").Funcs(template.FuncMap{
			"LocaleIntFmt": func(val int64) string {
				return p.Sprintf("%20.0d", val)
			},
			"LocaleTimeFmt": func(val time.Time) string {
				return val.Format("02/01/2006")
			},
			"sub": func(a, b int64) string {
				variation := a - b
				pctVariation := float64(a-b) / float64(b) * 100
				return p.Sprintf("%20.0d", variation) + p.Sprintf(" ( %6.2f%%)", pctVariation)
			},
		}).Parse(tmplStr))

	if err := tmpl.Execute(os.Stdout, last2days); err != nil {
		fmt.Println(err)
	}
}

func printPercentages(df *dataframe.DataFrame) {
	nrows := df.NRows()
	lastRow := df.Row(nrows-1, false, dataframe.SeriesName)
	secondLastRow := df.Row(nrows-2, false, dataframe.SeriesName)

	calcPct := func(row map[interface{}]interface{}, field string) float64 {
		return 100 * float64(lastRow[field].(int64)) / float64(row["totale_casi"].(int64))
	}

	today := map[string]float64{
		"mortality":     calcPct(lastRow, "deceduti"),
		"intensiveCare": calcPct(lastRow, "terapia_intensiva"),
		"hospitalized":  calcPct(lastRow, "totale_ospedalizzati"),
		"recovered":     calcPct(lastRow, "dimessi_guariti"),
	}
	yesterday := map[string]float64{
		"mortality":     calcPct(secondLastRow, "deceduti"),
		"intensiveCare": calcPct(secondLastRow, "terapia_intensiva"),
		"hospitalized":  calcPct(secondLastRow, "totale_ospedalizzati"),
		"recovered":     calcPct(secondLastRow, "dimessi_guariti"),
	}

	last2days := map[string]interface{}{
		"today":     today,
		"yesterday": yesterday,
	}

	tmplStr := `

Mortalitá: 		{{printf "%19.2f" .today.mortality}}% 	{{printf "%19.2f" .yesterday.mortality}}% 	{{sub .today.mortality .yesterday.mortality}}
Terapia intensiva: 	{{printf "%19.2f" .today.intensiveCare}}% 	{{printf "%19.2f" .yesterday.intensiveCare}}% 	{{sub .today.intensiveCare .yesterday.intensiveCare}}
Ricoverati: 		{{printf "%19.2f" .today.hospitalized}}% 	{{printf "%19.2f" .yesterday.hospitalized}}% 	{{sub .today.hospitalized .yesterday.hospitalized}}
Guariti: 		{{printf "%19.2f" .today.recovered}}% 	{{printf "%19.2f" .yesterday.recovered}}% 	{{sub .today.recovered .yesterday.recovered}}
`
	tmpl := template.Must(template.New("").Funcs(template.FuncMap{
		"sub": func(a, b float64) string {
			variation := a - b
			pctVariation := (a - b) / b * 100
			return fmt.Sprintf("%19.2f%%", variation) + fmt.Sprintf(" ( %6.2f%%)", pctVariation)
		},
	}).Parse(tmplStr))

	if err := tmpl.Execute(os.Stdout, last2days); err != nil {
		fmt.Println(err)
	}

}
