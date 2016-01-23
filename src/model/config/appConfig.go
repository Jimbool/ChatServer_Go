package config

type AppConfig struct {
	AppId   string //App唯一标识
	AppName string //App名称
	AppKey  string //App加密Key
}

func NewAppConfig(appId, appName, appKey string) *AppConfig {
	return &AppConfig{
		AppId:   appId,
		AppName: appName,
		AppKey:  appKey,
	}
}
