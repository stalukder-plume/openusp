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
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/n4-networks/openusp/internal/cwmp"
	"github.com/n4-networks/openusp/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CwmpDevice represents a TR-069 device
type CwmpDevice struct {
	DeviceId     string
	Manufacturer string
	OUI          string
	ProductClass string
	SerialNumber string
	SoftwareVersion string
	HardwareVersion string
	LastInformTime  time.Time
	ConnectionRequestURL string
	ParameterKey    string
	IsOnline        bool
	Parameters      map[string]cwmp.ParameterValueStruct
	mutex          sync.RWMutex
}

// CwmpManager handles TR-069 device management within the controller
type CwmpManager struct {
	devices    map[string]*CwmpDevice
	acsServer  *cwmp.AcsServer
	mutex      sync.RWMutex
	cfg        CwmpConfig
	dbH        *db.CwmpDb
}

// CwmpConfig holds CWMP configuration
type CwmpConfig struct {
	EnableACS           bool
	ACSPort            string
	ConnectionRequestPort string
	PeriodicInformInterval uint32
	ConnectionRequestAuth  string
}

// InitCwmp initializes the CWMP manager
func (c *Cntlr) InitCwmp() error {
	log.Println("Initializing CWMP Manager...")
	
	c.cwmpMgr = &CwmpManager{
		devices: make(map[string]*CwmpDevice),
		dbH:     &c.dbH,
	}
	
	// Load CWMP configuration
	if err := c.cwmpMgr.loadConfig(); err != nil {
		return fmt.Errorf("failed to load CWMP config: %w", err)
	}
	
	// Initialize ACS server if enabled
	if c.cwmpMgr.cfg.EnableACS {
		c.cwmpMgr.acsServer = &cwmp.AcsServer{}
		if err := c.cwmpMgr.acsServer.Init(); err != nil {
			return fmt.Errorf("failed to initialize ACS server: %w", err)
		}
		
		// Start ACS server in background
		go func() {
			if err := c.cwmpMgr.acsServer.Start(); err != nil {
				log.Printf("ACS server error: %v", err)
			}
		}()
	}
	
	log.Println("CWMP Manager initialized successfully")
	return nil
}

// loadConfig loads CWMP configuration from environment
func (cm *CwmpManager) loadConfig() error {
	// Configuration loading logic would be implemented here
	// For now, use defaults
	cm.cfg = CwmpConfig{
		EnableACS: true,
		ACSPort:   "7547",
		ConnectionRequestPort: "7548",
		PeriodicInformInterval: 300,
		ConnectionRequestAuth: "Basic",
	}
	return nil
}

// RegisterCwmpDevice registers a new TR-069 device
func (cm *CwmpManager) RegisterCwmpDevice(deviceInfo *cwmp.DeviceIdStruct, parameterList []cwmp.ParameterValueStruct) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	deviceId := fmt.Sprintf("cwmp:%s:%s:%s:%s",
		deviceInfo.Manufacturer,
		deviceInfo.OUI,
		deviceInfo.ProductClass,
		deviceInfo.SerialNumber)
	
	device := &CwmpDevice{
		DeviceId:       deviceId,
		Manufacturer:   deviceInfo.Manufacturer,
		OUI:           deviceInfo.OUI,
		ProductClass:  deviceInfo.ProductClass,
		SerialNumber:  deviceInfo.SerialNumber,
		LastInformTime: time.Now(),
		IsOnline:      true,
		Parameters:    make(map[string]cwmp.ParameterValueStruct),
	}
	
	// Store device parameters
	for _, param := range parameterList {
		device.Parameters[param.Name] = param
		
		// Extract important parameters
		switch param.Name {
		case "Device.DeviceInfo.SoftwareVersion":
			device.SoftwareVersion = param.Value
		case "Device.DeviceInfo.HardwareVersion":
			device.HardwareVersion = param.Value
		case "Device.ManagementServer.ConnectionRequestURL":
			device.ConnectionRequestURL = param.Value
		case "Device.ManagementServer.ParameterKey":
			device.ParameterKey = param.Value
		}
	}
	
	cm.devices[deviceId] = device
	log.Printf("Registered CWMP device: %s", deviceId)
	
	// Store device in database
	return cm.storeDeviceInDB(device)
}

