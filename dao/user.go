package dao

import "lostfound/model"

var user model.User

func OnlyUsername(username string) (*model.User, error) {
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func OnlyStudentNo(studentNo string) (*model.User, error) {
	if err := db.Where("studentNo = ?", studentNo).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func OnlyPhone(phone string) (*model.User, error) {
	if err := db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func OnlyEmail(email string) (*model.User, error) {
	if err := db.Where("email = ?", email).First(&user).Error; err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *model.User) error {
	if err := db.Create(user).Error; err != nil {
		return err
	}
	return nil
}
