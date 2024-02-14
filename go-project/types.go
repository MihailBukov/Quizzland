package main

import (
	"encoding/json"
)

const (
	Nuser = "nuser"
	Ruser = "ruser"
	Admin = "admin"
)

type Account struct {
	Id          uint    `json:"id" gorm:"primaryKey"`
	Username    string  `json:"username" gorm:"size:32;unique;not null"`
	Password    string  `json:"password" gorm:"size:100;not null"`
	FirstName   string  `json:"firstName" gorm:"size:32"`
	LastName    string  `json:"lastName" gorm:"size:32"`
	Email       string  `json:"email" gorm:"size:32;unique;not null"`
	Description string  `json:"description" gorm:"size:255"`
	Balance     float32 `json:"balance"`
	IsInGame    bool    `json:"isInGame"`
	Quizzes     []Quiz  `json:"quizzes" gorm:"foreignKey:OwnerId"`
	Role        string  `json:"role" gorm:"size:5"`
}

type AccountDto struct {
	Id          uint      `json:"id"`
	Username    string    `json:"username"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Description string    `json:"description"`
	Balance     float32   `json:"balance"`
	Quizzes     []QuizDto `json:"quizzes"`
}

type Question struct {
	Id                  uint     `json:"id" gorm:"primaryKey"`
	Text                string   `json:"text" gorm:"size:32"`
	Time                uint     `json:"time"`
	Answers             []Answer `json:"answers" gorm:"foreignKey:CorrespondingQuestionId"`
	CorrespondingQuizId uint     `json:"-"`
	CorrespondingQuiz   Quiz     `json:"correspondingQuiz" gorm:"foreignKey:CorrespondingQuizId;references:Id"`
}

type QuestionDto struct {
	Id      uint        `json:"id"`
	Text    string      `json:"text"`
	Time    uint        `json:"time"`
	Answers []AnswerDto `json:"answers"`
}

type Answer struct {
	Id                      uint     `json:"id" gorm:"primaryKey"`
	Points                  uint     `json:"points"`
	Text                    string   `json:"text" gorm:"size:255"`
	IsRight                 bool     `json:"isRight"`
	CorrespondingQuestionId uint     `json:"-"`
	CorrespondingQuestion   Question `json:"correspondingQuestion" gorm:"foreignKey:CorrespondingQuestionId;references:Id"`
}

type AnswerDto struct {
	Id   uint   `json:"id"`
	Text string `json:"text"`
}

type Quiz struct {
	Id          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"size:30"`
	Description string     `json:"description" gorm:"size:255"`
	Questions   []Question `json:"questions" gorm:"foreignKey:CorrespondingQuizId"`
	OwnerId     uint       `json:"-"`
	Owner       Account    `json:"owner" gorm:"foreignKey:OwnerId;references:Id"`
}

type QuizDto struct {
	Id          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
}

type Rating struct {
	Id                     uint    `json:"id" gorm:"primaryKey"`
	Value                  uint8   `json:"value"`
	CorrespondingProductId uint    `json:"-"`
	CorrespondingProduct   Product `json:"correspondingProduct" gorm:"foreignKey:CorrespondingProductId;references:Id"`
	OwnerId                uint    `json:"-"`
	Owner                  Account `json:"owner" gorm:"foreignKey:OwnerId;references:Id"`
}

type RatingDto struct {
	Id    uint   `json:"id"`
	Value uint8  `json:"value"`
	Owner string `json:"owner"`
}

type Comment struct {
	Id                     uint    `json:"id" gorm:"primaryKey"`
	Text                   string  `json:"text" gorm:"size:255"`
	CorrespondingProductId uint    `json:"-"`
	CorrespondingProduct   Product `json:"correspondingProduct" gorm:"foreignKey:CorrespondingProductId;references:Id"`
	OwnerId                uint    `json:"-"`
	Owner                  Account `json:"owner" gorm:"foreignKey:OwnerId;references:Id"`
}

type CommentDto struct {
	Id    uint   `json:"id"`
	Text  string `json:"text"`
	Owner string `json:"owner"`
}

type Product struct {
	Id       uint      `json:"id" gorm:"primaryKey"`
	ItemId   uint      `json:"itemId"`
	Item     Quiz      `json:"item" gorm:"foreignKey:ItemId;references:Id"`
	Price    float32   `json:"price"`
	Comments []Comment `json:"comments" gorm:"foreignKey:CorrespondingProductId"`
	Ratings  []Rating  `json:"ratings" gorm:"foreignKey:CorrespondingProductId"`
}

type ProductDto struct {
	Id       uint         `json:"id"`
	Item     QuizDto      `json:"item"`
	Price    float32      `json:"price"`
	Comments []CommentDto `json:"comments"`
	Ratings  []RatingDto  `json:"ratings"`
}

type Stat struct {
	Id         uint    `json:"id" gorm:"primaryKey"`
	PlayerId   uint    `json:"-"`
	Player     Account `json:"player" gorm:"foreignKey:PlayerId;references:Id"`
	GameId     uint    `json:"-"`
	ActiveGame Game    `json:"activeGame" gorm:"foreignKey:GameId;references:Id"`
	Score      uint    `json:"score"`
}

type StatDto struct {
	Id         uint   `json:"id"`
	PlayerName string `json:"playerName"`
	Score      uint   `json:"score"`
}

type Game struct {
	Id              uint    `json:"id" gorm:"primaryKey"`
	IsActive        bool    `json:"isActive"`
	IsInProgress    bool    `json:"isInProgress"`
	Code            string  `json:"code" gorm:"size:6"`
	CreatorId       uint    `json:"-"`
	Creator         Account `json:"creator" gorm:"foreignKey:CreatorId;references:Id"`
	Stats           []Stat  `json:"stats" gorm:"foreignKey:GameId"`
	QuizId          uint    `json:"-"`
	ActiveQuiz      Quiz    `json:"activeQuiz" gorm:"foreignKey:QuizId;references:Id"`
	CurrentQuestion uint    `json:"currentQuestion"`
}

type CreateAccountRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Role        string `json:"role"`
}

type ModifyAccountRequest struct {
	Username    string `json:"username"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Description string `json:"description"`
}

type DepositRequest struct {
	Amount float32
}

type CreateQuizRequest struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type ModifyQuizRequest struct {
	Id          uint       `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Questions   []Question `json:"questions"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateCommentRequest struct {
	ProductId uint   `json:"productId"`
	Text      string `json:"text"`
}

type ModifyCommentRequest struct {
	CommentId uint   `json:"commentId"`
	Text      string `json:"text"`
}

type CreateRatingRequest struct {
	ProductId uint  `json:"productId"`
	Rating    uint8 `json:"rating"`
}

type ModifyRatingRequest struct {
	RatingId uint  `json:"ratingId"`
	Rating   uint8 `json:"rating"`
}

type BuyQuizRequest struct {
	ProductId uint `json:"productId"`
}

type CreateGameRequest struct {
	QuizId uint `json:"quizId"`
}

type SellQuizRequest struct {
	QuizId uint    `json:"quizId"`
	Price  float32 `json:"price"`
}

type Event struct {
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type SendAnswerEvent struct {
	AnswerId   uint   `json:"answerId"`
	QuestionId string `json:"questionId"`
}

type NextRoundEvent struct {
	Stats    []StatDto   `json:"stats"`
	Question QuestionDto `json:"question"`
}
