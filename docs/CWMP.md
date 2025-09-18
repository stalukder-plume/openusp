# CWMP/TR-069 Integration

OpenUSP provides comprehensive support for CWMP (CPE WAN Management Protocol) as defined by TR-069 specification, enabling management of legacy TR-069 devices.

## CWMP Overview

CWMP (CPE WAN Management Protocol), also known as TR-069, is a technical specification that defines an application layer protocol for remote management of customer-premises equipment (CPE) connected to an Internet Protocol (IP) network.

### Key Features
- **Legacy Device Support**: Manage existing TR-069 devices
- **Protocol Translation**: Convert between CWMP and USP
- **Auto Configuration Server (ACS)**: Full ACS implementation
- **Firmware Management**: Remote firmware updates
- **Diagnostics**: Remote diagnostic capabilities
- **Event Notifications**: Real-time device events

## Architecture

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   ACS Server    │◄─────►│   HTTP/SOAP     │◄─────►│   CPE Device    │
│   (OpenUSP)     │       │   Transport     │       │   (TR-069)      │
└─────────────────┘       └─────────────────┘       └─────────────────┘
        │                         │                         │
        │                         │                         │
    ┌───▼───┐               ┌─────▼─────┐               ┌───▼───┐
    │ CWMP  │               │ HTTP/HTTPS│               │ CWMP  │
    │Server │               │   SOAP    │               │Client │
    │       │               │   XML     │               │       │
    └───────┘               └───────────┘               └───────┘
```

### Components

#### ACS (Auto Configuration Server)
- Manages CWMP devices
- Provides configuration management
- Handles firmware updates
- Processes diagnostic requests
- Manages device provisioning

#### CPE (Customer Premises Equipment)
- TR-069 enabled devices
- Connects to ACS for management
- Reports status and events
- Executes remote commands

## CWMP Protocol

### Connection Establishment
```
CPE                     ACS
 │                       │
 │   HTTP POST (Inform)  │
 ├──────────────────────►│
 │                       │
 │   HTTP 200 OK         │
 │◄──────────────────────┤
 │   (with SetParameterValues)
 │                       │
 │   HTTP POST (Response)│
 ├──────────────────────►│
 │                       │
 │   HTTP 204 No Content │
 │◄──────────────────────┤
```

### SOAP Message Structure
```xml
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope 
    xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"
    xmlns:cwmp="urn:dslforum-org:cwmp-1-2">
    <soap:Header>
        <cwmp:ID soap:mustUnderstand="1">12345</cwmp:ID>
    </soap:Header>
    <soap:Body>
        <cwmp:Inform>
            <DeviceId>
                <Manufacturer>Acme Corp</Manufacturer>
                <OUI>ABCDEF</OUI>
                <ProductClass>Router</ProductClass>
                <SerialNumber>123456789</SerialNumber>
            </DeviceId>
            <Event>
                <EventStruct>
                    <EventCode>1 BOOT</EventCode>
                    <CommandKey></CommandKey>
                </EventStruct>
            </Event>
            <MaxEnvelopes>1</MaxEnvelopes>
            <CurrentTime>2023-12-01T10:00:00Z</CurrentTime>
            <RetryCount>0</RetryCount>
            <ParameterList>
                <ParameterValueStruct>
                    <Name>Device.DeviceInfo.SoftwareVersion</Name>
                    <Value>1.2.3</Value>
                </ParameterValueStruct>
            </ParameterList>
        </cwmp:Inform>
    </soap:Body>
</soap:Envelope>
```

## Configuration

### ACS Server Configuration
```yaml
# CWMP ACS Configuration
cwmp:
  enabled: true
  port: 7547
  tls:
    enabled: true
    cert_file: "/certs/acs.crt"
    key_file: "/certs/acs.key"
  
  # Authentication
  auth:
    enabled: true
    username: "acs-user"
    password: "acs-password"
    
  # Connection settings
  connection:
    timeout: "30s"
    max_envelope_size: 65536
    soap_timeout: "60s"
    
  # Database settings  
  database:
    driver: "mongodb"
    connection_string: "mongodb://localhost:27017/cwmp"
    
  # Provisioning
  provisioning:
    auto_create_devices: true
    default_config_url: "http://acs.example.com/config/"
    firmware_url: "http://acs.example.com/firmware/"
