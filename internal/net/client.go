package net

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/AlexanderMac/faraway-chal/internal/constants"
	"github.com/AlexanderMac/faraway-chal/internal/pows"
	"github.com/AlexanderMac/faraway-chal/internal/utils"
)

type Client struct {
	serverAddr string
}

func NewClient(serverAddr string) *Client {
	return &Client{serverAddr}
}

func (client *Client) Start() {
	conn, err := net.Dial("tcp", client.serverAddr)
	utils.CheckError(err, fmt.Sprintf("Unable to connect to remote server %s", client.serverAddr))
	defer conn.Close()
	fmt.Printf("Connected to %s\n", client.serverAddr)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	var maxAttempts = 0
	var currentMessageId = constants.CHALLENGE_MESSAGE_ID
	var solutionMsg *SolutionMessage
	for {
		fmt.Printf("Sending message with Id: %v\n", currentMessageId)
		switch currentMessageId {
		case constants.CHALLENGE_MESSAGE_ID:
			if maxAttempts > 10 {
				return
			}
			maxAttempts++
			err := client.sendChallengeMessage(writer)
			utils.CheckError(err, "")
		case constants.SOLUTION_MESSAGE_ID:
			err := client.sendSolutionMessage(writer, solutionMsg)
			utils.CheckError(err, "")
		}

		msgId, err := reader.ReadByte()
		if utils.IsEOFError(err) {
			return
		}
		utils.CheckError(err, "")
		fmt.Printf("Got response message with Id: %d\n", msgId)

		switch msgId {
		case constants.CHALLENGE_MESSAGE_ID:
			currentMessageId = constants.SOLUTION_MESSAGE_ID
			solutionMsg, err = client.handleChallengeMessage(reader)
			utils.CheckError(err, "")
		case constants.GRANT_MESSAGE_ID:
			err = client.handleGrantMessage(reader)
			utils.CheckError(err, "")
		case constants.ERROR_MESSAGE_ID:
			err = client.handleErrorMessage(reader)
			utils.CheckError(err, "")
			currentMessageId = constants.CHALLENGE_MESSAGE_ID
			solutionMsg = nil
		default:
			fmt.Printf("Unrecognized message Id: %d. Resending >challenge message\n", msgId)
			currentMessageId = constants.CHALLENGE_MESSAGE_ID
			solutionMsg = nil
		}
	}
}

func (client *Client) sendChallengeMessage(writer *bufio.Writer) error {
	challengeMsg := &ChallengeMessage{}
	return utils.SendMessage(writer, constants.CHALLENGE_MESSAGE_ID, challengeMsg)
}

func (client *Client) handleChallengeMessage(reader *bufio.Reader) (*SolutionMessage, error) {
	challengeMsg, err := utils.ReadMessage[ChallengeMessage](reader)
	if err != nil {
		return nil, err
	}

	hashcash := new(pows.Hashcash)
	solution := hashcash.Solve(challengeMsg.Challenge, challengeMsg.Difficulty)

	solutionMsg := &SolutionMessage{
		Algorithm: challengeMsg.Algorithm,
		Challenge: challengeMsg.Challenge,
		Solution:  solution,
	}
	return solutionMsg, nil
}

func (client *Client) sendSolutionMessage(writer *bufio.Writer, solutionMsg *SolutionMessage) error {
	return utils.SendMessage(writer, constants.SOLUTION_MESSAGE_ID, solutionMsg)
}

func (client *Client) handleGrantMessage(reader *bufio.Reader) error {
	grantMsg, err := utils.ReadMessage[GrantMessage](reader)
	if err != nil {
		return err
	}

	fmt.Printf("Access Granted! Text: %s\n", grantMsg.Text)
	os.Exit(0)

	return nil
}

func (client *Client) handleErrorMessage(reader *bufio.Reader) error {
	errMsg, err := utils.ReadMessage[ErrorMessage](reader)
	if err != nil {
		return err
	}

	fmt.Printf("Got error from the server: %s\n", errMsg.Text)

	return nil
}
