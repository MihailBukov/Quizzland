package main

import (
	"errors"
	"unicode"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var store = sessions.NewCookieStore([]byte("secret-key"))

func Init() {

	store.Options = &sessions.Options{
	 Domain:   "127.0.0.1",
	 Path:     "/",
	 MaxAge:   3600 * 8, // 8 hours
	 HttpOnly: true,
 	}
}

func Login(body *LoginRequest) (uint, string, error) {
	acc, err := Db.GetAccountByUsername(body.Username)
	if err != nil {
		return 0, "", errors.New("invalid username or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(acc.Password), []byte(body.Password))
	if err != nil {
		return 0, "", errors.New("invalid username or password")
	}

	return acc.Id, acc.Role, nil
}

func CreateAccount(request *CreateAccountRequest, role string) error {
	for _, char := range request.Username {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return errors.New("username must be alphanumerical")
		}
	}

	if 3 > len(request.Username) && len(request.Username) < 33 {
		return errors.New("username must be between 4 and 32 characters")
	}

	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial bool
	for _, char := range request.Password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			return errors.New("password cannot contain spaces")
		}
	}

	if 5 > len(request.Password) && len(request.Password) > 20 {
		return errors.New("password must be between 6 and 20 characters")
	}

	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial {
		return errors.New("password must contain lower case upper case number and special character")
	}

	acc, err := Db.GetAccountByUsername(request.Username)
	if acc != nil {
		return errors.New("account with this username already exists")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	var hash []byte
	hash, err = bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var roleToAssign string
	if role == Admin {
		roleToAssign = request.Role
	} else {
		roleToAssign = Ruser
	}

	account := Account{
		Username:    request.Username,
		Password:    string(hash),
		FirstName:   request.FirstName,
		LastName:    request.LastName,
		Email:       request.Email,
		Description: request.Description,
		Balance:     0.0,
		Role:        roleToAssign,
	}

	if err := Db.PostAccount(&account); err != nil {
		return err
	}

	return nil
}

func GetAccount(username string) (*Account, error) {
	acc, err := Db.GetAccountByUsername(username)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

func GetAccountDto(username string) (*AccountDto, error) {
	acc, err := Db.GetAccountByUsername(username)

	if err != nil {
		return nil, err
	}

	return CreateAccountDto(acc), nil
}

func DeleteAccount(username string) error {
	_, err := GetAccount(username)
	if err != nil {
		return err
	}

	Db.DeleteAccountByUsername(username)

	return nil
}

func Deposit(body *DepositRequest, userId uint) error {
	if body.Amount <= 0 {
		return errors.New("amount must be positive number")
	}

	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	acc.Balance += body.Amount
	Db.PutAccount(acc)

	return nil
}

func BuyQuiz(body *BuyQuizRequest, userId uint) error {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	product, err := Db.GetProductById(body.ProductId)
	if err != nil {
		return err
	}

	if product.Item.Owner.Id == acc.Id {
		return errors.New("cannot buy quiz that you own")
	}

	if IsQuizExistingInAccount(acc, product.Item.Id) {
		return errors.New("account already owns the quiz")
	}

	if acc.Balance < product.Price {
		return errors.New("insufficient balance")
	}

	acc.Quizzes = append(acc.Quizzes, product.Item)

	Db.PutAccount(acc)

	return nil
}

func SellQuiz(body *SellQuizRequest, userId uint) error {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	if body.Price <= 0 {
		return errors.New("invalid price")
	}

	quiz, err := Db.GetQuizById(body.QuizId)
	if err != nil {
		return err
	}

	isForSale, err := Db.IsQuizForSale(quiz.Id)
	if err != nil {
		return err
	} else if isForSale {
		return errors.New("quiz is already listed for selling")
	}

	if !isQuizOwnedByAccount(acc, quiz) {
		return errors.New("account already owns the quiz")
	}

	product := Product{
		ItemId: quiz.Id,
		Item:   *quiz,
		Price:  body.Price,
	}

	err = Db.PutProduct(&product)
	if err != nil {
		return err
	}

	return nil
}

func EditAccount(username string, userId uint, role string, body ModifyAccountRequest) error {
	acc, err := Db.GetAccountByUsername(username)
	if err != nil {
		return err
	}

	if userId != acc.Id && role != Admin {
		return errors.New("you dont have permission to modify this account")
	}

	newAcc := Account{
		Id:          acc.Id,
		Username:    body.Username,
		Password:    acc.Password,
		FirstName:   body.FirstName,
		LastName:    body.LastName,
		Email:       body.Email,
		Description: body.Description,
		Balance:     acc.Balance,
		Role:        acc.Role,
		Quizzes:     acc.Quizzes,
	}

	err = Db.PutAccount(&newAcc)
	if err != nil {
		return err
	}

	return nil
}

func IsQuizExistingInAccount(acc *Account, quizId uint) bool {
	for _, val := range acc.Quizzes {
		if val.Id == quizId {
			return true
		}
	}

	return false
}

func isQuizOwnedByAccount(acc *Account, quiz *Quiz) bool {
	for _, val := range acc.Quizzes {
		if val.Id == quiz.Id && quiz.Owner.Id == acc.Id {
			return true
		}
	}

	return false
}

func CreateAccountDto(account *Account) *AccountDto {
	quizzes := make([]QuizDto, len(account.Quizzes))
	for _, quiz := range account.Quizzes {
		quizzes = append(quizzes, CreateQuizDto(&quiz))
	}
	
	return &AccountDto{
		Id:          account.Id,
		Username:    account.Username,
		FirstName:   account.FirstName,
		LastName:    account.LastName,
		Email:       account.Email,
		Description: account.Description,
		Balance:     account.Balance,
		Quizzes:     quizzes,
	}
}