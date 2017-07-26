package xmltools

import (
	"encoding/xml"
	"os"
)

// CoberturaClass report class item. Contains basic coverage info for analysis
type CoberturaClass struct {
	Name     string  `xml:"name,attr"`
	FileName string  `xml:"filename,attr"`
	LineRate float32 `xml:"line-rate,attr"`
}

const classNodeName string = "class"

// BuildCoberturaClasses creates an array of CoberturaClass for specified cobertura report
func BuildCoberturaClasses(xmlFileName string) ([]CoberturaClass, error) {
	var cobClasses []CoberturaClass
	xmlFile, err := os.Open(xmlFileName)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()
	decoder := xml.NewDecoder(xmlFile)
	for {
		t, _ := decoder.Token()
		if t == nil {
			break
		}
		switch se := t.(type) {
		case xml.StartElement:
			if se.Name.Local == classNodeName {
				var cobClass CoberturaClass
				decoder.DecodeElement(&cobClass, &se)
				cobClasses = append(cobClasses, cobClass)
			}
		default:
		}
	}
	return cobClasses, nil
}
