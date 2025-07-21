# ModelBench (Go Edition)

ModelBench is a tool for benchmarking machine learning models using Ollama. It helps you evaluate the performance of different models on your hardware. This version is written in Go.

## How Benchmark

Install ollama 
https://ollama.com/

```
ollama pull llama3.2:1b
ollama list
```

Now you can pull models and check which are available 
git clone https://github.com/johncave/modelbench
cd modelbench 
docker run --rm -v $(pwd):/app -w /app -e OLLAMA_URL=http://host.docker.internal:11434/api/chat golang:latest go run . run article --model llama3.2:1b --iterations 5
```
Run the benchmark using docker, otherwise run with Go on the local machine:
```
go run . run article --model llama3.2:1b --iterations 5
```