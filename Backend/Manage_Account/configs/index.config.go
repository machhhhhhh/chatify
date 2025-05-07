package configs

import "os"

type IServerSetting struct {
	Port    string
	BaseURL string
}

type IDatabaseSetting struct {
	Host         string
	User         string
	Password     string
	DatabaseName string
	DatabasePort string
	Timezone     string
}

type IAESSetting struct {
	AES_IV  string
	AES_KEY string
}

type IConfig struct {
	Name             string
	Version          string
	IsProductionMode bool
	ServerSetting    IServerSetting
	DatabaseSetting  IDatabaseSetting
	AESSetting       IAESSetting
}

var ENV IConfig = LoadConfig()

func GetEnv(key, default_value string) string {
	if value, is_exist := os.LookupEnv(key); is_exist == true {
		return value
	}
	return default_value
}

func LoadConfig() IConfig {
	return IConfig{
		Name:             "Loan Management",
		Version:          "1.0.0",
		IsProductionMode: GetEnv("ENV_MODE", "development") == "production",
		ServerSetting: IServerSetting{
			Port:    GetEnv("PORT", "3000"),
			BaseURL: GetEnv("SERVER_BASE_URL", "http://localhost:5555"),
		},
		DatabaseSetting: IDatabaseSetting{
			Host:         GetEnv("DATABASE_HOST", "127.0.0.1"),
			User:         GetEnv("DATABASE_USER", "admin"),
			Password:     GetEnv("DATABASE_PASSWORD", "P@ssw0rd#555"),
			DatabaseName: GetEnv("DATABASE_NAME", "CHATIFY"),
			DatabasePort: GetEnv("DATABASE_PORT", "5432"),
			Timezone:     GetEnv("DATABASE_TIMEZONE", "Asia/Bangkok"),
		},
		AESSetting: IAESSetting{
			AES_IV:  GetEnv("AES_IV", "4cneyoDet7Zrs3Wx"),
			AES_KEY: GetEnv("AES_KEY", "hiyt6nTEt6ASboHK0A4cneyoDet7Zrs3"),
		},
	}
}
