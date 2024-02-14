package main

import (
	"errors"
	"strings"
)

func DeleteComment(id uint, userId uint, role string) error {
	comment, err := Db.GetCommentById(id)
	if err != nil {
		return err
	}

	if comment.Owner.Id != userId && role != Admin {
		return errors.New("you dont have permission to modify this account")
	}

	if err := Db.DeleteCommentById(id); err != nil {
		return err
	}

	return nil
}

func GetCommentsForProduct(productId uint) ([]Comment, error) {
	comments, err := Db.GetCommentsByProductId(productId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func CreateComment(body *CreateCommentRequest, userId uint) error {
	if isStringNullOrEmptyOrBlank(body.Text) {
		return errors.New("comment is invalid")
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

	comment := Comment{
		Text:                   body.Text,
		CorrespondingProductId: body.ProductId,
		CorrespondingProduct:   *product,
		OwnerId:                acc.Id,
		Owner:                  *acc,
	}
	Db.PostComment(&comment)

	return nil
}

func isStringNullOrEmptyOrBlank(s string) bool {
	trimmed := strings.TrimSpace(s)
	return trimmed == ""
}

func ModifyComment(body *ModifyCommentRequest, userId uint, role string) error {
	if isStringNullOrEmptyOrBlank(body.Text) {
		return errors.New("comment is invalid")
	}

	username, err := Db.GetUsernameByAccountId(userId)
	if err != nil {
		return err
	}

	comment, err := Db.GetCommentById(body.CommentId)
	if err != nil {
		return err
	} else if username != comment.Owner.Username && role != Admin {
		return errors.New("you dont have permission to modify this account")
	}

	comment.Text = body.Text

	Db.PostComment(comment)

	return nil
}

func CreateCommentDto(comment *Comment) *CommentDto {
	return &CommentDto{
		Id:    comment.Id,
		Text:  comment.Text,
		Owner: comment.Owner.Username,
	}
}
