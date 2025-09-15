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

package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/n4-networks/openusp/pkg/cwmp"
)

// TR-069 CWMP API endpoints
const (
	CWMP_GET_DEVICES        = "/cwmp/devices/"
	CWMP_GET_DEVICE         = "/cwmp/device/{deviceId}"
	CWMP_GET_PARAMS         = "/cwmp/device/{deviceId}/params"
	CWMP_SET_PARAMS         = "/cwmp/device/{deviceId}/params"
	CWMP_REBOOT_DEVICE      = "/cwmp/device/{deviceId}/reboot"
	CWMP_FACTORY_RESET      = "/cwmp/device/{deviceId}/factory-reset"
	CWMP_GET_DEVICE_INFO    = "/cwmp/device/{deviceId}/info"
	CWMP_DOWNLOAD           = "/cwmp/device/{deviceId}/download"
	CWMP_UPLOAD             = "/cwmp/device/{deviceId}/upload"
	CWMP_CONNECTION_REQUEST = "/cwmp/device/{deviceId}/connection-request"
)

// CwmpDeviceInfo represents device information for API responses
type CwmpDeviceInfo struct {
	DeviceId         string            `json:"device_id"`
	Manufacturer     string            `json:"manufacturer"`
	OUI              string            `json:"oui"`
	ProductClass     string            `json:"product_class"`
	SerialNumber     string            `json:"serial_number"`
	SoftwareVersion  string            `json:"software_version"`
	HardwareVersion  string            `json:"hardware_version"`
	LastInformTime   string            `json:"last_inform_time"`
	IsOnline         bool              `json:"is_online"`
	ParameterCount   int               `json:"parameter_count"`
	ConnectionRequestURL string        `json:"connection_request_url"`
}

// CwmpParameterRequest represents parameter operation request
type CwmpParameterRequest struct {
	ParameterNames []string                      `json:"parameter_names,omitempty"`
	Parameters     []cwmp.ParameterValueStruct   `json:"parameters,omitempty"`
	ParameterKey   string                        `json:"parameter_key,omitempty"`
}

// CwmpRebootRequest represents reboot request
type CwmpRebootRequest struct {
	CommandKey string `json:"command_key"`
}

// CwmpDownloadRequest represents download request
type CwmpDownloadRequest struct {
	CommandKey     string `json:"command_key"`
	FileType       string `json:"file_type"`
	URL            string `json:"url"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	FileSize       uint32 `json:"file_size"`
	TargetFileName string `json:"target_filename"`
	DelaySeconds   uint32 `json:"delay_seconds"`
	SuccessURL     string `json:"success_url"`
	FailureURL     string `json:"failure_url"`
}

// CwmpUploadRequest represents upload request
type CwmpUploadRequest struct {
	CommandKey   string `json:"command_key"`
	FileType     string `json:"file_type"`
	URL          string `json:"url"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DelaySeconds uint32 `json:"delay_seconds"`
}

// setCwmpRoutesHandlers sets up CWMP API routes
func (as *ApiServer) setCwmpRoutesHandlers() {
	// Device management endpoints
	as.router.HandleFunc(CWMP_GET_DEVICES, as.getCwmpDevices).Methods("GET")
	as.router.HandleFunc(CWMP_GET_DEVICE, as.getCwmpDevice).Methods("GET")
	as.router.HandleFunc(CWMP_GET_DEVICE_INFO, as.getCwmpDeviceInfo).Methods("GET")
	
	// Parameter management endpoints
	as.router.HandleFunc(CWMP_GET_PARAMS, as.getCwmpParams).Methods("GET")
	as.router.HandleFunc(CWMP_SET_PARAMS, as.setCwmpParams).Methods("POST")
	
	// Device control endpoints
	as.router.HandleFunc(CWMP_REBOOT_DEVICE, as.rebootCwmpDevice).Methods("POST")
	as.router.HandleFunc(CWMP_FACTORY_RESET, as.factoryResetCwmpDevice).Methods("POST")
	as.router.HandleFunc(CWMP_CONNECTION_REQUEST, as.connectionRequestCwmpDevice).Methods("POST")
	
	// File transfer endpoints
	as.router.HandleFunc(CWMP_DOWNLOAD, as.downloadCwmpDevice).Methods("POST")
	as.router.HandleFunc(CWMP_UPLOAD, as.uploadCwmpDevice).Methods("POST")
}