```

### Device Registration
```go
type CWMPDevice struct {
    ID             string    `bson:"_id"`
    OUI            string    `bson:"oui"`
    ProductClass   string    `bson:"product_class"`
    SerialNumber   string    `bson:"serial_number"`
    Manufacturer   string    `bson:"manufacturer"`
    SoftwareVersion string   `bson:"software_version"`
    HardwareVersion string   `bson:"hardware_version"`
    ConnectionURL  string    `bson:"connection_url"`
    Username       string    `bson:"username"`
    Password       string    `bson:"password"`
    Status         string    `bson:"status"`
    LastInform     time.Time `bson:"last_inform"`
    ParameterKey   string    `bson:"parameter_key"`
    Created        time.Time `bson:"created"`
    Updated        time.Time `bson:"updated"`
}
```

## CWMP Methods

### Inform
Device reports status and events to ACS.

```go
func (s *CWMPServer) HandleInform(w http.ResponseWriter, r *http.Request) {
    // Parse SOAP message
    inform, err := s.parseInformMessage(r)
    if err != nil {
        s.sendFault(w, CWMPFaultInvalidArgs, "Invalid inform message")
        return
    }
    
    // Update device information
    device := &CWMPDevice{
        OUI:            inform.DeviceId.OUI,
        ProductClass:   inform.DeviceId.ProductClass,
        SerialNumber:   inform.DeviceId.SerialNumber,
        Manufacturer:   inform.DeviceId.Manufacturer,
        LastInform:     time.Now(),
        Status:         "online",
    }
    
    err = s.db.UpdateDevice(device)
    if err != nil {
        log.Errorf("Failed to update device: %v", err)
    }
    
    // Check for pending operations
    ops, err := s.db.GetPendingOperations(device.GetUniqueID())
    if err == nil && len(ops) > 0 {
        // Send first pending operation
        s.sendOperation(w, ops[0])
        return
    }
    
    // Send empty response
    s.sendEmptyResponse(w)
}
```

### GetParameterValues
Retrieve parameter values from device.

```go
type GetParameterValues struct {
    ParameterNames []string
}

func (s *CWMPServer) GetParameterValues(deviceID string, paramNames []string) error {
    op := &CWMPOperation{
        ID:         generateOperationID(),
        DeviceID:   deviceID,
        Type:       "GetParameterValues",
        Parameters: paramNames,
        Status:     "pending",
        Created:    time.Now(),
    }
    
    return s.db.SaveOperation(op)
}

// SOAP message generation
func (s *CWMPServer) generateGetParameterValuesSOAP(paramNames []string) string {
    return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" 
               xmlns:cwmp="urn:dslforum-org:cwmp-1-2">
    <soap:Header>
        <cwmp:ID soap:mustUnderstand="1">%s</cwmp:ID>
    </soap:Header>
    <soap:Body>
        <cwmp:GetParameterValues>
            <ParameterNames>
                %s
            </ParameterNames>
        </cwmp:GetParameterValues>
    </soap:Body>
</soap:Envelope>`, generateMessageID(), strings.Join(paramNames, "\n"))
}
```

### SetParameterValues
Modify parameter values on device.

```go
type SetParameterValues struct {
    ParameterList []ParameterValueStruct
    ParameterKey  string
}

type ParameterValueStruct struct {
    Name  string
    Value interface{}
    Type  string
}

func (s *CWMPServer) SetParameterValues(deviceID string, params []ParameterValueStruct, key string) error {
    op := &CWMPOperation{
        ID:           generateOperationID(),
        DeviceID:     deviceID,
        Type:         "SetParameterValues",
        Parameters:   params,
        ParameterKey: key,
        Status:       "pending",
        Created:      time.Now(),
    }
    
    return s.db.SaveOperation(op)
}

func (s *CWMPServer) generateSetParameterValuesSOAP(params []ParameterValueStruct, key string) string {
    paramList := ""
    for _, param := range params {
        paramList += fmt.Sprintf(`
            <ParameterValueStruct>
                <Name>%s</Name>
                <Value xsi:type="%s">%v</Value>
            </ParameterValueStruct>`, param.Name, param.Type, param.Value)
    }
    
    return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" 
               xmlns:cwmp="urn:dslforum-org:cwmp-1-2"
               xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
    <soap:Header>
        <cwmp:ID soap:mustUnderstand="1">%s</cwmp:ID>
    </soap:Header>
    <soap:Body>
        <cwmp:SetParameterValues>
            <ParameterList>%s</ParameterList>
            <ParameterKey>%s</ParameterKey>
        </cwmp:SetParameterValues>
    </soap:Body>
</soap:Envelope>`, generateMessageID(), paramList, key)
}
```

### GetParameterNames
Discover available parameters.

```go
func (s *CWMPServer) GetParameterNames(deviceID, path string, nextLevel bool) error {
    op := &CWMPOperation{
        ID:       generateOperationID(),
        DeviceID: deviceID,
        Type:     "GetParameterNames",
        Parameters: map[string]interface{}{
            "ParameterPath": path,
            "NextLevel":     nextLevel,
        },
        Status:  "pending",
        Created: time.Now(),
    }
    
    return s.db.SaveOperation(op)
}
```

### AddObject/DeleteObject
Manage object instances.

```go
func (s *CWMPServer) AddObject(deviceID, objectName, parameterKey string) error {
    op := &CWMPOperation{
        ID:           generateOperationID(),
        DeviceID:     deviceID,
        Type:         "AddObject",
        Parameters:   map[string]interface{}{"ObjectName": objectName},
        ParameterKey: parameterKey,
        Status:       "pending",
        Created:      time.Now(),
    }
    
    return s.db.SaveOperation(op)
}

