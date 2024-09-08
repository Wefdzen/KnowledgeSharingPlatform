package database

import (
	"fmt"
	"log"

	"wefdzen/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	//lauch config for Cfg
	var Cfg config.SettingDb
	config.LaunchCfgDb(&Cfg)
	//connect
	urlToDataBase := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", Cfg.PGuser, Cfg.PGpassword, Cfg.PGaddress, Cfg.PGPort, Cfg.PGdbname)
	db, err := gorm.Open(postgres.Open(urlToDataBase), &gorm.Config{})
	if err != nil {
		log.Fatal("can't open database")
		return nil, err
	}
	db.AutoMigrate(&Account{}) // если такой структуры небыло migrate will be create a new table
	// db.Create(&model.User{Login: "wefd2", Password: "1232", FullName: "Domitr V2"}) //add a new record
	return db, nil
}
