package config

type AppConfig struct {
	AWSRegion          string `env:"AWS_REGION, default=us-east-1"`
	AWSAccessKeyID     string `env:"AWS_AK, default="`
	AWSSecretAccessKey string `env:"AWS_AKS, default="`
	DbHost             string `env:"DB_HOST"`
	DbUser             string `env:"DB_USER, default=postgres"`
	DbPassword         string `env:"DB_PASSWORD, default=postgres"`
	DbName             string `env:"DB_NAME, default=postgres"`
	DbPort             string `env:"DB_PORT, default=5432"`
	EmailHost          string `env:"EMAIL_HOST, default=localhost"`
	EmailPort          string `env:"EMAIL_PORT, default=1025"`
	EmailUsername      string `env:"EMAIL_USERNAME, default="`
	EmailPassword      string `env:"EMAIL_PASSWORD, default="`
}
