package userUsecases

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/DioSurreal/Online-Shopping/modules/user"
	userPb "github.com/DioSurreal/Online-Shopping/modules/user/userPb"
	"github.com/DioSurreal/Online-Shopping/modules/user/userRepositories"
	"github.com/DioSurreal/Online-Shopping/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

type (
	UserUsecasesService interface{
		CreateUser(pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error)
		FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfile, error)
		AddUserMoney(pctx context.Context, req *user.CreateUserTransactionReq) (*user.UserSavingAccount, error)
		GetUserSavingAccount(pctx context.Context, userId string) (*user.UserSavingAccount, error)
		FindOneUserCredential(pctx context.Context, password, email string) (*userPb.UserProfile, error) 
		FindOneUserProfileToRefresh(pctx context.Context, userId string) (*userPb.UserProfile, error)
	}

	userUsecase struct {
		userRepository userRepositories.UserRepositoriesService
	}
)

func NewUserUsecase (userRepository userRepositories.UserRepositoriesService) UserUsecasesService {
	return &userUsecase{userRepository: userRepository}
} 

func (u *userUsecase) CreateUser(pctx context.Context, req *user.CreateUserReq) (*user.UserProfile, error) {
	if !u.userRepository.IsUniqueUser(pctx, req.Email, req.Username) {
		return nil, errors.New("error: email or username already exist")
	}

	// Hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("error: failed to hash password")
	}

	// Insert one user
	userId, err := u.userRepository.InsertOneUser(pctx, &user.User{
		Email:     req.Email,
		Password:  string(hashedPassword),
		Username:  req.Username,
		CreatedAt: utils.LocalTime(),
		UpdatedAt: utils.LocalTime(),
		UserRoles: []user.UserRole{
			{
				RoleTitle: "User",
				RoleCode:  0,
			},
		},
	})

	return u.FindOneUserProfile(pctx, userId.Hex())
}

func (u *userUsecase) FindOneUserProfile(pctx context.Context, userId string) (*user.UserProfile, error) {
	result, err := u.userRepository.FindOneUserProfile(pctx, userId)
	if err != nil {
		return nil, err
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &user.UserProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		CreatedAt: result.CreatedAt.In(loc),
		UpdatedAt: result.UpdatedAt.In(loc),
	}, nil
}

func (u *userUsecase) AddUserMoney(pctx context.Context, req *user.CreateUserTransactionReq) (*user.UserSavingAccount, error) {
	// Insert one user transaction
	if _, err := u.userRepository.InsertOneUserTranscation(pctx, &user.UserTransaction{
		UserId:  req.UserId,
		Amount:    req.Amount,
		CreatedAt: utils.LocalTime(),
	}); err != nil {
		return nil, err
	}

	// Get user saving account
	return u.userRepository.GetUserSavingAccount(pctx, req.UserId)
}

func (u *userUsecase) GetUserSavingAccount(pctx context.Context, userId string) (*user.UserSavingAccount, error) {
	return u.userRepository.GetUserSavingAccount(pctx, userId)
}

func (u *userUsecase) FindOneUserCredential(pctx context.Context, password, email string) (*userPb.UserProfile, error) {
	result, err := u.userRepository.FindOneUserCredential(pctx, email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err != nil {
		log.Printf("Error: FindOneUserCredential: %s", err.Error())
		return nil, errors.New("error: password is invalid")
	}

	roleCode := 0
	for _, v := range result.UserRoles {
		roleCode += v.RoleCode
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &userPb.UserProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}

func (u *userUsecase) FindOneUserProfileToRefresh(pctx context.Context, userId string) (*userPb.UserProfile, error) {
	result, err := u.userRepository.FindOneUserProfileToRefresh(pctx, userId)
	if err != nil {
		return nil, err
	}

	roleCode := 0
	for _, v := range result.UserRoles {
		roleCode += v.RoleCode
	}

	loc, _ := time.LoadLocation("Asia/Bangkok")

	return &userPb.UserProfile{
		Id:        result.Id.Hex(),
		Email:     result.Email,
		Username:  result.Username,
		RoleCode:  int32(roleCode),
		CreatedAt: result.CreatedAt.In(loc).String(),
		UpdatedAt: result.UpdatedAt.In(loc).String(),
	}, nil
}
