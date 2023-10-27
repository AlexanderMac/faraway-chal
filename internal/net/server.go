package net

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/AlexanderMac/faraway-chal/internal/constants"
	"github.com/AlexanderMac/faraway-chal/internal/utils"
)

type Server struct {
	srvAddr string
}

func NewServer(srvAddr string) *Server {
	return &Server{srvAddr}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", server.srvAddr)
	utils.CheckErrorWithMessage(err, fmt.Sprintf("Unable to start server %s", server.srvAddr))
	defer listener.Close()
	fmt.Printf("Listening on %s\n", server.srvAddr)

	for {
		conn, err := listener.Accept()
		utils.CheckError(err)
		go server.handleMessage(conn)
	}
}

func (server *Server) handleMessage(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected from %s\n", remoteAddr)

connected:
	for {
		message, err := utils.GobRead[Message](reader)
		switch {
		case err == io.EOF:
			break connected
		case err != nil:
			utils.CheckError(err)
		}
		fmt.Printf("Got message: %v\n", message)

		switch message.Code {
		case constants.INIT_MSG:
			// TODO: generate challenge
			message = &Message{
				Code: constants.CHALLENGE_MSG,
				Data: "challenge-data",
			}
			err = utils.GobWrite(writer, message)
			utils.CheckError(err)
		case constants.SOLUTION_MSG:
			// TODO: validate solution
			// TODO: get poem
			message = &Message{
				Code: constants.GRANT_MSG,
				Data: "grant-data",
			}
			err = utils.GobWrite(writer, message)
			utils.CheckError(err)
		default:
			message = &Message{
				Code: constants.ERROR_MSG,
				Data: "Unrecognized message code",
			}
			err = utils.GobWrite(writer, message)
			utils.CheckError(err)
		}
	}

	fmt.Printf("Client at %s disconnected\n", remoteAddr)
}
