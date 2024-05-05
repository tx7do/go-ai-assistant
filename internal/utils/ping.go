package utils

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	probing "github.com/prometheus-community/pro-bing"
)

var (
	regAverage = regexp.MustCompile("Average = (\\d+)ms|平均 = (\\d+)ms")
)

// Ping 使用内置库ping
// @param [in] ip 对端地址
// @param [in] count 测试次数
func Ping(host string, count int) float64 {
	pinger, err := probing.NewPinger(host)
	if err != nil {
		return -1
	}

	// Windows 必须设置，否则报错：socket: The requested protocol has not been configured into the system, or no implementation for it exists.
	pinger.SetPrivileged(true)
	pinger.Count = count

	// Blocks until finished.
	err = pinger.Run()
	if err != nil {
		return -1
	}

	// get send/receive/duplicate/rtt stats
	stats := pinger.Statistics()
	return float64(stats.AvgRtt.Milliseconds())
}

// SystemPing 调用系统自带ping程序
func SystemPing(host string, count int) float64 {
	strOut, err := exec.Command("ping", host, "-n", strconv.Itoa(count)).CombinedOutput()
	if err != nil {
		return -1
	}

	strUtf8 := ConvertByte2String(strOut, GB18030)
	if matches := regAverage.FindStringSubmatch(strUtf8); matches != nil && len(matches) == 3 {
		f, err := strconv.ParseFloat(matches[2], 64)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			return -1
		}
		//fmt.Println(f)
		return f
	}

	return -1
}
