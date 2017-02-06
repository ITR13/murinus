package main

import (
	"encoding/xml"
	"os"
)

var options Options

type Options struct {
	Character  uint8
	EdgeSlip   int
	BetterSlip int32
	AllKeys    *[]*Key `xml:"Keys>Key"`
}

func ReadOptions(path string, input *Input) {
	options = Options{0, EdgeSlipDefault, BetterSlipDefault, &input.allInputs}
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		e(err)
		defer file.Close()
		decoder := xml.NewDecoder(file)
		e(decoder.Decode(&options))
	}
}

func SaveOptions(path string, input *Input) {
	file, err := os.Create(path)
	e(err)
	defer file.Close()
	encoder := xml.NewEncoder(file)
	encoder.Indent("  ", "    ")

	e(encoder.Encode(&options))
}
