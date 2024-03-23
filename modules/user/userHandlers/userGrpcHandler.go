package userHandlers

import (
	"context"

	userPb "github.com/DioSurreal/Online-Shopping/modules/user/userPb"
	"github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
)
type(
	userGrpcHandler struct {
		userUsecase userUsecases.UserUsecasesService
		userPb.UnimplementedUserGrpcServiceServer
	}
)

func NewUserGrpcHandler(userUsecase userUsecases.UserUsecasesService) *userGrpcHandler {
	return &userGrpcHandler{userUsecase: userUsecase}
}

func (g *userGrpcHandler) CredentialSearch(ctx context.Context, req *userPb.CredentialSearchReq) (*userPb.UserProfile, error) {
	return nil,nil
}

func (g *userGrpcHandler) FindOneUserProfileToRefresh(ctx context.Context, req *userPb.FindOneUserProfileToRefreshReq) (*userPb.UserProfile, error) {
	return nil,nil
}

func (g *userGrpcHandler) GetUserSavingAccount(ctx context.Context, req *userPb.GetUserSavingAccountReq) (*userPb.GetUserSavingAccountRes, error) {
	return nil, nil
}