package v1

import "backend/api/v1/settings_api"

type ApiGroup struct {
	SettingsApi settings_api.SettingsApi
}

var ApiGroupApp = new(ApiGroup)
