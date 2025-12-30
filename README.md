# go-toolery

## 简介 (Introduction)

自研现代化的 Go 工具库，聚焦于数据结构、算法实现、工程辅助函数与常用工具组件，旨在实际工程项目中快速解决常见问题，减少重复造轮子

## 对比无参考意义，仅供个人学习使用

纯代码量：2400 + 行

## 手写传统算法 (Algorithmx)

- 手写二分查找 (BinarySearch)

- 手写快排 (QuickSort)

## 手写数据结构 (DataStructurex)

- 手写带头结点的双向循环链表 (Listx)

- 手写集合 (Setx)

- 手写栈 (Stackx)

- 手写双端队列 (DeQueuex)

- 手写堆 (PriorityQueuex)

- 手写二叉树 (Treex)

## 手写标准库辅助 (Standardx)

- 手写 Slice 辅助 (Slicex)

- 手写 Math 辅助 (Mathx)

- 手写 Rand 辅助 (Randx)

- 手写 Function 辅助

## 手写工程辅助 (Utilx)

- 手写最小并发度负载均衡算法

- 手写 Alias 采样

- 手写 BloomFilter 布隆过滤器

- 手写 PKCS7 数据填充

- 手写 Jaccard 相似度

- 手写 Slugify 函数

## 手写 JWT 认证 (Jwtx)

## 手写高性能 Log (Loggerx)

- 实现基本日志功能及终端打印控制、堆栈追踪、定时滚动、UDP 异步缓冲聚合日志等；
- 对比 Uber-Zap (2552 ns/op vs 2758 ns/op)，性能接近主流高性能库；
```shell
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./loggerx -bench=^Benchmark -run=^$ -count=1 -benchmem    
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/loggerx
cpu: Apple M1 Pro
BenchmarkMyLogger-10              453636              2552 ns/op             616 B/op          9 allocs/op
BenchmarkZap-10                   422684              2758 ns/op             529 B/op          7 allocs/op
BenchmarkZapSugar-10              408788              2956 ns/op             537 B/op          9 allocs/op
PASS
ok      github.com/yzletter/go-toolery/loggerx  4.062s 
```

## 手写 RPC 框架 (Rpcx)

- 实现参数序列化和反序列化，实现 RPC 调用功能；
- 对比 ByteDance-Sonic (1749 ns/op vs 701.3 ns/op)，性能接近主流高性能库；
```shell
yzletter@yangzhileideMacBook-Pro go-toolery % go test ./rpcx/serializer/test -bench=^Benchmark -run=^$ -count=1 -benchmem 
goos: darwin
goarch: arm64
pkg: github.com/yzletter/go-toolery/rpcx/serializer/test
cpu: Apple M1 Pro
BenchmarkBytedance-10            1680181               701.3 ns/op           523 B/op          6 allocs/op
BenchmarkGob-10                   115826              9968 ns/op            9048 B/op        187 allocs/op
BenchmarkMySerializer-10          671557              1749 ns/op            1816 B/op         60 allocs/op
PASS
ok      github.com/yzletter/go-toolery/rpcx/serializer/test     4.970s
```

## 