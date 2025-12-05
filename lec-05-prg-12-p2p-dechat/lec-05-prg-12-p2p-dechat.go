package main

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-zeromq/zmq4"
)

// 내 IP 가져오는 함수
func get_local_ip() string {
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

// Beacon Server
func runBeaconServer(localIP, port string) {
	pub := zmq4.NewPub(context.Background())
	defer pub.Close()
	addr := fmt.Sprintf("tcp://%s:%s", localIP, port)
	if err := pub.Listen(addr); err != nil {
		fmt.Printf("Beacon Bind 실패: %v\n", err)
		return
	}
	fmt.Printf("[Server] Beacon 활성화 (tcp://%s:%s)\n", localIP, port)

	for {
		time.Sleep(1 * time.Second)
		msg := fmt.Sprintf("NAMESERVER:%s", localIP)
		pub.Send(zmq4.NewMsgString(msg))
	}
}

// User Manager
func runUserManager(localIP, port string) {
	rep := zmq4.NewRep(context.Background())
	defer rep.Close()

	addr := fmt.Sprintf("tcp://%s:%s", localIP, port)
	if err := rep.Listen(addr); err != nil {
		fmt.Printf("UserDB Bind 실패: %v\n", err)
		return
	}
	fmt.Printf("[Server] 유저 DB 활성화 (tcp://%s:%s)\n", localIP, port)

	for {
		msg, err := rep.Recv()
		if err != nil {
			break
		}

		req := string(msg.Bytes())
		fmt.Printf("[Server] 유저 등록 요청 받음: %s\n", req)

		rep.Send(zmq4.NewMsgString("ok"))
	}
}

// Relay Server
func runRelayServer(localIP, pubPort, pullPort string) {
	pub := zmq4.NewPub(context.Background())
	defer pub.Close()
	pub.Listen(fmt.Sprintf("tcp://%s:%s", localIP, pubPort))

	pull := zmq4.NewPull(context.Background())
	defer pull.Close()
	pull.Listen(fmt.Sprintf("tcp://%s:%s", localIP, pullPort))

	fmt.Printf("[Server] Relay 서버 활성화 (PUB:%s, PULL:%s)\n", pubPort, pullPort)

	for {
		msg, err := pull.Recv()
		if err != nil {
			continue
		}

		content := string(msg.Bytes())
		fmt.Printf("p2p-relay:<==> %s\n", content)

		pub.Send(zmq4.NewMsgString("RELAY:" + content))
	}
}

func searchNameserver(ipMask, localIP, port string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sub := zmq4.NewSub(ctx)
	defer sub.Close()

	for i := 1; i < 255; i++ {
		target := fmt.Sprintf("tcp://%s.%d:%s", ipMask, i, port)
		sub.Dial(target)
	}
	sub.SetOption(zmq4.OptionSubscribe, "NAMESERVER")

	msgChan := make(chan string)
	go func() {
		msg, err := sub.Recv()
		if err == nil {
			msgChan <- string(msg.Bytes())
		}
	}()

	select {
	case res := <-msgChan:
		parts := strings.Split(res, ":")
		if len(parts) >= 2 && parts[0] == "NAMESERVER" {
			return parts[1]
		}
	case <-ctx.Done():
		return ""
	}
	return ""
}

func main() {
	fmt.Println("ZMQ 더티 P2P 구현 프로젝트")
	myIP := get_local_ip()
	fmt.Println("My IP:", myIP)

	userName := os.Args[1]

	mask := get_ip_mask(myIP)
	fmt.Println("IP Mask:", mask)

	portNameserver := "9001"
	portChatPub := "9002"
	portChatColl := "9003"
	portSubscribe := "9004"

	p2pServerIP := searchNameserver(mask, myIP, portNameserver)
	fmt.Println("P2P Server IP:", p2pServerIP)

	if p2pServerIP == "" {
		p2pServerIP = myIP
		fmt.Println("p2p server가 없으니까 내가 서버가 됨.")

		go runBeaconServer(myIP, portNameserver)
		go runUserManager(myIP, portSubscribe)
		go runRelayServer(myIP, portChatPub, portChatColl)

		time.Sleep(500 * time.Millisecond)
	} else {
		fmt.Printf("p2p server found at %s. I am a Client.\n", p2pServerIP)
	}

	fmt.Println("서버가 실행 중입니다...")
	req := zmq4.NewReq(context.Background())
	defer req.Close()
	req.Dial(fmt.Sprintf("tcp://%s:%s", p2pServerIP, portSubscribe))
	req.Send(zmq4.NewMsgString(fmt.Sprintf("%s:%s", myIP, userName)))

	if msg, err := req.Recv(); err == nil && string(msg.Bytes()) == "ok" {
		fmt.Println("p2p 서버에 유저 등록 완료")
	} else {
		fmt.Println("유저 등록 실패")
	}

	// SUB
	sub := zmq4.NewSub(context.Background())
	defer sub.Close()
	sub.Dial(fmt.Sprintf("tcp://%s:%s", p2pServerIP, portChatPub))
	sub.SetOption(zmq4.OptionSubscribe, "RELAY")

	// PUSH
	push := zmq4.NewPush(context.Background())
	defer push.Close()
	push.Dial(fmt.Sprintf("tcp://%s:%s", p2pServerIP, portChatColl))

	fmt.Println("starting autonomous message transmit and receive scenario.")

	// 수신용 채널
	msgChan := make(chan string)

	// 백그라운드에서 수신 대기
	go func() {
		for {
			msg, err := sub.Recv()
			if err != nil {
				continue
			}
			msgChan <- string(msg.Bytes())
		}
	}()

	rand.Seed(time.Now().UnixNano())

	for {
		select {
		case content := <-msgChan:
			if len(content) > 6 {
				parts := strings.Split(content, ":")
				if len(parts) >= 3 {
					fmt.Printf("p2p-recv::<<== %s:%s\n", parts[1], parts[2])
				} else {
					fmt.Printf("p2p-recv::<<== %s\n", content[6:])
				}
			}
		case <-time.After(100 * time.Millisecond):
			r := rand.Intn(100) + 1  // 100ms 타임아웃
			if r < 10 {
				time.Sleep(3 * time.Second)
				msg := fmt.Sprintf("(%s,%s:ON)", userName, myIP)
				push.Send(zmq4.NewMsgString(msg))
				fmt.Println("p2p-send::==>>", msg)
			} else if r > 90 {
				time.Sleep(3 * time.Second)
				msg := fmt.Sprintf("(%s,%s:OFF)", userName, myIP)
				push.Send(zmq4.NewMsgString(msg))
				fmt.Println("p2p-send::==>>", msg)
			}
		}
	}
}