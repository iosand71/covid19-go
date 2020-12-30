package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	dataframe "github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
)

func getData(URL string) string {
	var bodyString string
	resp, err := http.Get(URL)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
	}

	return bodyString
}

func loadDataFrame(csvStr string) (df *dataframe.DataFrame) {
	var ctx = context.Background()

	opts := imports.CSVLoadOptions{
		InferDataTypes: true,
		NilValue:       &[]string{"NA"}[0],
		DictateDataType: map[string]interface{}{
			"data": imports.Converter{
				ConcreteType:  time.Time{},
				ConverterFunc: convertToTime,
			},
			"ingressi_terapia_intensiva": imports.Converter{
				ConcreteType:  int64(0),
				ConverterFunc: convertToInt64,
			},
			"casi_testati": imports.Converter{
				ConcreteType:  int64(0),
				ConverterFunc: convertToInt64,
			},
		},
	}
	df, err := imports.LoadFromCSV(ctx, strings.NewReader(csvStr), opts)

	if err != nil {
		log.Fatal(err)
	}

	logDataframe(df)

	return df
}

func convertToTime(in interface{}) (interface{}, error) {
	return time.Parse(timeFormat, in.(string))
}

func convertToInt64(in interface{}) (interface{}, error) {
	if in == nil || in.(string) == "" {
		return nil, nil
	}
	return strconv.ParseInt(in.(string), 10, 64)
}

func logDataframe(df *dataframe.DataFrame) {

	iterator := df.ValuesIterator(dataframe.ValuesOptions{0, 1, true})

	df.Lock()
	for {
		row, vals, _ := iterator()
		if row == nil {
			break
		}
		log.Println(*row, vals)
	}
	df.Unlock()
}
