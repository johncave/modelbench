package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/johncave/modelbench/benchmarks" // Import for side effects (benchmark registration)
	"github.com/johncave/modelbench/pkg/benchmarkregistry"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "list":
		listBenchmarks()
	case "run":
		runBenchmark(os.Args[2:])
	case "help":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  modelbench list                     List available benchmarks")
	fmt.Println("  modelbench run <benchmark> [flags]  Run a specific benchmark")
	fmt.Println("\nExample:")
	fmt.Println("  modelbench run article --model llama2 --prompt \"Write about cats\" --iterations 3")
}

func listBenchmarks() {
	fmt.Println("Available benchmarks:")
	for _, name := range benchmarkregistry.List() {
		b := benchmarkregistry.Get(name)
		fmt.Printf("  - %-15s %s\n", name, b.Description())
	}
}

func runBenchmark(args []string) {
	if len(args) < 1 {
		fmt.Println("Error: benchmark name required")
		fmt.Println("\nUsage:")
		fmt.Println("  modelbench run <benchmark> [flags]")
		os.Exit(1)
	}

	benchName := args[0]
	b := benchmarkregistry.Get(benchName)
	if b == nil {
		fmt.Printf("Error: benchmark '%s' not found\n\n", benchName)
		listBenchmarks()
		os.Exit(1)
	}

	// Set up flags for this benchmark
	flags := flag.NewFlagSet(benchName, flag.ExitOnError)
	model := flags.String("model", "llama2", "Model to use for the benchmark")
	prompt := flags.String("prompt", "", "Prompt to use for the benchmark")
	iterations := flags.Int("iterations", 1, "Number of iterations to run")

	// Parse only the flags after the benchmark name
	if err := flags.Parse(args[1:]); err != nil {
		fmt.Printf("Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	// Set iterations
	b.SetIterations(*iterations)

	// Build args map
	benchArgs := map[string]string{
		"model": *model,
	}
	if *prompt != "" {
		benchArgs["prompt"] = *prompt
	}

	fmt.Printf("Running benchmark: %s\n", b.Name())
	fmt.Printf("Description: %s\n", b.Description())
	fmt.Printf("Iterations: %d\n\n", b.GetIterations())

	if err := b.Run(benchArgs); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
