package main

import (
	"context"
	"log"
	"time"

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
		BaseConfig: plugin.BaseConfig{
			EnableReporter: true,
			Reporter: &plugin.Reporter{
				Interval:  2 * time.Second,
				BatchSize: 2,
			},
		},
		ClientConfig: redis.ClientConfig{
			Address:  test.TestConfig.RedisAddress,
			Database: test.TestConfig.RedisDatabase,
			Username: test.TestConfig.RedisUsername,
			Password: test.TestConfig.RedisPassword,
		},
		ObjectTypes:       []types.ObjectType{types.ObjectTypePage, types.ObjectTypeDatabase},
		IncludeBlocks:     true,
		EnableProgress:    true,
		ProgressInterval:  2 * time.Second,
		ProgressBatchSize: 25,
		ContinueOnError:   true,
		MaxErrors:         10,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer redisPlugin.NotionClient.Close()

	// Export all content
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Printf("Starting export...")

	result, err := redisPlugin.Service.ExportAll(ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Exported %d objects successfully in %v", result.Success, time.Since(result.Start))
	log.Printf("Export stats: %+v", result)
}
