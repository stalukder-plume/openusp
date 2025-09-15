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

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/n4-networks/openusp/pkg/cwmp"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	acs := &cwmp.AcsServer{}
	
	log.Println("Initializing CWMP ACS Server...")
	if err := acs.Init(); err != nil {
		log.Fatalf("Failed to initialize ACS server: %v", err)
	}

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down CWMP ACS Server...")
		if err := acs.Stop(); err != nil {
			log.Printf("Error during shutdown: %v", err)
		}
		os.Exit(0)
	}()

	log.Println("Starting CWMP ACS Server...")
	if err := acs.Start(); err != nil {
		log.Fatalf("ACS server exited with error: %v", err)
	}
}
