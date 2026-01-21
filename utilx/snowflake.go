package utilx

import (
	"os"
	"strconv"
	"sync"
	"time"
)

const (
	Epoch int64 = 1742522400000 //2025-03-21 10:00:00的时间戳，单位毫秒
)

/**
原始snowflake生成的id为有符号的int64:
+--------------------------------------------------------------------------+
| 1 Bit Unused | 39 Bit Timestamp |  10 Bit NodeID  |   14 Bit Sequence ID |
+--------------------------------------------------------------------------+
Note:
1. 没有考虑process，不同的进程在同一时间（同一个毫秒内）有可能生成重复的ID
2. 该算法依赖机器时钟
*/

// 定义各字段所占 bit，根据自己公司的实际情况，修改以下常量。第一次上线之后，这些常量不能再修改
const (
	DataCenterBits = 4                                                // 数据中心占 4 个 bit
	MachineBits    = 6                                                // 机器 ID 占 6 个 bit
	SequenceBits   = 14                                               // 递增序号占 14 个 bit
	TimestampBits  = 63 - DataCenterBits - MachineBits - SequenceBits // 时间戳所占 bit
)

// 计算各字段最大能表示多少
const (
	MaxDataCenter = 1<<DataCenterBits - 1
	MaxMachineID  = 1<<MachineBits - 1
	MaxSequence   = 1<<SequenceBits - 1
)

// 左移量
const (
	TimestampShift  = DataCenterBits + MachineBits + SequenceBits
	DataCenterShift = MachineBits + SequenceBits
	MachineShift    = SequenceBits
)

// 检查配置合法性
func init() {
	if DataCenterBits <= 0 {
		panic("Snowflake数据中心至少于留1位")
	}
	if MachineBits <= 0 {
		panic("Snowflake机器至少于留1位")
	}
	if SequenceBits <= 0 {
		panic("Snowflake序列号至少于留1位")
	}
	if TimestampBits < 39 {
		panic("Snowflake时间戳至少需要留39位，保证17年内ID不会重复")
	}
}

var (
	sf    *Snowflake
	sonce sync.Once
)

type Snowflake struct {
	dataCenterID int64
	machineID    int64
	nowMilli     int64 // 当前的时间戳 - Epoch（毫秒）表示在哪个毫秒内
	seq          int64 // 当前序号
	mu           sync.Mutex
}

func NewSnowflake() *Snowflake {

	sonce.Do(func() {
		sf = new(Snowflake)
		dataCenter := os.Getenv("SNOWFLAKE_DATACENTER") // 从环境变量里读取数据中心（或者从zookeeper、etcd等地方读取）
		if len(dataCenter) == 0 {                       // 环境变量 SNOWFLAKE_DATACENTER 可以不配，默认数据中心编号为1
			sf.dataCenterID = 1
		} else {
			if dataCenterId, err := strconv.ParseInt(dataCenter, 10, 64); err == nil {
				if dataCenterId > MaxDataCenter {
					panic("数据中心 ID 超过最大值")
				}
				sf.dataCenterID = dataCenterId
			} else {
				panic("环境变量 SNOWFLAKE_DATACENTER 必须为数字")
			}
		}

		machine := os.Getenv("SNOWFLAKE_MACHINEID")
		if len(machine) == 0 { // 环境变量 SNOWFLAKE_MACHINEID 可以不配，默认机器编号为1
			sf.machineID = 1
		} else {
			if machineId, err := strconv.ParseInt(machine, 10, 64); err == nil {
				if machineId > MaxMachineID {
					panic("机器 ID 编号超过最大值")
				}
				sf.machineID = machineId
			} else {
				panic("环境变量 SNOWFLAKE_MACHINEID 必须为数字")
			}
		}
	})
	return sf
}

func (sf *Snowflake) GenerateID() int64 {
	// 上锁
	sf.mu.Lock()
	defer sf.mu.Unlock()

	// 获取当前序号，即在当前是毫秒内第几个 ID
	now := time.Now().UnixMilli() - Epoch
	if sf.nowMilli == now {
		sf.seq++
		sf.seq %= MaxSequence
	} else {
		sf.seq = 0
		sf.nowMilli = now
	}

	return now<<TimestampShift | sf.dataCenterID<<DataCenterShift | sf.machineID<<MachineShift | sf.seq
}
