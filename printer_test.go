package main

import (
	"os"
	"testing"
	"time"

	dataframe "github.com/rocketlaunchr/dataframe-go"
)

var (
	s0                       *dataframe.SeriesTime
	s1, s2, s3, s4, s5       *dataframe.SeriesInt64
	s6, s7, s8, s9, s10, s11 *dataframe.SeriesInt64
	df                       *dataframe.DataFrame
	oldStdout                *os.File
)

func init() {
	s0 = dataframe.NewSeriesTime("data", nil,
		time.Date(2020, time.December, 23, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.December, 24, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.December, 25, 0, 0, 0, 0, time.UTC),
		time.Date(2020, time.December, 26, 0, 0, 0, 0, time.UTC),
	)

	s1 = dataframe.NewSeriesInt64("totale_casi", nil, 1, 2, 3, 4)
	s2 = dataframe.NewSeriesInt64("nuovi_positivi", nil, 1, 2, 3, 4)
	s3 = dataframe.NewSeriesInt64("totale_positivi", nil, 1, 2, 3, 4)
	s4 = dataframe.NewSeriesInt64("variazione_totale_positivi", nil, 1, 2, 3, 4)
	s5 = dataframe.NewSeriesInt64("deceduti", nil, 1, 2, 3, 4)
	s6 = dataframe.NewSeriesInt64("terapia_intensiva", nil, 1, 2, 3, 4)
	s7 = dataframe.NewSeriesInt64("ingressi_terapia_intensiva", nil, 1, 2, 3, 4)
	s8 = dataframe.NewSeriesInt64("totale_ospedalizzati", nil, 1, 2, 3, 4)
	s9 = dataframe.NewSeriesInt64("dimessi_guariti", nil, 1, 2, 3, 4)
	s10 = dataframe.NewSeriesInt64("tamponi", nil, 1, 2, 3, 4)
	s11 = dataframe.NewSeriesInt64("casi_testati", nil, 1, 2, 3, 4)

	df = dataframe.NewDataFrame(s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11)

	oldStdout = os.Stdout
}

func disableStdout() {
	os.Stdout, _ = os.OpenFile(os.DevNull, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func enableStdout() {
	os.Stdout = oldStdout
}

func TestPrintSummary(t *testing.T) {
	disableStdout()

	err := printSummary(df)
	if err != nil {
		t.Error(err)
	}

	enableStdout()
}

func TestPrintPercentages(t *testing.T) {
	disableStdout()

	err := printPercentages(df)
	if err != nil {
		t.Error(err)
	}

	enableStdout()
}

func TestGetLastDays(t *testing.T) {
	first, second, third := getLastDays(df)
	if first == nil || second == nil || third == nil {
		t.Error("Unexpected null result from getLastDays")
	}
	if first["totale_casi"] != int64(4) {
		t.Error("Wrong value from getLastDays")
	}
}

func TestPreviousDate(t *testing.T) {
	cfg = Config{startDate: DateTime(time.Date(2020, time.December, 24, 1, 0, 0, 0, time.UTC))}
	TestPrintSummary(t)
	TestPrintPercentages(t)
}
