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

package cwmp

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// AcsConfig holds ACS server configuration
type AcsConfig struct {
	httpPort     string
	httpsPort    string
	isTlsEnabled bool
	certFile     string
	keyFile      string
	dbAddr       string
	sessionTimeout uint32
	informInterval uint32
	logLevel     string
}

// AcsServer represents the TR-069 ACS server
type AcsServer struct {
	cfg      AcsConfig
	dbClient *mongo.Client
	sessions map[string]*CwmpSession
	mutex    sync.RWMutex
	server   *http.Server
}

// CwmpSession represents a TR-069 CWMP session with a device
type CwmpSession struct {
	DeviceId     string
	SessionId    string
	CreatedTime  time.Time
	LastActivity time.Time
	HoldRequests bool
	MaxEnvelopes uint32
	State        SessionState
	PendingRPCs  []interface{}
	mutex        sync.RWMutex
}

type SessionState int

const (
	SessionStateNew SessionState = iota
	SessionStateInform
	SessionStateActive
	SessionStateClosed
)

// Init initializes the ACS server
func (acs *AcsServer) Init() error {
	log.Println("Initializing TR-069 ACS Server...")
	
	if err := acs.loadConfig(); err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if err := acs.connectDB(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	acs.sessions = make(map[string]*CwmpSession)
	
	// Initialize HTTP routes
	acs.initRoutes()
	
	log.Println("TR-069 ACS Server initialized successfully")
	return nil
}

// loadConfig loads configuration from environment variables
func (acs *AcsServer) loadConfig() error {
	if port, ok := os.LookupEnv("CWMP_HTTP_PORT"); ok {
		acs.cfg.httpPort = port
	} else {
		acs.cfg.httpPort = "7547"
	}

	if port, ok := os.LookupEnv("CWMP_HTTPS_PORT"); ok {
		acs.cfg.httpsPort = port
	} else {
		acs.cfg.httpsPort = "7548"
	}

	if tlsEnabled, ok := os.LookupEnv("CWMP_TLS_ENABLED"); ok {
		acs.cfg.isTlsEnabled = tlsEnabled == "true"
	} else {
		acs.cfg.isTlsEnabled = false
	}

	if cert, ok := os.LookupEnv("CWMP_CERT_FILE"); ok {
		acs.cfg.certFile = cert
	} else {
		acs.cfg.certFile = "server.crt"
	}

	if key, ok := os.LookupEnv("CWMP_KEY_FILE"); ok {
		acs.cfg.keyFile = key
	} else {
		acs.cfg.keyFile = "server.key"
	}

	if dbAddr, ok := os.LookupEnv("DB_ADDR"); ok {
		acs.cfg.dbAddr = dbAddr
	} else {
		acs.cfg.dbAddr = "localhost:27017"
	}

	if timeout, ok := os.LookupEnv("CWMP_SESSION_TIMEOUT"); ok {
		if t, err := strconv.ParseUint(timeout, 10, 32); err == nil {
			acs.cfg.sessionTimeout = uint32(t)
		}
	} else {
		acs.cfg.sessionTimeout = 30
	}

	if interval, ok := os.LookupEnv("CWMP_INFORM_INTERVAL"); ok {
		if i, err := strconv.ParseUint(interval, 10, 32); err == nil {
			acs.cfg.informInterval = uint32(i)
		}
	} else {
		acs.cfg.informInterval = 300
	}

	log.Printf("CWMP ACS Config: %+v", acs.cfg)
	return nil
}

// connectDB establishes database connection
func (acs *AcsServer) connectDB() error {
	// Database connection logic would be implemented here
	// For now, we'll use a placeholder
	log.Println("Connected to database for CWMP ACS")
	return nil
}

// initRoutes sets up HTTP routes for CWMP
func (acs *AcsServer) initRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", acs.handleCwmpRequest)
	mux.HandleFunc("/tr069", acs.handleCwmpRequest)
	mux.HandleFunc("/cwmp", acs.handleCwmpRequest)
	
	acs.server = &http.Server{
		Addr:    ":" + acs.cfg.httpPort,
		Handler: mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}

// Start starts the ACS server
func (acs *AcsServer) Start() error {
	log.Printf("Starting TR-069 ACS Server on port %s", acs.cfg.httpPort)
	
	if acs.cfg.isTlsEnabled {
		// Load TLS certificate
		cert, err := tls.LoadX509KeyPair(acs.cfg.certFile, acs.cfg.keyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}
		
		acs.server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		acs.server.Addr = ":" + acs.cfg.httpsPort
		return acs.server.ListenAndServeTLS("", "")
	}
	
	return acs.server.ListenAndServe()
}

// Stop gracefully stops the ACS server
func (acs *AcsServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return acs.server.Shutdown(ctx)
}

// handleCwmpRequest handles incoming CWMP SOAP requests
func (acs *AcsServer) handleCwmpRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received CWMP request from %s", r.RemoteAddr)
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Set SOAP headers
	w.Header().Set("Content-Type", "text/xml; charset=utf-8")
	w.Header().Set("SOAPAction", "")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Pragma", "no-cache")

	// Handle empty body (HTTP POST without SOAP content)
	if len(body) == 0 {
		log.Println("Received empty request body, sending empty response")
		acs.sendEmptyResponse(w)
		return
	}

	// Parse SOAP envelope
	var envelope SOAPEnvelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		log.Printf("Error parsing SOAP envelope: %v", err)
		acs.sendSOAPFault(w, FaultInvalidArguments, "Invalid SOAP envelope")
		return
	}

	// Route to appropriate handler based on SOAP body content
	response, err := acs.processSOAPRequest(&envelope, r)
	if err != nil {
		log.Printf("Error processing SOAP request: %v", err)
		acs.sendSOAPFault(w, FaultInternalError, err.Error())
		return
	}

	// Send response
	responseXML, err := xml.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Printf("Error marshaling response: %v", err)
		acs.sendSOAPFault(w, FaultInternalError, "Error creating response")
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	w.Write(responseXML)
}

