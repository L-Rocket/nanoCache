#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GO_OUT="${GO_OUT:-${ROOT_DIR}/perf_go.txt}"
CPP_OUT="${CPP_OUT:-${ROOT_DIR}/perf_cpp.txt}"
THREADS="${THREADS:-$(getconf _NPROCESSORS_ONLN)}"
OPS="${OPS:-300000}"
RUNS="${RUNS:-5}"

printf "[1/4] Running Go benchmarks (workers=%s, work items=%s, runs=%s)...\n" "${THREADS}" "${OPS}" "${RUNS}"
: > "${GO_OUT}"
for ((run = 1; run <= RUNS; run++)); do
    printf "RUN %s\n" "${run}" | tee -a "${GO_OUT}"
    GOMAXPROCS="${THREADS}" go test ./go-iml/cache -run '^$' -bench 'BenchmarkCache(Set|Get|ConcurrentSetGet)$' -benchmem -benchtime="${OPS}x" -count=1 | tee -a "${GO_OUT}"
done

printf "\n[2/4] Configuring C++ build...\n"
cmake -S "${ROOT_DIR}/cpp-impl" -B "${ROOT_DIR}/cpp-impl/build" -DCMAKE_BUILD_TYPE=Release >/dev/null

printf "[3/4] Building C++ benchmark target...\n"
cmake --build "${ROOT_DIR}/cpp-impl/build" --target benchmark_sharded_cache --parallel >/dev/null

printf "\n[4/4] Running C++ benchmarks (workers=%s, work items=%s, runs=%s)...\n" "${THREADS}" "${OPS}" "${RUNS}"
: > "${CPP_OUT}"
for ((run = 1; run <= RUNS; run++)); do
    printf "RUN %s\n" "${run}" | tee -a "${CPP_OUT}"
    "${ROOT_DIR}/cpp-impl/build/benchmark_sharded_cache" --ops "${OPS}" --threads "${THREADS}" | tee -a "${CPP_OUT}"
done

printf "\nNormalized throughput summary (mean and sample variation):\n"
awk '
function add(name, value) {
    count[name]++
    sum[name] += value
    sumsq[name] += value * value
}
function field_before(unit,    i) {
    for (i = 1; i <= NF; i++) {
        if ($i == unit) return $(i - 1)
    }
}
function assignment(name,    i, parts) {
    for (i = 1; i <= NF; i++) {
        split($i, parts, "=")
        if (parts[1] == name) return parts[2]
    }
}
function mean(name) { return sum[name] / count[name] }
function cv(name,    average, variance) {
    if (count[name] < 2) return 0
    average = mean(name)
    variance = (sumsq[name] - count[name] * average * average) / (count[name] - 1)
    return 100 * sqrt(variance) / average
}
$1 ~ /^BenchmarkCacheSet-/ { add("go_set", field_before("ops/s")) }
$1 ~ /^BenchmarkCacheGet-/ { add("go_get", field_before("ops/s")) }
$1 ~ /^BenchmarkCacheConcurrentSetGet-/ { add("go_pairs", field_before("pairs/s")) }
$1 == "cpp_set" { add("cpp_set", assignment("ops_per_second")) }
$1 == "cpp_get" { add("cpp_get", assignment("ops_per_second")) }
$1 == "cpp_concurrent_set_get" { add("cpp_pairs", assignment("pairs_per_second")) }
END {
    printf "%-22s %14s %10s %14s %10s %10s\n", "Scenario", "Go mean/s", "Go CV", "C++ mean/s", "C++ CV", "C++/Go"
    printf "%-22s %14.0f %9.2f%% %14.0f %9.2f%% %9.2fx\n", "Set (operations)", mean("go_set"), cv("go_set"), mean("cpp_set"), cv("cpp_set"), mean("cpp_set") / mean("go_set")
    printf "%-22s %14.0f %9.2f%% %14.0f %9.2f%% %9.2fx\n", "Get (operations)", mean("go_get"), cv("go_get"), mean("cpp_get"), cv("cpp_get"), mean("cpp_get") / mean("go_get")
    printf "%-22s %14.0f %9.2f%% %14.0f %9.2f%% %9.2fx\n", "Concurrent (pairs)", mean("go_pairs"), cv("go_pairs"), mean("cpp_pairs"), cv("cpp_pairs"), mean("cpp_pairs") / mean("go_pairs")
}
' "${GO_OUT}" "${CPP_OUT}"

printf "\nBenchmark outputs saved to:\n  - %s\n  - %s\n" "${GO_OUT}" "${CPP_OUT}"
