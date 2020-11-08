
binary-trees test, an experiment of optimization.

Go implementation rejected by https://benchmarksgame-team.pages.debian.net/benchmarksgame/ as not complying with dogmatic rules of allowed optimizations.

Test question from BenchmarksGame's gatekeeper Isaac Gouy was:

"_Do you think these are just ways to avoid using Go's memory management?_"


Some other implementations:

Fastest C program uses `Apache Portable Runtime` library **to avoid using C's memory management**.  
Fastest Rust program uses `typed_arena` external crate **to avoid Rust's memory management**.

It looks like Go lacks a typed arena library.
