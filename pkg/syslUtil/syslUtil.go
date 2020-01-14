package syslUtil

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/Joshcarp/sysl_testing/pkg/command"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func Parse(input, syslCommand string) (string, error) {
	// Declare a Memory filesystem and Logger
	fs := afero.NewMemMapFs()
	logger := logrus.New()

	// Replace any input files with /tmp.sysl
	re := regexp.MustCompile(`\w*\.sysl`)
	syslCommand = re.ReplaceAllString(syslCommand, "/tmp.sysl")

	f, err := fs.Create("/tmp.sysl")
	check(err)

	// Write String input to file
	_, e := f.Write([]byte(input))
	check(e)

	// Replace any output files with project.svg
	re = regexp.MustCompile(`(?m)(?:-o)\s"?([\S]+)`)
	syslCommand = re.ReplaceAllString(syslCommand, "-o project.svg")

	args, err := parseCommandLine(syslCommand)
	check(err)

	// Execute sysl
	command.Main2(args, fs, logger, command.Main3)

	output, err := afero.ReadFile(fs, "project.svg")

	return string(output), err
}

// parseCommandLine is from user laurent on stackoverflow to split commands into string slices
func parseCommandLine(cmd string) ([]string, error) {
	var args []string
	state := "start"
	current := ""
	quote := "\""
	escapeNext := true
	for i := 0; i < len(cmd); i++ {
		c := cmd[i]

		if state == "quotes" {
			if string(c) != quote {
				current += string(c)
			} else {
				args = append(args, current)
				current = ""
				state = "start"
			}
			continue
		}

		if escapeNext {
			current += string(c)
			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if c == '"' || c == '\'' {
			state = "quotes"
			quote = string(c)
			continue
		}

		if state == "arg" {
			if c == ' ' || c == '\t' {
				args = append(args, current)
				current = ""
				state = "start"
			} else {
				current += string(c)
			}
			continue
		}

		if c != ' ' && c != '\t' {
			state = "arg"
			current += string(c)
		}
	}

	if state == "quotes" {
		return []string{}, errors.New(fmt.Sprintf("Unclosed quote in cmd line: %s", cmd))
	}

	if current != "" {
		args = append(args, current)
	}

	return args, nil
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
