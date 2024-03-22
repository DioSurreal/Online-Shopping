package middlewareRepositories

type (
	MiddlewareRepositoriesService interface{}

	middlewareRepository struct{
		
	}
)

func NewMiddlewareRepository() MiddlewareRepositoriesService {
	return &middlewareRepository{}
}