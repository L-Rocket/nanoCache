# C++ and Go benchmark analysis

This analysis accompanies the normalized benchmark run from 2026-10-20 on an
Apple M4 (`darwin/arm64`). It explains what the measurements demonstrate and
keeps hypotheses separate from measured causes.

## Normalized result

Both implementations used 256 shards, a 10-second TTL, a one-second janitor
interval, 10 workers, identical precomputed key/value shapes, and 300,000 work
items per isolated process. Five runs produced:

| Scenario | C++ mean/s (CV) | Go mean/s (CV) | C++ / Go |
| --- | ---: | ---: | ---: |
| Set operations | 8,155,290 (12.35%) | 6,025,967 (3.71%) | 1.35x |
| Get operations | 14,808,860 (1.25%) | 8,873,410 (2.08%) | 1.67x |
| Concurrent Set+Get pairs | 15,396,280 (4.42%) | 11,752,032 (1.82%) | 1.31x |

The old 5.35x concurrent result was primarily a unit error: C++ reported the
Set and Get as two operations, while Go reported their pair as one benchmark
iteration. It was also affected by different shard counts and TTLs.

## Go profile evidence

A separate two-million-pair diagnostic run collected Go CPU and allocation
profiles. Profiling changes throughput, so its rate was not used in the table
above. The important CPU samples were:

| CPU function/category | Flat CPU samples |
| --- | ---: |
| `runtime.mapassign_faststr` | 17.5% |
| GC scanning (`runtime.scanobject`) | 10.2% |
| FNV-1a shard hashing | 9.0% |
| Memory release (`runtime.madvise`) | 8.5% |
| Clock reads (`runtime.walltime`) | 5.1% |
| Mutex locking, flat samples | 1.7% |

The allocation profile attributed 89% of allocated bytes to `shard.set`.
The normalized 300,000-pair benchmark reported 277 allocated bytes per pair.
Map insertion/growth and the live cache therefore create substantial heap work;
this is not just goroutine launch overhead.

The Go entry also retains the key twice: once as the map key and once in the
`entry.key` field. C++ keeps the key only in the `unordered_map` key. The extra
Go field is unused after insertion and increases the map value and the amount
of pointer-bearing live memory scanned by GC. This is a concrete implementation
asymmetry, although a short removal experiment was too noisy to assign it a
reliable percentage.

## GC control

Ten interleaved 300,000-pair runs compared the default collector with
`GOGC=off` on the same workload:

| Go mode | Mean pairs/s | CV |
| --- | ---: | ---: |
| Default GC | 9,934,200 | 6.62% |
| `GOGC=off` | 11,837,461 | 6.44% |

Disabling collection increased throughput by 19.2% in this diagnostic, which
is equivalent to roughly 16% less elapsed time. This demonstrates that GC is a
real contributor for the insertion-heavy workload, but not that GC explains
the entire language gap. `GOGC=off` also permits unbounded heap growth and is
not a production recommendation.

## Conclusion and limits

The measured Go-side costs are, in descending confidence:

1. Map insertion, growth, and the resulting heap allocation are the largest
   visible hot-path cost.
2. GC scanning and memory management measurably reduce throughput; the
   interleaved control quantified their effect for this workload.
3. Shard hashing and TTL clock reads are material CPU costs shared conceptually
   by both implementations, though their library/runtime costs differ.
4. Exclusive Set locking produces wait time, but 256 shards keep lock code from
   dominating Go's flat CPU samples in this unique-key workload.
5. Scheduler overhead is unlikely to explain much of this benchmark: only 10
   workers are launched once, so launch and scheduling costs are amortized over
   300,000 pairs.

This analysis does not isolate a single cause for the remaining C++ advantage.
A rigorous cross-runtime attribution would additionally require a C++ CPU and
allocation profile collected under the same sampling method and controlled
experiments for map capacity, entry layout, and clock calls.
