package server

import (
	"log"

	"github.com/DioSurreal/Online-Shopping/modules/item/itemHandlers"
	itemPb "github.com/DioSurreal/Online-Shopping/modules/item/itemPb"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemRepositories"
	"github.com/DioSurreal/Online-Shopping/modules/item/itemUsecases"
	"github.com/DioSurreal/Online-Shopping/pkg/grpccon"
)

func (s *server) itemService() {
	itemRepository := itemRepositories.NewItemRepository(s.db)
	itemUsecase := itemUsecases.NewItemUsecase(itemRepository)
	itemHttpHandler := itemHandlers.NewItemHttpHandler(s.cfg,itemUsecase)
    itemGrpcHandler := itemHandlers.NewItemGrpcHandler(itemUsecase)

	//grpc
	go func() {
		grpcServer, lis := grpccon.NewGrpcServer(&s.cfg.Jwt, s.cfg.Grpc.ItemUrl)

		itemPb.RegisterItemGrpcServiceServer(grpcServer, itemGrpcHandler)

		log.Printf("Item gRPC server listening on %s", s.cfg.Grpc.ItemUrl)
		grpcServer.Serve(lis)
	}()
	
	_  = itemHttpHandler
	_ = itemGrpcHandler

	item := s.app.Group("/item_v1")

	//Health Check
	item.GET("",s.healthCheckService)
}