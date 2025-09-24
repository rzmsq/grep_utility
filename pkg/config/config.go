package config

import "flag"

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

// NewFromFlags - Парсинг флагов
func NewFromFlags() *GrepConfig {
	cfg := &GrepConfig{}

	flag.IntVar(&cfg.A, "A", 0, "Print N string after match")
	flag.IntVar(&cfg.B, "B", 0, "Print N number before match")
	flag.IntVar(&cfg.C, "C", 0, "Print N number around match")
	flag.BoolVar(&cfg.Count, "c", false, "Print count")
	flag.BoolVar(&cfg.IgnoreCase, "i", false, "Ignore case")
	flag.BoolVar(&cfg.Invert, "v", false, "Invert match")
	flag.BoolVar(&cfg.Fixed, "F", false, "Fix match")
	flag.BoolVar(&cfg.LineNumber, "n", false, "Line number")
	flag.StringVar(&cfg.FilePath, "f", "", "File path")

	flag.Parse()

	var pattern string
	if len(flag.Args()) > 0 {
		pattern = flag.Arg(0)
	}

	cfg.Pattern = pattern

	return cfg
}
