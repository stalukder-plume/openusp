# USP (User Services Platform) Integration

OpenUSP provides comprehensive support for the USP protocol as defined by the Broadband Forum TR-369 specification.

## USP Overview

USP (User Services Platform) is a standardized protocol for managing, monitoring, upgrading, and controlling connected devices. It provides a bidirectional communication mechanism between Controllers and Agents.

### Key Features
- **Standardized Protocol**: Based on TR-369 specification
- **Multi-Transport Support**: WebSocket, MQTT, STOMP, CoAP
- **Security**: Built-in authentication and encryption
- **Extensibility**: Support for vendor-specific data models
- **Real-time Communication**: Event notifications and subscriptions

## Architecture

```
┌─────────────────┐       ┌─────────────────┐       ┌─────────────────┐
│   Controller    │◄─────►│      MTP        │◄─────►│     Agent       │
│   (OpenUSP)     │       │   (Transport)   │       │   (Device)      │
└─────────────────┘       └─────────────────┘       └─────────────────┘
        │                         │                         │
        │                         │                         │
    ┌───▼───┐               ┌─────▼─────┐               ┌───▼───┐
    │ USP   │               │WebSocket  │               │ USP   │
    │Message│               │   MQTT    │               │Message│
    │       │               │   STOMP   │               │       │
    └───────┘               │   CoAP    │               └───────┘
                            └───────────┘
```

### Components

#### Controller
- Manages and controls USP Agents
- Sends commands and receives responses
- Handles subscriptions and notifications
- Maintains device state and configuration

#### Agent
- Resides on managed devices
- Executes commands from Controllers
- Sends notifications and events
- Exposes device data model

#### Message Transport Protocol (MTP)
- Provides reliable message delivery
- Supports multiple transport mechanisms
- Handles connection management and security

## USP Messages

### Message Types

#### Request Messages
- **Get**: Retrieve parameter values
- **Set**: Modify parameter values
- **Add**: Create object instances
- **Delete**: Remove object instances
- **Operate**: Execute RPC methods
- **GetSupportedDM**: Get supported data model
- **GetInstances**: Get object instances
- **GetSupportedProtocol**: Get protocol information

#### Response Messages
- **GetResp**: Get operation response
- **SetResp**: Set operation response
- **AddResp**: Add operation response
- **DeleteResp**: Delete operation response
- **OperateResp**: Operate operation response
- **Error**: Error response

#### Notification Messages
- **Notify**: Event notifications
- **OnBoardRequest**: Device onboarding

### Message Structure

```protobuf
message Msg {
    Header header = 1;
    Body body = 2;
}

message Header {
    string msg_id = 1;
    MsgType msg_type = 2;
    // Additional header fields
}

message Body {
    oneof msg_body {
        Request request = 1;
        Response response = 2;
        Error error = 3;
        Notify notify = 4;
    }
}
```

## Transport Protocols

### WebSocket
```yaml
# WebSocket MTP Configuration
websocket:
  enabled: true
  port: 8080
  path: "/usp"
  tls:
    enabled: true
    cert_file: "/certs/server.crt"
    key_file: "/certs/server.key"
  ping_interval: 30s
  pong_timeout: 10s
```

**Connection Example:**
```javascript
const ws = new WebSocket('wss://controller.example.com:8080/usp');
ws.binaryType = 'arraybuffer';

ws.onopen = function() {
    console.log('USP WebSocket connection established');
};

ws.onmessage = function(event) {
    const uspMessage = parseUSPMessage(event.data);
    handleUSPMessage(uspMessage);
};
```

### MQTT
```yaml
# MQTT MTP Configuration
mqtt:
  enabled: true
  broker: "mqtt://broker.example.com:1883"
  client_id: "openusp-controller"
  username: "controller"
  password: "secret"
  tls:
    enabled: true
    ca_cert: "/certs/ca.crt"
  topics:
    request: "usp/controller/{agent_id}/request"
    response: "usp/controller/{agent_id}/response"
    notify: "usp/controller/+/notify"
```

**Topic Structure:**
- Request: `usp/controller/{agent_id}/request`
- Response: `usp/controller/{agent_id}/response`  
- Notify: `usp/controller/{agent_id}/notify`

### STOMP
```yaml
# STOMP MTP Configuration
stomp:
  enabled: true
  host: "stomp.example.com"
  port: 61613
  username: "controller"
  password: "secret"
  vhost: "/usp"
  destinations:
    request: "/queue/usp.request"
    response: "/queue/usp.response"
    notify: "/topic/usp.notify"
```

### CoAP
```yaml
# CoAP MTP Configuration
coap:
  enabled: true
  address: "0.0.0.0:5683"
  dtls:
    enabled: true
    psk_identity: "controller"
    psk_key: "secret-key"
  resources:
    usp: "/usp"
```

