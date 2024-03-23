package authHandlers

import (
	"context"

	authPb "github.com/DioSurreal/Online-Shopping/modules/auth/authPb"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authUsecases"
)

type(
	authGrpcHandler struct {
		authPb.UnimplementedAuthGrpcServiceServer
		authUsecase authUsecases.AuthUsecasesService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecases.AuthUsecasesService) *authGrpcHandler {
	return &authGrpcHandler{
		authUsecase: authUsecase,
	}
}

func (g *authGrpcHandler) AccessTokenSearch(ctx context.Context, req *authPb.AccessTokenSearchReq) (*authPb.AccessTokenSearchRes, error) {
	return nil,nil}

func (g *authGrpcHandler) RolesCount(ctx context.Context, req *authPb.RolesCountReq) (*authPb.RolesCountRes, error) {
	return nil,nil
}