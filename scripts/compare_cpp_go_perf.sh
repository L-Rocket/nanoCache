#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_OUT="${ROOT_DIR}/perf_go.txt"
CPP_OUT="${ROOT_DIR}/perf_cpp.txt"

printf "[1/4] Running Go benchmarks...\n"
go test ./go-iml/cache -run '^$' -bench 'BenchmarkCache(Set|Get|ConcurrentSetGet)$' -benchmem -count=3 | tee "${GO_OUT}"

printf "\n[2/4] Configuring C++ build...\n"
cmake -S "${ROOT_DIR}/cpp-impl" -B "${ROOT_DIR}/cpp-impl/build" -DCMAKE_BUILD_TYPE=Release >/dev/null

printf "[3/4] Building C++ benchmark target...\n"
cmake --build "${ROOT_DIR}/cpp-impl/build" --target benchmark_sharded_cache --parallel >/dev/null

printf "\n[4/4] Running C++ benchmarks...\n"
"${ROOT_DIR}/cpp-impl/build/benchmark_sharded_cache" --ops 300000 --threads 16 | tee "${CPP_OUT}"

printf "\nBenchmark outputs saved to:\n  - %s\n  - %s\n" "${GO_OUT}" "${CPP_OUT}"
