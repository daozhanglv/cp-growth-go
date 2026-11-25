# 维护与贡献指南

## 环境和检查

使用 [go.mod](./go.mod) 声明的 Go 版本或兼容版本，运行 `go version` 确认环境。维护检查在 64 位环境运行，首次下载模块需要联网。

```sh
sh scripts/check.sh
```

脚本可以从任意工作目录调用，也可以用 `GO=/path/to/go sh scripts/check.sh` 指定 Go 可执行文件。所有测试均禁用结果缓存，并设置超时。

自动检查覆盖：
- 维护范围内 Go 文件的 `gofmt` 格式校验与 `my-solutions/README.md` 索引同步检查（`sh scripts/sync-solutions.sh --check`）；
- `copypasta` 主模板包、`main/testutil`、`leetcode/testutil`、四个 `copypasta/treap` 包、`scripts` 工作流测试，以及 `my-solutions/` 下所有个人题解包的 `go vet` 和 `-race` 测试；
- CF10C、CF14D、CF20C 的内嵌样例。

它验证公共工具、模板、个人新增题解和部分代表性历史题解，不代表所有历史题解都已通过检查。修改其他历史题目或未纳入测试的模板时，还需运行对应的有限测试。CI 在上述代码、检查脚本或模块文件变动时触发，也支持手动触发。

不要默认运行全仓库测试：`copypasta/template` 中部分测试是联网生成器，有些会循环重试；部分题目测试依赖被忽略的本地 `.txt` 文件；无尽对拍需要单独运行并手动停止。

## 新增个人题解（`my-solutions/`）

使用脚本生成带独立目录的题解包（含 `README.md`、`solution.go`、`solution_test.go`）并自动刷新索引：

```sh
sh scripts/new-solution.sh Codeforces cf-1234a "CF 1234A" 800 "贪心" "https://codeforces.com/problemset/problem/1234/A"
```

修改笔记中的状态或标签后运行 `sh scripts/sync-solutions.sh` 即可同步 `my-solutions/README.md`。

## 测试历史题解

Codeforces 同一目录包含多道题目，建议只编译本次修改的文件：

```sh
go test -count=1 -timeout=30s main/1-99/14D.go main/1-99/14D_test.go
```

新题目的测试优先使用内嵌样例，避免依赖未提交的文件：

```go
testCases := [][2]string{
    {"4", "2"},
    {"5", "6"},
}
testutil.AssertEqualStringCase(t, testCases, 0, CF10C)
```

`AssertEqualCase` 支持按整行分隔的 `input` / `output` 或 `inputCopy` / `outputCopy` 标记，首个输入标记可以省略。数据中的 `myinput`、`output values` 等子串会保留。缺少输出标记、重复输出标记等格式错误会报告测试失败。如果数据本身包含完整的标记行，使用上面的 `AssertEqualStringCase` 显式传入输入和输出。

LeetCode 测试可用 `RunLeetCodeFunc` 直接传入样例，字符串参数保留外层引号：

```go
f := func(x int) int { return x * 2 }
if err := testutil.RunLeetCodeFunc(t, f,
    [][]string{{"2"}, {"-3"}},
    [][]string{{"4"}, {"-6"}},
); err != nil {
    t.Fatal(err)
}
```

`RunLeetCodeClassWithCase` 每组输入必须恰有两行：方法名数组和参数数组。输入与输出的样例数量必须相同。文件式函数测试每组有效行数为函数参数数量加返回值数量，因此不支持零参数且零返回值的函数。

使用文件式测试前，确认样例文件实际存在。`.txt`、`.in`、`.out` 默认被忽略；新提交尽量用内嵌样例，让其他人和 CI 能复现。

## 提交前

1. 对修改的 Go 文件运行 `gofmt -w`。
2. 为修复的问题补充可终止的回归测试，并运行对应测试和 `sh scripts/check.sh`。
3. 运行 `git diff --check`，检查新增文件和文档是否已纳入提交。
4. 在提交说明中描述具体触发条件、修复后的行为和验证范围。

保持题目链接和原作者署名。module 路径被大量内部导入引用，修改仓库地址时不要只改 `go.mod`。升级依赖或 Go 版本应单独说明原因并验证受影响的工具，避免和题解维护混在一起。
