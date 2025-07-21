package benchmarkregistry

import "github.com/johncave/modelbench/pkg/benchmark"

var benchmarks = map[string]benchmark.Benchmark{}

func Register(b benchmark.Benchmark) {
	benchmarks[b.Name()] = b
}

func Get(name string) benchmark.Benchmark {
	return benchmarks[name]
}

func List() []string {
	names := make([]string, 0, len(benchmarks))
	for name := range benchmarks {
		names = append(names, name)
	}
	return names
}
