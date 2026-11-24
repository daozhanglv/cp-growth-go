#!/bin/sh
# Create a new personal solution note under my-solutions/YYYY/MM/ and append it to my-solutions/README.md.
# Usage:
#   sh scripts/new-solution.sh <platform> <slug> [title] [difficulty] [topics]
# Example:
#   sh scripts/new-solution.sh Codeforces cf-1234a "CF 1234A Equalize Prices Again" 800 "贪心"
set -eu

if [ "$#" -lt 2 ]; then
    echo "Usage: sh scripts/new-solution.sh <platform> <slug> [title] [difficulty] [topics]" >&2
    exit 1
fi

cd "$(dirname "$0")/.."

PLATFORM=$1
SLUG=$2
TITLE=${3:-$SLUG}
DIFFICULTY=${4:-"-"}
TOPICS=${5:-"-"}
DATE=${DATE:-$(date +%Y-%m-%d)}
YEAR=$(printf '%s' "$DATE" | cut -d- -f1)
MONTH=$(printf '%s' "$DATE" | cut -d- -f2)

REL_PATH="$YEAR/$MONTH/$SLUG.md"
TARGET="my-solutions/$REL_PATH"

if [ -e "$TARGET" ]; then
    echo "Error: $TARGET already exists" >&2
    exit 1
fi

mkdir -p "$(dirname "$TARGET")"

cat > "$TARGET" <<EOF
# $TITLE

- 链接：
- 平台和难度：$PLATFORM / $DIFFICULTY
- 日期：$DATE
- 知识点：$TOPICS
- 状态：独立完成 / 看题解后完成 / 待复习

## 思路

## 卡住的地方

## 代码

\`\`\`go

\`\`\`

## 复杂度

- 时间：
- 空间：

## 学到了什么
EOF

printf '| %s | %s | %s | %s | %s | 进行中 | [%s](./%s) |\n' \
    "$DATE" "$PLATFORM" "$TITLE" "$DIFFICULTY" "$TOPICS" "$SLUG.md" "$REL_PATH" \
    >> my-solutions/README.md

echo "Created $TARGET and updated my-solutions/README.md"
