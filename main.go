package main

import (
	"context"
	"log"
	"os"

	"github.com/DioSurreal/Online-Shopping/config"
	"github.com/DioSurreal/Online-Shopping/pkg/database"
	"github.com/DioSurreal/Online-Shopping/server"
)

func main() {
	ctx := context.Background()
	_ = ctx

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error : .env path is required")
		}
		return os.Args[1]
	}())

	db := database.DbConn(ctx,&cfg)
	defer db.Disconnect(ctx)
	log.Println(db)

	server.Start(ctx,&cfg,db)

	
}