package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/jonfriesen/pglint/pkg/ecpg"
)

func main() {
	addSemicolons := flag.Bool("addSemicolons", true, "Automatically add semicolons to queries if not already present (will fail without a semicolon)")
	trimWhitespace := flag.Bool("trim", true, "Trim whitespace at head and tail of query")
	allowQuestionMarks := flag.Bool("questionmarks", true, "Determines whether or not the to fill question mark placeholders with sample data for linting")

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

	query := strings.Join(flag.Args(), " ")

	err = linter.CheckStatement(query)
	if err != nil {
		fmt.Println(err)
	}

}
