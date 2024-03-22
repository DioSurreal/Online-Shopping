package paymentHandlers

import (
	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/payment/paymentUsecases"
)

type(
	PaymentQueueHandlersService interface{}

    paymentQueueHandler struct {
		cfg *config.Config
		paymentUsecase paymentUsecases.PaymentUsecasesService
	}
)

func NewPaymentQueueHandler (cfg *config.Config,paymentUsecase paymentUsecases.PaymentUsecasesService) PaymentQueueHandlersService{
	return &paymentQueueHandler{
		cfg: cfg,
		paymentUsecase: paymentUsecase,
	}
}