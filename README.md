# go-toolery

## 简介 (Introduction)

自研现代化的 Go 工具库，聚焦于数据结构、算法实现、工程辅助函数与常用工具组件，旨在实际工程项目中快速解决常见问题，减少重复造轮子

## 对比无参考意义，仅供个人学习使用

纯代码量：3550 + 行

## 手写 Raft 分布式一致性算法 (Raftx)

  - 实现 Raft 角色与状态机、Leader 选举、日志复制与心跳、一致性校验与冲突处理、提交与应用日志等功能；

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

## 手写 JWT 认证 (Jwtx)

- 实现根据 payload 生成 JWT Token 和 校验 JWT Token 合法性的功能；

## 手写传统算法 (Algorithmx)

- 手写二分查找 (BinarySearch)
  - 对比  sort.Slice (560753667 ns/op vs 819695230 ns/op)，性能接近主流高性能库；
- 手写快排 (QuickSort)
  - 对比  sort.Slice (560753667 ns/op vs 819695230 ns/op)，性能接近主流高性能库；
  ```shell
  yzletter@yangzhileideMacBook-Pro go-toolery % go test ./algorithmx -bench=Benchmark -run=^$ -benchmem
  goos: darwin
  goarch: arm64
  pkg: github.com/yzletter/go-toolery/algorithmx
  cpu: Apple M1 Pro
  BenchmarkQuickSort-10                  2         560753667 ns/op        82002064 B/op       1002 allocs/op
  BenchmarkSliceSort-10                  2         819695230 ns/op        82058080 B/op       3004 allocs/op
  PASS
  ok      github.com/yzletter/go-toolery/algorithmx       4.738s
  ```

## 手写数据结构 (DataStructurex)

- 手写并发安全 Map (Concurrent_mapx)

  - 对比  sync.Map (479212142 ns/op vs 735929958 ns/op)，性能接近主流高性能库；

  ```shell
  yzletter@yangzhileideMacBook-Pro go-toolery % go test ./datastructurex/concurrent_mapx -run=^$ -bench=^Benchmark -benchtime=3s -count=1 -benchmem
  goos: darwin
  goarch: arm64
  pkg: github.com/yzletter/go-toolery/datastructurex/concurrent_mapx
  cpu: Apple M1 Pro
  BenchmarkMyMap-10             10         479212142 ns/op        429347919 B/op   6019169 allocs/op
  BenchmarkSyncMap-10            8         735929958 ns/op        523149185 B/op  13170668 allocs/op
  PASS
  ok      github.com/yzletter/go-toolery/datastructurex/concurrent_mapx   13.237s
  ```

- 手写带头结点的双向循环链表 (Listx)
- 手写集合 (Setx)
- 手写栈 (Stackx)
- 手写双端队列 (DeQueuex)
- 手写堆 (PriorityQueuex)
- 手写树 (Treex)
  - 手写二叉树
  - 手写 Trie 树

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

- 手写 Snowflake 雪花算法

  - 对比  bwmarrin/snowflake (53.28 ns/op vs 119.8 ns/op)，性能接近主流高性能库；

  ```shell
  yzletter@yangzhileideMacBook-Pro go-toolery % go test ./utilx -bench=^BenchmarkSnow -run=^$ -benchmem -count=1
  goos: darwin
  goarch: arm64
  pkg: github.com/yzletter/go-toolery/utilx
  cpu: Apple M1 Pro
  BenchmarkSnowflake-10                   21216313                53.28 ns/op            0 B/op          0 allocs/op
  BenchmarkSnowflakeByBwmarrin-10         10079784               119.8 ns/op            96 B/op          1 allocs/op
  PASS
  ok      github.com/yzletter/go-toolery/utilx    2.901s
  ```

  

  