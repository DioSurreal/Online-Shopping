package authHandlers

import "github.com/DioSurreal/Online-Shopping/modules/auth/authUsecases"

type(
	authGrpcHandler struct {
		authUsecase authUsecases.AuthUsecasesService
	}
)

func NewAuthGrpcHandler(authUsecase authUsecases.AuthUsecasesService) *authGrpcHandler {
	return &authGrpcHandler{authUsecase}
}