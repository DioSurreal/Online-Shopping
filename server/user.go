package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/user/userHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/user/userRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
)

func (s *server) userService() {
	userRepository := userRepositories.NewUserRepository(s.db)
	userUsecase := userUsecases.NewUserUsecase(userRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(s.cfg,userUsecase)
    userGrpcHandler := userHandlers.NewUserGrpcHandler(userUsecase)
	userQueueHandler := userHandlers.NewUserQueueHandler(s.cfg,userUsecase)

	_  = userHttpHandler
	_ = userGrpcHandler
	_ = userQueueHandler

	user := s.app.Group("/user_v1")

	//Health Check
	user.GET("",s.healthCheckService)
}