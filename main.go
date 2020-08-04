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
	"sync"
)

type certs struct {
	Url  string
	Cert string
}

var wg sync.WaitGroup

func main() {
	urls := readCsvFile("urls.csv")

	errs := make(chan certs)
	var results []certs
	fmt.Println(len(urls))
	wg.Add(len(urls))

	for _, v := range urls {
		go queryAsync(v[0], errs)
	}

	go func() {
		for v := range errs {
			fmt.Printf("Found bad cert for %s\n", v.Url)
			results = append(results, v)
		}
	}()

	wg.Wait()
	close(errs)

	file, err := json.MarshalIndent(results, "", " ")
	if err != nil {
		log.Fatalf("Error indenting - %s", err.Error())
	}

	_ = ioutil.WriteFile("output.json", file, 0644)
}

func queryAsync(url string, errs chan certs) {
	defer wg.Done()
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{}
	fmt.Printf("Running for %s\n", url)
	_, err := http.Get(url)
	if err != nil {
		a := certs{
			Url:  url,
			Cert: err.Error(),
		}
		errs <- a
	}
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
