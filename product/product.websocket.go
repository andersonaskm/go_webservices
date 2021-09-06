package product

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/net/websocket"
)

type message struct {
	Data string `json:"data"`
	Type string `json:"type"`
}

func productWebSocket(ws *websocket.Conn) {

	// cria um channel para lidar com a estrutura
	done := make(chan struct{})
	fmt.Println("Connection established")
	go func(c *websocket.Conn) {
		for {
			var msg message
			err := websocket.JSON.Receive(ws, &msg)
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Printf("Received message %s \r\n", msg.Data)
		}
		close(done)
	}(ws)
loop:
	for {
		select {
		case <-done:
			fmt.Println("Connection was closed")
			break loop
		default:
			products, errGetTopTenProducts := GetTopTenProducts()
			if errGetTopTenProducts != nil {
				log.Println(errGetTopTenProducts)
				break
			}
			errSend := websocket.JSON.Send(ws, products)
			if errSend != nil {
				log.Println(errSend)
				break
			}

			time.Sleep(10 * time.Second)
		}

	}
	fmt.Println("Closing the websocket")
	defer ws.Close()
}
