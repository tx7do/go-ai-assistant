package utils

import (
	"fmt"
	"testing"

	probing "github.com/prometheus-community/pro-bing"
)

func TestProPing(t *testing.T) {
	pinger, err := probing.NewPinger("baidu.com")
	if err != nil {
		fmt.Printf("new pinger error: %s\n", err.Error())
	}

	// Windows 必须设置，否则报错：socket: The requested protocol has not been configured into the system, or no implementation for it exists.
	pinger.SetPrivileged(true)

	pinger.Count = 3
	err = pinger.Run() // Blocks until finished.
	if err != nil {
		fmt.Printf("ping run error: %s\n", err.Error())
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats
	fmt.Println(stats)
}

func TestPing(t *testing.T) {
	fmt.Printf("ping www.baidu.com: %.2fms\n", Ping("www.baidu.com", 3))
}

func TestSystemPing(t *testing.T) {
	fmt.Printf("ping www.baidu.com: %.2fms\n", SystemPing("www.baidu.com", 3))
}