func (s *CWMPServer) DeleteObject(deviceID, objectName, parameterKey string) error {
    op := &CWMPOperation{
        ID:           generateOperationID(),
        DeviceID:     deviceID,
        Type:         "DeleteObject",
        Parameters:   map[string]interface{}{"ObjectName": objectName},
        ParameterKey: parameterKey,
        Status:       "pending",
        Created:      time.Now(),
    }
    
    return s.db.SaveOperation(op)
}
```

## Data Models

### TR-069 Data Model Structure
```
Device.
├── DeviceInfo.
│   ├── Manufacturer (string, R)
│   ├── ManufacturerOUI (string, R)
│   ├── ModelName (string, R)
│   ├── Description (string, R)
│   ├── ProductClass (string, R)
│   ├── SerialNumber (string, R)
│   ├── HardwareVersion (string, R)
│   ├── SoftwareVersion (string, R)
│   ├── ProvisioningCode (string, RW)
│   └── UpTime (unsignedInt, R)
├── ManagementServer.
│   ├── EnableCWMP (boolean, RW)
│   ├── URL (string, RW)
│   ├── Username (string, RW)
│   ├── Password (string, RW)
│   ├── PeriodicInformEnable (boolean, RW)
│   └── PeriodicInformInterval (unsignedInt, RW)
├── LAN.
│   ├── IPAddress (string, RW)
│   ├── SubnetMask (string, RW)
│   ├── DHCPServerEnable (boolean, RW)
│   └── DHCPServerConfigurable (boolean, R)
└── WiFi.
    ├── RadioNumberOfEntries (unsignedInt, R)
    ├── SSIDNumberOfEntries (unsignedInt, R)
    ├── Radio.{i}.
    │   ├── Enable (boolean, RW)
    │   ├── Status (string, R)
    │   ├── Channel (unsignedInt, RW)
    │   └── AutoChannelEnable (boolean, RW)
    └── SSID.{i}.
        ├── Enable (boolean, RW)
        ├── SSID (string, RW)
        ├── BeaconType (string, RW)
        └── BasicAuthenticationMode (string, RW)
```

### Parameter Storage
```go
type CWMPParameter struct {
    DeviceID    string      `bson:"device_id"`
    Path        string      `bson:"path"`
    Value       interface{} `bson:"value"`
    Type        string      `bson:"type"`
    Writable    bool        `bson:"writable"`
    LastUpdated time.Time   `bson:"last_updated"`
}

func (db *CWMPDatabase) SaveParameter(param *CWMPParameter) error {
    collection := db.client.Database("cwmp").Collection("parameters")
    
    filter := bson.M{
        "device_id": param.DeviceID,
        "path":      param.Path,
    }
    
    update := bson.M{
        "$set": bson.M{
            "value":        param.Value,
            "type":         param.Type,
            "writable":     param.Writable,
            "last_updated": time.Now(),
        },
    }
    
    _, err := collection.UpdateOne(context.Background(), filter, update, options.Update().SetUpsert(true))
    return err
}
```

## Session Management

### Session State
```go
type CWMPSession struct {
    ID           string            `bson:"_id"`
    DeviceID     string            `bson:"device_id"`
    Status       string            `bson:"status"` // active, completed, error
    Operations   []CWMPOperation   `bson:"operations"`
    CurrentOp    int               `bson:"current_op"`
    Started      time.Time         `bson:"started"`
    LastActivity time.Time         `bson:"last_activity"`
    Cookies      map[string]string `bson:"cookies"`
    MaxEnvelopes int               `bson:"max_envelopes"`
}

