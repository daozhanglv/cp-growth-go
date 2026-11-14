#!/bin/sh
# Run from any directory; GO may select a specific Go executable.
set -eu
cd "$(dirname "$0")/.."
GO=${GO:-go}

"$GO" vet ./copypasta
"$GO" test -race -count=1 -timeout=3m ./copypasta
"$GO" test -race -count=1 -timeout=2m ./main/testutil ./leetcode/testutil ./copypasta/treap/...
# Compile only these problem files to avoid unrelated drafts in the same package.
"$GO" test -count=1 -timeout=30s main/1-99/10C.go main/1-99/10C_test.go
"$GO" test -count=1 -timeout=30s main/1-99/14D.go main/1-99/14D_test.go
"$GO" test -count=1 -timeout=30s main/1-99/20C.go main/1-99/20C_test.go
