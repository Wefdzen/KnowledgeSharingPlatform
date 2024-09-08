package config

import (
	"encoding/json"
	"os"
)

// for config
func LaunchCfgDb(Cfg *SettingDb) {
	file, err := os.Open("config.cfg")
	if err != nil {
		panic("Can't open the file \"setting.cfg\"")
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	readSetting := make([]byte, fileInfo.Size())
	_, err = file.Read(readSetting)
	if err != nil {
		panic("can't read file")
	}

	err = json.Unmarshal(readSetting, &Cfg)
	if err != nil {
		panic("json err")
	}
}

type SettingDb struct {
	PGaddress  string
	PGpassword string
	PGuser     string
	PGdbname   string
	PGPort     string
}