// getCwmpDevices returns all CWMP devices
func (as *ApiServer) getCwmpDevices(w http.ResponseWriter, r *http.Request) {
	// Get query parameters for filtering
	manufacturer := r.URL.Query().Get("manufacturer")
	productClass := r.URL.Query().Get("product_class")
	onlineOnly := r.URL.Query().Get("online_only") == "true"
	
	// Get devices from controller (mock implementation)
	devices := []CwmpDeviceInfo{
		{
			DeviceId:        "cwmp:Example:123456:RG:ABC123",
			Manufacturer:    "Example",
			OUI:            "123456",
			ProductClass:   "RG",
			SerialNumber:   "ABC123",
			SoftwareVersion: "1.0.0",
			HardwareVersion: "1.0",
			LastInformTime:  "2023-12-01T10:00:00Z",
			IsOnline:       true,
			ParameterCount: 150,
			ConnectionRequestURL: "http://192.168.1.1:7547/",
		},
	}
	
	// Apply filters
	var filteredDevices []CwmpDeviceInfo
	for _, device := range devices {
		if manufacturer != "" && !strings.Contains(strings.ToLower(device.Manufacturer), strings.ToLower(manufacturer)) {
			continue
		}
		if productClass != "" && !strings.Contains(strings.ToLower(device.ProductClass), strings.ToLower(productClass)) {
			continue
		}
		if onlineOnly && !device.IsOnline {
			continue
		}
		filteredDevices = append(filteredDevices, device)
	}
	
	httpSendRes(w, filteredDevices, nil)
}

// getCwmpDevice returns specific CWMP device information
func (as *ApiServer) getCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// Mock device info (in real implementation, get from controller)
	device := CwmpDeviceInfo{
		DeviceId:        deviceId,
		Manufacturer:    "Example",
		OUI:            "123456",
		ProductClass:   "RG",
		SerialNumber:   "ABC123",
		SoftwareVersion: "1.0.0",
		HardwareVersion: "1.0",
		LastInformTime:  "2023-12-01T10:00:00Z",
		IsOnline:       true,
		ParameterCount: 150,
		ConnectionRequestURL: "http://192.168.1.1:7547/",
	}
	
	httpSendRes(w, device, nil)
}

// getCwmpDeviceInfo returns detailed device information
func (as *ApiServer) getCwmpDeviceInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// Get detailed device info including all parameters
	deviceInfo := map[string]interface{}{
		"device_id": deviceId,
		"basic_info": CwmpDeviceInfo{
			DeviceId:        deviceId,
			Manufacturer:    "Example",
			OUI:            "123456",
			ProductClass:   "RG",
			SerialNumber:   "ABC123",
			SoftwareVersion: "1.0.0",
			HardwareVersion: "1.0",
			LastInformTime:  "2023-12-01T10:00:00Z",
			IsOnline:       true,
			ParameterCount: 150,
			ConnectionRequestURL: "http://192.168.1.1:7547/",
		},
		"capabilities": []string{"Download", "Upload", "Reboot", "FactoryReset"},
		"statistics": map[string]interface{}{
			"uptime": "7 days, 3 hours",
			"memory_usage": "45%",
			"cpu_usage": "12%",
		},
	}
	
	httpSendRes(w, deviceInfo, nil)
}

