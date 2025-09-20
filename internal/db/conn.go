// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/n4-networks/openusp/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type dbCfg struct {
	serverAddr string
	name       string
	userName   string
	passwd     string
	timeout    int // in minute
}

var cfg dbCfg

func readConfigFromYAML() error {
	// Try to load configuration from YAML files
	// First try service-specific configs, then fall back to generic config
	yamlConfig, err := tryLoadConfig()
	if err != nil {
		log.Printf("Failed to load YAML configuration: %v", err)
		return err
	}

	// Map YAML config to legacy dbCfg struct
	cfg.serverAddr = fmt.Sprintf("%s:%d", yamlConfig.Database.Host, yamlConfig.Database.Port)
	cfg.userName = yamlConfig.Database.Username
	cfg.passwd = yamlConfig.Database.Password
	cfg.name = yamlConfig.Database.Name
	
	// Convert timeout from duration to minutes (legacy format)
	if yamlConfig.Database.Pool.Timeout > 0 {
		cfg.timeout = int(yamlConfig.Database.Pool.Timeout.Minutes())
	} else {
		cfg.timeout = 3 // Default 3 minutes
	}

	log.Printf("DB Config params: %+v\n", cfg)
	return nil
}

// tryLoadConfig attempts to load configuration from various sources
func tryLoadConfig() (*config.Config, error) {
	// Try service-specific configs first
	configFiles := []string{
		"./configs/apiserver.yaml",
		"./configs/controller.yaml", 
		"./configs/cli.yaml",
		"./configs/cwmpacs.yaml",
	}

	for _, configFile := range configFiles {
		if yamlConfig, err := config.LoadConfig(configFile); err == nil {
			return yamlConfig, nil
		}
	}

	// Fall back to generic config loading
	return config.LoadConfig("")
}

func Connect() (*mongo.Client, error) {
	if err := readConfigFromYAML(); err != nil {
		return nil, err
	}
	cred := options.Credential{Username: cfg.userName, Password: cfg.passwd}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + cfg.serverAddr).SetAuth(cred))
	if err != nil {
		return nil, err
	}
	timeout := time.Duration(cfg.timeout) * time.Minute
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	return client, err
}

func ConnectWithParams(addr string, user string, passwd string, timeout time.Duration) (*mongo.Client, error) {
	cred := options.Credential{Username: user, Password: passwd}
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://" + addr).SetAuth(cred))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err = client.Connect(ctx); err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	return client, err
}

func ConnectCache(addr string, timeout time.Duration) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	ctx, _ := context.WithTimeout(context.Background(), timeout)
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		log.Println("Error in connecting to redis", addr)
		return nil, err
	}
	fmt.Println(pong)
	return client, nil

}
