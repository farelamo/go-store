package controller

import (
	"errors"
	"log"
	"store/config"
	"store/constant"
	"store/helper"
	"store/model"
	"store/repository"
	authSchema "store/schema/auth"
	userSchema "store/schema/user"
	"store/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	RefreshToken(c *gin.Context)
	RevokeToken(c *gin.Context)
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type authController struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
}

func NewAuthControlller(authRepo repository.AuthRepository, userRepo repository.UserRepository) AuthController {
	return &authController{
		authRepo: authRepo,
		userRepo: userRepo,
	}
}

func (a *authController) Login(c *gin.Context) {
	var req authSchema.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, errMsg, nil, false)
		return
	}

	login, err := a.userRepo.GetByUsername(req.Username)
	if err != nil {
		utils.Response(c, 404, nil, nil, nil, "User not found", nil, false)
		return
	}

	if !utils.CheckHash(req.Password, login.Password) {
		utils.Response(c, 401, nil, nil, nil, "Invalid credentials", nil, false)
		return
	}

	var (
		tokenAccess  string
		tokenRefresh string
		now          = time.Now()
	)

	tokenRedis, err := a.authRepo.GetByKey(login.Username)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			tokenString, err := helper.GenerateToken(login.Username, now)
			if err != nil {
				utils.Response(c, 500, nil, nil, nil, "Failed to generate token", nil, false)
				return
			}

			accessUser := model.AuthRedis{Key: login.Username, Value: tokenString}
			if err := a.authRepo.SaveRedis(accessUser); err != nil {
				utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
				return
			}
			tokenAccess = tokenString

			tokenRefresh, err = helper.GenerateRefreshToken(login.Username, now)
			if err != nil {
				utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
				return
			}
		}
	} else {
		tokenAccess = tokenRedis
	}

	auth := authSchema.AuthRes{
		Username:              login.Username,
		AccessToken:           &tokenAccess,
		RefreshToken:          &tokenRefresh,
		Fullname:              login.FullName,
		Iat:                   now,
		AccessTokenExpiredAt:  now.Add(time.Duration(config.AccessTokenExpiration) * time.Minute),
		RefreshTokenExpiredAt: now.Add(time.Duration(config.RefreshTokenExpiration) * time.Hour),
	}

	utils.Response(c, 200, nil, nil, nil, "Login successful", auth, false)
}

func (a *authController) Register(c *gin.Context) {
	var req userSchema.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, errMsg, nil, false)
		return
	}

	_, err := a.userRepo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			hashedPassword, err := utils.HashString(req.Password)
			if err != nil {
				utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
				return
			}

			newUser := model.User{
				Username: req.Username,
				Password: hashedPassword,
				FullName: req.FullName,
			}

			if err := a.userRepo.Create(newUser); err != nil {
				utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
				return
			}

			utils.Response(c, 201, nil, nil, nil, "User registered successfully", nil, false)
			return
		}

		utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
		return
	}

	utils.Response(c, 400, nil, nil, nil, "User already exists", nil, false)
}

func (a *authController) RefreshToken(c *gin.Context) {
	var (
		req authSchema.RefreshTokenReq
		now = time.Now()
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, errMsg, nil, false)
		return
	}

	username, newResfreshToken, err := helper.GenerateRefreshTokenFromOld(req.RefreshToken, now)
	if err != nil {
		log.Println("Error generating refresh token:", err.Error())
		if errors.Is(err, constant.ErrJwtSigned) {
			utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
			return
		}
		utils.Response(c, 400, nil, nil, nil, err.Error(), nil, false)
		return
	}

	token, err := helper.GenerateToken(username, now)
	if err != nil {
		log.Println("Error generating access token:", err.Error())
		utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
		return
	}

	auth := model.AuthRedis{Key: username, Value: token}
	if err := a.authRepo.SaveRedis(auth); err != nil {
		log.Println("Error saving token to Redis:", err.Error())
		utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
		return
	}

	var res authSchema.AuthRes
	res.AccessToken = &token
	res.RefreshToken = &newResfreshToken
	res.Username = username
	res.Iat = now
	res.AccessTokenExpiredAt = now.Add(time.Duration(config.AccessTokenExpiration) * time.Minute)
	res.RefreshTokenExpiredAt = now.Add(time.Duration(config.RefreshTokenExpiration) * time.Hour)

	utils.Response(c, 200, nil, nil, nil, "Token refreshed successfully", res, false)
}

func (a *authController) RevokeToken(c *gin.Context) {
	var req authSchema.AuthRes
	if err := c.ShouldBindJSON(&req); err != nil {
		errMsg := utils.FormatValidationError(err)
		utils.Response(c, 400, nil, nil, nil, errMsg, nil, false)
		return
	}

	if err := a.authRepo.Delete(req.Username); err != nil {
		if errors.Is(err, redis.Nil) {
			utils.Response(c, 404, nil, nil, nil, "Token not found", nil, false)
			return
		}
		log.Println("Error deleting token from Redis:", err.Error())
		utils.Response(c, 500, nil, nil, nil, constant.ErrInternalServerError, nil, false)
		return
	}
	utils.Response(c, 200, nil, nil, nil, "Token revoked successfully", nil, false)
}
