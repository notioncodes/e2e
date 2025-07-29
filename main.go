package main

import (
	"context"
	"log"
	"time"

	"github.com/mateothegreat/go-multilog/multilog"
	"github.com/notioncodes/client"
	"github.com/notioncodes/plugin"
	redis "github.com/notioncodes/plugin-redis"
	"github.com/notioncodes/test"
	"github.com/notioncodes/types"
)

func main() {
	client, err := client.New(&client.Config{
		APIKey: test.TestConfig.NotionAPIKey,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// // defaultRedisConfig returns sensible defaults for Redis configuration.
	// func defaultRedisConfig() *ClientConfig {
	// 	return &ClientConfig{
	// 		Address:      "localhost:6379",
	// 		Database:     0,
	// 		KeyPrefix:    "notion",
	// 		KeySeparator: ":",
	// 		TTL:          24 * time.Hour,
	// 		PrettyJSON:   false,
	// 		IncludeMeta:  true,
	// 		Pipeline:     true,
	// 		BatchSize:    100,
	// 		MaxRetries:   3,
	// 		RetryBackoff: time.Second,
	// 	}
	// }

	redisPlugin, err := redis.NewPlugin(client, redis.Config{
		BaseConfig: plugin.BaseConfig{
			EnableReporter: true,
			Reporter: &plugin.Reporter{
				Interval: 3 * time.Second,
			},
		},
		ClientConfig: redis.ClientConfig{
			Address:  test.TestConfig.RedisAddress,
			Database: test.TestConfig.RedisDatabase,
			Username: test.TestConfig.RedisUsername,
			Password: test.TestConfig.RedisPassword,
		},
		Workers:           10,
		BatchSize:         10,
		Timeout:           30 * time.Second,
		ObjectTypes:       []types.ObjectType{types.ObjectTypePage, types.ObjectTypeDatabase},
		IncludeBlocks:     true,
		EnableProgress:    true,
		ProgressInterval:  2 * time.Second,
		ProgressBatchSize: 25,
		ContinueOnError:   true,
		MaxErrors:         10,
		Flush:             true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer redisPlugin.NotionClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if redisPlugin.Flush {
		err := redisPlugin.RedisClient.Flush()
		if err != nil {
			log.Fatal(err)
		}
	}

	result, err := redisPlugin.Service.ExportAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	multilog.Info("e2e", "export completed", map[string]interface{}{
		"duration": time.Since(result.Start),
		"success":  result.Success,
		"errors":   result.Errors,
	})
}
