package main

import (
	"encoding/json"
	"fmt"
)

// EventHandler is a function signature that is used to affect messages on the socket and triggered
// depending on the type
type EventHandler func(event Event, c *Client) error

const (
	EventSendRightAnswer = "send_right_answer"
	EventSendAnswer = "send_answer"
	EventNextRound  = "next_round"
	EventStartTimer = "start_timer"
)

// SendMessageHandler will send out a message to all other participants in the chat
func SendAnswerHandler(event Event, c *Client) error {
	// Marshal Payload into wanted format
	var sendAnswerEvent SendAnswerEvent
	if err := json.Unmarshal(event.Payload, &sendAnswerEvent); err != nil {
		return fmt.Errorf("bad payload in request: %v", err)
	}

	answer, err := Db.GetAnswerById(sendAnswerEvent.AnswerId)
	if err != nil {
		return nil
	} else if answer.IsRight {
		AddPointsToPlayer(c.userId, c.gameCode, answer.Points)
	}

	return nil
}

func NextRoundSend(question QuestionDto, stats []StatDto, gameCode string) error {
	var broadMessage NextRoundEvent
	broadMessage.Stats = stats
	broadMessage.Question = question

	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to marshal broadcast message: %v", err)
	}

	// Place payload into an Event
	var outgoingEvent Event
	outgoingEvent.Payload = data
	outgoingEvent.Type = EventNextRound
	// Broadcast to all other Clients
	for client := range manager.clients {
		// Only send to clients inside the same chatroom
		if client.gameCode == gameCode {
			client.egress <- outgoingEvent
		}

	}
	return nil
}
