# Compress buffer resize, Phase 1 evidence

Raw benchmark output for the `compress` package, collected on the foxy fork.

- Run URL: https://github.com/foxyan-ru/ch-go/actions/runs/39734081865
- Commit under test: `d09ef4452753a38eb55327539bfd00105df5aa44` (branch `bench/compress`)
- Baseline upstream commit: `b77c7939d345eb0ba1388c2d2cd1a104f5355bac`
- Runner: ubuntu-24.04, linux/amd64, AMD EPYC 9V45 96-Core Processor
- Go: 1.26
- benchstat: `golang.org/x/perf/cmd/benchstat@v0.0.0-20260929162123-406019bb8b68`
- Command: `go test -run '^$' -bench 'BenchmarkReader_Read|BenchmarkWriter_Compress' -benchmem -count=10 ./compress`

`old.txt` is the baseline build, `new.txt` is the build with the `growfast` tag, which resizes
the output buffer without re-zeroing the region that the codec overwrites. Both were produced in
the same job, on the same runner.

`benchstat.txt` is the comparison of the two.
