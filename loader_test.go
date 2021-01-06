package main

import (
	"log"
	"testing"
)

const (
	countryData = `
data,stato,ricoverati_con_sintomi,terapia_intensiva,totale_ospedalizzati,isolamento_domiciliare,totale_positivi,variazione_totale_positivi,nuovi_positivi,dimessi_guariti,deceduti,casi_da_sospetto_diagnostico,casi_da_screening,totale_casi,tamponi,casi_testati,note,ingressi_terapia_intensiva,note_test,note_casi
2020-12-28T17:00:00,ITA,23932,2565,26497,548724,575221,-6539,8585,1408686,72370,,,2056277,26114818,14685718,,167,,
2020-12-29T17:00:00,ITA,23662,2549,26211,542517,568728,-6493,11224,1425730,73029,,,2067487,26243558,14731420,,256,,
2020-12-30T17:00:00,ITA,23566,2528,26094,538301,564395,-4333,16202,1445690,73604,,,2083689,26412603,14795168,,175,,
2020-12-31T17:00:00,ITA,23151,2555,25706,544190,569896,5501,23477,1463111,74159,,,2107166,26598607,14871966,,202,,

`
	regionalData = `
data,stato,codice_regione,denominazione_regione,lat,long,ricoverati_con_sintomi,terapia_intensiva,totale_ospedalizzati,isolamento_domiciliare,totale_positivi,variazione_totale_positivi,nuovi_positivi,dimessi_guariti,deceduti,casi_da_sospetto_diagnostico,casi_da_screening,totale_casi,tamponi,casi_testati,note,ingressi_terapia_intensiva,note_test,note_casi
2020-12-31T17:00:00,ITA,01,Piemonte,45.0732745,7.680687483,2895,190,3085,25172,28257,-702,1367,161649,7922,,,197828,1682529,984809,,7,,
2020-12-31T17:00:00,ITA,16,Puglia,41.12559576,16.86736689,1490,129,1619,51383,53002,262,1661,35490,2472,,,90964,1044314,653907,,11,,
2020-12-31T17:00:00,ITA,20,Sardegna,39.21531192,9.110616306,486,46,532,15921,16453,128,368,13913,747,,,31113,482520,404939,Si segnala il decesso dei seguenti pz: 1 uomo 79 aa residente nella Città Metropolitana di Cagliari; 1 donna 101 aa residente nella Città Metropolitana di Cagliari; 1 uomo 88 aa residente nella Provincia di Sassari; 1 uomo 72 aa residente nella Provincia di Sassari. ,3,,
2020-12-31T17:00:00,ITA,19,Sicilia,38.11569725,13.362356699999998,1069,171,1240,32628,33868,481,1299,57364,2412,,,93644,1219132,812545,,13,,
2020-12-31T17:00:00,ITA,09,Toscana,43.76923077,11.25588885,838,150,988,8690,9678,-59,632,106977,3673,,,120328,1883593,1061504,,3,"Positivi diagnosticati solo con test antigenico rapido: in questo momento non è procedura adottatata da Regione Toscana, pertanto il valore è pari a zero",
`
)

func TestLoadDataFrame(t *testing.T) {
	disableStdout()
	df := loadDataFrame(countryData)
	enableStdout()

	df.MustNameToColumn("data")
	df.MustNameToColumn("ingressi_terapia_intensiva")
	df.MustNameToColumn("casi_testati")
}

func TestFilterByRegion(t *testing.T) {
	df := loadDataFrame(regionalData)
	filterByRegion(df, "Piemonte")
	if df.NRows() != 1 {
		t.Error("Wrong filter by region")
	}
}

func TestSetDateTime(t *testing.T) {
	var date DateTime
	err := date.Set("2020-12-25")
	if err != nil {
		t.Error("There'a problem with iso date parsing")
	}
	log.Println(date.String())
}
