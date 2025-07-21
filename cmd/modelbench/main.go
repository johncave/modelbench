package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	_ "github.com/johncave/modelbench/benchmarks"
	"github.com/johncave/modelbench/pkg/benchmarkregistry"
)

func main() {
	list := flag.Bool("list", false, "List available benchmarks")
	benchName := flag.String("bench", "", "Benchmark to run")
	flag.Parse()

	if *list {
		fmt.Println("Available benchmarks:")
		for _, name := range benchmarkregistry.List() {
			fmt.Println("-", name)
		}
		return
	}

	if *benchName == "" {
		fmt.Println("Please specify a benchmark to run with --bench or use --list to see available benchmarks.")
		os.Exit(1)
	}

	args := make(map[string]string)
	for _, arg := range flag.Args() {
		parts := strings.SplitN(arg, "=", 2)
		if len(parts) == 2 {
			args[parts[0]] = parts[1]
		}
	}

	b := benchmarkregistry.Get(*benchName)
	if b == nil {
		fmt.Printf("Benchmark '%s' not found.\n", *benchName)
		os.Exit(1)
	}

	fmt.Printf("Running benchmark: %s\n", b.Name())
	if err := b.Run(args); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
