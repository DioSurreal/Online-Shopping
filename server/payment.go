package server

import (
	"github.com/DioSurreal/Online-Shopping/modules/payment/paymentHandlers"
	"github.com/DioSurreal/Online-Shopping/modules/payment/paymentRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/payment/paymentUsecases"
)

func (s *server) paymentService() {
	paymentRepository := paymentRepositories.NewPaymentRepository(s.db)
	paymentUsecase := paymentUsecases.NewPaymentUsecase(paymentRepository)
	paymentHttpHandler := paymentHandlers.NewPaymentHttpHandler(s.cfg,paymentUsecase)

	_  = paymentHttpHandler
	

	payment := s.app.Group("/payment_v1")

	//Health Check
	payment.GET("",s.healthCheckService)

	payment.POST("/payment/buy", paymentHttpHandler.BuyItem, s.middleware.JwtAuthorization)
	payment.POST("/payment/sell", paymentHttpHandler.SellItem, s.middleware.JwtAuthorization)
}