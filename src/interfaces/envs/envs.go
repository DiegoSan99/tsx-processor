package envs

import (
	"context"

	"github.com/DiegoSan99/transaction-processor/src/config"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func WithEnvs(ctx context.Context, cfg *config.AppConfig) context.Context {
	_ = godotenv.Load()
	envconfig.Process(ctx, cfg)
	ctx = context.WithValue(ctx, "envs", cfg)
	return ctx
}