// GetCwmpDevice retrieves a CWMP device by ID
func (cm *CwmpManager) GetCwmpDevice(deviceId string) (*CwmpDevice, error) {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	if device, exists := cm.devices[deviceId]; exists {
		return device, nil
	}
	
	return nil, fmt.Errorf("device not found: %s", deviceId)
}

// GetAllCwmpDevices returns all registered CWMP devices
func (cm *CwmpManager) GetAllCwmpDevices() []*CwmpDevice {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	devices := make([]*CwmpDevice, 0, len(cm.devices))
	for _, device := range cm.devices {
		devices = append(devices, device)
	}
	
	return devices
}

// GetParameterValues requests parameter values from a CWMP device
func (cm *CwmpManager) GetParameterValues(deviceId string, parameterNames []string) error {
	device, err := cm.GetCwmpDevice(deviceId)
	if err != nil {
		return err
	}
	
	if !device.IsOnline {
		return fmt.Errorf("device is offline: %s", deviceId)
	}
	
	if cm.acsServer != nil {
		return cm.acsServer.GetParameterValues(deviceId, parameterNames)
	}
	
	return fmt.Errorf("ACS server not available")
}

// SetParameterValues sets parameter values on a CWMP device
func (cm *CwmpManager) SetParameterValues(deviceId string, parameters []cwmp.ParameterValueStruct, parameterKey string) error {
	device, err := cm.GetCwmpDevice(deviceId)
	if err != nil {
		return err
	}
	
	if !device.IsOnline {
		return fmt.Errorf("device is offline: %s", deviceId)
	}
	
	if cm.acsServer != nil {
		return cm.acsServer.SetParameterValues(deviceId, parameters, parameterKey)
	}
	
	return fmt.Errorf("ACS server not available")
}

// RebootCwmpDevice reboots a CWMP device
func (cm *CwmpManager) RebootCwmpDevice(deviceId string, commandKey string) error {
	device, err := cm.GetCwmpDevice(deviceId)
	if err != nil {
		return err
	}
	
	if !device.IsOnline {
		return fmt.Errorf("device is offline: %s", deviceId)
	}
	
	if cm.acsServer != nil {
		return cm.acsServer.RebootDevice(deviceId, commandKey)
	}
	
	return fmt.Errorf("ACS server not available")
}

// UpdateDeviceStatus updates device online status
func (cm *CwmpManager) UpdateDeviceStatus(deviceId string, isOnline bool) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	if device, exists := cm.devices[deviceId]; exists {
		device.IsOnline = isOnline
		if isOnline {
			device.LastInformTime = time.Now()
		}
		log.Printf("Device %s status updated: online=%v", deviceId, isOnline)
		return nil
	}
	
	return fmt.Errorf("device not found: %s", deviceId)
}

// UpdateDeviceParameters updates device parameters after receiving response
func (cm *CwmpManager) UpdateDeviceParameters(deviceId string, parameters []cwmp.ParameterValueStruct) error {
	device, err := cm.GetCwmpDevice(deviceId)
	if err != nil {
		return err
	}
	
	device.mutex.Lock()
	defer device.mutex.Unlock()
	
	for _, param := range parameters {
		device.Parameters[param.Name] = param
	}
	
	log.Printf("Updated parameters for device %s: %d parameters", deviceId, len(parameters))
	
	// Store updated parameters in database
	return cm.updateDeviceParametersInDB(deviceId, parameters)
}

