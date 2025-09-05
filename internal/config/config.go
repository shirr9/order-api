package config

type Config struct {
	PostgresDb `yaml:"postgres"`
}

// PostgresDb password must have right name from .env, it must be connected with config.yaml
type PostgresDb struct {
	Username string `yaml:"user"`
	Password string `yaml:"password" env:"USER_PASSWORD" env-required:"true"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	DbName   string `yaml:"db_name"`
	SSlMode  string `yaml:"sslmode" env-default:"disable"`
}
