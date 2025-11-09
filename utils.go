package main

import (
	"net"
	"strings"
)

// 내 IP 가져오는 함수
func get_local_ip() string {
	// 구글 DNS 서버 TEST
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "127.0.0.1"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

// IP 대역폭 확인
func get_ip_mask(ip string) string {
	parts := strings.Split(ip, ".")
	if len(parts) < 3 {
		return ip
	}
	return strings.Join(parts[:3], ".")
}