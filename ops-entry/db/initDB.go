package db

import "ops-entry/db/configManager"

func InitDb() error {
	err := configManager.Init()
	if err != nil {
		return err
	}
	return nil
}
