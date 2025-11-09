package main

import "fmt"

func main() {
    fmt.Println("ZMQ 구현 프로젝트")
    // 프린트 로직 출력해보기
	// TODO: IP 확인 로직 추가하기

	myIP := get_local_ip()
    fmt.Println("My IP:", myIP)
    
    mask := get_ip_mask(myIP)
    fmt.Println("IP Mask:", mask)
}