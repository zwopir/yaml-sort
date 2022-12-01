package main

import (
	"github.com/alecthomas/kong"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
	"strings"
)

var CLI struct {
	Sort struct {
	} `cmd:"" embed:"" help:"Sort YAML"`
	Infile  string `name:"infile" help:"input file. Defaults to stdin"`
	Outfile string `name:"outfile" help:"output file. Defaults to stdout"`
}

func main() {
	kongCTX := kong.Parse(&CLI)
	// set infile
	infile := os.Stdin
	if !(strings.EqualFold(CLI.Infile, "-") || strings.EqualFold(CLI.Infile, "")) {
		var err error
		if infile, err = os.Open(CLI.Infile); err != nil {
			kongCTX.Errorf("failed to open input file %s: %s", err)
			os.Exit(1)
		}
		defer func() {
			_ = infile.Close()
		}()
	}
	// set outfile
	outfile := os.Stdout
	if !(strings.EqualFold(CLI.Outfile, "-") || strings.EqualFold(CLI.Outfile, "")) {
		var err error
		if outfile, err = os.Create(CLI.Outfile); err != nil {
			kongCTX.Errorf("failed to create output file %s: %s", err)
			os.Exit(1)
		}
		defer func() {
			_ = outfile.Close()
		}()
	}
	var in yaml.MapSlice
	dec := yaml.NewDecoder(infile)
	dec.SetStrict(true)
	if err := dec.Decode(&in); err != nil {
		kongCTX.Errorf("failed to decode input yaml: %s", err)
		os.Exit(1)
	}
	out := sortYAML(in)
	if err := yaml.NewEncoder(outfile).Encode(out); err != nil {
		kongCTX.Errorf("failed to encode sorted yaml: %s", err)
		os.Exit(1)
	}
}

type sortedYAML []yaml.MapItem

func (s sortedYAML) Len() int {
	return len(s)
}

func (s sortedYAML) Less(i, j int) bool {
	iStr, ok := s[i].Key.(string)
	if !ok {
		panic("key is not string assertable")
	}
	jStr, ok := s[j].Key.(string)
	if !ok {
		panic("key is not string assertable")
	}
	return strings.Compare(iStr, jStr) < 0
}

func (s sortedYAML) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sortYAML(in yaml.MapSlice) sortedYAML {
	sort.Sort(sortedYAML(in))
	for _, v := range in {
		if obj, ok := v.Value.(yaml.MapSlice); ok {
			sortedObj := sortYAML(obj)
			v.Value = sortedObj
		}
	}
	return sortedYAML(in)
}
