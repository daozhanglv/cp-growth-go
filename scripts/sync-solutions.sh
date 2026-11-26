#!/bin/sh
# Regenerate (or verify with --check) the statistics and solution index table
# between <!-- SOLUTIONS:START --> and <!-- SOLUTIONS:END --> in my-solutions/README.md.
set -eu

CD_ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
cd "$CD_ROOT"

CHECK_ONLY=0
if [ "${1:-}" = "--check" ]; then
    CHECK_ONLY=1
fi

SOLUTIONS_ROOT=${SOLUTIONS_ROOT:-my-solutions}
README_FILE="$SOLUTIONS_ROOT/README.md"

if [ ! -f "$README_FILE" ]; then
    echo "Error: $README_FILE not found" >&2
    exit 1
fi

TMP_DIR=$(mktemp -d)
trap 'rm -rf "$TMP_DIR"' EXIT

ROWS_RAW="$TMP_DIR/rows_raw.tsv"
: > "$ROWS_RAW"

find "$SOLUTIONS_ROOT" -mindepth 3 -maxdepth 4 -name '*.md' 2>/dev/null | sort | while IFS= read -r file; do
    rel=${file#"$SOLUTIONS_ROOT"/}
    case "$rel" in
        [0-9][0-9][0-9][0-9]/[0-9][0-9]/*/README.md|[0-9][0-9][0-9][0-9]/[0-9][0-9]/*.md) ;;
        *) continue ;;
    esac

    dir_part=$(dirname "$file")
    has_go=0
    if [ "$(basename "$file")" = "README.md" ] && [ -f "$dir_part/solution.go" ]; then
        has_go=1
    fi

    awk -v rel="$rel" -v has_go="$has_go" '
    function trim(s) {
        sub(/^[ \t\r]+/, "", s)
        sub(/[ \t\r]+$/, "", s)
        return s
    }
    function esc(s) {
        gsub(/\|/, "\\|", s)
        return s
    }
    BEGIN {
        title = ""; link = ""; platform = "-"; diff = "-"; date = "-"; topics = "-"; status = "进行中"
    }
    title == "" && /^#[ \t]+/ {
        s = $0
        sub(/^#[ \t]+/, "", s)
        title = trim(s)
        next
    }
    # Match the entire full-width colon even when awk treats UTF-8 as bytes.
    /^- 链接(：|:)/ {
        s = $0; sub(/^- 链接(：|:)/, "", s); link = trim(s); next
    }
    /^- 平台和难度(：|:)/ {
        s = $0; sub(/^- 平台和难度(：|:)/, "", s); s = trim(s)
        n = split(s, parts, "/")
        if (n >= 1) platform = trim(parts[1])
        if (n >= 2) diff = trim(parts[2])
        next
    }
    /^- 平台(：|:)/ {
        s = $0; sub(/^- 平台(：|:)/, "", s); v = trim(s); if (v != "") platform = v; next
    }
    /^- 难度(：|:)/ {
        s = $0; sub(/^- 难度(：|:)/, "", s); v = trim(s); if (v != "") diff = v; next
    }
    /^- 日期(：|:)/ {
        s = $0; sub(/^- 日期(：|:)/, "", s); v = trim(s); if (v != "") date = v; next
    }
    /^- 知识点(：|:)/ {
        s = $0; sub(/^- 知识点(：|:)/, "", s); v = trim(s); if (v != "") topics = v; next
    }
    /^- 状态(：|:)/ {
        s = $0; sub(/^- 状态(：|:)/, "", s); v = trim(s); if (v != "") status = v; next
    }
    END {
        if (title == "") title = rel
        display_title = esc(title)
        if (link != "") {
            display_title = "[" display_title "](" link ")"
        }
        entry_links = "[笔记](./" rel ")"
        if (has_go == 1) {
            code_rel = rel
            sub(/README\.md$/, "solution.go", code_rel)
            entry_links = entry_links " · [代码](./" code_rel ")"
        }
        printf "%s\t%s\t%s\t| %s | %s | %s | %s | %s | %s | %s |\n",
            date, rel, status,
            esc(date), esc(platform), display_title, esc(diff), esc(topics), esc(status), entry_links
    }
    ' "$file" >> "$ROWS_RAW"
done

SORTED_ROWS="$TMP_DIR/rows_sorted.tsv"
LC_ALL=C sort -t '	' -k1,1r -k2,2 "$ROWS_RAW" > "$SORTED_ROWS"

GENERATED_BLOCK="$TMP_DIR/block.md"
awk -F '\t' '
BEGIN {
    total = 0; indep = 0; editorial = 0; review = 0; wip = 0
}
NF >= 4 {
    total++
    st = $3
    if (st == "独立完成") indep++
    else if (st == "看题解后完成") editorial++
    else if (st == "待复习") review++
    else wip++
    # Strip the three metadata fields and keep any tabs inside the Markdown row.
    row = $0
    sub(/^[^\t]*\t[^\t]*\t[^\t]*\t/, "", row)
    rows[total] = row
}
END {
    printf "<!-- SOLUTIONS:START -->\n"
    printf "共记录 **%d** 题（独立完成 %d · 看题解后完成 %d · 待复习 %d · 进行中 %d）。\n\n",
        total, indep, editorial, review, wip
    printf "| 日期 | 平台 | 题目 | 难度 | 知识点 | 状态 | 笔记 / 代码 |\n"
    printf "| --- | --- | --- | --- | --- | --- | --- |\n"
    for (i = 1; i <= total; i++) {
        print rows[i]
    }
    printf "<!-- SOLUTIONS:END -->\n"
}
' "$SORTED_ROWS" > "$GENERATED_BLOCK"

NEW_README="$TMP_DIR/README.md"
awk -v block_file="$GENERATED_BLOCK" '
BEGIN {
    in_block = 0
    replaced = 0
}
/^<!-- SOLUTIONS:START -->$/ {
    while ((getline line < block_file) > 0) {
        print line
    }
    close(block_file)
    in_block = 1
    replaced = 1
    next
}
/^<!-- SOLUTIONS:END -->$/ {
    in_block = 0
    next
}
!in_block {
    print
}
END {
    if (!replaced) {
        print ""
        while ((getline line < block_file) > 0) {
            print line
        }
        close(block_file)
    }
}
' "$README_FILE" > "$NEW_README"

if [ "$CHECK_ONLY" -eq 1 ]; then
    if ! cmp -s "$README_FILE" "$NEW_README"; then
        echo "Error: $README_FILE is out of sync with solution notes." >&2
        echo "Run: sh scripts/sync-solutions.sh" >&2
        exit 1
    fi
    exit 0
fi

cp "$NEW_README" "$README_FILE"
