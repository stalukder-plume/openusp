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

package cntlr

import (
	"log"
	"strconv"

	"github.com/n4-networks/openusp/pkg/config"
)

const (
	_                  = iota
	SERVER_MODE_NORMAL = iota
	SERVER_MODE_TLS
	SERVER_MODE_NORMAL_AND_TLS
)

type grpcCfg struct {
	port string
}

type cacheCfg struct {
	serverAddr string
}

type uspCfg struct {
	endpointId        string
	protoVersion      string
	protoVersionCheck bool
}

type cntlrCfg struct {
	cache cacheCfg
	grpc  grpcCfg
	usp   uspCfg
}

func (c *Cntlr) loadConfig() error {
	// Load YAML configuration - try to find controller.yaml specifically
	cfg, err := config.LoadConfig("./configs/controller.yaml")
	if err != nil {
		log.Printf("Error loading YAML configuration: %v", err)
		return err
	}
	
	c.config = cfg

	// Map YAML config to legacy cntlrCfg struct for backward compatibility
	c.cfg.cache.serverAddr = cfg.GetCacheAddress()
	c.cfg.grpc.port = strconv.Itoa(cfg.Protocols.GRPC.Port)
	c.cfg.usp.endpointId = cfg.Security.USP.ControllerEndpointID
	c.cfg.usp.protoVersion = cfg.Security.USP.ProtocolVersion
	c.cfg.usp.protoVersionCheck = cfg.Security.USP.VersionCheck

	log.Printf("CNTLR Config params: %+v\n", c.cfg)

	return nil
}
