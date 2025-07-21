# ModelBench (Go Edition)

ModelBench is a tool for benchmarking machine learning models using Ollama. It helps you evaluate the performance of different models on your hardware. This version is written in Go.

## Prerequisites

- Go 1.18+
- [Ollama](https://ollama.ai/) installed and running
- Sufficient disk space for model downloads

## Installation

1. Clone this repository
2. Build the executable:
```bash
go build
```

## Usage

To run a model benchmark:
```bash
./modelbench test --model llama2 --prompt-size 1000
```

Options:
- `--model`: Name of the Ollama model to test (e.g., llama2, mistral, codellama)
- `--prompt-size`: Size of the prompt in characters (default: 1000)
- `--iterations`: Number of test iterations (default: 3)
- `--warmup`: Number of warmup runs (default: 1)
