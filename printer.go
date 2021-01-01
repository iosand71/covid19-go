package main

import (
	"context"
	"fmt"
	"os"
	"text/template"
	"time"

	dataframe "github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/utils"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func printAvailableRegions(df *dataframe.DataFrame) {
	var denominazioneRegione = df.MustNameToColumn("denominazione_regione")
	iter := df.Series[denominazioneRegione].ValuesIterator()
	set := make(map[string]bool)

	fmt.Println("\navailable regions:")
	for {
		row, vals, _ := iter()
		if row == nil {
			break
		}
		set[vals.(string)] = true
	}

	for key := range set {
		fmt.Printf("%v, ", key)
	}
	fmt.Println()
}

func getLastDays(df *dataframe.DataFrame) (map[interface{}]interface{}, map[interface{}]interface{}, map[interface{}]interface{}) {
	ctx := context.Background()

	nrows := df.NRows() - 1
	dateColumn := df.MustNameToColumn("data")
	startDate := time.Time(cfg.startDate)
	maxResults := 1

	if !cfg.startDate.isZero() {
		found, _ := utils.Search(ctx, df.Series[dateColumn], startDate, startDate.Add(time.Hour*24), utils.SearchOptions{Max: &maxResults})
		if len(found) == 0 {
			fmt.Println("Nessun dato alla data specificata")
			os.Exit(1)
		}
		nrows = found[0]
	}

	lastRow := df.Row(nrows, false, dataframe.SeriesName)
	secondLastRow := df.Row(nrows-1, false, dataframe.SeriesName)
	thirdLastRow := df.Row(nrows-2, false, dataframe.SeriesName)

	return lastRow, secondLastRow, thirdLastRow
}

func printSummary(df *dataframe.DataFrame) {
	lastRow, secondLastRow, thirdLastRow := getLastDays(df)
	last2days := map[string]interface{}{
		"today":        lastRow,
		"yesterday":    secondLastRow,
		"threeDaysAgo": thirdLastRow,
	}

	tmplStr := `
Aggiornamento: {{LocaleTimeFmt .today.data}}             Ultimi              Precedenti 	   	  Differenza (%)
-------------------------------------------------------------------------------------------------------

Totale casi: 		{{LocaleIntFmt .today.totale_casi}}	{{LocaleIntFmt .yesterday.totale_casi}}	{{pctVar .today.totale_casi .yesterday.totale_casi}}
Nuovi casi: 		{{LocaleIntFmt .today.nuovi_positivi}} 	{{LocaleIntFmt .yesterday.nuovi_positivi}}	{{pctVar .today.nuovi_positivi .yesterday.nuovi_positivi}}
Totale positivi: 	{{LocaleIntFmt .today.totale_positivi}} 	{{LocaleIntFmt .yesterday.totale_positivi}}	{{pctVar .today.totale_positivi .yesterday.totale_positivi}}
Variazione positivi: 	{{ LocaleIntFmt .today.variazione_totale_positivi }} 	{{ LocaleIntFmt .yesterday.variazione_totale_positivi }}	{{pctVar .today.variazione_totale_positivi .yesterday.variazione_totale_positivi}}
Totale decessi: 	{{LocaleIntFmt .today.deceduti}} 	{{LocaleIntFmt .yesterday.deceduti}}	{{pctVar .today.deceduti .yesterday.deceduti}}
Variazione decessi: 	{{sub .today.deceduti .yesterday.deceduti}} 	{{sub .yesterday.deceduti .threeDaysAgo.deceduti}}
Terapia intensiva: 	{{LocaleIntFmt .today.terapia_intensiva}} 	{{LocaleIntFmt .yesterday.terapia_intensiva}}	{{pctVar .today.terapia_intensiva .yesterday.terapia_intensiva}}
Ingressi in intensiva: 	{{LocaleIntFmt .today.ingressi_terapia_intensiva}} 	{{LocaleIntFmt .yesterday.ingressi_terapia_intensiva}}	{{pctVar .today.ingressi_terapia_intensiva .yesterday.ingressi_terapia_intensiva}}
Ospedalizzati: 		{{LocaleIntFmt .today.totale_ospedalizzati}} 	{{LocaleIntFmt .yesterday.totale_ospedalizzati}}	{{pctVar .today.totale_ospedalizzati .yesterday.totale_ospedalizzati}}
Dimessi: 		{{LocaleIntFmt .today.dimessi_guariti}} 	{{LocaleIntFmt .yesterday.dimessi_guariti}}	{{pctVar .today.dimessi_guariti .yesterday.dimessi_guariti}}
Totale tamponi: 	{{LocaleIntFmt .today.tamponi}} 	{{LocaleIntFmt .yesterday.tamponi}}	{{pctVar .today.tamponi .yesterday.tamponi}}
Totale testati: 	{{LocaleIntFmt .today.casi_testati}} 	{{LocaleIntFmt .yesterday.casi_testati}}	{{pctVar .today.casi_testati .yesterday.casi_testati}}
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
				return p.Sprintf("%20.0d", a-b)
			},
			"pctVar": func(a, b int64) string {
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
	lastRow, secondLastRow, _ := getLastDays(df)

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
