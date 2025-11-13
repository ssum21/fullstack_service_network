package main

import (
	"context"
	"fmt"
	"time"
	"github.com/go-zeromq/zmq4" // ZMQ 라이브러리 설치
)

// Beacon Server (여기 있다고 알리기 역할)
func runBeaconServer(localIP, port string) {
	pub := zmq4.NewPub(context.Background()) // PUB 소켓 생성
	defer pub.Close()
	addr := fmt.Sprintf("tcp://%s:%s", localIP, port) // Bind
	if err := pub.Listen(addr); err != nil {
		fmt.Printf("Beacon Bind 실패: %v\n", err)
		return
	}
	fmt.Printf("[Server] Beacon 활성화 (tcp://%s:%s)\n", localIP, port)

	// 무한루프 돌면서 메시지 전송
	for {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("NAMESERVER:%s", localIP)
		pub.Send(zmq4.NewMsgString(msg))
	}
}