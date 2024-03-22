package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/auth/authHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authUsecases"
)


func (s *server) authService() {
	authRepository := authRepositories.NewAuthRepository(s.db)
	authUsecase := authUsecases.NewAuthUsecase(authRepository)
	authHttpHandler := authHandlers.NewAuthHttpHandler(s.cfg,authUsecase)
    authGrpcHandler := authHandlers.NewAuthGrpcHandler(authUsecase)

	_  = authHttpHandler
	_ = authGrpcHandler

	auth := s.app.Group("/auth_v1")

	//Health Check
	auth.GET("",s.healthCheckService)
}