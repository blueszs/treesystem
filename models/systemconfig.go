package models

type SystemConfig struct {
	Old_Config_Key string `json:"old_config_key"`
	Config_Key     string `json:"config_key"`
	Config_Value   string `json:"config_value"`
}
