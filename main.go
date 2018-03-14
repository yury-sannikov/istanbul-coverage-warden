package main

import (
	"os"

	flags "github.com/jessevdk/go-flags"
)

// Options contains general options
type Options struct {
}

var options Options

var parser = flags.NewParser(&options, flags.Default)

func main() {

	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}

// func compare(previousCoberturaData []xmltools.CoberturaClass, newCoberturaMap coberturaClassMap) bool {
// 	var result = true
// 	for _, item := range previousCoberturaData {
// 		newItem, found := newCoberturaMap[item.FileName]
// 		if !found {
// 			fmt.Println("Unable to find coverage infrormation for ", item.FileName)
// 			continue
// 		}
// 		if item.LineRate > newItem.LineRate {
// 			fmt.Printf("Code coverage dropped for %s from %.2f%% to %.2f%%\n", item.FileName, item.LineRate*100.0, newItem.LineRate*100.0)
// 			result = false
// 		}
// 	}
// 	return result
// }
