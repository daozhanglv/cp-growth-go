#!/bin/sh
# Create a testable personal solution package + note under my-solutions/YYYY/MM/<slug>/
# and automatically sync my-solutions/README.md.
#
# Usage:
#   sh scripts/new-solution.sh <platform> <slug> [title] [difficulty] [topics] [link]
# Example:
#   sh scripts/new-solution.sh Codeforces cf-1234a "CF 1234A Equalize Prices Again" 800 "贪心" "https://codeforces.com/problemset/problem/1234/A"
set -eu

if [ "$#" -lt 2 ]; then
    echo "Usage: sh scripts/new-solution.sh <platform> <slug> [title] [difficulty] [topics] [link]" >&2
    exit 1
fi

CD_ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$CD_ROOT"

SOLUTIONS_ROOT=${SOLUTIONS_ROOT:-my-solutions}
PLATFORM=$1
SLUG=$2
TITLE=${3:-$SLUG}
DIFFICULTY=${4:-"-"}
TOPICS=${5:-"-"}
LINK=${6:-""}
DATE=${DATE:-$(date +%Y-%m-%d)}
YEAR=$(printf '%s' "$DATE" | cut -d- -f1)
MONTH=$(printf '%s' "$DATE" | cut -d- -f2)

REL_DIR="$YEAR/$MONTH/$SLUG"
TARGET_DIR="$SOLUTIONS_ROOT/$REL_DIR"

if [ -e "$TARGET_DIR" ]; then
    echo "Error: $TARGET_DIR already exists" >&2
    exit 1
fi

mkdir -p "$TARGET_DIR"

cat > "$TARGET_DIR/README.md" <<EOF
# $TITLE

- 链接：$LINK
- 平台：$PLATFORM
- 难度：$DIFFICULTY
- 日期：$DATE
- 知识点：$TOPICS
- 状态：进行中

> 状态可选：\`进行中\` / \`独立完成\` / \`看题解后完成\` / \`待复习\`。修改上方元数据后运行 \`sh scripts/sync-solutions.sh\` 即可自动刷新索引与统计。

## 思路

## 卡住的地方

## 复杂度

- 时间：
- 空间：

## 学到了什么

## 代码与测试

- 解答：[solution.go](./solution.go)
- 测试：[solution_test.go](./solution_test.go)
EOF

cat > "$TARGET_DIR/solution.go" <<'EOF'
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func solve(in io.Reader, _w io.Writer) {
	r := bufio.NewReader(in)
	out := bufio.NewWriter(_w)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(r, &n); err != nil {
		return
	}
	fmt.Fprintln(out, n)
}

func main() {
	solve(os.Stdin, os.Stdout)
}
EOF

cat > "$TARGET_DIR/solution_test.go" <<'EOF'
package main

import (
	"testing"

	"github.com/EndlessCheng/codeforces-go/main/testutil"
)

func TestSolve(t *testing.T) {
	testCases := [][2]string{
		{"1", "1"},
	}
	testutil.AssertEqualStringCase(t, testCases, 0, solve)
}
EOF

SOLUTIONS_ROOT="$SOLUTIONS_ROOT" sh scripts/sync-solutions.sh

echo "Created $TARGET_DIR (README.md, solution.go, solution_test.go) and synced $SOLUTIONS_ROOT/README.md"
