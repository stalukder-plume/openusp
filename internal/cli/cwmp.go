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
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	"github.com/n4-networks/openusp/internal/cwmp"
)

// CWMP CLI commands and help text
const (
	showCwmpDevicesHelp    = "show cwmp devices [manufacturer] [product_class] - List all CWMP/TR-069 devices"
	showCwmpDeviceHelp     = "show cwmp device <device_id> - Show specific CWMP device information"
	getCwmpParamsHelp      = "get cwmp params <device_id> <param1> [param2] ... - Get parameter values from CWMP device"
	setCwmpParamsHelp      = "set cwmp params <device_id> <param=value> [param2=value2] ... - Set parameter values on CWMP device"
	rebootCwmpDeviceHelp   = "reboot cwmp device <device_id> [command_key] - Reboot CWMP device"
	factoryResetCwmpDeviceHelp = "factory-reset cwmp device <device_id> - Factory reset CWMP device"
	downloadCwmpFileHelp   = "download cwmp file <device_id> <url> <file_type> [target_filename] - Download file to CWMP device"
	uploadCwmpFileHelp     = "upload cwmp file <device_id> <url> <file_type> - Upload file from CWMP device"
	connectionRequestHelp  = "connection-request cwmp <device_id> - Send connection request to CWMP device"
)

// registerNounsCwmp registers CWMP-related CLI commands
func (cli *Cli) registerNounsCwmp() {
	cwmpCmds := []noun{
		{"show", "cwmp", showCwmpDevicesHelp, cli.showCwmpDevices},
		{"show.cwmp", "devices", showCwmpDevicesHelp, cli.showCwmpDevices},
		{"show.cwmp", "device", showCwmpDeviceHelp, cli.showCwmpDevice},
		{"get", "cwmp", getCwmpParamsHelp, cli.getCwmpParams},
		{"get.cwmp", "params", getCwmpParamsHelp, cli.getCwmpParams},
		{"set", "cwmp", setCwmpParamsHelp, cli.setCwmpParams},
		{"set.cwmp", "params", setCwmpParamsHelp, cli.setCwmpParams},
		{"reboot", "cwmp", rebootCwmpDeviceHelp, cli.rebootCwmpDevice},
		{"reboot.cwmp", "device", rebootCwmpDeviceHelp, cli.rebootCwmpDevice},
		{"factory-reset", "cwmp", factoryResetCwmpDeviceHelp, cli.factoryResetCwmpDevice},
		{"factory-reset.cwmp", "device", factoryResetCwmpDeviceHelp, cli.factoryResetCwmpDevice},
		{"download", "cwmp", downloadCwmpFileHelp, cli.downloadCwmpFile},
		{"download.cwmp", "file", downloadCwmpFileHelp, cli.downloadCwmpFile},
		{"upload", "cwmp", uploadCwmpFileHelp, cli.uploadCwmpFile},
		{"upload.cwmp", "file", uploadCwmpFileHelp, cli.uploadCwmpFile},
		{"connection-request", "cwmp", connectionRequestHelp, cli.connectionRequestCwmp},
	}
	cli.registerNouns(cwmpCmds)
}

