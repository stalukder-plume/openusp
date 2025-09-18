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
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	CwmpDeviceCollection   = "cwmpdevices"
	CwmpSessionCollection  = "cwmpsessions"
	CwmpParameterCollection = "cwmpparams"
	CwmpFileTransferCollection = "cwmpfiles"
)

// CwmpDevice represents a TR-069 device in the database
type CwmpDevice struct {
	ID                string            `bson:"_id" json:"id"`
	OUI               string            `bson:"oui" json:"oui"`
	ProductClass      string            `bson:"product_class" json:"product_class"`
	SerialNumber      string            `bson:"serial_number" json:"serial_number"`
	ManufacturerOUI   string            `bson:"manufacturer_oui" json:"manufacturer_oui"`
	Manufacturer      string            `bson:"manufacturer" json:"manufacturer"`
	ModelName         string            `bson:"model_name" json:"model_name"`
	Description       string            `bson:"description" json:"description"`
	ProductClass2     string            `bson:"product_class2" json:"product_class2"`
	HardwareVersion   string            `bson:"hardware_version" json:"hardware_version"`
	SoftwareVersion   string            `bson:"software_version" json:"software_version"`
	SpecVersion       string            `bson:"spec_version" json:"spec_version"`
	ProvisioningCode  string            `bson:"provisioning_code" json:"provisioning_code"`
	ConnectionRequestURL string         `bson:"connection_request_url" json:"connection_request_url"`
	ConnectionRequestUsername string    `bson:"connection_request_username" json:"connection_request_username"`
	ConnectionRequestPassword string    `bson:"connection_request_password" json:"connection_request_password"`
	PeriodicInformEnable bool           `bson:"periodic_inform_enable" json:"periodic_inform_enable"`
	PeriodicInformInterval int          `bson:"periodic_inform_interval" json:"periodic_inform_interval"`
	LastInform        time.Time         `bson:"last_inform" json:"last_inform"`
	LastBootstrap     time.Time         `bson:"last_bootstrap" json:"last_bootstrap"`
	CurrentTime       time.Time         `bson:"current_time" json:"current_time"`
	UpTime           int               `bson:"up_time" json:"up_time"`
	IPAddress        string            `bson:"ip_address" json:"ip_address"`
	Tags             []string          `bson:"tags" json:"tags"`
	Parameters       map[string]string `bson:"parameters" json:"parameters"`
	Events           []DeviceEvent     `bson:"events" json:"events"`
	CreatedAt        time.Time         `bson:"created_at" json:"created_at"`
	UpdatedAt        time.Time         `bson:"updated_at" json:"updated_at"`
}

// CwmpSession represents an active CWMP session
type CwmpSession struct {
	ID                string    `bson:"_id" json:"id"`
	DeviceID          string    `bson:"device_id" json:"device_id"`
	SessionID         string    `bson:"session_id" json:"session_id"`
	State             string    `bson:"state" json:"state"`
	CurrentRPCMethod  string    `bson:"current_rpc_method" json:"current_rpc_method"`
	PendingRPCs       []string  `bson:"pending_rpcs" json:"pending_rpcs"`
	LastActivity      time.Time `bson:"last_activity" json:"last_activity"`
	ConnectionRequestURL string `bson:"connection_request_url" json:"connection_request_url"`
	CreatedAt         time.Time `bson:"created_at" json:"created_at"`
}

// CwmpParameter represents a TR-069 device parameter
type CwmpParameter struct {
	ID         string    `bson:"_id" json:"id"`
	DeviceID   string    `bson:"device_id" json:"device_id"`
	Path       string    `bson:"path" json:"path"`
	Value      string    `bson:"value" json:"value"`
	Type       string    `bson:"type" json:"type"`
	Writable   bool      `bson:"writable" json:"writable"`
	LastUpdate time.Time `bson:"last_update" json:"last_update"`
}

