package main

import "fmt"

func main() {
    fmt.Println("ZMQ 더티 P2P 구현 프로젝트")
    // 프린트 로직 출력해보기
	// TODO: IP 확인 로직 추가하기

	myIP := get_local_ip()
    fmt.Println("My IP:", myIP)
    
    mask := get_ip_mask(myIP)
    fmt.Println("IP Mask:", mask)

	go runBeaconServer(myIP, "5555")
	go runUserManager(myIP, "5560")
	go runRelayServer(myIP, "5570", "5580")

	fmt.Println("서버가 실행 중입니당...")
	// 무한루프 돌면서 서버 유지시켜야 한다는데 이유 찾아보기
	select {}
}