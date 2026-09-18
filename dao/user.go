package dao

import (
	"lostfound/model"
	"lostfound/pkg/errcode"
)

func OnlyUsername(username string) (*model.User, error) {
	var user model.User
	if err := db.Where("user_name = ?", username).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, errcode.ErrUserExist
}

func OnlyStudentNo(studentNo string) (*model.User, error) {
	var user model.User
	if err := db.Where("student_no = ?", studentNo).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, errcode.ErrStudentIDRegistered
}

func OnlyPhone(phone string) (*model.User, error) {
	var user model.User
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, errcode.ErrPhoneRegistered
}

func OnlyEmail(email string) (*model.User, error) {
	var user model.User
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, errcode.ErrEmailRegistered
}

func CreateUser(user *model.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
