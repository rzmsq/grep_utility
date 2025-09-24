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

// newGrepData - конструктор для GrepData
func newGrepData(config cfg.GrepConfig) *Data {
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
	data := newGrepData(*config)
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

		if isMatch(line, config, regex) == true {
			data.matchCount++
			if data.Config.B > 0 || data.Config.C > 0 {
				for _, preLine := range data.beforeBuffer {
					fmt.Println(preLine)
				}
				clear(data.beforeBuffer)
			}
			printLine(data, line)
			data.lastPrinted = data.lineNumber
			data.linesToPrintAfterMatch = data.Config.A
			continue
		}

		if data.linesToPrintAfterMatch > 0 {
			if data.lineNumber > data.lastPrinted {
				printLine(data, line)
				data.lastPrinted = data.lineNumber
				data.linesToPrintAfterMatch--
			} else {
				data.linesToPrintAfterMatch--
			}
		}

		if data.Config.B > 0 || data.Config.C > 0 {
			data.beforeBuffer = append(data.beforeBuffer, line)
		}
	}
	return nil
}

// isMatch - проверяет, соответствует ли строка шаблону
func isMatch(line string, config *cfg.GrepConfig, regex *regexp.Regexp) bool {
	// ... (реализуем логику с regexp, fixed и ignore case) ...
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

// printLine - выводит найденные строки и контекст
func printLine(data *Data, line string) {
	if data.Config.LineNumber {
		fmt.Printf("%d %s\n", data.lineNumber, line)
	} else {
		fmt.Println(line)
	}
}
