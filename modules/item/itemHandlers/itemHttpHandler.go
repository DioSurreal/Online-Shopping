package itemHandlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/modules/item"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/request"
	"github.com/DioSurreal/Online-Shopping/pkg/response"
	"github.com/labstack/echo/v4"
)
type(
	ItemHttpHandlersService interface{
		CreateItem(c echo.Context) error
		FindOneItem(c echo.Context) error
		FindManyItems(c echo.Context) error
		EditItem(c echo.Context) error
		EnableOrDisableItem(c echo.Context) error
	}

	itemHttpHandler struct {
		cfg  *config.Config
		itemUsecase itemUsecases.ItemUsecasesService
	}
)

func NewItemHttpHandler(cfg *config.Config,itemUsecase itemUsecases.ItemUsecasesService) ItemHttpHandlersService {
	return &itemHttpHandler{
		cfg: cfg,
		itemUsecase: itemUsecase}
}

func (h *itemHttpHandler) CreateItem(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.CreateItemReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.CreateItem(ctx, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttpHandler) FindOneItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	res, err := h.itemUsecase.FindOneItem(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusCreated, res)
}

func (h *itemHttpHandler) FindManyItems(c echo.Context) error {
	ctx := context.Background()

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemSearchReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.FindManyItems(ctx, h.cfg.Paginate.ItemNextPageBasedUrl, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EditItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	wrapper := request.ContextWrapper(c)

	req := new(item.ItemUpdateReq)

	if err := wrapper.Bind(req); err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	res, err := h.itemUsecase.EditItem(ctx, itemId, req)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, res)
}

func (h *itemHttpHandler) EnableOrDisableItem(c echo.Context) error {
	ctx := context.Background()

	itemId := strings.TrimPrefix(c.Param("item_id"), "item:")

	res, err := h.itemUsecase.EnableOrDisableItem(ctx, itemId)
	if err != nil {
		return response.ErrResponse(c, http.StatusBadRequest, err.Error())
	}

	return response.SuccessResponse(c, http.StatusOK, map[string]any{
		"message": fmt.Sprintf("item_id: %s is successfully is activated to: %v", itemId, res),
	})
}