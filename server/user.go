package server

import (
	"log"

	userPb "github.com/DioSurreal/Online-Shopping/modules/user/userPb"
	"github.com/DioSurreal/Online-Shopping/modules/user/userHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/user/userRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/user/userUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/grpccon"
)

func (s *server) userService() {
	userRepository := userRepositories.NewUserRepository(s.db)
	userUsecase := userUsecases.NewUserUsecase(userRepository)
	userHttpHandler := userHandlers.NewUserHttpHandler(s.cfg,userUsecase)
    userGrpcHandler := userHandlers.NewUserGrpcHandler(userUsecase)
	userQueueHandler := userHandlers.NewUserQueueHandler(s.cfg,userUsecase)

	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.UserUrl)

		userPb.RegisterUserGrpcServiceServer(grpcServer, userGrpcHandler)

		log.Printf("User gRPC server listening on %s", s.cfg.Grpc.UserUrl)
		grpcServer.Serve(lis)
	}()
	_  = userHttpHandler
	_ = userGrpcHandler
	_ = userQueueHandler

	user := s.app.Group("/user_v1")

	//Health Check
	user.GET("",s.healthCheckService)
}