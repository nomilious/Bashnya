package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

var lines = make(map[string]int, 0)

func open_file(name string, alternative *os.File, write bool) *os.File {
	var reader *os.File
	if name != "" {
		var f *os.File
		var err error
		if write {
			f, err = os.Create(name)
		} else {
			f, err = os.Open(name)
		}
		if err != nil {
			fmt.Println("ERRRRRROOOOOORRR")
			os.Exit(-1)
		}
		reader = f
	} else {
		reader = alternative
	}
	return reader
}

func uniq(rules map[string]interface{}, in_file string, out_file string) {
	var reader = open_file(in_file, os.Stdin, false)
	defer reader.Close()
	read_all(reader,
		*(rules["i"].(*bool)),
		*(rules["f"].(*int)),
		*(rules["s"].(*int)))

	var outputer = open_file(out_file, os.Stdout, true)
	defer outputer.Close()
	output_all(outputer,
		*(rules["c"].(*bool)),
		*(rules["d"].(*bool)),
		*(rules["u"].(*bool)))
}
func output_all(out *os.File, count bool, dupplicates bool, uniq bool) {
	if dupplicates { // deleete unique
		for key, value := range lines {
			if value == 1 {
				delete(lines, key)
			}
		}
	}
	if uniq { // delete dupplicates
		for key, value := range lines {
			if value > 1 {
				delete(lines, key)
			}
		}
	}
	for key, value := range lines {
		if count {
			io.WriteString(out, strconv.Itoa(value))
			io.WriteString(out, " ")
		}
		io.WriteString(out, key)
		io.WriteString(out, "\n")
	}
}
func main() {
	flags := map[string]interface{}{
		"f": 0, "s": 0,
		"d": false, "i": false, "u": false, "c": false,
	}
	for flagName, flagValue := range flags {
		switch flagValue.(type) {
		case int:
			flags[flagName] = flag.Int(flagName, flagValue.(int), "")
		case bool:
			flags[flagName] = flag.Bool(flagName, flagValue.(bool), "a bool")
		}
	}
	flag.Parse()

	args := flag.Args()
	var in_file, out_file string
	if len(args) == 1 {
		in_file = args[0]
	} else if len(flag.Args()) == 2 {
		in_file, out_file = args[0], args[1]
	}

	uniq(flags, in_file, out_file)
}
func add_to_map(line string, caseIgnore bool, num_chars int) {
	var key string
	var value int = 0
	for k, v := range lines { // check if exists in map
		if len(k) < num_chars || len(line) < num_chars { // if too short
			continue
		}
		if caseIgnore {
			if strings.ToLower(k[num_chars:]) == strings.ToLower(line[num_chars:]) {
				value = v
				key = k
			}
		} else if k[num_chars:] == line[num_chars:] {
			value = v
			key = k
		}
	}
	if value == 0 { // add new key
		lines[line] = 1
	} else {
		lines[key]++
	}

}
func read_all(reader *os.File, caseIgnore bool, num_fields int, num_chars int) {
	// if f skip, if s skip, but memorise the original
	b := make([]byte, 1)
	line := ""
	num_field := 0
	for {
		_, err := reader.Read(b)
		if err == io.EOF {
			break // End of input
		} else if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		// handle read data
		if string(b) != "\n" {
			line += string(b)
		} else {
			num_field++
			if num_field > num_fields {
				add_to_map(line, caseIgnore, num_chars)
			}
			// set to default
			line = ""
		}
	}
	// if there is no '\n' at the end of the file
	if line != "" && num_field > num_fields {
		add_to_map(line, caseIgnore, num_chars)
	}
}
