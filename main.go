/*
Copyright © 2024 Catizard <1185032459@qq.com>
*/
package main

import (
	"github.com/Catizard/lampghost/cmd"
	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/service"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/charmbracelet/log"
)

func main() {
	log.Default().SetLevel(log.DebugLevel)
	db := sqlite.NewDB(config.GetDSN())
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	service.InitService(db)
	cmd.Execute()
}
