package net

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"

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
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected from %s\n", remoteAddr)

connected:
	for {
		message, err := rw.ReadString('\n')
		switch {
		case err == io.EOF:
			break connected
		case err != nil:
			utils.CheckError(err)
		}

		message = strings.Trim(message, "\n ")
		parsed := strings.Split(message, ":")
		code, data := parsed[0], parsed[1]
		fmt.Printf("Got message with code: %s, and data: %s\n", code, data)

		switch code {
		case constants.INIT_MSG:
			// TODO: generate challenge
			message = fmt.Sprintf("%s:%s", constants.CHALLENGE_MSG, "challenge-data")
			_, err = rw.WriteString(message + "\n")
			utils.CheckError(err)
		case constants.SOLUTION_MSG:
			// TODO: validate solution
			// TODO: get poem
			message = fmt.Sprintf("%s:%s", constants.GRANT_MSG, "grant-data")
			_, err = rw.WriteString(message + "\n")
			utils.CheckError(err)
		default:
			_, err = rw.WriteString("Unrecognized message code\n")
			utils.CheckError(err)
		}

		err = rw.Flush()
		if err != nil {
			fmt.Printf("Unable to send message: %s\n", err)
			break
		}
	}

	fmt.Printf("Client at %s disconnected\n", remoteAddr)
}
