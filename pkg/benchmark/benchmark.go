package benchmark

type Benchmark interface {
	Name() string
	Description() string
	Run(args map[string]string) error
	SetIterations(n int)
	GetIterations() int
}
