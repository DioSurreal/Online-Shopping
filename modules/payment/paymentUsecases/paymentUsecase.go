package paymentUsecases

import "github.com/DioSurreal/Online-Shopping/modules/payment/paymentRepositories"

type(
	PaymentUsecasesService interface{}

    paymentUsecase struct {
		paymentRepository paymentRepositories.PaymentRepositoriesService
	}
)

func NewPaymentUsecase (paymentRepository paymentRepositories.PaymentRepositoriesService) PaymentUsecasesService{
	return &paymentUsecase{
		paymentRepository: paymentRepository,}
}
