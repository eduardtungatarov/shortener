package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/eduardtungatarov/shortener/internal/analyzer"
)

func main() {
	singlechecker.Main(analyzer.HardStopAnalyzer)
}
