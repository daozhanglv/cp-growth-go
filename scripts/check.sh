#!/bin/sh
# Run from any directory; GO may select a specific Go executable.
set -eu
cd "$(dirname "$0")/.."
GO=${GO:-go}
GOFMT=${GOFMT:-gofmt}

# 1. Check formatting on maintained packages, tests, and personal solutions.
UNFORMATTED=$(
    "$GOFMT" -l \
        copypasta/*_test.go \
        copypasta/berlekamp_massey.go \
        copypasta/cartesian_tree.go \
        copypasta/convex_hull_trick.go \
        copypasta/li_chao_tree.go \
        copypasta/skew_heap.go \
        copypasta/skip_list.go \
        copypasta/suffix_array.go \
        copypasta/van_emde_boas.go \
        copypasta/wavelet_tree.go \
        copypasta/treap \
        main/testutil \
        leetcode/testutil \
        scripts \
        my-solutions
)
if [ -n "$UNFORMATTED" ]; then
    echo "The following maintained Go files need 'gofmt -w':" >&2
    printf '%s\n' "$UNFORMATTED" >&2
    exit 1
fi

# 2. Ensure my-solutions/README.md index matches solution notes.
sh scripts/sync-solutions.sh --check

# 3. Vet and test templates, shared utilities, and workflow scripts.
"$GO" vet ./copypasta ./main/testutil ./leetcode/testutil ./copypasta/treap/... ./scripts/...
"$GO" test -race -count=1 -timeout=3m ./copypasta
"$GO" test -race -count=1 -timeout=2m ./main/testutil ./leetcode/testutil ./copypasta/treap/... ./scripts/...

# 4. Automatically vet and test any personal Go solutions under my-solutions/.
if [ -n "$(find my-solutions -name '*.go' -print -quit 2>/dev/null)" ]; then
    "$GO" vet ./my-solutions/...
    "$GO" test -race -count=1 -timeout=2m ./my-solutions/...
fi

# 5. Compile only these representative problem files to avoid unrelated drafts in the same package.
"$GO" test -count=1 -timeout=30s main/1-99/10C.go main/1-99/10C_test.go
"$GO" test -count=1 -timeout=30s main/1-99/14D.go main/1-99/14D_test.go
"$GO" test -count=1 -timeout=30s main/1-99/20C.go main/1-99/20C_test.go
