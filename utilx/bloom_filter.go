package utilx

import (
	"bufio"
	"encoding/gob"

	"github.com/dgryski/go-farm" //google FarmHash是CityHash的继承者
	"github.com/yzletter/go-toolery/errx"
 
	"math/rand"
	"os"
)

type BloomFilter struct {
	Data      []byte   // 用 Byte 切片表示底层 Bit
	BitCount  uint     // 底层的 Bit 数
	HashSeeds []uint32 // 哈希采用的种子
}

// NewBloomFilter 构造函数, 传入哈希次数和底层数组长度, 长度需要是 8 的整数倍
func NewBloomFilter(hashCount int, length int) *BloomFilter {
	arrLen := length / 8 // 转为 Byte 数组的长度
	data := make([]byte, arrLen)

	// 构造哈希种子
	hashSeeds := make([]uint32, hashCount)
	for i := 0; i < hashCount; i++ {
		hashSeeds[i] = rand.Uint32()
	}

	return &BloomFilter{
		BitCount:  uint(length),
		HashSeeds: hashSeeds,
		Data:      data,
	}
}

// 检查该 Bit 位是否为 1
func (filter *BloomFilter) checkBit(index uint) bool {
	a := index / 8
	b := index % 8
	c := uint(1 << b)

	res := uint(filter.Data[a]) & c

	return res == c
}

// 将该 Bit 位置为 1
func (filter *BloomFilter) setBit(index uint) {
	a := index / 8
	b := index % 8
	c := uint(1 << b)

	filter.Data[a] |= byte(c)
}

// Add 将字符串添加到过滤器
func (filter *BloomFilter) Add(str string) {
	// 哈希 N 次
	for _, seed := range filter.HashSeeds {
		hashed := farm.Hash32WithSeed([]byte(str), seed) // 进行哈希
		index := uint(hashed) % filter.BitCount          // 底层下标
		filter.setBit(index)
	}
}

// Exists 判断过滤器中是否有字符串: 注意过滤器中没有一定没有，但过滤器中有不一定有
func (filter *BloomFilter) Exists(str string) bool {
	// 哈希 N 次
	for _, seed := range filter.HashSeeds {
		hashed := farm.Hash32WithSeed([]byte(str), seed) // 进行哈希
		index := uint(hashed) % filter.BitCount          // 底层下标
		if !filter.checkBit(index) {
			return false
		}
	}

	return true
}

// BloomFilter 在内存中, 进程 kill 前要进行持久化

// Dump 输出到文件
func (filter *BloomFilter) Dump(fileOut string) error {
	// 打开文件
	fout, err := os.OpenFile(fileOut, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o666)
	if err != nil {
		return errx.ErrBloomFilterDumpFailed
	}
	defer fout.Close()

	// 构造 Writer
	writer := bufio.NewWriter(fout)
	defer writer.Flush()

	// Encode
	encoder := gob.NewEncoder(writer)
	err = encoder.Encode(*filter)
	if err != nil {
		return errx.ErrBloomFilterDumpFailed
	}
	return nil
}

// LoadBloomFilter 从文件读取
func LoadBloomFilter(fileIn string) (*BloomFilter, error) {
	// 打开文件
	fin, err := os.Open(fileIn)
	if err != nil {
		return nil, errx.ErrBloomFilterLoadailed
	}

	// 构造 Reader
	reader := bufio.NewReader(fin)

	// Decode
	var filter *BloomFilter
	encoder := gob.NewDecoder(reader)
	err = encoder.Decode(&filter)
	if err != nil {
		return nil, errx.ErrBloomFilterLoadailed
	}

	return filter, nil
}
