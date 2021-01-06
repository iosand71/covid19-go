package main

import (
	"testing"
)

func TestLoadDataFrame(t *testing.T) {
	csvString := `
data,stato,ricoverati_con_sintomi,terapia_intensiva,totale_ospedalizzati,isolamento_domiciliare,totale_positivi,variazione_totale_positivi,nuovi_positivi,dimessi_guariti,deceduti,casi_da_sospetto_diagnostico,casi_da_screening,totale_casi,tamponi,casi_testati,note,ingressi_terapia_intensiva,note_test,note_casi
2020-12-28T17:00:00,ITA,23932,2565,26497,548724,575221,-6539,8585,1408686,72370,,,2056277,26114818,14685718,,167,,
2020-12-29T17:00:00,ITA,23662,2549,26211,542517,568728,-6493,11224,1425730,73029,,,2067487,26243558,14731420,,256,,
2020-12-30T17:00:00,ITA,23566,2528,26094,538301,564395,-4333,16202,1445690,73604,,,2083689,26412603,14795168,,175,,
2020-12-31T17:00:00,ITA,23151,2555,25706,544190,569896,5501,23477,1463111,74159,,,2107166,26598607,14871966,,202,,

`
	disableStdout()
	df := loadDataFrame(csvString)
	enableStdout()

	df.MustNameToColumn("data")
	df.MustNameToColumn("ingressi_terapia_intensiva")
	df.MustNameToColumn("casi_testati")
}