// showCwmpDevices displays all CWMP devices
func (cli *Cli) showCwmpDevices(c *ishell.Context) {
	// Build query parameters
	queryParams := ""
	if len(c.Args) > 0 {
		params := make([]string, 0)
		for i, arg := range c.Args {
			switch i {
			case 0:
				if arg != "all" {
					params = append(params, "manufacturer="+arg)
				}
			case 1:
				params = append(params, "product_class="+arg)
			}
		}
		if len(params) > 0 {
			queryParams = "?" + strings.Join(params, "&")
		}
	}

	url := cli.cfg.apiServerAddr + "/cwmp/devices/" + queryParams
	data, err := cli.restGet(url)
	if err != nil {
		c.Printf("Error getting CWMP devices: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var devices []map[string]interface{}
	if err := json.Unmarshal(data, &devices); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	if len(devices) == 0 {
		c.Println("No CWMP devices found")
		cli.lastCmdErr = nil
		return
	}

	// Display device information
	c.Printf("Found %d CWMP device(s):\n", len(devices))
	c.Println("==========================================")
	
	for i, device := range devices {
		c.Printf("Device #%d:\n", i+1)
		c.Printf("  Device ID        : %v\n", device["device_id"])
		c.Printf("  Manufacturer     : %v\n", device["manufacturer"])
		c.Printf("  Product Class    : %v\n", device["product_class"])
		c.Printf("  Serial Number    : %v\n", device["serial_number"])
		c.Printf("  Software Version : %v\n", device["software_version"])
		c.Printf("  Online Status    : %v\n", device["is_online"])
		c.Printf("  Last Inform      : %v\n", device["last_inform_time"])
		c.Printf("  Parameters       : %v\n", device["parameter_count"])
		c.Println("------------------------------------------")
	}
	
	cli.lastCmdErr = nil
}

// showCwmpDevice displays specific CWMP device information
func (cli *Cli) showCwmpDevice(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Error: Device ID required")
		c.Println(showCwmpDeviceHelp)
		cli.lastCmdErr = errors.New("device ID required")
		return
	}

	deviceId := c.Args[0]
	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/info"
	data, err := cli.restGet(url)
	if err != nil {
		c.Printf("Error getting CWMP device info: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var deviceInfo map[string]interface{}
	if err := json.Unmarshal(data, &deviceInfo); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	// Display detailed device information
	c.Printf("CWMP Device Information for: %s\n", deviceId)
	c.Println("==========================================")
	
	if basicInfo, ok := deviceInfo["basic_info"].(map[string]interface{}); ok {
		c.Printf("Manufacturer     : %v\n", basicInfo["manufacturer"])
		c.Printf("OUI              : %v\n", basicInfo["oui"])
		c.Printf("Product Class    : %v\n", basicInfo["product_class"])
		c.Printf("Serial Number    : %v\n", basicInfo["serial_number"])
		c.Printf("Software Version : %v\n", basicInfo["software_version"])
		c.Printf("Hardware Version : %v\n", basicInfo["hardware_version"])
		c.Printf("Online Status    : %v\n", basicInfo["is_online"])
		c.Printf("Last Inform Time : %v\n", basicInfo["last_inform_time"])
		c.Printf("Connection URL   : %v\n", basicInfo["connection_request_url"])
		c.Printf("Parameter Count  : %v\n", basicInfo["parameter_count"])
	}

	if capabilities, ok := deviceInfo["capabilities"].([]interface{}); ok {
		c.Printf("Capabilities     : %v\n", capabilities)
	}

	if stats, ok := deviceInfo["statistics"].(map[string]interface{}); ok {
		c.Println("\nDevice Statistics:")
		for key, value := range stats {
			c.Printf("  %-15s: %v\n", key, value)
		}
	}

	cli.lastCmdErr = nil
}

// getCwmpParams gets parameter values from CWMP device
func (cli *Cli) getCwmpParams(c *ishell.Context) {
	if len(c.Args) < 2 {
		c.Println("Error: Device ID and parameter names required")
		c.Println(getCwmpParamsHelp)
		cli.lastCmdErr = errors.New("device ID and parameters required")
		return
	}

	deviceId := c.Args[0]
	paramNames := c.Args[1:]

	// Build query string
	queryParams := make([]string, len(paramNames))
	for i, param := range paramNames {
		queryParams[i] = "param=" + param
	}
	queryString := "?" + strings.Join(queryParams, "&")

	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/params" + queryString
	data, err := cli.restGet(url)
	if err != nil {
		c.Printf("Error getting CWMP parameters: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Parameters for device %s:\n", deviceId)
	c.Println("==========================================")

	if params, ok := response["parameters"].([]interface{}); ok {
		for _, p := range params {
			if param, ok := p.(map[string]interface{}); ok {
				c.Printf("%-50s : %v (%v)\n", param["Name"], param["Value"], param["Type"])
			}
		}
	}

	if timestamp, ok := response["timestamp"]; ok {
		c.Printf("\nRetrieved at: %v\n", timestamp)
	}

	cli.lastCmdErr = nil
}

// setCwmpParams sets parameter values on CWMP device
func (cli *Cli) setCwmpParams(c *ishell.Context) {
	if len(c.Args) < 2 {
		c.Println("Error: Device ID and parameter=value pairs required")
		c.Println(setCwmpParamsHelp)
		cli.lastCmdErr = errors.New("device ID and parameters required")
		return
	}

	deviceId := c.Args[0]
	paramPairs := c.Args[1:]

	// Parse parameter=value pairs
	var parameters []cwmp.ParameterValueStruct
	for _, pair := range paramPairs {
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			c.Printf("Error: Invalid parameter format '%s'. Use param=value\n", pair)
			cli.lastCmdErr = errors.New("invalid parameter format")
			return
		}
		parameters = append(parameters, cwmp.ParameterValueStruct{
			Name:  parts[0],
			Value: parts[1],
			Type:  "string", // Default type
		})
	}

	// Create request body
	requestBody := map[string]interface{}{
		"parameters":    parameters,
		"parameter_key": fmt.Sprintf("CLI_%d", len(parameters)),
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.Printf("Error creating request: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/params"
	data, err := cli.restPost(url, jsonData)
	if err != nil {
		c.Printf("Error setting CWMP parameters: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Set parameters result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}

// rebootCwmpDevice reboots a CWMP device
func (cli *Cli) rebootCwmpDevice(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Error: Device ID required")
		c.Println(rebootCwmpDeviceHelp)
		cli.lastCmdErr = errors.New("device ID required")
		return
	}

	deviceId := c.Args[0]
	commandKey := "CLI_REBOOT"
	if len(c.Args) > 1 {
		commandKey = c.Args[1]
	}

	requestBody := map[string]interface{}{
		"command_key": commandKey,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.Printf("Error creating request: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/reboot"
	data, err := cli.restPost(url, jsonData)
	if err != nil {
		c.Printf("Error rebooting CWMP device: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Reboot command result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}

// factoryResetCwmpDevice performs factory reset on CWMP device
func (cli *Cli) factoryResetCwmpDevice(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Error: Device ID required")
		c.Println(factoryResetCwmpDeviceHelp)
		cli.lastCmdErr = errors.New("device ID required")
		return
	}

	deviceId := c.Args[0]

	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/factory-reset"
	data, err := cli.restPost(url, []byte("{}"))
	if err != nil {
		c.Printf("Error factory resetting CWMP device: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Factory reset command result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}

// downloadCwmpFile downloads file to CWMP device
func (cli *Cli) downloadCwmpFile(c *ishell.Context) {
	if len(c.Args) < 3 {
		c.Println("Error: Device ID, URL, and file type required")
		c.Println(downloadCwmpFileHelp)
		cli.lastCmdErr = errors.New("device ID, URL, and file type required")
		return
	}

	deviceId := c.Args[0]
	url := c.Args[1]
	fileType := c.Args[2]
	targetFilename := ""
	if len(c.Args) > 3 {
		targetFilename = c.Args[3]
	}

	requestBody := map[string]interface{}{
		"command_key":      "CLI_DOWNLOAD",
		"file_type":       fileType,
		"url":            url,
		"target_filename": targetFilename,
		"delay_seconds":   0,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.Printf("Error creating request: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	apiUrl := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/download"
	data, err := cli.restPost(apiUrl, jsonData)
	if err != nil {
		c.Printf("Error downloading to CWMP device: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Download command result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}

// uploadCwmpFile uploads file from CWMP device
func (cli *Cli) uploadCwmpFile(c *ishell.Context) {
	if len(c.Args) < 3 {
		c.Println("Error: Device ID, URL, and file type required")
		c.Println(uploadCwmpFileHelp)
		cli.lastCmdErr = errors.New("device ID, URL, and file type required")
		return
	}

	deviceId := c.Args[0]
	url := c.Args[1]
	fileType := c.Args[2]

	requestBody := map[string]interface{}{
		"command_key":    "CLI_UPLOAD",
		"file_type":     fileType,
		"url":          url,
		"delay_seconds": 0,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		c.Printf("Error creating request: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	apiUrl := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/upload"
	data, err := cli.restPost(apiUrl, jsonData)
	if err != nil {
		c.Printf("Error uploading from CWMP device: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Upload command result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}

// connectionRequestCwmp sends connection request to CWMP device
func (cli *Cli) connectionRequestCwmp(c *ishell.Context) {
	if len(c.Args) < 1 {
		c.Println("Error: Device ID required")
		c.Println(connectionRequestHelp)
		cli.lastCmdErr = errors.New("device ID required")
		return
	}

	deviceId := c.Args[0]

	url := cli.cfg.apiServerAddr + "/cwmp/device/" + deviceId + "/connection-request"
	data, err := cli.restPost(url, []byte("{}"))
	if err != nil {
		c.Printf("Error sending connection request to CWMP device: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	var response map[string]interface{}
	if err := json.Unmarshal(data, &response); err != nil {
		c.Printf("Error parsing response: %v\n", err)
		cli.lastCmdErr = err
		return
	}

	c.Printf("Connection request result: %v\n", response["status"])
	c.Printf("Message: %v\n", response["message"])
	
	cli.lastCmdErr = nil
}