/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package main

import (
	"log"

	"github.com/Catizard/lampghost/cmd"
	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/Catizard/lampghost/internal/sqlite/service"
)

func main() {
	db := sqlite.NewDB(config.GetDSN())
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service.InitService(db)
	cmd.Execute()
}
