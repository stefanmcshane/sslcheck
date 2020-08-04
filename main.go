package main

import (
	"crypto/tls"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type certs struct {
	Url  string
	Cert string
}

func main() {
	urls := readCsvFile("urls.csv")

	var errs []certs

	for _, v := range urls {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{}
		_, err := http.Get(v[0])
		if err != nil {
			a := certs{
				Url:  v[0],
				Cert: err.Error(),
			}
			errs = append(errs, a)
			fmt.Println(err.Error())
		}
	}
	fmt.Println(errs)
	file, err := json.MarshalIndent(errs, "", " ")
	if err != nil {
		log.Fatalf("Error indenting - %s", err.Error())
	}

	fmt.Println(string(file))

	_ = ioutil.WriteFile("output.json", file, 0644)
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}
