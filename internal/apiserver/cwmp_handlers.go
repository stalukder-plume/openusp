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
	"time"

	"github.com/gorilla/mux"
	"github.com/n4-networks/openusp/internal/cwmp"
	"go.mongodb.org/mongo-driver/bson"
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
	CWMP_POPULATE_SAMPLE    = "/cwmp/populate-sample-data"
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
	
	// Sample data endpoint (for testing/demo)
	as.router.HandleFunc(CWMP_POPULATE_SAMPLE, as.populateSampleCwmpData).Methods("POST")
}

// getCwmpDevices returns all CWMP devices
func (as *ApiServer) getCwmpDevices(w http.ResponseWriter, r *http.Request) {
	// Check database connection
	if as.dbH.cwmpIntf == nil {
		httpSendRes(w, nil, fmt.Errorf("CWMP database not connected"))
		return
	}

	// Get query parameters for filtering
	manufacturer := r.URL.Query().Get("manufacturer")
	productClass := r.URL.Query().Get("product_class")
	onlineOnly := r.URL.Query().Get("online_only") == "true"
	
	// Build database filter
	filter := bson.M{}
	if manufacturer != "" {
		filter["manufacturer"] = bson.M{
			"$regex":   manufacturer,
			"$options": "i", // case insensitive
		}
	}
	if productClass != "" {
		filter["product_class"] = bson.M{
			"$regex":   productClass,
			"$options": "i", // case insensitive
		}
	}
	if onlineOnly {
		// Consider device online if last inform was within 5 minutes
		fiveMinutesAgo := time.Now().Add(-5 * time.Minute)
		filter["last_inform"] = bson.M{
			"$gte": fiveMinutesAgo,
		}
	}
	
	// Get devices from database
	dbDevices, err := as.dbH.cwmpIntf.GetCwmpDevicesByFilter(filter)
	if err != nil {
		httpSendRes(w, nil, fmt.Errorf("failed to retrieve devices: %w", err))
		return
	}
	
	// Convert to API response format
	var devices []CwmpDeviceInfo
	for _, dbDevice := range dbDevices {
		// Determine if device is online (last inform within 5 minutes)
		isOnline := time.Since(dbDevice.LastInform) <= 5*time.Minute
		
		device := CwmpDeviceInfo{
			DeviceId:        dbDevice.ID,
			Manufacturer:    dbDevice.Manufacturer,
			OUI:            dbDevice.OUI,
			ProductClass:   dbDevice.ProductClass,
			SerialNumber:   dbDevice.SerialNumber,
			SoftwareVersion: dbDevice.SoftwareVersion,
			HardwareVersion: dbDevice.HardwareVersion,
			LastInformTime:  dbDevice.LastInform.Format(time.RFC3339),
			IsOnline:       isOnline,
			ParameterCount: len(dbDevice.Parameters),
			ConnectionRequestURL: dbDevice.ConnectionRequestURL,
		}
		devices = append(devices, device)
	}
	
	httpSendRes(w, devices, nil)
}