// storeDeviceInDB stores device information in database
func (cm *CwmpManager) storeDeviceInDB(device *CwmpDevice) error {
	if cm.dbH == nil {
		return fmt.Errorf("database not initialized")
	}
	
	// Convert to database model
	dbDevice := &db.CwmpDevice{
		ID:                      device.DeviceId,
		OUI:                     device.OUI,
		ProductClass:           device.ProductClass,
		SerialNumber:           device.SerialNumber,
		Manufacturer:           device.Manufacturer,
		HardwareVersion:        device.HardwareVersion,
		SoftwareVersion:        device.SoftwareVersion,
		ConnectionRequestURL:   device.ConnectionRequestURL,
		LastInform:             device.LastInformTime,
		IPAddress:              "", // Set by session
		Parameters:             make(map[string]string),
		CreatedAt:              time.Now(),
		UpdatedAt:              time.Now(),
	}
	
	// Convert parameters
	for key, param := range device.Parameters {
		dbDevice.Parameters[key] = param.Value
	}
	
	// Insert into database
	collection := cm.dbH.GetCwmpDeviceCollection()
	ctx := context.Background()
	_, err := collection.InsertOne(ctx, dbDevice)
	if err != nil {
		log.Printf("Error storing CWMP device in database: %v", err)
		return err
	}
	
	log.Printf("Stored CWMP device in database: %s", device.DeviceId)
	return nil
}

// updateDeviceParametersInDB updates device parameters in database
func (cm *CwmpManager) updateDeviceParametersInDB(deviceId string, parameters []cwmp.ParameterValueStruct) error {
	if cm.dbH == nil {
		return fmt.Errorf("database not initialized")
	}
	
	collection := cm.dbH.GetCwmpParameterCollection()
	ctx := context.Background()
	
	for _, param := range parameters {
		dbParam := &db.CwmpParameter{
			ID:         fmt.Sprintf("%s_%s", deviceId, param.Name),
			DeviceID:   deviceId,
			Path:       param.Name,
			Value:      param.Value,
			Type:       param.Type,
			Writable:   true, // Default writable
			LastUpdate: time.Now(),
		}
		
		// Upsert parameter
		filter := bson.M{"device_id": deviceId, "path": param.Name}
		update := bson.M{"$set": dbParam}
		opts := options.Update().SetUpsert(true)
		_, err := collection.UpdateOne(ctx, filter, update, opts)
		if err != nil {
			log.Printf("Error updating parameter %s for device %s: %v", param.Name, deviceId, err)
			return err
		}
	}
	
	log.Printf("Updated CWMP device parameters in database: %s", deviceId)
	return nil
}

// MonitorCwmpDevices monitors CWMP device status and handles timeouts
func (cm *CwmpManager) MonitorCwmpDevices() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for range ticker.C {
		cm.checkDeviceTimeouts()
	}
}

// checkDeviceTimeouts checks for device timeouts and marks them offline
func (cm *CwmpManager) checkDeviceTimeouts() {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	
	timeout := time.Duration(cm.cfg.PeriodicInformInterval*2) * time.Second
	now := time.Now()
	
	for deviceId, device := range cm.devices {
		if device.IsOnline && now.Sub(device.LastInformTime) > timeout {
			device.IsOnline = false
			log.Printf("Device marked offline due to timeout: %s", deviceId)
		}
	}
}

// GetCwmpDevicesByFilter returns devices matching filter criteria
func (cm *CwmpManager) GetCwmpDevicesByFilter(manufacturer, productClass string) []*CwmpDevice {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	var filtered []*CwmpDevice
	for _, device := range cm.devices {
		match := true
		
		if manufacturer != "" && !strings.Contains(strings.ToLower(device.Manufacturer), strings.ToLower(manufacturer)) {
			match = false
		}
		
		if productClass != "" && !strings.Contains(strings.ToLower(device.ProductClass), strings.ToLower(productClass)) {
			match = false
		}
		
		if match {
			filtered = append(filtered, device)
		}
	}
	
	return filtered
}

// GetCwmpDeviceCount returns the total number of registered CWMP devices
func (cm *CwmpManager) GetCwmpDeviceCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	return len(cm.devices)
}

// GetOnlineCwmpDeviceCount returns the number of online CWMP devices
func (cm *CwmpManager) GetOnlineCwmpDeviceCount() int {
	cm.mutex.RLock()
	defer cm.mutex.RUnlock()
	
	count := 0
	for _, device := range cm.devices {
		if device.IsOnline {
			count++
		}
	}
	return count
}