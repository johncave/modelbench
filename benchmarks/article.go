package benchmarks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/johncave/modelbench/pkg/benchmarkregistry"
)

type ArticleBenchmark struct {
	iterations int
}

type Request struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Response struct {
	Model           string    `json:"model"`
	CreatedAt       time.Time `json:"created_at"`
	Message         Message   `json:"message"`
	Done            bool      `json:"done"`
	TotalDuration   int64     `json:"total_duration"`
	LoadDuration    int       `json:"load_duration"`
	PromptEvalCount int       `json:"prompt_eval_count"`
	EvalCount       int       `json:"eval_count"`
	EvalDuration    int64     `json:"eval_duration"`
}

type BenchmarkResult struct {
	TotalTime        time.Duration
	TokenCount       int
	TokensPerSec     float64
	LoadDuration     time.Duration
	EvalDuration     time.Duration
	PromptTokenCount int
	OutputTokenCount int
}

const defaultOllamaURL = "http://host.docker.internal:11434/api/chat"

func (a *ArticleBenchmark) Name() string        { return "article" }
func (a *ArticleBenchmark) Description() string { return "Generate a detailed article using Ollama" }
func (a *ArticleBenchmark) SetIterations(n int) { a.iterations = n }
func (a *ArticleBenchmark) GetIterations() int  { return a.iterations }

func (a *ArticleBenchmark) Run(args map[string]string) error {
	model, ok := args["model"]
	if !ok {
		model = "llama3:latest" // default model
	}

	prompt, ok := args["prompt"]
	if !ok {
		prompt = "Write a detailed encyclopedia article about ancient history."
	}

	var totalTokens int
	var totalDuration time.Duration

	for i := 0; i < a.iterations; i++ {
		fmt.Printf("\nIteration %d/%d:\n", i+1, a.iterations)

		start := time.Now()
		msg := Message{
			Role:    "user",
			Content: prompt,
		}
		req := Request{
			Model:    model,
			Stream:   false,
			Messages: []Message{msg},
		}
		resp, err := talkToOllama(defaultOllamaURL, req)
		if err != nil {
			return err
		}

		iterDuration := time.Since(start)
		totalDuration += iterDuration
		totalTokens += resp.EvalCount + resp.PromptEvalCount

		result := calculateMetrics(resp, iterDuration)
		printIterationResults(i+1, a.iterations, result)
	}

	// Print final summary
	fmt.Printf("\nFinal Results (over %d iterations):\n", a.iterations)
	fmt.Printf("Average Time:      %.2fs\n", totalDuration.Seconds()/float64(a.iterations))
	fmt.Printf("Total Tokens:      %d\n", totalTokens)
	fmt.Printf("Average Speed:     %.2f tokens/sec\n", float64(totalTokens)/totalDuration.Seconds())

	return nil
}

func calculateMetrics(resp *Response, totalTime time.Duration) BenchmarkResult {
	outputTokens := resp.EvalCount
	promptTokens := resp.PromptEvalCount
	totalTokens := outputTokens + promptTokens

	return BenchmarkResult{
		TotalTime:        totalTime,
		TokenCount:       totalTokens,
		TokensPerSec:     float64(outputTokens) / totalTime.Seconds(),
		LoadDuration:     time.Duration(resp.LoadDuration) * time.Microsecond,
		EvalDuration:     time.Duration(resp.EvalDuration) * time.Microsecond,
		PromptTokenCount: promptTokens,
		OutputTokenCount: outputTokens,
	}
}

func printIterationResults(iter, total int, result BenchmarkResult) {
	fmt.Printf("Time:             %.2fs\n", result.TotalTime.Seconds())
	fmt.Printf("Tokens:           %d (%d prompt, %d generated)\n",
		result.TokenCount,
		result.PromptTokenCount,
		result.OutputTokenCount,
	)
	fmt.Printf("Speed:            %.2f tokens/sec\n", result.TokensPerSec)
}

func talkToOllama(url string, ollamaReq Request) (*Response, error) {
	js, err := json.Marshal(&ollamaReq)
	if err != nil {
		return nil, err
	}
	client := http.Client{}
	httpReq, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(js))
	if err != nil {
		return nil, err
	}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()
	ollamaResp := Response{}
	err = json.NewDecoder(httpResp.Body).Decode(&ollamaResp)
	return &ollamaResp, err
}

func init() {
	benchmarkregistry.Register(&ArticleBenchmark{iterations: 1})
}
