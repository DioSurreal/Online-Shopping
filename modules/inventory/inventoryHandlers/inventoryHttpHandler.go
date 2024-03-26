package inventoryHandlers

import (
	"context"
	"net/http"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/inventory"
	"github.com/DioSurreal/Online-Shopping/modules/inventory/inventoryUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/request"
	"github.com/DioSurreal/Online-Shopping/pkg/response"
	"github.com/labstack/echo/v4"
)

type (
	InventoryHttpHandlersService interface{
		FindUserItems(c echo.Context) error
	}

	inventoryHttpHandler struct {
		cfg              *config.Config
		inventoryUsecase inventoryUsecases.InventoryUsecasesService
	}
)

func NewInventoryHttpHandler(cfg *config.Config, inventoryUsecase inventoryUsecases.InventoryUsecasesService) InventoryHttpHandlersService {
	return &inventoryHttpHandler{
		cfg:              cfg,
		inventoryUsecase: inventoryUsecase}
}

func (h *inventoryHttpHandler) FindUserItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(inventory.InventorySearchReq)
	userId := c.Param("user_id")

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.inventoryUsecase.FindUserItems(ctx, h.cfg, userId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}
