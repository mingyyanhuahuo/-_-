package service

import (
	"lostfound/dao"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/hashpassword"
	"lostfound/pkg/jwtutil"
	"time"
)

func Register(body *model.RegisterBody) (int64, error) {
	if body.Username == "" || body.Password == "" || body.Nickname == "" || body.StudentNo == "" || body.Phone == "" || body.Email == "" {
		return 0, errcode.ErrBadRequest
	}

	if _, err := dao.OnlyUsername(body.Username); err != nil {
		return 0, errcode.ErrUserExist
	}
	if _, err := dao.OnlyStudentNo(body.StudentNo); err != nil {
		return 0, errcode.ErrStudentIDRegistered
	}
	if _, err := dao.OnlyPhone(body.Phone); err != nil {
		return 0, errcode.ErrPhoneRegistered
	}
	if _, err := dao.OnlyEmail(body.Email); err != nil {
		return 0, errcode.ErrEmailRegistered
	}

	HashPassword, err := hashpassword.Hash(body.Password)
	if err != nil {
		return 0, err
	}

	user := model.User{
		UserName:  body.Username,
		PassHash:  HashPassword,
		NickName:  body.Nickname,
		StudentNo: body.StudentNo,
		Phone:     body.Phone,
		Email:     body.Email,
		Role:      "student",
		CreatedAt: time.Now(),
	}

	if err := dao.CreateUser(&user); err != nil {
		return 0, err
	}

	return int64(user.ID), nil
}

func Login(username, password string) (model.LoginResponse, error) {
	var LoginReponse model.LoginResponse
	user, err := dao.OnlyUsername(username)
	if err != nil {
		return LoginReponse, err
	}
	if user == nil {
		return LoginReponse, errcode.ErrUserPwdWrong
	}

	if err := hashpassword.CheckHash(user.PassHash, password); err != nil {
		return LoginReponse, errcode.ErrUserPwdWrong
	} else {
		access, err := jwtutil.GenerateAccessToken(user.ID, user.Role)
		if err != nil {
			return LoginReponse, err
		}
		refresh, err := jwtutil.GenerateRefreshToken(user.ID, user.Role)
		if err != nil {
			return LoginReponse, err
		}

		LoginReponse = model.LoginResponse{
			AccessToken:  access,
			RefreshToken: refresh,
			ExpiresIn:    7200,
			UserInfo: model.UserInfo{
				UserID:     user.ID,
				UserName:   user.UserName,
				NickName:   user.NickName,
				Avatar:     user.Avatar,
				StudentNo:  user.StudentNo,
				Phone:      user.Phone,
				Email:      user.Email,
				Role:       user.Role,
				CreateTime: user.CreatedAt,
			},
		}
		return LoginReponse, nil
	}

}
