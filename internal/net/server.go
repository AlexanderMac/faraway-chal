package net

import (
	"bufio"
	"fmt"
	"net"
	"sync"

	"github.com/AlexanderMac/faraway-chal/internal/constants"
	"github.com/AlexanderMac/faraway-chal/internal/data"
	"github.com/AlexanderMac/faraway-chal/internal/pows"
	"github.com/AlexanderMac/faraway-chal/internal/utils"
)

type Server struct {
	serverAddr string
	mut        sync.RWMutex
	challenges map[string]bool
}

func NewServer(serverAddr string) *Server {
	return &Server{
		serverAddr: serverAddr,
		challenges: make(map[string]bool),
	}
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", server.serverAddr)
	utils.CheckError(err, fmt.Sprintf("Unable to start server %s", server.serverAddr))
	defer listener.Close()
	fmt.Printf("Listening on %s\n", server.serverAddr)

	for {
		conn, err := listener.Accept()
		utils.CheckError(err, "")
		go server.handleMessage(conn)
	}
}

func (server *Server) handleMessage(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	clientAddr := conn.RemoteAddr().String()
	fmt.Printf("Client connected from %s\n", clientAddr)
	defer func() {
		fmt.Printf("Client at %s disconnected\n", clientAddr)
	}()

	for {
		msgId, err := reader.ReadByte()
		if utils.IsEOFError(err) {
			return
		}
		utils.CheckError(err, "")
		fmt.Printf("Got message with Id: %d\n", msgId)

		switch msgId {
		case constants.CHALLENGE_MESSAGE_ID:
			err = server.handleChallengeMessage(reader, writer, clientAddr)
			if utils.IsEOFError(err) {
				return
			}
			utils.CheckError(err, "")
		case constants.SOLUTION_MESSAGE_ID:
			err = server.handleSolutionMessage(reader, writer)
			if utils.IsEOFError(err) {
				return
			}
			utils.CheckError(err, "")
		default:
			err = server.sendErrorMessage(writer, "Unrecognized message id")
			utils.CheckError(err, "")
		}
	}
}

func (server *Server) handleChallengeMessage(reader *bufio.Reader, writer *bufio.Writer, clientAddr string) error {
	_, err := utils.ReadMessage[ChallengeMessage](reader)
	if err != nil {
		return err
	}

	challenge := utils.CreateChallenge(clientAddr)
	server.mut.Lock()
	server.challenges[challenge] = true
	server.mut.Unlock()

	challengeMsg := &ChallengeMessage{
		Algorithm:  "Hashcash",
		Challenge:  challenge,
		Difficulty: constants.DIFFICULTY,
	}
	return utils.SendMessage[ChallengeMessage](writer, constants.CHALLENGE_MESSAGE_ID, challengeMsg)
}

func (server *Server) handleSolutionMessage(reader *bufio.Reader, writer *bufio.Writer) error {
	solutionMsg, err := utils.ReadMessage[SolutionMessage](reader)
	if err != nil {
		return err
	}

	server.mut.RLock()
	_, ok := server.challenges[solutionMsg.Challenge]
	server.mut.RUnlock()
	if !ok {
		return server.sendErrorMessage(writer, "Unknown challenge: "+solutionMsg.Challenge)
	}

	var hashcash pows.Hashcash
	valid, err := hashcash.Validate(solutionMsg.Challenge, solutionMsg.Solution, constants.DIFFICULTY)
	if err != nil {
		return err
	}
	if !valid {
		return server.sendErrorMessage(writer, "Incorrect solution: "+solutionMsg.Solution)
	}

	server.mut.Lock()
	delete(server.challenges, solutionMsg.Challenge)
	server.mut.Unlock()

	poem, err := data.ReadPoem()
	if err != nil {
		return err
	}

	grantMsg := &GrantMessage{
		Text: poem,
	}
	return utils.SendMessage(writer, constants.GRANT_MESSAGE_ID, grantMsg)
}

func (server *Server) sendErrorMessage(writer *bufio.Writer, text string) error {
	errMsg := &ErrorMessage{
		Text: text,
	}
	return utils.SendMessage[ErrorMessage](writer, constants.ERROR_MESSAGE_ID, errMsg)
}
