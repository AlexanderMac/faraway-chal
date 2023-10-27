package net

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

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
	fmt.Printf("Connected to %s...\n", client.srvAddr)

	var outMsg = Message{
		Code: constants.INIT_MSG,
	}
	for {
		inMsg, err := client.handleMessage(conn, &outMsg)
		utils.CheckError(err)

		switch inMsg.Code {
		case constants.CHALLENGE_MSG:
			outMsg = Message{
				Code: constants.SOLUTION_MSG,
				Data: "solution",
			}
		case constants.GRANT_MSG:
			fmt.Printf("Granted access, the response data: %s\n", inMsg.Data)
			os.Exit(0)
		default:
			fmt.Printf("Unrecognized message code: %s. Resending >init\n", inMsg.Code)
			outMsg = Message{
				Code: constants.INIT_MSG,
				Data: "",
			}
		}
	}
}

func (client *Client) handleMessage(conn net.Conn, outMsg *Message) (*Message, error) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	message := fmt.Sprintf("%s:%s", outMsg.Code, outMsg.Data)
	_, err := rw.WriteString(message + "\n")
	if err != nil {
		return &Message{}, err
	}
	err = rw.Flush()
	if err != nil {
		return &Message{}, err
	}

	message, err = rw.ReadString('\n')
	if err != nil {
		return &Message{}, err
	}

	message = strings.Trim(message, "\n ")
	parsed := strings.Split(message, ":")
	inMsg := Message{
		Code: parsed[0],
		Data: parsed[1],
	}

	return &inMsg, nil
}
