package paymentHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/payment/paymentUsecases"
)

type(
	PaymentHttpHandlersService interface{}

    paymentHttpHandler struct {
		cfg *config.Config
		paymentUsecase paymentUsecases.PaymentUsecasesService
	}
)

func NewPaymentHttpHandler (cfg *config.Config,paymentUsecase paymentUsecases.PaymentUsecasesService) PaymentHttpHandlersService{
	return &paymentHttpHandler{
		cfg: cfg,
		paymentUsecase: paymentUsecase,}
}