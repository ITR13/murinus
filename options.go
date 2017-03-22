/*
   This file is part of Murinus.

   Murinus is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Murinus is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Murinus.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"encoding/xml"
	"os"
)

var options Options

type Options struct {
	CharacterP1 uint8
	CharacterP2 uint8
	EdgeSlip    int
	BetterSlip  int32
	ShowDivert  uint8
	showDivert  bool
	AllKeys     *[]*Key `xml:"Keys>Key"`
}

func ReadOptions(path string, input *Input) {
	options = Options{0, 2, EdgeSlipDefault, BetterSlipDefault,
		ShowDivertDefault, ShowDivertDefault != 0, &input.allInputs}
	if path == "" {
		return
	}
	if _, err := os.Stat(path); err == nil {
		file, err := os.Open(path)
		e(err)
		defer file.Close()
		decoder := xml.NewDecoder(file)
		e(decoder.Decode(&options))
	}
	options.showDivert = options.ShowDivert != 0
}

func SaveOptions(path string, input *Input) {
	file, err := os.Create(path)
	e(err)
	defer file.Close()
	encoder := xml.NewEncoder(file)
	encoder.Indent("  ", "    ")

	e(encoder.Encode(&options))
}
