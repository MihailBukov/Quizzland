package main

import (
	"errors"
	"math/rand"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	rand.New(rand.NewSource(time.Now().UnixNano()))
	randomString := make([]byte, length)

	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}

func CreateGame(userId uint, quizId uint) (string, error) {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return "", err
	}

	if acc.IsInGame {
		return "", errors.New("cannot create a game user is already in an active one")
	}

	quiz, err := Db.GetQuizById(quizId)
	if err != nil {
		return "", err
	}

	stats := make([]Stat, 1)
	stats = append(stats, Stat{
		PlayerId: acc.Id,
		Player:   *acc,
		Score:    0,
	})

	var code string
	for {
		code = GenerateRandomString(6)
		if _, err := Db.GetGameByCode(code); err == nil {
			break
		}
	}

	game := Game{
		IsActive:        true,
		IsInProgress:    false,
		Code:            code,
		CreatorId:       acc.Id,
		Creator:         *acc,
		Stats:           stats,
		QuizId:          quiz.Id,
		ActiveQuiz:      *quiz,
		CurrentQuestion: 0,
	}
	err = Db.SaveGame(&game)
	if err != nil {
		return "", err
	}

	return game.Code, nil
}

func JoinGame(gameCode string, userId uint) error {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	if acc.IsInGame {
		return errors.New("cannot join a game user is already in an active one")
	}

	game, err := Db.GetGameByCode(gameCode)
	if err != nil {
		return err
	} else if game.IsInProgress {
		return errors.New("cannot join game in progress")
	}

	game.Stats = append(game.Stats, Stat{
		PlayerId: acc.Id,
		Player:   *acc,
		Score:    0,
	})

	err = Db.SaveGame(game)
	if err != nil {
		return err
	}

	return nil
}

// Should be run using gourotine
func StartGame(gameCode string, userId uint) error {
	game, err := Db.GetGameByCode(gameCode)
	if err != nil {
		return err
	} else if game.IsInProgress {
		return errors.New("cannot start game in progress")
	} else if game.Creator.Id != userId {
		return errors.New("cannot start game you are not the creator of")
	}

	game.IsInProgress = true
	err = Db.SaveGame(game)
	if err != nil {
		return err
	}

	for _, question := range game.ActiveQuiz.Questions {
		stats := make([]StatDto, 4)
		for _, stat := range game.Stats {
			stats = append(stats, *createStatDto(stat))
		}
		NextRoundSend(*CreateQuestionDto(question), stats, game.Code)
		timer := time.NewTimer(time.Duration(question.Time) * time.Second)
		<-timer.C
	}

	game.IsInProgress = false
	game.IsActive = false
	err = Db.SaveGame(game)
	if err != nil {
		return err
	}

	return nil
}

func AddPointsToPlayer(userId uint, gameCode string, points uint) error {
	game, err := Db.GetGameByCode(gameCode)
	if err != nil {
		return err
	}

	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	for i := range game.Stats {
		if game.Stats[i].Player.Id == acc.Id {
			game.Stats[i].Score += points
			break
		}
	}

	return nil
}

func createStatDto(stat Stat) *StatDto {
	return &StatDto{
		Id:         stat.Id,
		PlayerName: stat.Player.Username,
		Score:      stat.Score,
	}
}
