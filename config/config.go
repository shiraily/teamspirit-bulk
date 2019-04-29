package config

var Cfg config

type config struct {
	GoogleSheets struct {
		SheetID       string `yaml:"sheet_id"`
		SheetWorkTime string `yaml:"sheet_work_time"`
		SheetSetting  string `yaml:"sheet_setting"`
		ClientSecret  string `yaml:"client_secret"`
		TimeDayStart  string `yaml:"time_day_start"`
	} `yaml:"google_sheets"`
	TeamSpirit struct {
		Domain   string
		User     string
		Password string
	}
}
