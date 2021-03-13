package ecpg

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"
)

// Config handles all of the configurable adjustments to consider when evaluating a query
// DependencyChecker is optional if you have another method to check if ecpg is present
type Config struct {
	// Will automatically add a semicolon to the end of the statement if not present
	AddSemiColon bool
	// Will trim whitespace on incoming statement
	TrimWhiteSpace bool
	// Allow question mark as placeholder for compatibility reasons.
	QuestionMarks bool
	// Checks if ecpg cli tool dependency is present (optional)
	DependencyChecker func(string) bool
}

// ECPG is a object for accesing the ecpg tool
type ECPG struct {
	config *Config
}

const depECPG = "ecpg --help"

// NewECPG provides an ECPG object and checks if ecpg is installed
// in this environment. Throws error if ecpg is not present.
func NewECPG(c *Config) (*ECPG, error) {

	if c.DependencyChecker == nil {
		c.DependencyChecker = checkDependencies
	}

	if !c.DependencyChecker(depECPG) {
		return nil, errors.New("ecpg is required to use this library, please see README.config.md for ECPG library requirements")
	}

	return &ECPG{
		config: c,
	}, nil

}

// CheckStatement checks if a statement is valid PostgreSQL syntax
// Empty checks are done on the stmt string
func (e *ECPG) CheckStatement(stmt string) []error {
	if stmt == "" {
		return []error{errors.New("statement string is empty")}
	}

	if e.config.TrimWhiteSpace {
		stmt = strings.TrimSpace(stmt)
	}

	if e.config.QuestionMarks {
		stmt = strings.Replace(stmt, "?", "\"placeholder\"", -1)
	}

	if e.config.AddSemiColon && !strings.HasSuffix(stmt, ";") {
		stmt = fmt.Sprintf("%s;", stmt)
	}

	return executeECPGCommand(stmt)

}

// executeECPGCommand executes a sql statement against the ecpg tool
func executeECPGCommand(stmt string) []error {
	collectedErrorStrings := []error{}

	splitStatement := strings.SplitAfter(stmt, ";")
	for _, element := range splitStatement {
		if element == "" {
			continue
		}
		if !strings.HasPrefix(element, "EXEC SQL") {
			element = fmt.Sprintf("EXEC SQL %s", element)
		}
		element = strings.Replace(element, "\n", "", -1)

		args := []string{"-o", "-", "-"}
		cmd := exec.Command("ecpg", args...)

		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, element)
		}()

		out, err := cmd.CombinedOutput()
		if err != nil {
			collectedErrorStrings = append(collectedErrorStrings, errors.New(pullErrorData(string(out))))
		}
	}
	if len(collectedErrorStrings) != 0 {
		return collectedErrorStrings
	}
	return nil
}

// pullErrorData grabs the first line of the out text (the actual error)
// then removes the common, unhelpful location indicator and returns
func pullErrorData(out string) string {
	outs := strings.Split(out, "\n")

	for _, v := range outs {
		if strings.HasPrefix(v, "stdin:1:") {
			// remove `stdin:1:` & return
			return strings.TrimSpace(strings.Replace(v, "stdin:1:", "", 1))
		}
	}

	return ""
}

// checkDependencies executes a command and verifies output
func checkDependencies(p string) bool {

	cmd := exec.Command("sh", "-c", fmt.Sprintf("command -v %v", p))

	err := cmd.Run()
	if err != nil {
		return false
	}

	return true
}
