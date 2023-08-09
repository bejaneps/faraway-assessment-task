package cache

import (
	"context"
	"fmt"

	pkgRedis "github.com/redis/go-redis/v9"
)

type redis struct {
	c *pkgRedis.Client
}

// SRandMember executes SRANDMEMBER command in Redis
func (r *redis) SRandMember(ctx context.Context, key string) (string, error) {
	res, err := r.c.SRandMember(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to execute SRANDMEMBER: %w", err)
	}

	return res, nil
}

func (r *redis) Close() error {
	return r.c.Close()
}

func (r *redis) client() interface{} {
	return r.c
}
