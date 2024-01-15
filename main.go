package main

import (
	"log"

	"github.com/KurbanowS/Animal/config"
	"github.com/KurbanowS/Animal/internal/api"
	"github.com/KurbanowS/Animal/internal/store"
	"github.com/KurbanowS/Animal/internal/store/pgx"
	"github.com/KurbanowS/Animal/internal/utils"
	"github.com/gin-gonic/gin"
)

func main() {
	defer utils.InitLogs().Close()
	config.LoadConfig()
	defer store.Init().(*pgx.PgxStore).Close()

	routes := gin.Default()
	api.Routes(routes)
	if err := routes.Run(":8000"); err != nil {
		log.Fatal(err)
	}
}
