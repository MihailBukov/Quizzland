package main

import "errors"

func CreateQuiz(body *CreateQuizRequest, userId uint) error {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	quiz := Quiz{
		Name:        body.Name,
		Description: body.Description,
		Questions:   body.Questions,
		OwnerId:     acc.Id,
		Owner:       *acc,
	}

	err = Db.PostQuiz(&quiz)
	if err != nil {
		return err
	}

	return nil
}

func GetQuizzesForSale() ([]ProductDto, error) {
	products, err := Db.GetProducts()
	if err != nil {
		return nil, err
	}

	productsDto := make([]ProductDto, len(products))
	for _, product := range products {
		productsDto = append(productsDto, *CreateProductDto(&product))
	}

	return productsDto, nil
}

func DeleteQuiz(id int) error {
	if err := Db.DeleteQuizById(id); err != nil {
		return err
	}

	return nil
}

func ModifyQuiz(body *ModifyQuizRequest, userId uint, role string) error {
	acc, err := Db.GetAccountById(userId)
	if err != nil {
		return err
	}

	if !IsQuizExistingInAccount(acc, body.Id) && role != Admin {
		return errors.New("you do not have permission to modify this resource")
	}

	quiz := Quiz{
		Id:          body.Id,
		Name:        body.Name,
		Description: body.Description,
		Questions:   body.Questions,
		OwnerId:     acc.Id,
		Owner:       *acc,
	}

	err = Db.PostQuiz(&quiz)
	if err != nil {
		return err
	}

	return nil
}

func createAnswerDto(answer Answer) *AnswerDto {
	return &AnswerDto{
		Id:   answer.Id,
		Text: answer.Text,
	}
}

func CreateQuestionDto(Question Question) *QuestionDto {
	answers := make([]AnswerDto, 4)
	for _, answer := range Question.Answers {
		answers = append(answers, *createAnswerDto(answer))
	}

	return &QuestionDto{
		Id:      Question.Id,
		Text:    Question.Text,
		Time:    Question.Time,
		Answers: answers,
	}
}

func CreateQuizDto(quiz *Quiz) QuizDto {
	return QuizDto{
		Id:          quiz.Id,
		Name:        quiz.Name,
		Description: quiz.Description,
		Owner:       quiz.Owner.Username,
	}
}

func CreateProductDto(product *Product) *ProductDto {
	ratings := make([]RatingDto, len(product.Ratings))
	for _, rating := range product.Ratings {
		ratings = append(ratings, *CreateRatingDto(&rating))
	}

	comments := make([]CommentDto, len(product.Comments))
	for _, comment := range product.Comments {
		comments = append(comments, *CreateCommentDto(&comment))
	}

	return &ProductDto{
		Id:       product.Id,
		Item:     CreateQuizDto(&product.Item),
		Price:    product.Price,
		Comments: comments,
		Ratings:  ratings,
	}
}
