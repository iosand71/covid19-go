package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	cli "github.com/jawher/mow.cli"
)

const (
	// ItalyDataURL CSV data at national level
	ItalyDataURL = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-andamento-nazionale/dpc-covid19-ita-andamento-nazionale.csv"
	// RegionsDataURL CSV data at regional level
	RegionsDataURL = "https://github.com/pcm-dpc/COVID-19/raw/master/dati-regioni/dpc-covid19-ita-regioni.csv"
)

var cfg Config

func main() {
	log.SetFlags(0)
	log.SetOutput(ioutil.Discard) // discard debugging output as default

	app := cli.App("covid19", "Daily stats for covid19 in Italy.")

	app.Version("v version", "covid19 0.0.1")
	app.Spec = "[-r [-a]] [-d] [-D]"

	app.StringOptPtr(&cfg.Region, "r region", "", "Specify a region")
	app.BoolOptPtr(&cfg.printRegions, "a availables", false, "Print available regions")
	app.VarOpt("d date", &cfg.startDate, "Date in yyyy-mm-dd format")
	app.BoolOptPtr(&cfg.debug, "D debug", false, "Debug mode")

	app.Action = func() {
		mainAction()
	}

	app.Run(os.Args)
}

func mainAction() {
	fmt.Println("Covid19 - dati sintetici")

	if cfg.debug == true {
		log.SetOutput(os.Stderr)
	}

	var csv string
	if cfg.Region != "" {
		csv = getData(RegionsDataURL)
		fmt.Printf("regione selezionata: %v\n", cfg.Region)
	} else {
		csv = getData(ItalyDataURL)
	}

	df := loadDataFrame(csv)

	if cfg.Region != "" {
		if cfg.printRegions == true {
			printAvailableRegions(df)
		}
		filterByRegion(df, cfg.Region)
	}

	printSummary(df)
	printPercentages(df)
}
