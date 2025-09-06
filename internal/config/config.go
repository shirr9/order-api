package config

type Config struct {
	Env        string `yaml:"env"`
	LogPath    string `yaml:"log_path"`
	PostgresDb `yaml:"postgres"`
}

// PostgresDb password must have correct name from .env, it must be connected with config.yaml
type PostgresDb struct {
	Username string `yaml:"user"`
	Password string `yaml:"password" env:"USER_PASSWORD" env-required:"true"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	SSlMode  string `yaml:"sslmode" env-default:"disable"`
}