// processSOAPRequest processes different types of SOAP requests
func (acs *AcsServer) processSOAPRequest(envelope *SOAPEnvelope, r *http.Request) (*SOAPEnvelope, error) {
	// Create response envelope
	response := &SOAPEnvelope{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		CwmpNS: "urn:dslforum-org:cwmp-1-2",
		XsiNS:  "http://www.w3.org/2001/XMLSchema-instance",
		XsdNS:  "http://www.w3.org/2001/XMLSchema",
		Header: &SOAPHeader{},
		Body:   SOAPBody{},
	}

	// Extract body content and determine request type
	bodyBytes, err := xml.Marshal(envelope.Body.Content)
	if err != nil {
		return nil, fmt.Errorf("error marshaling body content: %w", err)
	}

	// Check for Inform method
	if strings.Contains(string(bodyBytes), "Inform") {
		return acs.handleInform(envelope, response, r)
	}

	// Check for GetParameterValuesResponse
	if strings.Contains(string(bodyBytes), "GetParameterValuesResponse") {
		return acs.handleGetParameterValuesResponse(envelope, response)
	}

	// Check for SetParameterValuesResponse
	if strings.Contains(string(bodyBytes), "SetParameterValuesResponse") {
		return acs.handleSetParameterValuesResponse(envelope, response)
	}

	// Default: send empty response
	return response, nil
}

// handleInform handles CWMP Inform requests
func (acs *AcsServer) handleInform(envelope *SOAPEnvelope, response *SOAPEnvelope, r *http.Request) (*SOAPEnvelope, error) {
	log.Println("Processing Inform request")

	// Parse Inform message
	var inform Inform
	bodyBytes, _ := xml.Marshal(envelope.Body.Content)
	if err := xml.Unmarshal(bodyBytes, &inform); err != nil {
		return nil, fmt.Errorf("error parsing Inform message: %w", err)
	}

	// Create or update session
	deviceId := fmt.Sprintf("%s-%s-%s-%s", 
		inform.DeviceId.Manufacturer,
		inform.DeviceId.OUI,
		inform.DeviceId.ProductClass,
		inform.DeviceId.SerialNumber)

	session := acs.getOrCreateSession(deviceId)
	session.State = SessionStateInform
	session.LastActivity = time.Now()

	// Log device information
	log.Printf("Device connected: %s (Events: %v)", deviceId, inform.Event)

	// Store device parameters in database (implementation needed)
	// acs.storeDeviceParameters(deviceId, inform.ParameterList)

	// Create InformResponse
	informResponse := &InformResponse{
		MaxEnvelopes: 1,
	}

	response.Body.Content = informResponse
	response.Header.NoMoreRequests = true

	return response, nil
}

