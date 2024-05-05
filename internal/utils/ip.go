package utils

import "net"

func CheckIPAddress(ip string) bool {
	return net.ParseIP(ip) != nil
}

func CheckCIDR(network string) bool {
	_, _, err := net.ParseCIDR(network)
	if err != nil {
		return false
	}
	return true
}

func CidrToIpRange(cidr string) ([]string, error) {
	ip, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	var ips []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); inc(ip) {
		ips = append(ips, ip.String())
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
