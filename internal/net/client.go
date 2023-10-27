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

func RunClient(srvAddr string) {
	conn, err := net.Dial("tcp", srvAddr)
	utils.CheckErrorWithMessage(err, fmt.Sprintf("Unable to connect to remote server %s", srvAddr))
	defer conn.Close()
	fmt.Printf("Connected to %s...\n", srvAddr)

	var outMsg = Message{
		Code: constants.INIT_MSG,
	}
	for {
		inMsg, err := handleServerMessage(conn, outMsg)
		utils.CheckErrorWithMessage(err, "Error on handle message")

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

func handleServerMessage(conn net.Conn, outMsg Message) (Message, error) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	message := fmt.Sprintf("%s:%s", outMsg.Code, outMsg.Data)
	_, err := rw.WriteString(message + "\n")
	if err != nil {
		return Message{}, err
	}
	rw.Flush()

	message, err = rw.ReadString('\n')
	if err != nil {
		return Message{}, err
	}

	message = strings.Trim(message, "\n ")
	parsed := strings.Split(message, ":")
	inMsg := Message{
		Code: parsed[0],
		Data: parsed[1],
	}

	return inMsg, nil
}
