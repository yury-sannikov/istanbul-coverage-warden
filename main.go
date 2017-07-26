package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yury-sannikov/istanbul-coverage-warden/xmltools"
)

type coberturaClassMap map[string]xmltools.CoberturaClass

var currentCoberturaFile = flag.String("current", "cobertura-coverage.xml", "Current cobertura file path")
var previousCoberturaFile = flag.String("previous", "cobertura-coverage-prev.xml", "Previous cobertura file path")

func main() {
	flag.Parse()

	currentClasses, err := xmltools.BuildCoberturaClasses(*currentCoberturaFile)
	if err != nil {
		fmt.Println("Error opening current cobertura report file:", err)
		os.Exit(2)
	}
	prevClasses, err := xmltools.BuildCoberturaClasses(*previousCoberturaFile)
	if err != nil {
		fmt.Println("Error opening previous cobertura report file:", err)
		os.Exit(2)
	}

	currentClassMap := make(coberturaClassMap)
	for _, item := range currentClasses {
		currentClassMap[item.FileName] = item
	}
	result := compare(prevClasses, currentClassMap)

	if !result {
		fmt.Println("Code coverage drop has been detected.")
		os.Exit(1)
	}
}

func compare(previousCoberturaData []xmltools.CoberturaClass, newCoberturaMap coberturaClassMap) bool {
	var result = true
	for _, item := range previousCoberturaData {
		newItem, found := newCoberturaMap[item.FileName]
		if !found {
			fmt.Println("Unable to find coverage infrormation for ", item.FileName)
			continue
		}
		if item.LineRate > newItem.LineRate {
			fmt.Printf("Code coverage dropped for %s from %.2f%% to %.2f%%\n", item.FileName, item.LineRate*100.0, newItem.LineRate*100.0)
			result = false
		}
	}
	return result
}
