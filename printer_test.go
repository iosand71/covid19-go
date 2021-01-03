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
)

func init() {
	s0 = dataframe.NewSeriesTime("data", nil, time.Now(), time.Now(), time.Now())

	s1 = dataframe.NewSeriesInt64("totale_casi", nil, 1, 2, 3)
	s2 = dataframe.NewSeriesInt64("nuovi_positivi", nil, 1, 2, 3)
	s3 = dataframe.NewSeriesInt64("totale_positivi", nil, 1, 2, 3)
	s4 = dataframe.NewSeriesInt64("variazione_totale_positivi", nil, 1, 2, 3)
	s5 = dataframe.NewSeriesInt64("deceduti", nil, 1, 2, 3)
	s6 = dataframe.NewSeriesInt64("terapia_intensiva", nil, 1, 2, 3)
	s7 = dataframe.NewSeriesInt64("ingressi_terapia_intensiva", nil, nil, 2, 3)
	s8 = dataframe.NewSeriesInt64("totale_ospedalizzati", nil, 1, 2, 3)
	s9 = dataframe.NewSeriesInt64("dimessi_guariti", nil, 1, 2, 3)
	s10 = dataframe.NewSeriesInt64("tamponi", nil, 1, 2, 3)
	s11 = dataframe.NewSeriesInt64("casi_testati", nil, 1, 2, 3)

	df = dataframe.NewDataFrame(s0, s1, s2, s3, s4, s5, s6, s7, s8, s9, s10, s11)

	os.Stdout, _ = os.OpenFile(os.DevNull, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func TestPrintSummary(t *testing.T) {
	err := printSummary(df)
	if err != nil {
		t.Error(err)
	}
}

func TestPrintPercentages(t *testing.T) {
	err := printPercentages(df)
	if err != nil {
		t.Error(err)
	}
}