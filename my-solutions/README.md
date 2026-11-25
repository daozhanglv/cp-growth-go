# 我的题解和笔记

这里只放我自己写的题解、复盘笔记和可独立编译测试的 Go 代码。新增一道题时，运行脚本会自动创建 `年份/月份/题目/` 目录（含 `README.md`、`solution.go`、`solution_test.go`）并刷新下方索引：

```sh
sh scripts/new-solution.sh Codeforces cf-1234a "CF 1234A Equalize Prices Again" 800 "贪心" "https://codeforces.com/problemset/problem/1234/A"
```

做完题或修改笔记头部的状态（`进行中` / `独立完成` / `看题解后完成` / `待复习`）、知识点后，只需运行：

```sh
sh scripts/sync-solutions.sh
```

即可从各篇笔记自动重新生成统计与总表，无需手动维护两处文件。`sh scripts/check.sh` 会自动检查总表同步状态，并对 `my-solutions/` 下的所有 Go 解答执行 `go vet` 和测试。

<!-- SOLUTIONS:START -->
共记录 **0** 题（独立完成 0 · 看题解后完成 0 · 待复习 0 · 进行中 0）。

| 日期 | 平台 | 题目 | 难度 | 知识点 | 状态 | 笔记 / 代码 |
| --- | --- | --- | --- | --- | --- | --- |
<!-- SOLUTIONS:END -->