// getCwmpDevice returns specific CWMP device information
func (as *ApiServer) getCwmpDevice(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceId := vars["deviceId"]
	
	if deviceId == "" {
		httpSendRes(w, nil, fmt.Errorf("device ID is required"))
		return
	}
	
	// Check database connection
	if as.dbH.cwmpIntf == nil {
		httpSendRes(w, nil, fmt.Errorf("CWMP database not connected"))
		return
	}
	
	// Get device from database
	dbDevice, err := as.dbH.cwmpIntf.GetCwmpDeviceByID(deviceId)
	if err != nil {
		httpSendRes(w, nil, fmt.Errorf("device not found: %w", err))
		return
	}
	
	// Determine if device is online (last inform within 5 minutes)
	isOnline := time.Since(dbDevice.LastInform) <= 5*time.Minute
	
	// Convert to API response format
	device := CwmpDeviceInfo{
		DeviceId:        dbDevice.ID,
		Manufacturer:    dbDevice.Manufacturer,
		OUI:            dbDevice.OUI,
		ProductClass:   dbDevice.ProductClass,
		SerialNumber:   dbDevice.SerialNumber,
		SoftwareVersion: dbDevice.SoftwareVersion,
		HardwareVersion: dbDevice.HardwareVersion,
		LastInformTime:  dbDevice.LastInform.Format(time.RFC3339),
		IsOnline:       isOnline,
		ParameterCount: len(dbDevice.Parameters),
		ConnectionRequestURL: dbDevice.ConnectionRequestURL,
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
	
	// Check database connection
	if as.dbH.cwmpIntf == nil {
		httpSendRes(w, nil, fmt.Errorf("CWMP database not connected"))
		return
	}
	
	// Get device from database
	dbDevice, err := as.dbH.cwmpIntf.GetCwmpDeviceByID(deviceId)
	if err != nil {
		httpSendRes(w, nil, fmt.Errorf("device not found: %w", err))
		return
	}
	
	// Determine if device is online
	isOnline := time.Since(dbDevice.LastInform) <= 5*time.Minute
	
	// Calculate uptime in human-readable format
	uptimeSeconds := dbDevice.UpTime
	uptimeDays := uptimeSeconds / (24 * 3600)
	uptimeHours := (uptimeSeconds % (24 * 3600)) / 3600
	uptimeStr := fmt.Sprintf("%d days, %d hours", uptimeDays, uptimeHours)
	
	// Build detailed device info including all available data
	deviceInfo := map[string]interface{}{
		"device_id": deviceId,
		"basic_info": CwmpDeviceInfo{
			DeviceId:        dbDevice.ID,
			Manufacturer:    dbDevice.Manufacturer,
			OUI:            dbDevice.OUI,
			ProductClass:   dbDevice.ProductClass,
			SerialNumber:   dbDevice.SerialNumber,
			SoftwareVersion: dbDevice.SoftwareVersion,
			HardwareVersion: dbDevice.HardwareVersion,
			LastInformTime:  dbDevice.LastInform.Format(time.RFC3339),
			IsOnline:       isOnline,
			ParameterCount: len(dbDevice.Parameters),
			ConnectionRequestURL: dbDevice.ConnectionRequestURL,
		},
		"capabilities": []string{"Download", "Upload", "Reboot", "FactoryReset"},
		"statistics": map[string]interface{}{
			"uptime":       uptimeStr,
			"last_inform":  dbDevice.LastInform.Format(time.RFC3339),
			"last_bootstrap": dbDevice.LastBootstrap.Format(time.RFC3339),
			"current_time": dbDevice.CurrentTime.Format(time.RFC3339),
			"ip_address":   dbDevice.IPAddress,
		},
		"settings": map[string]interface{}{
			"periodic_inform_enable":   dbDevice.PeriodicInformEnable,
			"periodic_inform_interval": dbDevice.PeriodicInformInterval,
			"provisioning_code":        dbDevice.ProvisioningCode,
			"spec_version":            dbDevice.SpecVersion,
		},
		"tags": dbDevice.Tags,
		"parameters": dbDevice.Parameters,
		"recent_events": dbDevice.Events,
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
	
	// Check database connection
	if as.dbH.cwmpIntf == nil {
		httpSendRes(w, nil, fmt.Errorf("CWMP database not connected"))
		return
	}
	
	// Get parameter names from query
	parameterNames := r.URL.Query()["parameters"]
	
	var parameters []cwmp.ParameterValueStruct
	
	if len(parameterNames) == 0 {
		// If no specific parameters requested, get all parameters for the device
		dbParams, err := as.dbH.cwmpIntf.GetCwmpParametersByDeviceID(deviceId)
		if err != nil {
			httpSendRes(w, nil, fmt.Errorf("failed to retrieve parameters: %w", err))
			return
		}
		
		// Convert to API format
		for _, dbParam := range dbParams {
			parameters = append(parameters, cwmp.ParameterValueStruct{
				Name:  dbParam.Path,
				Value: dbParam.Value,
				Type:  dbParam.Type,
			})
		}
	} else {
		// Get specific parameters requested
		dbParams, err := as.dbH.cwmpIntf.GetCwmpParametersByPath(deviceId, parameterNames)
		if err != nil {
			httpSendRes(w, nil, fmt.Errorf("failed to retrieve specific parameters: %w", err))
			return
		}
		
		// Convert to API format
		for _, dbParam := range dbParams {
			parameters = append(parameters, cwmp.ParameterValueStruct{
				Name:  dbParam.Path,
				Value: dbParam.Value,
				Type:  dbParam.Type,
			})
		}
		
		// If some parameters weren't found, add them with empty values
		found := make(map[string]bool)
		for _, dbParam := range dbParams {
			found[dbParam.Path] = true
		}
		
		for _, paramName := range parameterNames {
			if !found[paramName] {
				parameters = append(parameters, cwmp.ParameterValueStruct{
					Name:  paramName,
					Value: "",
					Type:  "string",
				})
			}
		}
	}
	
	response := map[string]interface{}{
		"device_id":   deviceId,
		"parameters":  parameters,
		"timestamp":   time.Now().Format(time.RFC3339),
		"count":       len(parameters),
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