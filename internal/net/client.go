package net

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/AlexanderMac/faraway-chal/internal/constants"
	"github.com/AlexanderMac/faraway-chal/internal/utils"
)

type Client struct {
	srvAddr string
}

func NewClient(srvAddr string) *Client {
	return &Client{srvAddr}
}

func (client *Client) Start() {
	conn, err := net.Dial("tcp", client.srvAddr)
	utils.CheckErrorWithMessage(err, fmt.Sprintf("Unable to connect to remote server %s", client.srvAddr))
	defer conn.Close()
	fmt.Printf("Connected to %s\n", client.srvAddr)

	var outputMsg = Message{
		Code: constants.INIT_MSG,
	}
	for {
		fmt.Printf("Sending message: %v\n", outputMsg)
		inputMsg, err := client.handleMessage(conn, &outputMsg)
		utils.CheckError(err)

		fmt.Printf("Got message: %v\n", inputMsg)
		switch inputMsg.Code {
		case constants.CHALLENGE_MSG:
			outputMsg = Message{
				Code: constants.SOLUTION_MSG,
				Data: "solution",
			}
		case constants.GRANT_MSG:
			fmt.Printf("Granted access, the response data: %s\n", inputMsg.Data)
			os.Exit(0)
		default:
			fmt.Printf("Unrecognized message code: %s. Resending >init message\n", inputMsg.Code)
			outputMsg = Message{
				Code: constants.INIT_MSG,
				Data: "",
			}
		}
	}
}

func (client *Client) handleMessage(conn net.Conn, outputMsg *Message) (*Message, error) {
	w := bufio.NewWriter(conn)
	err := utils.GobWrite(w, outputMsg)
	utils.CheckError(err)

	inputMsg, err := utils.GobRead[Message](bufio.NewReader(conn))
	utils.CheckError(err)

	return inputMsg, nil
}
