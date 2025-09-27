package grep

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	cfg "grep_utility/pkg/config"
)

// Data хранит данные, необходимые для обработки
type Data struct {
	Config                 cfg.GrepConfig
	beforeBuffer           []string
	linesToPrintAfterMatch int
	matchCount             int
	lineNumber             int
	lastPrinted            int
}

// newData - конструктор для Data
func newData(config cfg.GrepConfig) *Data {
	return &Data{
		Config:                 config,
		beforeBuffer:           make([]string, 0, config.B),
		linesToPrintAfterMatch: 0,
		matchCount:             0,
		lineNumber:             0,
		lastPrinted:            -1,
	}
}

// RunGrep - основной цикл обработки
func RunGrep(config *cfg.GrepConfig, reader io.Reader) error {
	data := newData(*config)
	scanner := bufio.NewScanner(reader)

	pattern := data.Config.Pattern
	if data.Config.IgnoreCase && !data.Config.Fixed {
		pattern = "(?i)" + pattern
	}

	regex, errCompile := regexp.Compile(pattern)
	if errCompile != nil {
		return errCompile
	}

	for scanner.Scan() {
		data.lineNumber++
		line := scanner.Text()

		if isMatch(line, config, regex) {
			data.matchCount++

			if data.Config.Count {
				continue
			}

			if data.Config.B > 0 || data.Config.C > 0 {
				i := max(data.Config.B, data.Config.C)
				if i >= len(data.beforeBuffer) {
					i = 0
				}
				for i < len(data.beforeBuffer) {
					printLine(data, data.beforeBuffer[i], i+1)
					i++
				}
				data.beforeBuffer = nil
			}
			printFindLine(data, line)
			data.lastPrinted = data.lineNumber
			data.linesToPrintAfterMatch = max(data.Config.A, data.Config.C)
			continue
		}

		if data.linesToPrintAfterMatch <= 0 && (data.Config.B > 0 || data.Config.C > 0) {
			data.beforeBuffer = append(data.beforeBuffer, line)
		}

		if data.linesToPrintAfterMatch > 0 {
			if data.lineNumber > data.lastPrinted {
				printLine(data, line, data.lineNumber)
				data.lastPrinted = data.lineNumber
				data.linesToPrintAfterMatch--
			} else {
				data.linesToPrintAfterMatch--
			}
		}
	}

	if scanner.Err() != nil {
		return scanner.Err()
	}

	if data.Config.Count {
		fmt.Println(data.matchCount)
	}

	return nil
}

// isMatch - проверяет, соответствует ли строка шаблону
func isMatch(line string, config *cfg.GrepConfig, regex *regexp.Regexp) bool {
	var match bool
	if config.Fixed {
		if config.IgnoreCase {
			match = strings.Contains(strings.ToLower(line), strings.ToLower(config.Pattern))
		} else {
			match = strings.Contains(line, config.Pattern)
		}
	} else {
		match = regex.MatchString(line)
	}

	if config.Invert {
		return !match
	}
	return match
}

// printFindLine - выводит найденную строку
func printFindLine(data *Data, line string) {
	if data.Config.LineNumber {
		fmt.Printf("-> %d %s <-\n", data.lineNumber, line)
	} else {
		fmt.Println("-> " + line + " <-")
	}
}

// printLine - выводит строку
func printLine(data *Data, line string, lineNumber int) {
	if data.Config.LineNumber {
		fmt.Printf("%d %s\n", lineNumber, line)
	} else {
		fmt.Println(line)
	}
}