// getCwmpParams gets parameter values from CWMP device
func (as *ApiServer) getCwmpParams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// Get parameter names from query or body
	parameterNames := r.URL.Query()["param"]
	if len(parameterNames) == 0 {
		// If no specific parameters requested, return commonly requested ones
		parameterNames = []string{
			"Device.DeviceInfo.SoftwareVersion",
			"Device.DeviceInfo.HardwareVersion",
			"Device.DeviceInfo.ManufacturerOUI",
			"Device.DeviceInfo.SerialNumber",
		}
	}
	
	// Mock response (in real implementation, get from controller)
	parameters := []cwmp.ParameterValueStruct{
		{Name: "Device.DeviceInfo.SoftwareVersion", Value: "1.0.0", Type: "string"},
		{Name: "Device.DeviceInfo.HardwareVersion", Value: "1.0", Type: "string"},
		{Name: "Device.DeviceInfo.ManufacturerOUI", Value: "123456", Type: "string"},
		{Name: "Device.DeviceInfo.SerialNumber", Value: "ABC123", Type: "string"},
	}
	
	response := map[string]interface{}{
		"device_id":   deviceId,
		"parameters": parameters,
		"timestamp":  "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// setCwmpParams sets parameter values on CWMP device
func (as *ApiServer) setCwmpParams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	var req CwmpParameterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpSendRes(w, nil, fmt.Errorf("invalid request body: %w", err))
		return
	}
	
	if len(req.Parameters) == 0 {
		httpSendRes(w, nil, fmt.Errorf("parameters are required"))
		return
	}
	
	// In real implementation, send to controller
	// err := as.controller.SetCwmpParameters(deviceId, req.Parameters, req.ParameterKey)
	
	response := map[string]interface{}{
		"device_id":     deviceId,
		"status":       "success",
		"message":      fmt.Sprintf("Set %d parameters", len(req.Parameters)),
		"parameter_key": req.ParameterKey,
		"timestamp":    "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// rebootCwmpDevice reboots a CWMP device
func (as *ApiServer) rebootCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	var req CwmpRebootRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpSendRes(w, nil, fmt.Errorf("invalid request body: %w", err))
		return
	}
	
	// In real implementation, send to controller
	// err := as.controller.RebootCwmpDevice(deviceId, req.CommandKey)
	
	response := map[string]interface{}{
		"device_id":    deviceId,
		"status":      "success",
		"message":     "Reboot command sent",
		"command_key": req.CommandKey,
		"timestamp":   "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// factoryResetCwmpDevice performs factory reset on CWMP device
func (as *ApiServer) factoryResetCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// In real implementation, send to controller
	// err := as.controller.FactoryResetCwmpDevice(deviceId)
	
	response := map[string]interface{}{
		"device_id": deviceId,
		"status":   "success",
		"message":  "Factory reset command sent",
		"timestamp": "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// connectionRequestCwmpDevice initiates connection request to CWMP device
func (as *ApiServer) connectionRequestCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// In real implementation, send connection request to device
	// err := as.controller.SendConnectionRequest(deviceId)
	
	response := map[string]interface{}{
		"device_id": deviceId,
		"status":   "success",
		"message":  "Connection request sent",
		"timestamp": "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// downloadCwmpDevice initiates download to CWMP device
func (as *ApiServer) downloadCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	var req CwmpDownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpSendRes(w, nil, fmt.Errorf("invalid request body: %w", err))
		return
	}
	
	if req.URL == "" || req.FileType == "" {
		httpSendRes(w, nil, fmt.Errorf("URL and file_type are required"))
		return
	}
	
	// In real implementation, send to controller
	// err := as.controller.DownloadToCwmpDevice(deviceId, req)
	
	response := map[string]interface{}{
		"device_id":    deviceId,
		"status":      "success",
		"message":     "Download command sent",
		"command_key": req.CommandKey,
		"file_type":   req.FileType,
		"url":        req.URL,
		"timestamp":   "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}

// uploadCwmpDevice initiates upload from CWMP device
func (as *ApiServer) uploadCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	var req CwmpUploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpSendRes(w, nil, fmt.Errorf("invalid request body: %w", err))
		return
	}
	
	if req.URL == "" || req.FileType == "" {
		httpSendRes(w, nil, fmt.Errorf("URL and file_type are required"))
		return
	}
	
	// In real implementation, send to controller
	// err := as.controller.UploadFromCwmpDevice(deviceId, req)
	
	response := map[string]interface{}{
		"device_id":    deviceId,
		"status":      "success",
		"message":     "Upload command sent",
		"command_key": req.CommandKey,
		"file_type":   req.FileType,
		"url":        req.URL,
		"timestamp":   "2023-12-01T10:00:00Z",
	}
	
	httpSendRes(w, response, nil)
}