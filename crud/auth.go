package crud

import (
	"fmt"

	model "github.com/Prokuma/ProkumaLabAccount-Backend/models"
)

func CreateUser(user *model.User) error {
	err := DB.Create(user).Error

	if err != nil {
		fmt.Println("User could not create: ", err)
		return err
	}

	return nil
}

func GetUser(userId string) (model.User, error) {
	var user model.User
	err := DB.Where(&model.User{UserId: userId}).First(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func GetUserFromEmail(email string) (model.User, error) {
	var user model.User
	err := DB.Where(&model.User{Email: email}).First(&user).Error

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func UpdateUser(user *model.User) error {
	err := DB.Save(user).Error

	if err != nil {
		fmt.Println("User could not update: ", err)
		return err
	}

	return nil
}
