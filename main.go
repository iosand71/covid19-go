package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	cli "github.com/jawher/mow.cli"
)

const (
	timeFormat = "2006-01-02T15:04:05"
	dateFormat = "2006-01-02"
)

// DateTime alias
type DateTime time.Time

//Set date using custom format
func (d *DateTime) Set(v string) error {
	parsed, err := time.Parse(dateFormat, v)
	if err != nil {
		return err
	}
	*d = DateTime(parsed)
	return nil
}

func (d *DateTime) isZero() bool {
	date := time.Time(*d)
	return date.IsZero()
}

func (d *DateTime) String() string {
	date := time.Time(*d)
	return date.Format(dateFormat)
}

// Config holder for cli configs
type Config struct {
	Region       string
	startDate    DateTime
	printRegions bool
	debug        bool
}

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
	app.Spec = "[-r [-a]] [-d] [-V]"

	app.StringOptPtr(&cfg.Region, "r region", "", "Specify a region")
	app.BoolOptPtr(&cfg.printRegions, "a availables", false, "Print available regions")
	app.VarOpt("d date", &cfg.startDate, "Date in yyyy-mm-dd format")
	app.BoolOptPtr(&cfg.debug, "V verbose", false, "Verbose mode (debugging)")

	app.Action = func() {
		mainAction()
	}

	app.Run(os.Args)
}

func mainAction() {
	fmt.Println("Covid19 - dati sintetici")
	fmt.Println()

	var csv string
	if cfg.debug == true {
		log.SetOutput(os.Stderr)
	}

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
