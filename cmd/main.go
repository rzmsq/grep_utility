package main

import (
	"io"
	"log"
	"os"

	cfg "grep_utility/pkg/config"
	"grep_utility/pkg/grep"
)

// TODO: Отрефакторить main
// TODO: Обработка ошибок и выход с кодом ошибки
func main() {
	// 1. Парсинг флагов
	// ... (используем пакет flag) ...
	grepConfig := cfg.NewFromFlags()

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
	err := grep.RunGrep(grepConfig, reader)
	if err != nil {
		return
	}
}