## Data Models

### Supported Data Models
- **Device:2.14** - Core device data model
- **WiFi** - WiFi configuration and statistics
- **Ethernet** - Ethernet interface management
- **IP** - IP configuration and routing
- **Firewall** - Security and filtering rules
- **NAT** - Network Address Translation
- **DHCP** - DHCP client and server
- **DNS** - DNS configuration
- **Time** - Time synchronization

### Data Model Structure
```
Device.
├── DeviceInfo.
│   ├── Manufacturer (string, R)
│   ├── ModelName (string, R)
│   ├── SoftwareVersion (string, R)
│   └── SerialNumber (string, R)
├── WiFi.
│   ├── RadioNumberOfEntries (int, R)
│   ├── SSIDNumberOfEntries (int, R)
│   ├── Radio.{i}.
│   │   ├── Enable (boolean, RW)
│   │   ├── Status (string, R)
│   │   ├── Channel (int, RW)
│   │   └── OperatingFrequencyBand (string, RW)
│   └── SSID.{i}.
│       ├── Enable (boolean, RW)
│       ├── SSID (string, RW)
│       ├── Passphrase (string, RW)
│       └── Status (string, R)
└── Ethernet.
    ├── InterfaceNumberOfEntries (int, R)
    └── Interface.{i}.
        ├── Enable (boolean, RW)
        ├── Status (string, R)
        ├── MACAddress (string, R)
        └── MaxBitRate (int, RW)
```

## Operations

### Get Operation
Retrieve parameter values from devices.

```go
// Get single parameter
func (c *Controller) GetParameter(agentID, path string) (*Parameter, error) {
    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_GET,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Get{
                        Get: &usp.Get{
                            ParamPaths: []string{path},
                        },
                    },
                },
            },
        },
    }
    return c.sendRequest(agentID, msg)
}

// Get multiple parameters
func (c *Controller) GetParameters(agentID string, paths []string) ([]*Parameter, error) {
    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_GET,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Get{
                        Get: &usp.Get{
                            ParamPaths: paths,
                        },
                    },
                },
            },
        },
    }
    return c.sendRequest(agentID, msg)
}
```

### Set Operation
Modify parameter values on devices.

```go
func (c *Controller) SetParameters(agentID string, params map[string]interface{}) error {
    updateObjs := make([]*usp.Set_UpdateObject, 0, len(params))
    
    for path, value := range params {
        updateObjs = append(updateObjs, &usp.Set_UpdateObject{
            ObjPath: path,
            ParamSettings: []*usp.Set_UpdateParamSetting{
                {
                    Param: path,
                    Value: fmt.Sprintf("%v", value),
                },
            },
        })
    }

    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_SET,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Set{
                        Set: &usp.Set{
                            AllowPartialFalse: false,
                            UpdateObjs:        updateObjs,
                        },
                    },
                },
            },
        },
    }
    
    return c.sendRequest(agentID, msg)
}
```

### Add Operation
Create new object instances.

```go
func (c *Controller) AddObject(agentID, objPath string, params map[string]interface{}) (string, error) {
    paramSettings := make([]*usp.Add_CreateParamSetting, 0, len(params))
    
    for param, value := range params {
        paramSettings = append(paramSettings, &usp.Add_CreateParamSetting{
            Param: param,
            Value: fmt.Sprintf("%v", value),
        })
    }

    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_ADD,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Add{
                        Add: &usp.Add{
                            AllowPartialFalse: false,
                            CreateObjs: []*usp.Add_CreateObject{
                                {
                                    ObjPath:       objPath,
                                    ParamSettings: paramSettings,
                                },
                            },
                        },
                    },
                },
            },
        },
    }
    
    return c.sendRequest(agentID, msg)
}
```

### Delete Operation
Remove object instances.

```go
func (c *Controller) DeleteObject(agentID, objPath string) error {
    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_DELETE,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Delete{
                        Delete: &usp.Delete{
                            AllowPartialFalse: false,
                            ObjPaths:          []string{objPath},
                        },
                    },
                },
            },
        },
    }
    
    return c.sendRequest(agentID, msg)
}
```

## Event Handling

### Subscriptions
```go
type EventSubscription struct {
    ID        string
    AgentID   string
    EventType string
    Path      string
    Handler   func(*Event)
}

func (c *Controller) Subscribe(agentID, eventType, path string, handler func(*Event)) (*EventSubscription, error) {
    subscription := &EventSubscription{
        ID:        generateSubscriptionID(),
        AgentID:   agentID,
        EventType: eventType,
        Path:      path,
        Handler:   handler,
    }
    
    // Send subscription request to agent
    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_NOTIFY,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Notify{
                        Notify: &usp.Notify{
                            SubscriptionId: subscription.ID,
                            Send: &usp.Notify_Event{
                                Event: &usp.Event{
                                    ObjPath:   path,
                                    EventName: eventType,
                                },
                            },
                        },
                    },
                },
            },
        },
    }
    
    err := c.sendRequest(agentID, msg)
    if err != nil {
        return nil, err
    }
    
    c.subscriptions[subscription.ID] = subscription
    return subscription, nil
}
```

