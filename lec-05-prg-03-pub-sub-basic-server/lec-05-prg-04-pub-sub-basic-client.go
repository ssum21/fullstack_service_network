//
//  Weather update client.
//  Connects SUB socket to tcp://localhost:5556
//  Collects weather updates and finds avg temp in zipcode
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	zmq.SetRetryAfterEINTR(true)

	//  Socket to talk to server
	fmt.Println("Collecting updates from weather server...")
	subscriber, _ := zmq.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5556")

	//  Subscribe to zipcode, default is NYC, 10001
	filter := "10001 "
	if len(os.Args) > 1 {
		filter = os.Args[1] + " "
	}
	subscriber.SetSubscribe(filter)
	update_nbr := 0
	for update_nbr < 20 {

		msg, _ := subscriber.Recv(0)

		if msgs := strings.Fields(msg); len(msgs) >= 3 {
			temperature, _ := strconv.Atoi(msgs[1])

			fmt.Printf("Receive temperature for zipcode '%s' was %d F\n",
				strings.TrimSpace(filter), temperature)

			update_nbr++
		}
	}
}
