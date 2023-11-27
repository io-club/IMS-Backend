package util

import (
	"errors"
	"log"
	"net"
	"strings"
)

// 获取本机网卡 IP(内网 ip)
func GetLocalIPWithHardware() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP 地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return "", err
	}
	// 取第一个非 lo 的网卡 IP
	for _, addr = range addrs {
		// 这个网络地址是 IP 地址：ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过 IPV6
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				log.Printf("ip: %s", ipv4)
				return ipv4, nil
			}
		}
	}
	return "", errors.New("ERR_NO_LOCAL_IP_FOUND")
}

func GetLocalIPWithNet() (string, error) {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	localAddr := conn.LocalAddr().String()
	idx := strings.LastIndex(localAddr, ":")
	return localAddr[0:idx], nil
}
