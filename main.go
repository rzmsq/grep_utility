package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// GrepConfig хранит все настройки, полученные из флагов
type GrepConfig struct {
	Pattern    string
	FilePath   string
	A, B, C    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNumber bool
}

// GrepData хранит данные, необходимые для обработки
type GrepData struct {
	Config                 GrepConfig
	beforeBuffer           []string
	linesToPrintAfterMatch int
	matchCount             int
	lineNumber             int
	lastPrinted            int
}

// newGrepData - конструктор для GrepData
func newGrepData(config GrepConfig) *GrepData {
	return &GrepData{
		Config:                 config,
		beforeBuffer:           make([]string, 0, config.B),
		linesToPrintAfterMatch: 0,
		matchCount:             0,
		lineNumber:             0,
		lastPrinted:            -1,
	}
}

func setFlag(grepConfig *GrepConfig) {
	flag.IntVar(&grepConfig.A, "A", 0, "Print N string after match")
	flag.IntVar(&grepConfig.B, "B", 0, "Print N number before match")
	flag.IntVar(&grepConfig.C, "C", 0, "Print N number around match")
	flag.BoolVar(&grepConfig.Count, "c", false, "Print count")
	flag.BoolVar(&grepConfig.IgnoreCase, "i", false, "Print count")
	flag.BoolVar(&grepConfig.Invert, "v", false, "Invert match")
	flag.BoolVar(&grepConfig.Fixed, "fix", false, "Fix match")
	flag.BoolVar(&grepConfig.LineNumber, "n", false, "Line number")
	flag.StringVar(&grepConfig.FilePath, "f", "", "Line number")
}

func main() {
	// 1. Парсинг флагов
	// ... (используем пакет flag) ...
	var grepConfig GrepConfig
	setFlag(&grepConfig)
	flag.Parse()

	// 2. Чтение файла или STDIN
	// ... (открываем файл или os.Stdin) ...
	var reader io.Reader
	if grepConfig.FilePath == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(grepConfig.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer func(file *os.File) {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)
		reader = file
	}

	// 3. Запуск основного цикла
	// ... (runGrep(&config, reader)) ...
	runGrep(&grepConfig, reader)
}

// runGrep - основной цикл обработки
func runGrep(config *GrepConfig, reader io.Reader) {
	data := newGrepData(*config)
	scanner := bufio.NewScanner(reader)

	// ... (цикл `for scanner.Scan()`) ...
	// В цикле:
	// - Считываем строку
	// - Проверяем на совпадение
	// - Обновляем буфер beforeBuffer
	// - Если совпадение, вызываем функцию для вывода
	for scanner.Scan() {
		data.lineNumber++
		line := scanner.Text()

		if isMatch(line, config) {
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
}

// isMatch - проверяет, соответствует ли строка шаблону
func isMatch(line string, config *GrepConfig) bool {
	// ... (реализуем логику с regexp, fixed и ignore case) ...
	var regex *regexp.Regexp
	if config.Fixed {
		regex = regexp.MustCompile(config.Pattern)
	} else {
		regex = regexp.MustCompile(config.Pattern + ".?")
	}
	return regex.MatchString(line)
}

// printLine - выводит найденные строки и контекст
func printLine(data *GrepData, line string) {
	if data.Config.LineNumber {
		fmt.Printf("%d %s\n", data.lineNumber, line)
	} else {
		fmt.Println(line)
	}
}
