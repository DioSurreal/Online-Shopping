package inventory

import (
	"github.com/DioSurreal/Online-Shopping/modules/item"
	"github.com/DioSurreal/Online-Shopping/modules/model"
)



type (
	UpdateInventoryReq struct {
		PlayerId string `json:"player_id" validate:"required,max=64"`
		ItemId   string `json:"item_id" validate:"required,max=64"`
	}

	ItemInInventory struct {
		InventoryId string `json:"inventory_id"`
		PlayerId    string `json:"player_id"`
		*item.ItemShowCase
	}

	InventorySearchReq struct {
		model.PaginateReq
	}

	RollbackPlayerInventoryReq struct {
		InventoryId string `json:"inventory_id"`
		PlayerId    string `json:"player_id"`
		ItemId      string `json:"item_id"`
	}
)