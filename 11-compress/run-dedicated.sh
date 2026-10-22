#!/usr/bin/env bash
# Benchmark the ch-go compress buffer resize (Phase 1, P1) on a dedicated host.
# Run on a fresh Linux x86_64 instance. Results land in $HOME/ch-go-bench and $HOME/machine.txt.
set -euo pipefail

GO_VERSION="1.26.0"
BENCHSTAT_VERSION="v0.0.0-20260929162123-406019bb8b68"
COUNT="${COUNT:-20}"
BENCH='BenchmarkReader_Read|BenchmarkWriter_Compress'
OUT="$HOME/ch-go-bench"

echo "==> machine info"
{
  grep -m1 'model name' /proc/cpuinfo || true
  echo "cores: $(nproc)"
  uname -a
} | tee "$HOME/machine.txt"

export PATH="$PATH:/usr/local/go/bin"

if ! command -v go >/dev/null 2>&1; then
  echo "==> installing Go ${GO_VERSION}"
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tgz
  sudo rm -rf /usr/local/go
  sudo tar -C /usr/local -xzf /tmp/go.tgz
  rm -f /tmp/go.tgz
fi
echo "==> $(go version)"

if ! command -v git >/dev/null 2>&1; then
  echo "==> installing git"
  if command -v dnf >/dev/null 2>&1; then
    sudo dnf install -y git
  else
    sudo apt-get update -qq && sudo apt-get install -y -qq git
  fi
fi

mkdir -p "$OUT"
cd "$OUT"
if [ ! -d ch-go ]; then
  git clone https://github.com/foxyan-ru/ch-go.git
fi
cd ch-go
git checkout bench/compress
echo "==> at commit $(git rev-parse HEAD)"

echo "==> installing benchstat ${BENCHSTAT_VERSION}"
go install "golang.org/x/perf/cmd/benchstat@${BENCHSTAT_VERSION}"
export PATH="$PATH:$(go env GOPATH)/bin"

echo "==> checking both variants build and pass tests"
go test ./compress ./proto
go test -tags growfast ./compress ./proto

echo "==> baseline run (count=${COUNT})"
go test -run '^$' -bench "$BENCH" -benchmem -count="$COUNT" ./compress | tee "$OUT/old.txt"

echo "==> variant run -tags growfast (count=${COUNT})"
go test -run '^$' -bench "$BENCH" -benchmem -count="$COUNT" -tags growfast ./compress | tee "$OUT/new.txt"

echo ""
echo "==================== benchstat ===================="
benchstat "$OUT/old.txt" "$OUT/new.txt" | tee "$OUT/benchstat.txt"
echo ""
echo "Files ready:"
echo "  $HOME/machine.txt"
echo "  $OUT/old.txt"
echo "  $OUT/new.txt"
echo "  $OUT/benchstat.txt"
