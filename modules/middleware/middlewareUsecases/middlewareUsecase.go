package middlewareUsecases

import (
	"errors"
	"log"

	"github.com/DioSurreal/Online-Shopping/config"
	middlewareRepository "github.com/DioSurreal/Online-Shopping/modules/middleware/middlewareRepositories"
	"github.com/DioSurreal/Online-Shopping/pkg/jwtauth"
	"github.com/DioSurreal/Online-Shopping/pkg/rbac"
	"github.com/labstack/echo/v4"
)

type (
	MiddlewareUsecasesService interface {
		JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error)
		UserIdParamValidation(c echo.Context) (echo.Context, error)
		RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error)
	}

	middlewareUsecase struct {
		middlewareRepository middlewareRepository.MiddlewareRepositoriesService
	}
)

func NewMiddlewareUsecase(middlewareRepository middlewareRepository.MiddlewareRepositoriesService) MiddlewareUsecasesService {
	return &middlewareUsecase{middlewareRepository}
}

func (u *middlewareUsecase) JwtAuthorization(c echo.Context, cfg *config.Config, accessToken string) (echo.Context, error) {
	ctx := c.Request().Context()

	claims, err := jwtauth.ParseToken(cfg.Jwt.AccessSecretKey, accessToken)
	if err != nil {
		return nil, err
	}

	if err := u.middlewareRepository.AccessTokenSearch(ctx, cfg.Grpc.AuthUrl, accessToken); err != nil {
		return nil, err
	}

	c.Set("User_id", claims.UserId)
	c.Set("role_code", claims.RoleCode)

	return c, nil
}

func (u *middlewareUsecase) RbacAuthorization(c echo.Context, cfg *config.Config, expected []int) (echo.Context, error) {
	ctx := c.Request().Context()

	UserRoleCode := c.Get("role_code").(int)

	rolesCount, err := u.middlewareRepository.RolesCount(ctx, cfg.Grpc.AuthUrl)
	if err != nil {
		return nil, err
	}

	UserRoleBinary := rbac.IntToBinary(UserRoleCode, int(rolesCount))

	for i := 0; i < int(rolesCount); i++ {
		if UserRoleBinary[i]&expected[i] == 1 {
			return c, nil
		}
	}

	return nil, errors.New("error: permission denied")
}

func (u *middlewareUsecase) UserIdParamValidation(c echo.Context) (echo.Context, error) {
	UserIdReq := c.Param("User_id")
	UserIdToken := c.Get("User_id").(string)

	if UserIdToken == "" {
		log.Printf("Error: User_id not found")
		return nil, errors.New("error: User_id is required")
	}

	if UserIdToken != UserIdReq {
		log.Printf("Error: User_id not match, User_id_req: %s, User_id_token: %s", UserIdReq, UserIdToken)
		return nil, errors.New("error: User_id not match")
	}

	return c, nil
}
