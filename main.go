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

	redisPlugin, err := redis.NewPlugin(client, redis.Config{
		Config: plugin.Config{
			EnableReporter: true,
			Reporter: &plugin.Reporter{
				Interval: 3 * time.Second,
			},
		},
		ClientConfig: redis.ClientConfig{
			Address:      test.TestConfig.RedisAddress,
			Database:     test.TestConfig.RedisDatabase,
			Username:     test.TestConfig.RedisUsername,
			Password:     test.TestConfig.RedisPassword,
			KeyPrefix:    "test",
			KeySeparator: ":",
			TTL:          24 * time.Hour,
			PrettyJSON:   true,
			IncludeMeta:  true,
			Pipeline:     true,
			BatchSize:    100,
			MaxRetries:   3,
			RetryBackoff: time.Second,
		},
		Common: plugin.CommonSettings{
			Workers:         1,
			RuntimeTimeout:  30 * time.Second,
			RequestTimeout:  10 * time.Second,
			RequestDelay:    1 * time.Second,
			ContinueOnError: false,
		},
		Content: plugin.ContentSettings{
			Flush: true,
			Types: []types.ObjectType{
				types.ObjectTypePage,
				types.ObjectTypeDatabase,
				types.ObjectTypeBlock,
			},
			Databases: plugin.DatabaseSettings{
				Pages:  true,
				Blocks: true,
			},
			Pages: plugin.PageSettings{
				Blocks: true,
				// Comments:    true,
				// Attachments: true,
			},
			Blocks: plugin.BlockSettings{
				Children: true,
			},
		},
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

	result, err := redisPlugin.Service.Export(ctx)
	// result, err := redisPlugin.Service.ExportDatabase(ctx, types.DatabaseID("23fd7342e57181668a2ee373221477ad"), true)
	if err != nil {
		log.Fatal(err)
	}

	multilog.Info("e2e", "export completed", map[string]interface{}{
		"duration": time.Since(result.Start),
		"success":  result.Success,
		"errors":   result.Errors,
	})
}
