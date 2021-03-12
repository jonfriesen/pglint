package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jonfriesen/pglint/pkg/ecpg"
)

func main() {
	addSemicolons := flag.Bool("addSemicolons", true, "Automatically add semicolons to queries if not already present (will fail without a semicolon)")
	trimWhitespace := flag.Bool("trim", true, "Trim whitespace at head and tail of query")
	allowQuestionMarks := flag.Bool("questionmarks", true, "Determines whether or not the to fill question mark placeholders with sample data for linting")
	inputFile := ""
	flag.StringVar(&inputFile, "input", "", "What file to read in to analyze, just reads from end of command if empty")

	flag.Parse()

	config := ecpg.Config{
		AddSemiColon:   *addSemicolons,
		TrimWhiteSpace: *trimWhitespace,
		QuestionMarks:  *allowQuestionMarks,
	}

	linter, err := ecpg.NewECPG(&config)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	query := ""
	if inputFile == "" {
		query = strings.Join(flag.Args(), " ")
	} else {
		file, err := ioutil.ReadFile(inputFile)
		if err != nil {
			fmt.Printf("ERROR: %s", err)
			os.Exit(1)
		}
		query = string(file)
	}

	errArray := linter.CheckStatement(query)
	if len(errArray) != 0 {
		for _, err := range errArray {
			fmt.Printf("ERROR: %s", err)
		}
		os.Exit(1)
	}

}
