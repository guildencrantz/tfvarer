package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/hashicorp/hcl/hcl/printer"
	jsonParser "github.com/hashicorp/hcl/json/parser"
	"github.com/imdario/mergo"
)

func usage() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal("Unable to get the executable path!", err)
	}

	fmt.Printf("Usage: %s [flags]\n", ex)
	fmt.Println("\nMerge JSON and HCL configurations, allowing output in either HCL or JSON format")
	fmt.Println("Configs are processed in the order they are specified, reguardless of format.\n")

	fmt.Println("Flags: ")
	w := tabwriter.NewWriter(os.Stdout, 1, 8, 1, '\t', 0)
	defer w.Flush()

	fmt.Fprintf(w, "\t%s\t%s\n", "-json <path to json file>", "Path to a file containing JSON formatted input config. Multiple can be specified.")
	fmt.Fprintf(w, "\t%s\t%s\n", "-hcl <path to hcl file>", "Path to a file containing HCL formatted input config. Multiple can be specified.")
	fmt.Fprintf(w, "\t%s\t%s\n", "-output [hcl | json]", "Whether to output the calculated config in 'hcl' or 'json' (default 'hcl')")
	fmt.Fprintf(w, "\t%s\t%s\n", "-overwrite [true | false]", "If overwrite is true the value from the _last_ config will be used for conflicting keys, if false the _first_ value is preserved (default 'true')")
}

func main() {
	inputs, output := parseFlags(os.Args)
	if len(inputs) == 0 {
		fmt.Println("You must specify at least one input file\n")
		usage()
		os.Exit(1)
	}

	//var final interface{}
	var final map[string]interface{}
	for _, input := range inputs {
		mergo.Merge(&final, input.Unmarshal(), output.overwrite)
	}

	fmt.Println(output.Format(final))
}

func hclFormatter(v interface{}) string {
	j := jsonFormatter(v)

	ast, err := jsonParser.Parse([]byte(j))
	if err != nil {
		log.Fatal("Unable to parse JSON", err)
	}

	bb := new(bytes.Buffer)
	if printer.Fprint(bb, ast) != nil {
		log.Fatal("Unable to format HCL", err)
	}

	return bb.String()
}

func jsonFormatter(i interface{}) string {
	b, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		log.Fatal("Unable to format JSON", err)
	}

	return string(b)
}
