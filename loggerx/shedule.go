package loggerx

import "time"

var (
	loc, _ = time.LoadLocation("Asia/Shanghai") // 设置时区
)

func Schedule(work func(), hour, minute, second int) {
	const period = 86400
	ticker := time.NewTicker(1) //time.After，time.Ticker，time.Timer，time.Sleep都可以互相替换
	<-ticker.C
	defer ticker.Stop()

	now := time.Now()                                                                     // 当前时间
	target := time.Date(now.Year(), now.Month(), now.Day(), hour, minute, second, 0, loc) // 目标时间

	// 计算时间差 diff
	diff := (target.Unix()%period - now.Unix()%period + period) % period

	if diff == 0 {
		go work()
	} else {
		// 下一次执行时间
		ticker.Reset(time.Duration(diff) * time.Second)
		<-ticker.C
		go work()
	}

	// 定期执行
	ticker.Reset(time.Duration(period) * time.Second)
	for {
		<-ticker.C
		go work()
	}
}
