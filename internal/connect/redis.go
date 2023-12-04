package connect

import (
	"context"
	"crypto/tls"
	"fmt"
	
	"github.com/go-redis/redis/v8"
	
	"github.com/creamsensation/cp/env"
	"github.com/creamsensation/cp/internal/config"
	"github.com/creamsensation/cp/internal/style"
)

func Redis(config config.Cache) *redis.Client {
	var client *redis.Client
	if env.Development() {
		client = redis.NewClient(
			&redis.Options{
				Addr:     config.Address,
				Password: config.Password,
				DB:       config.Db,
			},
		)
	}
	if env.Production() {
		client = redis.NewClient(
			&redis.Options{
				Addr:     config.Address,
				Password: config.Password,
				DB:       config.Db,
				TLSConfig: &tls.Config{
					MinVersion: tls.VersionTLS12,
				},
			},
		)
	}
	ping := client.Ping(context.Background())
	if ping.Err() != nil {
		fmt.Printf("Redis: %s\n", style.RedColor.Render("ERROR"))
		fmt.Println(style.RedColor.Render(ping.Err().Error()))
	}
	if ping.Err() == nil {
		fmt.Printf("Redis: %s\n", style.EmeraldColor.Render("CONNECTED"))
	}
	return client
}
