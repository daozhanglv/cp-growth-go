package scripts_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestSolutionWorkflowScripts(t *testing.T) {
	repoRoot, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}

	// Create temporary solutions root inside repoRoot so generated solution_test.go
	// shares the repository's go.mod and can be compiled and tested with `go test`.
	tmpRoot, err := os.MkdirTemp(repoRoot, ".tmp-solutions-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpRoot)

	relRoot := filepath.Base(tmpRoot)
	readmePath := filepath.Join(tmpRoot, "README.md")
	if err := os.WriteFile(readmePath, []byte("# Test Solutions\n\n<!-- SOLUTIONS:START -->\n<!-- SOLUTIONS:END -->\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	// Exercise byte-oriented awk even on machines that default to a UTF-8 locale.
	t.Setenv("LC_ALL", "C")

	runSh := func(envExtra []string, args ...string) string {
		t.Helper()
		cmd := exec.Command("sh", args...)
		cmd.Dir = repoRoot
		cmd.Env = append(os.Environ(), append([]string{"SOLUTIONS_ROOT=" + relRoot}, envExtra...)...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("sh %v failed: %v\n%s", args, err, out)
		}
		return string(out)
	}

	// 1. Create a solution with a pipe and a tab in its topics.
	runSh([]string{"DATE=2026-11-25"}, "scripts/new-solution.sh",
		"Codeforces", "cf-100a", "CF 100A Carpeting", "1100", "贪心 |\t数学", "https://codeforces.com/problemset/problem/100/A")

	// 2. Verify sync-solutions.sh --check passes immediately after creation.
	runSh(nil, "scripts/sync-solutions.sh", "--check")

	readmeBytes, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	readme := string(readmeBytes)
	if !utf8.ValidString(readme) {
		t.Fatalf("index contains invalid UTF-8:\n%s", readme)
	}
	wantRow := "| 2026-11-25 | Codeforces | [CF 100A Carpeting](https://codeforces.com/problemset/problem/100/A) | 1100 | 贪心 \\|\t数学 | 进行中 |"
	if !strings.Contains(readme, wantRow) {
		t.Fatalf("expected complete row with escaped pipe and intact metadata, got:\n%s", readme)
	}
	if !strings.Contains(readme, "共记录 **1** 题（独立完成 0 · 看题解后完成 0 · 待复习 0 · 进行中 1）") {
		t.Fatalf("unexpected stats line:\n%s", readme)
	}

	// 3. Verify generated Go solution compiles and passes `go test`.
	goBin := os.Getenv("GO")
	if goBin == "" {
		goBin = "go"
	}
	pkgPattern := "./" + filepath.ToSlash(filepath.Join(relRoot, "2026/11/cf-100a"))
	testCmd := exec.Command(goBin, "test", "-count=1", pkgPattern)
	testCmd.Dir = repoRoot
	if out, err := testCmd.CombinedOutput(); err != nil {
		t.Fatalf("generated solution_test.go failed: %v\n%s", err, out)
	}

	// 4. Update status in note README.md -> verify --check detects drift, then re-sync updates stats.
	notePath := filepath.Join(tmpRoot, "2026/11/cf-100a/README.md")
	noteBytes, err := os.ReadFile(notePath)
	if err != nil {
		t.Fatal(err)
	}
	updatedNote := strings.NewReplacer(
		"- 状态：进行中", "- 状态: 独立完成",
		"- 链接：https://codeforces.com/problemset/problem/100/A", "- 链接：",
	).Replace(string(noteBytes))
	if err := os.WriteFile(notePath, []byte(updatedNote), 0o644); err != nil {
		t.Fatal(err)
	}

	checkCmd := exec.Command("sh", "scripts/sync-solutions.sh", "--check")
	checkCmd.Dir = repoRoot
	checkCmd.Env = append(os.Environ(), "SOLUTIONS_ROOT="+relRoot)
	if err := checkCmd.Run(); err == nil {
		t.Fatal("expected sync-solutions.sh --check to fail when README is stale")
	}

	runSh(nil, "scripts/sync-solutions.sh")
	runSh(nil, "scripts/sync-solutions.sh", "--check")

	readmeBytes, err = os.ReadFile(readmePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(readmeBytes), "共记录 **1** 题（独立完成 1 · 看题解后完成 0 · 待复习 0 · 进行中 0）") {
		t.Fatalf("stats did not update after sync:\n%s", readmeBytes)
	}
	if !strings.Contains(string(readmeBytes), "| CF 100A Carpeting |") {
		t.Fatalf("expected an unlinked title when the link is empty:\n%s", readmeBytes)
	}
}