### Event Types
- **ValueChange**: Parameter value changed
- **ObjectCreation**: Object instance created
- **ObjectDeletion**: Object instance deleted
- **OperationComplete**: Asynchronous operation completed
- **OnBoardRequest**: Device requesting onboarding
- **Boot**: Device boot notification
- **Periodic**: Periodic event notification

## Security

### Authentication
```yaml
# Certificate-based authentication
authentication:
  type: "certificate"
  ca_cert: "/certs/ca.crt"
  server_cert: "/certs/server.crt"
  server_key: "/certs/server.key"
  client_verification: true

# Shared secret authentication
authentication:
  type: "shared_secret"
  algorithm: "HMAC-SHA256"
  secret: "shared-secret-key"
```

### Encryption
All USP messages are encrypted using TLS/DTLS transport encryption:
- **WebSocket**: TLS 1.3
- **MQTT**: TLS 1.3 with client certificates
- **STOMP**: TLS 1.3
- **CoAP**: DTLS 1.2

### Message Integrity
USP messages include integrity protection through:
- Transport-level encryption
- Message-level signatures (optional)
- Sequence number validation

## Error Handling

### Error Types
```go
const (
    USPErrorMessageFormat    = 7000
    USPErrorMessageStructure = 7001
    USPErrorUnsupportedParam = 7002
    USPErrorInvalidArguments = 7003
    USPErrorResourcesExceeded = 7004
    USPErrorPermissionDenied = 7005
    USPErrorInvalidConfig    = 7006
    USPErrorInvalidPath      = 7007
    USPErrorParameterReadOnly = 7008
    USPErrorValueConflict    = 7009
    USPErrorOperationFailure = 7010
)
```

### Error Response
```go
func (c *Controller) handleError(msgID string, errorCode uint32, errorMessage string) *usp.Msg {
    return &usp.Msg{
        Header: &usp.Header{
            MsgId:   msgID,
            MsgType: usp.Header_ERROR,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Error{
                Error: &usp.Error{
                    ErrCode: errorCode,
                    ErrMsg:  errorMessage,
                },
            },
        },
    }
}
```

## Performance Optimization

### Connection Pooling
```go
type ConnectionPool struct {
    connections map[string]*Connection
    mu          sync.RWMutex
    maxConn     int
}

func (p *ConnectionPool) GetConnection(agentID string) (*Connection, error) {
    p.mu.RLock()
    conn, exists := p.connections[agentID]
    p.mu.RUnlock()
    
    if exists && conn.IsHealthy() {
        return conn, nil
    }
    
    return p.createConnection(agentID)
}
```

### Message Batching
```go
func (c *Controller) BatchOperations(agentID string, operations []Operation) error {
    msg := &usp.Msg{
        Header: &usp.Header{
            MsgId:   generateMsgID(),
            MsgType: usp.Header_GET,
        },
        Body: &usp.Body{
            MsgBody: &usp.Body_Request{
                Request: &usp.Request{
                    ReqType: &usp.Request_Get{
                        Get: &usp.Get{
                            ParamPaths: extractPaths(operations),
                        },
                    },
                },
            },
        },
    }
    
    return c.sendRequest(agentID, msg)
}
```

## Testing

### Unit Tests
```go
func TestUSPController_GetParameter(t *testing.T) {
    controller := NewTestController()
    agent := NewTestAgent()
    
    // Test successful get
    param, err := controller.GetParameter(agent.ID, "Device.DeviceInfo.Manufacturer")
    assert.NoError(t, err)
    assert.Equal(t, "Test Manufacturer", param.Value)
    
    // Test invalid path
    _, err = controller.GetParameter(agent.ID, "Invalid.Path")
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "invalid path")
}
```

### Integration Tests
```go
func TestUSPIntegration_EndToEnd(t *testing.T) {
    // Start test environment
    controller := startTestController()
    agent := startTestAgent()
    defer cleanup(controller, agent)
    
    // Test device onboarding
    err := controller.OnboardDevice(agent.EndpointID)
    assert.NoError(t, err)
    
    // Test parameter operations
    err = controller.SetParameter(agent.ID, "Device.WiFi.Radio.1.Enable", true)
    assert.NoError(t, err)
    
    param, err := controller.GetParameter(agent.ID, "Device.WiFi.Radio.1.Enable")
    assert.NoError(t, err)
    assert.True(t, param.Value.(bool))
}
```

This USP integration guide provides comprehensive information for implementing and using USP protocol features in the OpenUSP platform.
