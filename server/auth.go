package server

import (
	"log"

authPb "github.com/DioSurreal/Online-Shopping/modules/auth/authPb"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/auth/authUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/grpccon"
)


func (s *server) authService() {
	authRepository := authRepositories.NewAuthRepository(s.db)
	authUsecase := authUsecase.NewAuthUsecase(authRepository)
	authHttpHandler := authHandlers.NewAuthHttpHandler(s.cfg,authUsecase)
    authGrpcHandler := authHandlers.NewAuthGrpcHandler(authUsecase)


	//Grpc
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.AuthUrl)

		authPb.RegisterAuthGrpcServiceServer(grpcServer, authGrpcHandler)

		log.Printf("Auth gRPC server listening on %s", s.cfg.Grpc.AuthUrl)
		grpcServer.Serve(lis)
	}()
	_ = authGrpcHandler

	auth := s.app.Group("/auth_v1")

	//Health Check
	auth.GET("",s.healthCheckService)

	auth.GET("/test/:user_id", s.healthCheckService)
	auth.POST("/auth/login", authHttpHandler.Login)
	auth.POST("/auth/refresh-token", authHttpHandler.RefreshToken)
	auth.POST("/auth/logout", authHttpHandler.Login)
}