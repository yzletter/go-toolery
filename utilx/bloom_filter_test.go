package utilx_test

import (
	"testing"

	"github.com/yzletter/go-toolery/utilx"
)

func TestBloomFilter(t *testing.T) {
	strs := []string{"天青色等烟雨", "而我在等你", "炊烟袅袅升起", "隔江千万里"}
	filter := utilx.NewBloomFilter(8, 1<<20) // 8 次哈希, 底层 1M 个 Bit

	// 将前两个字符串存入
	filter.Add(strs[0])
	filter.Add(strs[1])

	// 判断前两个存在，后两个不存在
	if !filter.Exists(strs[0]) {
		t.Fail()
	}
	if !filter.Exists(strs[1]) {
		t.Fail()
	}
	if filter.Exists(strs[2]) {
		t.Fail()
	}
	if filter.Exists(strs[3]) {
		t.Fail()
	}

	// 测试文件导入导出功能
	filePath := "./bloom_filter.bin"
	_ = filter.Dump(filePath)
	newFilter, _ := utilx.LoadBloomFilter(filePath)
	if newFilter == nil {
		t.Fail()
	} else {
		// 继续判断前两个存在，后两个不存在
		if !newFilter.Exists(strs[0]) {
			t.Fail()
		}
		if !newFilter.Exists(strs[1]) {
			t.Fail()
		}
		if newFilter.Exists(strs[2]) {
			t.Fail()
		}
		if newFilter.Exists(strs[3]) {
			t.Fail()
		}
	}
}

// go test -v ./utilx -run=^TestBloomFilter$ -count=1
/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test -v ./utilx -run=^TestBloomFilter$ -count=1
=== RUN   TestBloomFilter
--- PASS: TestBloomFilter (0.00s)
PASS
ok      github.com/yzletter/go-toolery/utilx    1.096s
*/