func (s *CWMPServer) CreateSession(deviceID string, cookies map[string]string) *CWMPSession {
    session := &CWMPSession{
        ID:           generateSessionID(),
        DeviceID:     deviceID,
        Status:       "active",
        Operations:   []CWMPOperation{},
        CurrentOp:    0,
        Started:      time.Now(),
        LastActivity: time.Now(),
        Cookies:      cookies,
        MaxEnvelopes: 1,
    }
    
    s.sessions[session.ID] = session
    return session
}
```

### HTTP Cookie Management
```go
func (s *CWMPServer) setCWMPCookies(w http.ResponseWriter, sessionID string) {
    cookie := &http.Cookie{
        Name:     "CWMPSESSIONID",
        Value:    sessionID,
        Path:     "/",
        HttpOnly: true,
        Secure:   s.config.TLS.Enabled,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   3600, // 1 hour
    }
    http.SetCookie(w, cookie)
}

func (s *CWMPServer) getCWMPSession(r *http.Request) (*CWMPSession, error) {
    cookie, err := r.Cookie("CWMPSESSIONID")
    if err != nil {
        return nil, err
    }
    
    session, exists := s.sessions[cookie.Value]
    if !exists {
        return nil, fmt.Errorf("session not found")
    }
    
    return session, nil
}
```

## Event Handling

### Event Types
```go
const (
    EventBootstrap      = "0 BOOTSTRAP"
    EventBoot          = "1 BOOT" 
    EventPeriodic      = "2 PERIODIC"
    EventScheduled     = "3 SCHEDULED"
    EventValueChange   = "4 VALUE CHANGE"
    EventKicked        = "5 KICKED"
    EventConnectionRequest = "6 CONNECTION REQUEST"
    EventTransferComplete  = "7 TRANSFER COMPLETE"
    EventDiagnosticsComplete = "8 DIAGNOSTICS COMPLETE"
    EventRequestDownload   = "9 REQUEST DOWNLOAD"
    EventAutonomousTransferComplete = "10 AUTONOMOUS TRANSFER COMPLETE"
    EventDUStateChangeComplete = "11 DU STATE CHANGE COMPLETE"
    EventAutonomousDUStateChangeComplete = "12 AUTONOMOUS DU STATE CHANGE COMPLETE"
)
```

### Event Processing
```go
func (s *CWMPServer) processEvent(deviceID string, event EventStruct) error {
    log.Infof("Processing event %s for device %s", event.EventCode, deviceID)
    
    switch event.EventCode {
    case EventBoot:
        return s.handleBootEvent(deviceID, event)
    case EventPeriodic:
        return s.handlePeriodicEvent(deviceID, event)
    case EventValueChange:
        return s.handleValueChangeEvent(deviceID, event)
    case EventTransferComplete:
        return s.handleTransferCompleteEvent(deviceID, event)
    default:
        log.Warnf("Unhandled event type: %s", event.EventCode)
    }
    
    return nil
}

func (s *CWMPServer) handleBootEvent(deviceID string, event EventStruct) error {
    // Update device status
    err := s.db.UpdateDeviceStatus(deviceID, "online")
    if err != nil {
        return err
    }
    
    // Trigger auto-provisioning if enabled
    if s.config.Provisioning.AutoProvision {
        return s.triggerProvisioning(deviceID)
    }
    
    return nil
}
```

## Firmware Management

### Download
```go
type Download struct {
    CommandKey   string
    FileType     string
    URL          string
    Username     string
    Password     string
    FileSize     uint32
    TargetFileName string
    DelaySeconds uint32
    SuccessURL   string
    FailureURL   string
}

func (s *CWMPServer) ScheduleDownload(deviceID string, download *Download) error {
    op := &CWMPOperation{
        ID:         generateOperationID(),
        DeviceID:   deviceID,
        Type:       "Download",
        Parameters: download,
        Status:     "pending",
        Created:    time.Now(),
    }
    
    return s.db.SaveOperation(op)
}

