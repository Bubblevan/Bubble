package config

type Logger struct {
	Level        string `yaml:"level"`
	Prefix       string `yaml:"prefix"`
	Director     string `yaml:"director"`
	ShowLine     string `yaml:"show_line"`
	LogInConsole bool   `yaml:"log_in_console"`
}
