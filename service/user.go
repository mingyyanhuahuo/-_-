package service

import (
	"lostfound/dao"
	"lostfound/model"
	"lostfound/pkg/errcode"
	"lostfound/pkg/hashpassword"
	"lostfound/pkg/jwtutil"
	"lostfound/pkg/redisdb"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Register(body *model.RegisterBody) (int64, error) {
	if body.Username == "" || body.Password == "" || body.Nickname == "" || body.StudentNo == "" || body.Phone == "" || body.Email == "" {
		return 0, errcode.ErrBadRequest
	}

	if _, err := dao.OnlyUsername(body.Username); err != nil {
		return 0, err
	}
	if _, err := dao.OnlyStudentNo(body.StudentNo); err != nil {
		return 0, err
	}
	if _, err := dao.OnlyPhone(body.Phone); err != nil {
		return 0, err
	}
	if _, err := dao.OnlyEmail(body.Email); err != nil {
		return 0, err
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

func Login(username string, password string) (model.LoginResponse, error) {
	var LoginReponse model.LoginResponse
	user, err := dao.OnlyUsername(username)
	if err != nil && err != errcode.ErrUserExist {
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

func RefreshToken(refreshToken string) (model.RefreshResponse, error) {
	var Response model.RefreshResponse
	Cliam, err := jwtutil.ParseToken(refreshToken, jwtutil.TokenTypeRefresh)
	if err != nil {
		return Response, err
	}

	access, err := jwtutil.GenerateAccessToken(Cliam.UserID, Cliam.Role)
	if err != nil {
		return Response, err
	}
	refresh, err := jwtutil.GenerateRefreshToken(Cliam.UserID, Cliam.Role)
	if err != nil {
		return Response, err
	}

	if err := redisdb.RemveToken(refreshToken, time.Until(Cliam.ExpiresAt.Time)); err != nil {
		return Response, err
	}

	Response = model.RefreshResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    7200,
	}
	return Response, nil
}

func Logout(refreshToken string) error {
	Cliam, err := jwtutil.ParseToken(refreshToken, jwtutil.TokenTypeRefresh)
	if err != nil {
		return err
	}

	dration := time.Until(Cliam.ExpiresAt.Time)
	if dration <= 0 {
		return nil
	}
	if err := redisdb.RemveToken(refreshToken, dration); err != nil {
		return err
	}
	return nil
}

func GetMe(userid uint) (model.UserInfo, error) {
	var userinfo model.UserInfo
	user, err := dao.IDtoUser(userid)
	if err != nil {
		return userinfo, err
	}

	userinfo = model.UserInfo{
		UserID:     user.ID,
		UserName:   user.UserName,
		NickName:   user.NickName,
		Avatar:     user.Avatar,
		StudentNo:  user.StudentNo,
		Phone:      user.Phone,
		Email:      user.Email,
		Role:       user.Role,
		CreateTime: user.CreatedAt,
	}

	return userinfo, nil
}

func UpdatePassaard(oldPwd string, newPwd string, ID uint, refreshToken string) error {
	user, err := dao.IDtoUser(ID)
	if err != nil {
		return err
	}

	newHashedPassword, err := hashpassword.Hash(newPwd)
	if err != nil {
		return err
	}

	switch err := hashpassword.CheckHash(user.PassHash, oldPwd); err {
	case bcrypt.ErrMismatchedHashAndPassword:
		return errcode.ErrOldPwdWrong
	case nil:
		if err := dao.UpdatePassword(user, newHashedPassword); err != nil {
			return err
		}
		Cliam, err := jwtutil.ParseToken(refreshToken, jwtutil.TokenTypeRefresh)
		if err != nil {
			return err
		}

		dration := time.Until(Cliam.ExpiresAt.Time)
		if dration <= 0 {
			return nil
		}
		if err := redisdb.RemveToken(refreshToken, dration); err != nil {
			return err
		}
		return nil
	default:
		return err
	}
}