// handleGetParameterValuesResponse handles response from device
func (acs *AcsServer) handleGetParameterValuesResponse(envelope *SOAPEnvelope, response *SOAPEnvelope) (*SOAPEnvelope, error) {
	log.Println("Processing GetParameterValuesResponse")
	
	// Parse response and store in database
	var getParamResponse GetParameterValuesResponse
	bodyBytes, _ := xml.Marshal(envelope.Body.Content)
	if err := xml.Unmarshal(bodyBytes, &getParamResponse); err != nil {
		return nil, fmt.Errorf("error parsing GetParameterValuesResponse: %w", err)
	}

	log.Printf("Received parameters: %v", getParamResponse.ParameterList)
	
	response.Header.NoMoreRequests = true
	return response, nil
}

// handleSetParameterValuesResponse handles response from device
func (acs *AcsServer) handleSetParameterValuesResponse(envelope *SOAPEnvelope, response *SOAPEnvelope) (*SOAPEnvelope, error) {
	log.Println("Processing SetParameterValuesResponse")
	
	var setParamResponse SetParameterValuesResponse
	bodyBytes, _ := xml.Marshal(envelope.Body.Content)
	if err := xml.Unmarshal(bodyBytes, &setParamResponse); err != nil {
		return nil, fmt.Errorf("error parsing SetParameterValuesResponse: %w", err)
	}

	log.Printf("Set parameter status: %d", setParamResponse.Status)
	
	response.Header.NoMoreRequests = true
	return response, nil
}

// getOrCreateSession gets existing session or creates new one
func (acs *AcsServer) getOrCreateSession(deviceId string) *CwmpSession {
	acs.mutex.Lock()
	defer acs.mutex.Unlock()

	if session, exists := acs.sessions[deviceId]; exists {
		return session
	}

	session := &CwmpSession{
		DeviceId:     deviceId,
		SessionId:    fmt.Sprintf("session-%d", time.Now().Unix()),
		CreatedTime:  time.Now(),
		LastActivity: time.Now(),
		State:        SessionStateNew,
		MaxEnvelopes: 1,
		PendingRPCs:  make([]interface{}, 0),
	}

	acs.sessions[deviceId] = session
	log.Printf("Created new session for device: %s", deviceId)
	
	return session
}

// sendEmptyResponse sends an empty SOAP response
func (acs *AcsServer) sendEmptyResponse(w http.ResponseWriter) {
	response := `<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Header/>
	<soap:Body/>
</soap:Envelope>`
	
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(response))
}

// sendSOAPFault sends a SOAP fault response
func (acs *AcsServer) sendSOAPFault(w http.ResponseWriter, faultCode uint32, faultString string) {
	fault := &SOAPEnvelope{
		SoapNS: "http://schemas.xmlsoap.org/soap/envelope/",
		CwmpNS: "urn:dslforum-org:cwmp-1-2",
		Body: SOAPBody{
			Fault: &SOAPFault{
				FaultCode:   "Client",
				FaultString: faultString,
				Detail: &FaultDetail{
					CWMPFault: &CWMPFault{
						FaultCode:   faultCode,
						FaultString: faultString,
					},
				},
			},
		},
	}

	faultXML, err := xml.MarshalIndent(fault, "", "  ")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(xml.Header))
	w.Write(faultXML)
}

// SendRPC sends an RPC request to a device
func (acs *AcsServer) SendRPC(deviceId string, rpc interface{}) error {
	acs.mutex.RLock()
	session, exists := acs.sessions[deviceId]
	acs.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("no active session for device: %s", deviceId)
	}

	session.mutex.Lock()
	session.PendingRPCs = append(session.PendingRPCs, rpc)
	session.mutex.Unlock()

	log.Printf("Queued RPC for device %s: %T", deviceId, rpc)
	return nil
}

// GetParameterValues requests parameter values from a device
func (acs *AcsServer) GetParameterValues(deviceId string, parameterNames []string) error {
	rpc := &GetParameterValues{
		ParameterNames: parameterNames,
	}
	return acs.SendRPC(deviceId, rpc)
}

// SetParameterValues sets parameter values on a device
func (acs *AcsServer) SetParameterValues(deviceId string, parameters []ParameterValueStruct, parameterKey string) error {
	rpc := &SetParameterValues{
		ParameterList: parameters,
		ParameterKey:  parameterKey,
	}
	return acs.SendRPC(deviceId, rpc)
}

// RebootDevice sends a reboot command to a device
func (acs *AcsServer) RebootDevice(deviceId string, commandKey string) error {
	rpc := &Reboot{
		CommandKey: commandKey,
	}
	return acs.SendRPC(deviceId, rpc)
}