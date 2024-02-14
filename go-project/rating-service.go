package main

import (
	"errors"
)

func DeleteRating(id uint, userId uint, role string) error {
	rating, err := Db.GetRatingById(id)
	if err != nil {
		return err
	}

	if rating.Owner.Id != userId && role != Admin {
		return errors.New("you dont have permission to modify this account")
	}

	if err := Db.DeleteRatingById(id); err != nil {
		return err
	}

	return nil
}

func GetRatingsForProduct(productId uint) ([]Rating, error) {
	ratings, err := Db.GetRatingsByProductId(productId)
	if err != nil {
		return nil, err
	}

	return ratings, nil
}

func CreateRating(body *CreateRatingRequest, userId uint) error {
	if isRatingInvalid(body.Rating) {
		return errors.New("rating is invalid")
	}
	username, err := Db.GetUsernameByAccountId(userId)
	if err != nil {
		return err
	}

	acc, err := GetAccount(username)
	if err != nil {
		return err
	}

	product, err := Db.GetProductById(body.ProductId)
	if err != nil {
		return err
	}

	rating := Rating{
		Value:                  body.Rating,
		CorrespondingProductId: body.ProductId,
		CorrespondingProduct:   *product,
		OwnerId:                acc.Id,
		Owner:                  *acc,
	}
	Db.PostRating(&rating)

	return nil
}

func isRatingInvalid(rating uint8) bool {
	return rating < 1 || rating > 5
}

func ModifyRating(body *ModifyRatingRequest, userId uint, role string) error {
	if isRatingInvalid(body.Rating) {
		return errors.New("rating is invalid")
	}

	username, err := Db.GetUsernameByAccountId(userId)
	if err != nil {
		return err
	}

	rating, err := Db.GetRatingById(body.RatingId)
	if err != nil {
		return err
	} else if username != rating.Owner.Username && role != Admin {
		return errors.New("you dont have permission to modify this account")
	}

	rating.Value = uint8(body.Rating)

	Db.PostRating(rating)

	return nil
}

func CreateRatingDto(rating *Rating) *RatingDto {
	return &RatingDto{
		Id:    rating.Id,
		Value: rating.Value,
		Owner: rating.Owner.Username,
	}
}