// CwmpFileTransfer represents a file transfer operation
type CwmpFileTransfer struct {
	ID           string    `bson:"_id" json:"id"`
	DeviceID     string    `bson:"device_id" json:"device_id"`
	CommandKey   string    `bson:"command_key" json:"command_key"`
	FileType     string    `bson:"file_type" json:"file_type"`
	URL          string    `bson:"url" json:"url"`
	Username     string    `bson:"username" json:"username"`
	Password     string    `bson:"password" json:"password"`
	FileSize     int64     `bson:"file_size" json:"file_size"`
	TargetFileName string  `bson:"target_file_name" json:"target_file_name"`
	DelaySeconds int       `bson:"delay_seconds" json:"delay_seconds"`
	SuccessURL   string    `bson:"success_url" json:"success_url"`
	FailureURL   string    `bson:"failure_url" json:"failure_url"`
	Status       string    `bson:"status" json:"status"`
	StartTime    time.Time `bson:"start_time" json:"start_time"`
	CompleteTime time.Time `bson:"complete_time" json:"complete_time"`
	FaultCode    string    `bson:"fault_code" json:"fault_code"`
	FaultString  string    `bson:"fault_string" json:"fault_string"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
}

// DeviceEvent represents an event from a TR-069 device
type DeviceEvent struct {
	EventCode  string    `bson:"event_code" json:"event_code"`
	CommandKey string    `bson:"command_key" json:"command_key"`
	Timestamp  time.Time `bson:"timestamp" json:"timestamp"`
}

// CwmpDb extends UspDb with TR-069 specific collections
type CwmpDb struct {
	UspDb
	cwmpDeviceColl   *mongo.Collection
	cwmpSessionColl  *mongo.Collection
	cwmpParamColl    *mongo.Collection
	cwmpFileColl     *mongo.Collection
}

// InitCwmp initializes CWMP collections and creates indexes
func (c *CwmpDb) InitCwmp(client *mongo.Client) error {
	if client == nil {
		return errors.New("DB is not connected, please try again...")
	}

	// Initialize USP collections first
	if err := c.UspDb.Init(client); err != nil {
		return err
	}

	dbName := cfg.name
	c.cwmpDeviceColl = client.Database(dbName).Collection(CwmpDeviceCollection)
	c.cwmpSessionColl = client.Database(dbName).Collection(CwmpSessionCollection)
	c.cwmpParamColl = client.Database(dbName).Collection(CwmpParameterCollection)
	c.cwmpFileColl = client.Database(dbName).Collection(CwmpFileTransferCollection)

	// Create indexes for better performance
	return c.createCwmpIndexes()
}

// createCwmpIndexes creates necessary indexes for CWMP collections
func (c *CwmpDb) createCwmpIndexes() error {
	ctx := context.Background()

	// Device collection indexes
	deviceIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "oui", Value: 1}, {Key: "serial_number", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "last_inform", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "ip_address", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "manufacturer", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "model_name", Value: 1}},
		},
	}

	// Session collection indexes  
	sessionIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "device_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "session_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "last_activity", Value: -1}},
		},
		{
			Keys: bson.D{{Key: "state", Value: 1}},
		},
	}

	// Parameter collection indexes
	parameterIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "device_id", Value: 1}, {Key: "path", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "path", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "last_update", Value: -1}},
		},
	}

	// File transfer collection indexes
	fileIndexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "device_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "command_key", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
		},
	}

	// Create indexes
	if _, err := c.cwmpDeviceColl.Indexes().CreateMany(ctx, deviceIndexes); err != nil {
		return err
	}
	if _, err := c.cwmpSessionColl.Indexes().CreateMany(ctx, sessionIndexes); err != nil {
		return err
	}
	if _, err := c.cwmpParamColl.Indexes().CreateMany(ctx, parameterIndexes); err != nil {
		return err
	}
	if _, err := c.cwmpFileColl.Indexes().CreateMany(ctx, fileIndexes); err != nil {
		return err
	}

	return nil
}

// DeleteCwmpCollection drops a CWMP collection
func (c *CwmpDb) DeleteCwmpCollection(collName string) error {
	var err error
	ctx := context.Background()
	
	switch collName {
	case CwmpDeviceCollection:
		err = c.cwmpDeviceColl.Drop(ctx)
	case CwmpSessionCollection:
		err = c.cwmpSessionColl.Drop(ctx)
	case CwmpParameterCollection:
		err = c.cwmpParamColl.Drop(ctx)
	case CwmpFileTransferCollection:
		err = c.cwmpFileColl.Drop(ctx)
	default:
		err = errors.New("Invalid CWMP collection name: " + collName)
	}
	return err
}

// GetCwmpDeviceCollection returns the CWMP device collection
func (c *CwmpDb) GetCwmpDeviceCollection() *mongo.Collection {
	return c.cwmpDeviceColl
}

// GetCwmpSessionCollection returns the CWMP session collection
func (c *CwmpDb) GetCwmpSessionCollection() *mongo.Collection {
	return c.cwmpSessionColl
}

// GetCwmpParameterCollection returns the CWMP parameter collection
func (c *CwmpDb) GetCwmpParameterCollection() *mongo.Collection {
	return c.cwmpParamColl
}

// GetCwmpFileTransferCollection returns the CWMP file transfer collection
func (c *CwmpDb) GetCwmpFileTransferCollection() *mongo.Collection {
	return c.cwmpFileColl
}

// GetAllCwmpDevices retrieves all CWMP devices from the database
func (c *CwmpDb) GetAllCwmpDevices() ([]CwmpDevice, error) {
	if c.cwmpDeviceColl == nil {
		return nil, errors.New("CWMP device collection not initialized")
	}

	ctx := context.Background()
	cursor, err := c.cwmpDeviceColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var devices []CwmpDevice
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// GetCwmpDeviceByID retrieves a specific CWMP device by ID
func (c *CwmpDb) GetCwmpDeviceByID(deviceID string) (*CwmpDevice, error) {
	if c.cwmpDeviceColl == nil {
		return nil, errors.New("CWMP device collection not initialized")
	}

	ctx := context.Background()
	var device CwmpDevice
	err := c.cwmpDeviceColl.FindOne(ctx, bson.M{"_id": deviceID}).Decode(&device)
	if err != nil {
		return nil, err
	}

	return &device, nil
}

// GetCwmpDevicesByFilter retrieves CWMP devices based on filter criteria
func (c *CwmpDb) GetCwmpDevicesByFilter(filter bson.M) ([]CwmpDevice, error) {
	if c.cwmpDeviceColl == nil {
		return nil, errors.New("CWMP device collection not initialized")
	}

	ctx := context.Background()
	cursor, err := c.cwmpDeviceColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var devices []CwmpDevice
	if err = cursor.All(ctx, &devices); err != nil {
		return nil, err
	}

	return devices, nil
}

// GetCwmpParametersByDeviceID retrieves parameters for a specific CWMP device
func (c *CwmpDb) GetCwmpParametersByDeviceID(deviceID string) ([]CwmpParameter, error) {
	if c.cwmpParamColl == nil {
		return nil, errors.New("CWMP parameter collection not initialized")
	}

	ctx := context.Background()
	cursor, err := c.cwmpParamColl.Find(ctx, bson.M{"device_id": deviceID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var parameters []CwmpParameter
	if err = cursor.All(ctx, &parameters); err != nil {
		return nil, err
	}

	return parameters, nil
}

// GetCwmpParametersByPath retrieves specific parameters by path for a device
func (c *CwmpDb) GetCwmpParametersByPath(deviceID string, paths []string) ([]CwmpParameter, error) {
	if c.cwmpParamColl == nil {
		return nil, errors.New("CWMP parameter collection not initialized")
	}

	ctx := context.Background()
	filter := bson.M{
		"device_id": deviceID,
		"path":      bson.M{"$in": paths},
	}

	cursor, err := c.cwmpParamColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var parameters []CwmpParameter
	if err = cursor.All(ctx, &parameters); err != nil {
		return nil, err
	}

	return parameters, nil
}

// UpsertCwmpDevice inserts or updates a CWMP device
func (c *CwmpDb) UpsertCwmpDevice(device *CwmpDevice) error {
	if c.cwmpDeviceColl == nil {
		return errors.New("CWMP device collection not initialized")
	}

	ctx := context.Background()
	device.UpdatedAt = time.Now()
	
	opts := options.Replace().SetUpsert(true)
	_, err := c.cwmpDeviceColl.ReplaceOne(ctx, bson.M{"_id": device.ID}, device, opts)
	
	return err
}

// UpsertCwmpParameters inserts or updates CWMP parameters
func (c *CwmpDb) UpsertCwmpParameters(parameters []CwmpParameter) error {
	if c.cwmpParamColl == nil {
		return errors.New("CWMP parameter collection not initialized")
	}

	if len(parameters) == 0 {
		return nil
	}

	ctx := context.Background()
	var operations []mongo.WriteModel

	for _, param := range parameters {
		param.LastUpdate = time.Now()
		
		filter := bson.M{
			"device_id": param.DeviceID,
			"path":      param.Path,
		}
		
		replacement := bson.M{
			"device_id":   param.DeviceID,
			"path":        param.Path,
			"value":       param.Value,
			"type":        param.Type,
			"writable":    param.Writable,
			"last_update": param.LastUpdate,
		}
		
		operation := mongo.NewReplaceOneModel().SetFilter(filter).SetReplacement(replacement).SetUpsert(true)
		operations = append(operations, operation)
	}

	_, err := c.cwmpParamColl.BulkWrite(ctx, operations)
	return err
}