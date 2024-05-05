package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomIP(t *testing.T) {
	for i := 0; i < 100; i++ {
		fmt.Println(RandomIP())
	}
}

func TestCheckIPAddress(t *testing.T) {
	assert.True(t, CheckIPAddress("10.40.210.253"))
	assert.False(t, CheckIPAddress("1000.40.210.253"))

	assert.True(t, CheckIPAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334"))
	assert.False(t, CheckIPAddress("2001:0db8:85a3:0000:0000:8a2e:0370:7334:3445"))
}

func TestCheckCIDR(t *testing.T) {
	assert.True(t, CheckCIDR("100.64.0.0/16"))
	assert.True(t, CheckCIDR("100.64.0.0/20"))
	assert.False(t, CheckCIDR("100.64.0.0/64"))
}

func TestCidrToIpRange(t *testing.T) {
	fmt.Println(CidrToIpRange("1.1.1.0/24"))
}

func TestGetExePath(t *testing.T) {
	fmt.Println(GetExePath())
}

func TestRandomPort(t *testing.T) {
	fmt.Println(RandomPort())
}

func TestGetOutBoundIP(t *testing.T) {
	fmt.Println(GetOutBoundIP())
}
