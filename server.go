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

// User Manager (유저 등록 요청 받기)
func runUserManager(localIP, port string) {
	rep := zmq4.NewRep(context.Background()) // REP 소켓 생성
	defer rep.Close()

	addr := fmt.Sprintf("tcp://%s:%s", localIP, port)
	if err := rep.Listen(addr); err != nil {
		fmt.Printf("UserDB Bind 실패: %v\n", err)
		return
	}
	fmt.Printf("[Server] 유저 DB 활성화 (tcp://%s:%s)\n", localIP, port)

	for {
		msg, err := rep.Recv()
		if err != nil { break }
		
		req := string(msg.Bytes())
		fmt.Printf("[Server] 유저 등록 요청 받음: %s\n", req)
		
		rep.Send(zmq4.NewMsgString("ok")) // 간단히 ok 응답으로 수행 확인
	}
}

// Relay Server (중계 서버)
func runRelayServer(localIP, pubPort, pullPort string) {
	pub := zmq4.NewPub(context.Background()) // PUB: 메시지 뿌리기용
	defer pub.Close()
	pub.Listen(fmt.Sprintf("tcp://%s:%s", localIP, pubPort))

	pull := zmq4.NewPull(context.Background()) // PULL: 메시지 받기용
	defer pull.Close()
	pull.Listen(fmt.Sprintf("tcp://%s:%s", localIP, pullPort))

	fmt.Printf("[Server] Relay 서버 활성화 (PUB:%s, PULL:%s)\n", pubPort, pullPort)

	for {
		msg, err := pull.Recv()
		if err != nil { continue }
		
		content := string(msg.Bytes())
		fmt.Printf("p2p-relay:<==> %s\n", content)
		
		pub.Send(zmq4.NewMsgString("RELAY:" + content)) // 모두에게 다시 쏘기
	}
}