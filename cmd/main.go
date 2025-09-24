package main

import (
	"fmt"
	cfg "grep_utility/pkg/config"
	"grep_utility/pkg/grep"
	"io"
	"log"
	"os"
)

func main() {
	err := runApp()
	if err != nil {
		_, errStdErr := fmt.Fprintf(os.Stderr, "%s\n", err)
		if errStdErr != nil {
			log.Fatal(errStdErr)
		}
		os.Exit(1)
	}
}
func runApp() error {
	// 1. Парсинг флагов
	grepConfig := cfg.NewFromFlags()

	// 2. Чтение файла или STDIN
	var reader io.Reader
	if grepConfig.FilePath == "" {
		reader = os.Stdin
	} else {
		file, err := os.Open(grepConfig.FilePath)
		if err != nil {
			return fmt.Errorf("error opening file: %s", err)
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
	err := grep.RunGrep(grepConfig, reader)
	if err != nil {
		return fmt.Errorf("error running Grep: %s", err)
	}
	return nil
}