func (s *CWMPServer) generateDownloadSOAP(download *Download) string {
    return fmt.Sprintf(`
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" 
               xmlns:cwmp="urn:dslforum-org:cwmp-1-2">
    <soap:Header>
        <cwmp:ID soap:mustUnderstand="1">%s</cwmp:ID>
    </soap:Header>
    <soap:Body>
        <cwmp:Download>
            <CommandKey>%s</CommandKey>
            <FileType>%s</FileType>
            <URL>%s</URL>
            <Username>%s</Username>
            <Password>%s</Password>
            <FileSize>%d</FileSize>
            <TargetFileName>%s</TargetFileName>
            <DelaySeconds>%d</DelaySeconds>
            <SuccessURL>%s</SuccessURL>
            <FailureURL>%s</FailureURL>
        </cwmp:Download>
    </soap:Body>
</soap:Envelope>`, 
    generateMessageID(),
    download.CommandKey,
    download.FileType,
    download.URL,
    download.Username,
    download.Password,
    download.FileSize,
    download.TargetFileName,
    download.DelaySeconds,
    download.SuccessURL,
    download.FailureURL)
}
```

## Protocol Translation

### CWMP to USP Translation
```go
type ProtocolTranslator struct {
    cwmpServer *CWMPServer
    uspController *USPController
}

func (t *ProtocolTranslator) TranslateCWMPToUSP(deviceID string, cwmpOp *CWMPOperation) error {
    switch cwmpOp.Type {
    case "GetParameterValues":
        paths := cwmpOp.Parameters.([]string)
        return t.uspController.GetParameters(deviceID, paths)
        
    case "SetParameterValues":
        params := cwmpOp.Parameters.([]ParameterValueStruct)
        uspParams := make(map[string]interface{})
        for _, p := range params {
            uspParams[p.Name] = p.Value
        }
        return t.uspController.SetParameters(deviceID, uspParams)
        
    case "AddObject":
        objPath := cwmpOp.Parameters.(map[string]interface{})["ObjectName"].(string)
        return t.uspController.AddObject(deviceID, objPath, nil)
        
    case "DeleteObject":
        objPath := cwmpOp.Parameters.(map[string]interface{})["ObjectName"].(string)
        return t.uspController.DeleteObject(deviceID, objPath)
    }
    
    return fmt.Errorf("unsupported operation: %s", cwmpOp.Type)
}
```

## Testing

### Unit Tests
```go
func TestCWMPServer_HandleInform(t *testing.T) {
    server := NewTestCWMPServer()
    
    informXML := `
<?xml version="1.0" encoding="UTF-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/" 
               xmlns:cwmp="urn:dslforum-org:cwmp-1-2">
    <soap:Body>
        <cwmp:Inform>
            <DeviceId>
                <Manufacturer>Test Corp</Manufacturer>
                <OUI>ABCDEF</OUI>
                <ProductClass>Router</ProductClass>
                <SerialNumber>123456789</SerialNumber>
            </DeviceId>
            <Event>
                <EventStruct>
                    <EventCode>1 BOOT</EventCode>
                    <CommandKey></CommandKey>
                </EventStruct>
            </Event>
            <MaxEnvelopes>1</MaxEnvelopes>
            <CurrentTime>2023-12-01T10:00:00Z</CurrentTime>
            <RetryCount>0</RetryCount>
        </cwmp:Inform>
    </soap:Body>
</soap:Envelope>`
    
    req := httptest.NewRequest("POST", "/", strings.NewReader(informXML))
    w := httptest.NewRecorder()
    
    server.HandleInform(w, req)
    
    assert.Equal(t, 200, w.Code)
    
    // Verify device was created/updated
    device, err := server.db.GetDevice("ABCDEF-Router-123456789")
    assert.NoError(t, err)
    assert.Equal(t, "Test Corp", device.Manufacturer)
    assert.Equal(t, "online", device.Status)
}
```

### Integration Tests
```go
func TestCWMPIntegration_EndToEnd(t *testing.T) {
    // Start CWMP server
    server := startTestCWMPServer()
    defer server.Stop()
    
    // Simulate CPE device
    cpe := NewTestCPEDevice()
    
    // Test inform -> response cycle
    err := cpe.SendInform(server.URL, "1 BOOT")
    assert.NoError(t, err)
    
    // Test parameter operations
    err = server.SetParameter(cpe.DeviceID, "Device.WiFi.Radio.1.Enable", "true")
    assert.NoError(t, err)
    
    // Wait for next inform with parameter change
    response, err := cpe.WaitForNextInform()
    assert.NoError(t, err)
    assert.Contains(t, response.ParameterList, "Device.WiFi.Radio.1.Enable")
}
```

This CWMP/TR-069 integration guide provides comprehensive information for implementing and using CWMP protocol features in the OpenUSP platform.