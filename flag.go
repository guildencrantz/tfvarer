package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/hashicorp/hcl"
	"github.com/imdario/mergo"
)

type input struct {
	file         string
	Unmarshaller func([]byte, interface{}) error
}

func (i input) Unmarshal() interface{} {
	b, err := ioutil.ReadFile(i.file)
	if err != nil {
		log.Fatalf("Unable to read %q\n%e", i.file, err)
	}

	var ret interface{}
	i.Unmarshaller(b, &ret)

	return ret
}

type output struct {
	overwrite func(*mergo.Config)
	Format    func(interface{}) string
}

func parseFlags(args []string) ([]input, output) {
	outputter := output{
		overwrite: mergo.WithOverride,
		Format:    hclFormatter,
	}
	inputs := make([]input, 0, (len(args)-1)/2)

	for i, s := range args {
		if i%2 == 0 {
			continue
		}
		if s[0] != '-' {
			log.Fatalf("Expected a flag for arg %d, got %q", i, s)
		}

		var v string
		if len(args) > (i + 1) {
			v = args[i+1]
		}

		fn := s[1:]
		switch fn {
		case "hcl":
			inputs = append(inputs, input{
				file:         v,
				Unmarshaller: hcl.Unmarshal,
			})
		case "json":
			inputs = append(inputs, input{
				file:         v,
				Unmarshaller: json.Unmarshal,
			})
		case "output":
			switch v {
			case "json":
				outputter.Format = jsonFormatter
			case "hcl":
				outputter.Format = hclFormatter
			default:
				usage()
				log.Fatalf("Unknown output flag %q. Must be one of 'hcl' or 'json'", v)
			}
		case "overwrite":
			switch v {
			case "true":
				outputter.overwrite = mergo.WithOverride
			case "false":
				outputter.overwrite = func(c *mergo.Config) { c.Overwrite = false }
			default:
				usage()
				log.Fatalf("Overwrite must be either 'true' or 'false'. Received %q", v)
			}
		case "h":
			fallthrough
		case "help":
			usage()
			os.Exit(0)
		default:
			usage()
			log.Fatalf("Unknown flag type %q. Must be one of 'hcl', 'json', or 'output'", fn)
		}

		if v == "" {
			log.Fatalf("No value for flag %q!", fn)
		}
	}

	return inputs, outputter
}
