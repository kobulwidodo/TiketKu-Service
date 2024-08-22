package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	Nil = redis.Nil
)

type Interface interface {
	BatchLock(ctx context.Context, lockKey []string, ttl time.Duration) (bool, error)
	BatchReleaseLock(ctx context.Context, lockKey []string) error
}

type TLSConfig struct {
	Enabled            bool
	InsecureSkipVerify bool
}

type Config struct {
	Protocol string
	Host     string
	Port     string
	Username string
	Password string
	TLS      TLSConfig
}

type cache struct {
	conf Config
	rdb  *redis.Client
}

func Init(cfg Config) Interface {
	c := &cache{
		conf: cfg,
	}
	c.connect(context.Background())
	return c
}

func (c *cache) connect(ctx context.Context) {
	redisOpts := redis.Options{
		Network:  c.conf.Protocol,
		Addr:     fmt.Sprintf("%s:%s", c.conf.Host, c.conf.Port),
		Username: c.conf.Username,
		Password: c.conf.Password,
	}

	if c.conf.TLS.Enabled {
		redisOpts.TLSConfig = &tls.Config{
			InsecureSkipVerify: c.conf.TLS.InsecureSkipVerify,
		}
	}

	client := redis.NewClient(&redisOpts)

	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("[FATAL] cannot connect to redis on address @%s:%v, with error: %s", c.conf.Host, c.conf.Port, err)
	}
	c.rdb = client
	log.Printf("REDIS: Address @%s:%v", c.conf.Host, c.conf.Port)
}

func (c *cache) BatchLock(ctx context.Context, lockKey []string, ttl time.Duration) (bool, error) {
	pipe := c.rdb.Pipeline()

	lockValues := make(map[string]string)
	for _, key := range lockKey {
		lockVal := fmt.Sprintf("%d", time.Now().UnixNano())
		pipe.SetNX(ctx, key, lockVal, ttl)
		lockValues[key] = lockVal
	}

	cmds, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	lockedKeys := []string{}
	for i, cmd := range cmds {
		if cmd.(*redis.BoolCmd).Val() {
			lockedKeys = append(lockedKeys, lockKey[i])
			continue
		}
		for _, lockKey := range lockedKeys {
			c.rdb.Del(ctx, lockKey)
		}
		return false, nil
	}

	return true, nil
}

func (c *cache) BatchReleaseLock(ctx context.Context, lockKey []string) error {
	pipe := c.rdb.Pipeline()

	for _, key := range lockKey {
		pipe.Del(ctx, key)
	}

	_, err := pipe.Exec(ctx)
	return err
}
