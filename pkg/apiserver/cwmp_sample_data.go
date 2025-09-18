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
	"fmt"
	"net/http"
	"time"
	"github.com/n4-networks/openusp/pkg/db"
)

// Helper function to create sample CWMP data for testing
func (as *ApiServer) createSampleCwmpData() error {
	if as.dbH.cwmpIntf == nil {
		return fmt.Errorf("CWMP database not connected")
	}

	// Create sample devices
	sampleDevices := []*db.CwmpDevice{
		{
			ID:                "cwmp:Broadcom:001122:RG:BCM001",
			OUI:               "001122",
			ProductClass:      "RG",
			SerialNumber:      "BCM001",
			Manufacturer:      "Broadcom",
			ModelName:         "BCM63138",
			HardwareVersion:   "1.0",
			SoftwareVersion:   "2.1.0-beta",
			SpecVersion:       "1.4",
			ConnectionRequestURL: "http://192.168.1.1:7547/",
			PeriodicInformEnable: true,
			PeriodicInformInterval: 300,
			LastInform:        time.Now().Add(-2 * time.Minute), // Online
			LastBootstrap:     time.Now().Add(-24 * time.Hour),
			CurrentTime:       time.Now(),
			UpTime:           86400, // 1 day in seconds
			IPAddress:        "192.168.1.1",
			Tags:             []string{"residential", "broadcom"},
			Parameters: map[string]string{
				"Device.DeviceInfo.SoftwareVersion":   "2.1.0-beta",
				"Device.DeviceInfo.HardwareVersion":   "1.0",
				"Device.DeviceInfo.SerialNumber":      "BCM001",
				"Device.DeviceInfo.ManufacturerOUI":   "001122",
				"Device.WiFi.Radio.1.Enable":          "true",
				"Device.WiFi.Radio.1.Channel":         "6",
				"Device.WiFi.SSID.1.SSID":            "HomeNetwork",
			},
			Events: []db.DeviceEvent{
				{
					EventCode:  "1 BOOT",
					CommandKey: "",
					Timestamp:  time.Now().Add(-24 * time.Hour),
				},
			},
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:                "cwmp:Netgear:334455:RG:R6300",
			OUI:               "334455",
			ProductClass:      "RG",
			SerialNumber:      "R6300_001",
			Manufacturer:      "Netgear",
			ModelName:         "R6300",
			HardwareVersion:   "2.0",
			SoftwareVersion:   "1.5.2",
			SpecVersion:       "1.4",
			ConnectionRequestURL: "http://192.168.100.1:7547/",
			PeriodicInformEnable: true,
			PeriodicInformInterval: 600,
			LastInform:        time.Now().Add(-10 * time.Minute), // Offline
			LastBootstrap:     time.Now().Add(-48 * time.Hour),
			CurrentTime:       time.Now(),
			UpTime:           172800, // 2 days in seconds
			IPAddress:        "192.168.100.1",
			Tags:             []string{"residential", "netgear"},
			Parameters: map[string]string{
				"Device.DeviceInfo.SoftwareVersion":   "1.5.2",
				"Device.DeviceInfo.HardwareVersion":   "2.0",
				"Device.DeviceInfo.SerialNumber":      "R6300_001",
				"Device.DeviceInfo.ManufacturerOUI":   "334455",
				"Device.WiFi.Radio.1.Enable":          "true",
				"Device.WiFi.Radio.1.Channel":         "11",
				"Device.WiFi.SSID.1.SSID":            "NetgearHome",
			},
			Events: []db.DeviceEvent{
				{
					EventCode:  "1 BOOT",
					CommandKey: "",
					Timestamp:  time.Now().Add(-48 * time.Hour),
				},
			},
			CreatedAt: time.Now().Add(-60 * 24 * time.Hour),
			UpdatedAt: time.Now().Add(-10 * time.Minute),
		},
	}

	// Insert sample devices
	for _, device := range sampleDevices {
		if err := as.dbH.cwmpIntf.UpsertCwmpDevice(device); err != nil {
			return fmt.Errorf("failed to insert sample device %s: %w", device.ID, err)
		}
	}

	// Create sample parameters
	sampleParameters := []db.CwmpParameter{
		{
			ID:         "param_bcm_sw",
			DeviceID:   "cwmp:Broadcom:001122:RG:BCM001",
			Path:       "Device.DeviceInfo.SoftwareVersion",
			Value:      "2.1.0-beta",
			Type:       "string",
			Writable:   false,
			LastUpdate: time.Now(),
		},
		{
			ID:         "param_bcm_hw",
			DeviceID:   "cwmp:Broadcom:001122:RG:BCM001",
			Path:       "Device.DeviceInfo.HardwareVersion",
			Value:      "1.0",
			Type:       "string",
			Writable:   false,
			LastUpdate: time.Now(),
		},
		{
			ID:         "param_bcm_wifi_enable",
			DeviceID:   "cwmp:Broadcom:001122:RG:BCM001",
			Path:       "Device.WiFi.Radio.1.Enable",
			Value:      "true",
			Type:       "boolean",
			Writable:   true,
			LastUpdate: time.Now(),
		},
		{
			ID:         "param_ng_sw",
			DeviceID:   "cwmp:Netgear:334455:RG:R6300",
			Path:       "Device.DeviceInfo.SoftwareVersion",
			Value:      "1.5.2",
			Type:       "string",
			Writable:   false,
			LastUpdate: time.Now(),
		},
		{
			ID:         "param_ng_wifi_enable",
			DeviceID:   "cwmp:Netgear:334455:RG:R6300",
			Path:       "Device.WiFi.Radio.1.Enable",
			Value:      "true",
			Type:       "boolean",
			Writable:   true,
			LastUpdate: time.Now(),
		},
	}

	// Insert sample parameters
	if err := as.dbH.cwmpIntf.UpsertCwmpParameters(sampleParameters); err != nil {
		return fmt.Errorf("failed to insert sample parameters: %w", err)
	}

	return nil
}

// API endpoint to populate sample data (for testing/demo purposes)
func (as *ApiServer) populateSampleCwmpData(w http.ResponseWriter, r *http.Request) {
	if err := as.createSampleCwmpData(); err != nil {
		httpSendRes(w, nil, fmt.Errorf("failed to create sample data: %w", err))
		return
	}

	response := map[string]interface{}{
		"message": "Sample CWMP data created successfully",
		"devices_created": 2,
		"parameters_created": 5,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	httpSendRes(w, response, nil)
}