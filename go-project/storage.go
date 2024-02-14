package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	DeleteRatingById(id uint) error
	GetRatingById(id uint) (*Rating, error)
	GetRatingsByProductId(id uint) ([]Rating, error)
	PostRating(rating *Rating) error

	DeleteCommentById(id uint) error
	GetCommentById(id uint) (*Comment, error)
	GetCommentsByProductId(id uint) ([]Comment, error)
	PostComment(comment *Comment) error

	GetAnswersByQuestionId(questionId int) ([]Answer, error)
	CreateAnswerForQuestion(answer *Answer) error
	DeleteAnswerById(id int) error
	GetAnswerById(id uint) (*Answer, error)

	CreateQuestion(question *Question) error
	DeleteQuestionById(id int) error
	PutQuestion(question *Question) error
	GetQuestionsByQuizId(quizId int) ([]Question, error)

	GetUsernameByAccountId(id uint) (string, error)
	PostAccount(account *Account) error
	DeleteAccountByUsername(username string) error
	PutAccount(account *Account) error
	GetAccountById(id uint) (*Account, error)
	GetAccountByUsername(username string) (*Account, error)

	GetProducts() ([]Product, error)
	GetProductById(id uint) (*Product, error)
	PutProduct(product *Product) error
	DeleteProductById(id int) error
	IsQuizForSale(quizId uint) (bool, error)

	GetQuizById(id uint) (*Quiz, error)
	GetQuizzesByOwnerId(id int) ([]Quiz, error)
	PutQuiz(quiz *Quiz) error
	PostQuiz(quiz *Quiz) error
	DeleteQuizById(id int) error

	GetGameById(id uint) (*Game, error)
	SaveGame(game *Game) error
	GetGameByCode(code string) (*Game, error)
}

type MySqlStore struct {
	db *gorm.DB
}

var Db MySqlStore

func (s *MySqlStore) SaveGame(game *Game) error {
	if err := s.db.Save(game).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetGameById(id uint) (*Game, error) {
	var game Game

	if err := s.db.First(&game, id).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *MySqlStore) GetGameByCode(code string) (*Game, error) {
	var game Game

	if err := s.db.Where("code = ?", code).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (s *MySqlStore) GetRatingById(id uint) (*Rating, error) {
	var rating Rating

	if err := s.db.Where("id = ?", id).Find(&rating).Error; err != nil {
		return nil, err
	}

	return &rating, nil
}

func (s *MySqlStore) DeleteRatingById(id uint) error {
	if err := s.db.Where("id = ?", id).Delete(&Rating{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetRatingsByProductId(id uint) ([]Rating, error) {
	var ratings []Rating

	if err := s.db.Where("corresponding_product_id = ?", id).Find(&ratings).Error; err != nil {
		return nil, err
	}

	return ratings, nil
}

func (s *MySqlStore) PostRating(rating *Rating) error {
	if err := s.db.Save(rating).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetCommentById(id uint) (*Comment, error) {
	var comment Comment

	if err := s.db.Where("id = ?", id).Find(&comment).Error; err != nil {
		return nil, err
	}

	return &comment, nil
}

func (s *MySqlStore) GetUsernameByAccountId(id uint) (string, error) {
	var username string

	if err := s.db.Model(Account{}).Select("username").Where("id = ?", id).Scan(&username).Error; err != nil {
		return "", err
	}

	return username, nil
}

func (s *MySqlStore) DeleteCommentById(id uint) error {
	if err := s.db.Where("id = ?", id).Delete(&Comment{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetCommentsByProductId(id uint) ([]Comment, error) {
	var comments []Comment

	if err := s.db.Where("corresponding_product_id = ?", id).Find(&comments).Error; err != nil {
		return nil, err
	}

	return comments, nil
}

func (s *MySqlStore) PostComment(comment *Comment) error {
	if err := s.db.Save(comment).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetAnswersByQuestionId(questionId int) ([]Answer, error) {
	var answers []Answer

	if err := s.db.Where("corresponding_question_id = ?", questionId).Find(&answers).Error; err != nil {
		return nil, err
	}

	return answers, nil
}

func (s *MySqlStore) CreateAnswerForQuestion(answer *Answer) error {
	if err := s.db.Save(answer).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) DeleteAnswerById(id int) error {
	if err := s.db.Where("id = ?", id).Delete(&Answer{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetAnswerById(id uint) (*Answer, error) {
	var answer Answer

	if err := s.db.Where("id = ?", id).Find(&answer).Error; err != nil {
		return nil, err
	}

	return &answer, nil
}

func (s *MySqlStore) CreateQuestion(question *Question) error {
	if err := s.db.Save(question).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) DeleteQuestionById(id int) error {
	if err := s.db.Where("id = ?", id).Delete(&Question{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) PutQuestion(question *Question) error {
	if err := s.db.Save(question).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) GetQuestionsByQuizId(quizId int) ([]Question, error) {
	var questions []Question

	if err := s.db.Where("corresponding_quiz_id = ?", quizId).Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}

func (s *MySqlStore) PostAccount(account *Account) error {
	if err := s.db.Create(account).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) DeleteAccountByUsername(username string) error {
	if err := s.db.Where("username = ?", username).Delete(&Account{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) PutAccount(account *Account) error {
	if err := s.db.Save(account).Error; err != nil {
		return err
	}

	return nil
}


func (s *MySqlStore) GetAccountById(id uint) (*Account, error) {
	var account Account

	if err := s.db.Where("id = ?", id).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *MySqlStore) GetAccountByUsername(username string) (*Account, error) {
	var account Account

	if err := s.db.Where("username = ?", username).First(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (s *MySqlStore) GetProducts() ([]Product, error) {
	var products []Product
	if err := s.db.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (s *MySqlStore) GetProductById(id uint) (*Product, error) {
	var product Product

	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *MySqlStore) PutProduct(product *Product) error {
	if err := s.db.Save(product).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) DeleteProductById(id int) error {
	if err := s.db.Where("id = ?", id).Delete(&Product{}).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) IsQuizForSale(quizId uint) (bool, error) {
	var count int64

    if err := s.db.Model(&Product{}).
        Where("item_id = ?", quizId).
        Count(&count).Error; err != nil {
        return false, err
    }

    return count > 0, nil
}

func (s *MySqlStore) GetQuizById(id uint) (*Quiz, error) {
	var quiz Quiz

	if err := s.db.First(&quiz, id).Error; err != nil {
		return nil, err
	}

	return &quiz, nil
}

func (s *MySqlStore) GetQuizzesByOwnerId(id int) ([]Quiz, error) {
	var quizzes []Quiz

	if err := s.db.Where("owner_id = ?", id).Find(&quizzes).Error; err != nil {
		return nil, err
	}

	return quizzes, nil
}


func (s *MySqlStore) PutQuiz(quiz *Quiz) error {
	if err := s.db.Save(quiz).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) PostQuiz(quiz *Quiz) error {
	if err := s.db.Create(quiz).Error; err != nil {
		return err
	}

	return nil
}

func (s *MySqlStore) DeleteQuizById(id int) error {
	if err := s.db.Where("id = ?", id).Delete(&Quiz{}).Error; err != nil {
		return err
	}

	return nil
}

func NewMySqlStore() error {
	database, err := gorm.Open(mysql.Open("root:parola@tcp(127.0.0.1:3306)/quizzland?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return err
	}

	database.AutoMigrate(&Account{}, &Product{}, &Question{}, &Answer{}, &Quiz{}, &Rating{}, &Comment{}, &Stat{}, &Game{})

	Db = MySqlStore{db: database}

	return nil
}
