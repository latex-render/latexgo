// Copyright ©2020 The go-latex Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	genFonts()
	genKerns()
}

func genFonts() {
	buf := new(bytes.Buffer)
	cmd := exec.Command("python", "-c", pyFonts)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run python script: %+v", err)
	}

	var db []struct {
		FontName string  `json:"font_name"`
		FontType string  `json:"font_type"`
		Math     bool    `json:"math"`
		Size     float64 `json:"size"`
		Symbol   string  `json:"symbol"`
		Metrics  Metrics `json:"metrics"`
	}

	err = json.NewDecoder(buf).Decode(&db)
	if err != nil {
		log.Fatalf("could not decode json: %+v", err)
	}

	out, err := os.Create("fakebackend_fonts_gen.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(out, `// Autogenerated. DO NOT EDIT.

package fakebackend

import "github.com/go-latex/latex/tex"

func init() {
	fontsDb = dbFonts{
`)

	for _, v := range db {
		fmt.Fprintf(
			out,
			"fontKey{symbol: %q, math: %v, font: tex.Font{Name:%q, Size:%v, Type:%q}}:",
			v.Symbol, v.Math, v.FontName, v.Size, v.FontType,
		)

		fmt.Fprintf(
			out,
			"tex.Metrics{Advance: %g, Height: %g, Width: %g, XMin: %g, XMax: %g, YMin: %g, YMax: %g, Iceberg: %g, Slanted: %v},\n",
			v.Metrics.Advance,
			v.Metrics.Height,
			v.Metrics.Width,
			v.Metrics.XMin,
			v.Metrics.XMax,
			v.Metrics.YMin,
			v.Metrics.YMax,
			v.Metrics.Iceberg,
			v.Metrics.Slanted,
		)
	}

	fmt.Fprintf(out, "\t}\n}\n")

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

type Metrics struct {
	Advance float64 `json:"advance"`
	Height  float64 `json:"height"`
	Width   float64 `json:"width"`
	XMin    float64 `json:"xmin"`
	XMax    float64 `json:"xmax"`
	YMin    float64 `json:"ymin"`
	YMax    float64 `json:"ymax"`
	Iceberg float64 `json:"iceberg"`
	Slanted bool    `json:"slanted"`
}

const pyFonts = `
import sys
import string
import json
import matplotlib.mathtext as mtex

dejavu = mtex.DejaVuSansFonts(mtex.FontProperties(), mtex.MathtextBackendPdf())

math = [True, False]
sizes = [10,12]
fonts = [
	("default", "regular"),
	("default", "rm"),
	("it", "it")
	]
symbols = list(string.ascii_letters) + ["\\"+k for k in mtex.tex2uni.keys()] + [
	"é",
	" ",
	]

db = []
for math in [True, False]:
    for font in fonts:
        fontName = font[0]
        fontClass = font[1]
        for size in sizes:
            for sym in symbols:
                m = dejavu.get_metrics(fontName, fontClass, sym, size, 72, math)
                db.append({
					'font_name': font[0], 'font_type': font[1],
                    'math':math, 'size': size,
                    'symbol': sym,
                    'metrics': {
                        'advance': m.advance,
                        'height': m.height,
                        'width': m.width,
                        'xmin': m.xmin,
                        'xmax': m.xmax,
                        'ymin': m.ymin,
                        'ymax': m.ymax,
                        'iceberg': m.iceberg,
                        'slanted': m.slanted,
                        },
                    })
				#print(f"sym: {sym}=> {m}")

with open("testdata/metrics-dejavu-sans.json", "w") as f:
    json.dump(db, f, indent='  ')
    pass
json.dump(db, sys.stdout)
sys.stdout.flush()
`

func genKerns() {
	buf := new(bytes.Buffer)
	cmd := exec.Command("python", "-c", pyKerns)
	cmd.Stdout = buf

	err := cmd.Run()
	if err != nil {
		log.Fatalf("could not run python script: %+v", err)
	}

	var db []struct {
		FontName string  `json:"font_name"`
		FontType string  `json:"font_type"`
		Size     float64 `json:"size"`
		Symbol1  string  `json:"sym1"`
		Symbol2  string  `json:"sym2"`
		Kern     float64 `json:"kern"`
	}

	err = json.NewDecoder(buf).Decode(&db)
	if err != nil {
		log.Fatalf("could not decode json: %+v", err)
	}

	out, err := os.Create("fakebackend_kerns_gen.go")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(out, `// Autogenerated. DO NOT EDIT.

package fakebackend

import "github.com/go-latex/latex/tex"

func init() {
	kernsDb = dbKerns{
`)

	for _, v := range db {
		fmt.Fprintf(
			out,
			"kernKey{font: tex.Font{Name:%q, Size:%v, Type:%q}, s1: %q, s2: %q}: %g,\n",
			v.FontName, v.Size, v.FontType, v.Symbol1, v.Symbol2, v.Kern,
		)
	}

	fmt.Fprintf(out, "\t}\n}\n")

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

const pyKerns = `
import sys
import string
import json
import matplotlib.mathtext as mtex

dejavu = mtex.DejaVuSansFonts(mtex.FontProperties(), mtex.MathtextBackendPdf())

dpi = 72
sizes = [12]
fonts = [
	("default", "regular"),
	("default", "rm"),
	("it", "it")
	]
symbols = [
	("A", "V"),
	("A", "é"),
	("V", "é"),
	("é", "é"),
	("f", "i"),
	("A", r"\sigma"),
	("a", r"\sigma"),
	("é", r"\sigma"),
	(r"\sum", r"\sigma"),
]

db = []
for font in fonts:
    fontName = font[0]
    fontClass = font[1]
    for size in sizes:
        for sym in symbols:
            kern = dejavu.get_kern(fontName, fontClass, sym[0], size, fontName, fontClass, sym[1], size, dpi)
            db.append({
				'font_name': font[0], 'font_type': font[1], 'size': size,
                'sym1': sym[0],
                'sym2': sym[1],
                'kern': kern,
                })

            kern = dejavu.get_kern(fontName, fontClass, sym[1], size, fontName, fontClass, sym[0], size, dpi)
            db.append({
				'font_name': font[0], 'font_type': font[1], 'size': size,
                'sym1': sym[1],
                'sym2': sym[0],
                'kern': kern,
                })

with open("testdata/kerns-dejavu-sans.json", "w") as f:
    json.dump(db, f, indent='  ')
    pass
json.dump(db, sys.stdout)
sys.stdout.flush()
`
