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

package cli

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/abiosoft/ishell"
	"github.com/go-stomp/stomp"
	"github.com/n4-networks/openusp/pkg/config"
)

type cliCfg struct {
	apiServerAddr string
	stompAddr     string
	agentId       string
	histFile      string
	connTimeout   time.Duration
	logSetting    string
	authName      string
	authPasswd    string
}
type restHandler struct {
	client *http.Client
}

type shHandler struct {
	shell    *ishell.Shell
	histFile string
	cmds     map[string]*ishell.Cmd
}

type stompHandler struct {
	client *stomp.Conn
}

type Cli struct {
	cfg        cliCfg
	config     *config.Config
	agent      agentInfo
	stomp      stompHandler
	sh         shHandler
	rest       restHandler
	lastCmdErr error
}

func (cli *Cli) GetLastCmdErr() error {
	return cli.lastCmdErr
}

func (cli *Cli) ClearLastCmdErr() {
	cli.lastCmdErr = nil
}

func (cli *Cli) Init() error {

	if err := cli.loadConfig(); err != nil {
		log.Println("Could not configure CLI, err:", err)
		return err
	}

	// Initialize logging
	if err := cli.loggingInit(); err != nil {
		log.Println("Logging settings could not be applied")
	}

	// Initialization rest client
	if err := cli.restInit(); err != nil {
		log.Println("Could not initialize rest client:", err)
		return err
	}

	// Initialization of Agent Parameters
	if err := cli.initCliWithAgentParams(); err != nil {
		log.Println("Could not set agent information:", err)
	}
	log.Println("CLI version:", getVer())

	// Initialize shell
	cli.sh.shell = ishell.New()

	// Set default Prompt
	cli.sh.shell.SetPrompt("OpenUsp-Cli>> ")
	cli.sh.histFile = "history"
	cli.sh.shell.SetHistoryPath(cli.sh.histFile)

	// Initialize shell Cmds
	cli.sh.cmds = make(map[string]*ishell.Cmd)

	// Initialize shell Cmds
	cli.sh.cmds = make(map[string]*ishell.Cmd)

	// Register verb cmds
	cli.registerVerbs()

	// MTP and DB
	cli.registerNounsMtp()
	cli.registerNounsDb()
	cli.registerNounsStomp()

	// CLI related
	cli.registerNounsHistory()
	cli.registerNounsLogging()
	cli.registerNounsVersion()

	// Agent
	cli.registerNounsAgent()

	// Device Model
	cli.registerNounsDevice()
	cli.registerNounsDevInfo()
	cli.registerNounsBridging()
	cli.registerNounsDhcpv4()
	cli.registerNounsEth()
	cli.registerNounsIp()
	cli.registerNounsNat()
	cli.registerNounsWiFi()
	cli.registerNounsTime()
	cli.registerNounsNw()

	// Basic low level
	cli.registerNounsDatamodel()
	cli.registerNounsCommand()
	cli.registerNounsParam()
	cli.registerNounsInstance()

	return nil
}

func (cli *Cli) Run() {
	cli.sh.shell.Println("**************************************************************")
	cli.sh.shell.Println("                          OpenUsp Cli")
	cli.sh.shell.Println("**************************************************************")
	cli.sh.shell.Run()
}

func (cli *Cli) loadConfig() error {
	// Load YAML configuration - try to find cli.yaml specifically
	cfg, err := config.LoadConfig("./configs/cli.yaml")
	if err != nil {
		log.Printf("Error loading YAML configuration: %v", err)
		return err
	}
	
	cli.config = cfg

	// Map YAML config to legacy cliCfg struct for backward compatibility
	cli.cfg.apiServerAddr = fmt.Sprintf("http://%s", cfg.GetHTTPAddress())
	cli.cfg.stompAddr = cfg.GetStompAddress()
	cli.cfg.connTimeout = cfg.Database.Pool.Timeout
	cli.cfg.histFile = "history" // Default history file
	cli.cfg.logSetting = cfg.Logging.Level
	cli.cfg.authName = cfg.Security.Auth.Username
	cli.cfg.authPasswd = cfg.Security.Auth.Password
	
	// Agent ID from USP config
	if cfg.Security.USP.AgentID != "" {
		cli.cfg.agentId = cfg.Security.USP.AgentID
	} else {
		log.Println("Default Agent Endpoint ID is not defined, please configure through cli anytime")
	}

	log.Printf("Cli Config params: %+v\n", cli.cfg)
	return nil
}

func (cli *Cli) ProcessCmd(args string) error {
	log.Println("Running cli command in non interactive mode")
	log.Println("Processing cmd:", args)
	tok := strings.Split(args, " ")
	return cli.sh.shell.Process(tok...)
}

func (cli *Cli) SetOut(writer io.Writer) error {
	cli.sh.shell.SetOut(writer)
	return nil
}

// IsConnectedToDb checks if CLI has a valid REST connection to the API server
func (cli *Cli) IsConnectedToDb() bool {
	// Check if the REST client is configured and can reach the API server
	return cli.rest.client != nil && cli.cfg.apiServerAddr != ""
}

// IsConnectedToMtp checks if CLI has a valid MTP (STOMP) connection
func (cli *Cli) IsConnectedToMtp() bool {
	// Check if STOMP client is connected
	return cli.stomp.client != nil
}
