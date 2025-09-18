# Security Guide

This document outlines security considerations, best practices, and implementation details for the OpenUSP platform.

## Security Architecture

### Defense in Depth
OpenUSP implements multiple layers of security controls:

```
┌─────────────────────────────────────────────────────────────────┐
│                       Network Security                         │
├─────────────────────────────────────────────────────────────────┤
│  Firewall  │  DDoS Protection  │  VPN  │  Network Segmentation │
└─────────────────┬───────────────────────┬───────────────────────┘
                  │                       │
┌─────────────────▼───────────────────────▼───────────────────────┐
│                    Transport Security                          │
├─────────────────────────────────────────────────────────────────┤
│  TLS 1.3  │  Certificate Management  │  Perfect Forward Secrecy │
└─────────────────┬───────────────────────┬───────────────────────┘
                  │                       │
┌─────────────────▼───────────────────────▼───────────────────────┐
│                 Application Security                           │
├─────────────────────────────────────────────────────────────────┤
│  OAuth2/OIDC  │  RBAC  │  Input Validation  │  Rate Limiting   │
└─────────────────┬───────────────────────┬───────────────────────┘
                  │                       │
┌─────────────────▼───────────────────────▼───────────────────────┐
│                    Data Security                               │
├─────────────────────────────────────────────────────────────────┤
│  Encryption at Rest  │  Data Masking  │  Audit Logging         │
└─────────────────────────────────────────────────────────────────┘
```

## Authentication & Authorization

### Multi-Factor Authentication (MFA)
```yaml
# MFA Configuration
authentication:
  mfa:
    enabled: true
    providers:
      - totp          # Time-based One-Time Password
      - sms           # SMS verification
      - email         # Email verification
    backup_codes:
      enabled: true
      count: 10
    grace_period: 300  # 5 minutes
```

```go
type MFAConfig struct {
    Enabled       bool     `yaml:"enabled"`
    Providers     []string `yaml:"providers"`
    BackupCodes   bool     `yaml:"backup_codes"`
    GracePeriod   int      `yaml:"grace_period"`
    RequiredRoles []string `yaml:"required_roles"`
}

func (auth *AuthService) EnableMFA(userID string, provider string) (*MFASetup, error) {
    user, err := auth.db.GetUser(userID)
    if err != nil {
        return nil, err
    }
    
    switch provider {
    case "totp":
        return auth.setupTOTP(user)
    case "sms":
        return auth.setupSMS(user)
    default:
        return nil, fmt.Errorf("unsupported MFA provider: %s", provider)
    }
}

func (auth *AuthService) setupTOTP(user *User) (*MFASetup, error) {
    secret := make([]byte, 20)
    if _, err := rand.Read(secret); err != nil {
        return nil, err
    }
    
    secretBase32 := base32.StdEncoding.EncodeToString(secret)
    
    // Generate QR code URL
    qrURL := fmt.Sprintf("otpauth://totp/OpenUSP:%s?secret=%s&issuer=OpenUSP",
        user.Email, secretBase32)
    
    // Store encrypted secret
    encryptedSecret, err := auth.encrypt(secretBase32)
    if err != nil {
        return nil, err
    }
    
    user.MFASecret = encryptedSecret
    user.MFAEnabled = false // Will be enabled after verification
    
    err = auth.db.UpdateUser(user)
    if err != nil {
        return nil, err
    }
    
    return &MFASetup{
        Secret: secretBase32,
        QRCode: qrURL,
    }, nil
}
```

### OAuth2 Integration
```yaml
oauth2:
  providers:
    google:
      client_id: "google-client-id"
      client_secret: "google-client-secret"
      redirect_uri: "https://openusp.com/auth/google/callback"
      scopes: ["openid", "email", "profile"]
      
    azure:
      client_id: "azure-client-id"
      client_secret: "azure-client-secret"
      tenant_id: "azure-tenant-id"
      redirect_uri: "https://openusp.com/auth/azure/callback"
      
    okta:
      domain: "company.okta.com"
      client_id: "okta-client-id"
      client_secret: "okta-client-secret"
      redirect_uri: "https://openusp.com/auth/okta/callback"
```

### Role-Based Access Control (RBAC)
```go
type Role struct {
    Name        string       `json:"name"`
    Description string       `json:"description"`
    Permissions []Permission `json:"permissions"`
    IsSystem    bool         `json:"is_system"`
}

type Permission struct {
    Resource string   `json:"resource"`  // device, user, system, api
    Actions  []string `json:"actions"`   // create, read, update, delete
    Scope    string   `json:"scope"`     // global, group, own
}

// Predefined roles
var SystemRoles = map[string]Role{
    "super_admin": {
        Name:        "Super Administrator",
        Description: "Full system access",
        Permissions: []Permission{
            {Resource: "*", Actions: []string{"*"}, Scope: "global"},
        },
        IsSystem: true,
    },
    "admin": {
        Name:        "Administrator",
        Description: "Administrative access to devices and users",
        Permissions: []Permission{
            {Resource: "device", Actions: []string{"create", "read", "update", "delete"}, Scope: "global"},
            {Resource: "user", Actions: []string{"create", "read", "update", "delete"}, Scope: "group"},
            {Resource: "system", Actions: []string{"read"}, Scope: "global"},
        },
        IsSystem: true,
    },
    "operator": {
        Name:        "Operator",
        Description: "Device management access",
        Permissions: []Permission{
            {Resource: "device", Actions: []string{"read", "update"}, Scope: "group"},
            {Resource: "user", Actions: []string{"read"}, Scope: "own"},
        },
        IsSystem: true,
    },
    "viewer": {
        Name:        "Viewer",
        Description: "Read-only access",
        Permissions: []Permission{
            {Resource: "device", Actions: []string{"read"}, Scope: "group"},
            {Resource: "user", Actions: []string{"read"}, Scope: "own"},
        },
        IsSystem: true,
    },
}

func (auth *AuthService) CheckPermission(userID, resource, action string, resourceID string) bool {
    user, err := auth.db.GetUser(userID)
    if err != nil {
        return false
    }
    
    role, exists := SystemRoles[user.Role]
    if !exists {
        return false
    }
    
    for _, perm := range role.Permissions {
        if perm.Resource == "*" || perm.Resource == resource {
            if contains(perm.Actions, "*") || contains(perm.Actions, action) {
                if auth.checkScope(user, perm.Scope, resource, resourceID) {
                    return true
                }
            }
        }
    }
    
    return false
}
```

## Transport Security

### TLS Configuration
```yaml
# TLS Configuration for API Server
tls:
  enabled: true
  version: "1.3"
  cert_file: "/certs/api-server.crt"
  key_file: "/certs/api-server.key"
  ca_file: "/certs/ca.crt"
  
  # Client certificate verification
  client_auth: "require"
  client_ca_file: "/certs/client-ca.crt"
  
  # Cipher suites (TLS 1.3)
  cipher_suites:
    - TLS_AES_256_GCM_SHA384
    - TLS_AES_128_GCM_SHA256
    - TLS_CHACHA20_POLY1305_SHA256
  
  # Security headers
  security_headers:
    strict_transport_security: "max-age=31536000; includeSubDomains"
    x_frame_options: "DENY"
    x_content_type_options: "nosniff"
    x_xss_protection: "1; mode=block"
    content_security_policy: "default-src 'self'"
```

### Certificate Management
```go
type CertificateManager struct {
    storage   CertStorage
    acme      ACMEClient
    scheduler *cron.Cron
}

func NewCertificateManager(config *CertConfig) *CertificateManager {
    cm := &CertificateManager{
        storage:   NewFileCertStorage(config.CertDir),
        acme:      NewACMEClient(config.ACME),
        scheduler: cron.New(),
    }
    
    // Schedule certificate renewal check
    cm.scheduler.AddFunc("@daily", cm.checkAndRenewCertificates)
    cm.scheduler.Start()
    
    return cm
}

func (cm *CertificateManager) checkAndRenewCertificates() {
    certs, err := cm.storage.ListCertificates()
    if err != nil {
        log.Errorf("Failed to list certificates: %v", err)
        return
    }
    
    for _, cert := range certs {
        if cert.ExpiresWithin(30 * 24 * time.Hour) {
            log.Infof("Renewing certificate for %s", cert.Subject)
            err := cm.renewCertificate(cert)
            if err != nil {
                log.Errorf("Failed to renew certificate for %s: %v", cert.Subject, err)
                // Send alert
                cm.sendAlert("certificate_renewal_failed", cert.Subject, err)
            }
        }
    }
}

func (cm *CertificateManager) renewCertificate(cert *Certificate) error {
    // Request new certificate from ACME provider
    newCert, err := cm.acme.RenewCertificate(cert.Domain)
    if err != nil {
        return fmt.Errorf("ACME renewal failed: %w", err)
    }
    
    // Validate certificate
    if err := cm.validateCertificate(newCert); err != nil {
        return fmt.Errorf("certificate validation failed: %w", err)
    }
    
    // Store new certificate
    if err := cm.storage.StoreCertificate(newCert); err != nil {
        return fmt.Errorf("failed to store certificate: %w", err)
    }
    
    // Reload server configuration
    if err := cm.reloadServerConfig(); err != nil {
        return fmt.Errorf("failed to reload server config: %w", err)
    }
    
    log.Infof("Successfully renewed certificate for %s", cert.Domain)
    return nil
}
```

### Device Authentication
```go
type DeviceAuthenticator struct {
    certValidator CertificateValidator
    keyManager    KeyManager
    trustStore    TrustStore
}

func (da *DeviceAuthenticator) AuthenticateDevice(cert *x509.Certificate) (*DeviceIdentity, error) {
    // Validate certificate chain
    if err := da.certValidator.ValidateChain(cert); err != nil {
        return nil, fmt.Errorf("certificate validation failed: %w", err)
    }
    
    // Check certificate revocation
    if revoked, err := da.trustStore.IsRevoked(cert); err != nil {
        return nil, fmt.Errorf("revocation check failed: %w", err)
    } else if revoked {
        return nil, fmt.Errorf("certificate is revoked")
    }
    
    // Extract device identity from certificate
    identity, err := da.extractDeviceIdentity(cert)
    if err != nil {
        return nil, fmt.Errorf("failed to extract device identity: %w", err)
    }
    
    // Verify device is authorized
    if !da.isDeviceAuthorized(identity) {
        return nil, fmt.Errorf("device not authorized")
    }
    
    return identity, nil
}

func (da *DeviceAuthenticator) extractDeviceIdentity(cert *x509.Certificate) (*DeviceIdentity, error) {
    // Parse certificate subject for device information
    subject := cert.Subject
    
    // Extract OUI from Organization
    oui := subject.Organization
    if len(oui) == 0 {
        return nil, fmt.Errorf("missing OUI in certificate")
    }
    
    // Extract serial number from Common Name
    serialNumber := subject.CommonName
    if serialNumber == "" {
        return nil, fmt.Errorf("missing serial number in certificate")
    }
    
    // Extract product class from Organizational Unit
    var productClass string
    if len(subject.OrganizationalUnit) > 0 {
        productClass = subject.OrganizationalUnit[0]
    }
    
    return &DeviceIdentity{
        OUI:          oui[0],
        SerialNumber: serialNumber,
        ProductClass: productClass,
        Certificate:  cert,
    }, nil
}
```

## Input Validation & Sanitization

### Request Validation
```go
type Validator struct {
    validate *validator.Validate
    rules    map[string]ValidationRule
}

type ValidationRule struct {
    Required bool
    Type     string
    MinLen   int
    MaxLen   int
    Pattern  *regexp.Regexp
    Enum     []string
}

func NewValidator() *Validator {
    v := &Validator{
        validate: validator.New(),
        rules:    make(map[string]ValidationRule),
    }
    
    // Register custom validation rules
    v.validate.RegisterValidation("endpoint_id", validateEndpointID)
    v.validate.RegisterValidation("parameter_path", validateParameterPath)
    v.validate.RegisterValidation("parameter_value", validateParameterValue)
    
    return v
}

func validateEndpointID(fl validator.FieldLevel) bool {
    endpointID := fl.Field().String()
    
    // USP Endpoint ID format: proto::oui-prod_class-serial_number
    pattern := `^[a-z]+::[A-F0-9]{6}-[A-Za-z0-9\-_]+-[A-Za-z0-9\-_]+$`
    matched, _ := regexp.MatchString(pattern, endpointID)
    return matched
}

func validateParameterPath(fl validator.FieldLevel) bool {
    path := fl.Field().String()
    
    // Data model path format
    pattern := `^Device\.([A-Za-z0-9]+\.)*[A-Za-z0-9]+$`
    matched, _ := regexp.MatchString(pattern, path)
    return matched
}

func (v *Validator) ValidateRequest(req interface{}) error {
    if err := v.validate.Struct(req); err != nil {
        var validationErrors []ValidationError
        
        for _, err := range err.(validator.ValidationErrors) {
            validationErrors = append(validationErrors, ValidationError{
                Field:   err.Field(),
                Tag:     err.Tag(),
                Value:   err.Value(),
                Message: getValidationMessage(err),
            })
        }
        
        return &ValidationErrorCollection{Errors: validationErrors}
    }
    
    return nil
}

type SetParameterRequest struct {
    AgentID    string                `json:"agent_id" validate:"required,mongodb"`
    Parameters []ParameterSetting    `json:"parameters" validate:"required,min=1,max=100"`
}

type ParameterSetting struct {
    Path  string      `json:"path" validate:"required,parameter_path"`
    Value interface{} `json:"value" validate:"required"`
}
```

### SQL Injection Prevention
```go
// Use parameterized queries for database operations
func (db *Database) GetAgentByEndpointID(endpointID string) (*Agent, error) {
    // MongoDB query with proper escaping
    filter := bson.M{"endpoint_id": endpointID}
    
    var agent Agent
    err := db.agents.FindOne(context.Background(), filter).Decode(&agent)
    if err != nil {
        if err == mongo.ErrNoDocuments {
            return nil, ErrAgentNotFound
        }
        return nil, err
    }
    
    return &agent, nil
}

// For any external database queries, use prepared statements
func (db *Database) GetParameterHistory(deviceID, path string, limit int) ([]*ParameterHistory, error) {
    query := `
        SELECT path, value, timestamp, operation_id 
        FROM parameter_history 
        WHERE device_id = $1 AND path = $2 
        ORDER BY timestamp DESC 
        LIMIT $3`
    
    rows, err := db.conn.Query(query, deviceID, path, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    var history []*ParameterHistory
    for rows.Next() {
        var h ParameterHistory
        err := rows.Scan(&h.Path, &h.Value, &h.Timestamp, &h.OperationID)
        if err != nil {
            return nil, err
        }
        history = append(history, &h)
    }
    
    return history, rows.Err()
}
```

### Cross-Site Scripting (XSS) Prevention
```go
import "html/template"
import "github.com/microcosm-cc/bluemonday"

type TemplateData struct {
    Title       template.HTML
    Content     template.HTML
    UserInput   string // Raw user input - will be escaped
}

func (h *HTTPHandler) renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    // Use strict sanitization policy
    p := bluemonday.StrictPolicy()
    
    // Create template with custom functions
    t := template.New(tmpl).Funcs(template.FuncMap{
        "sanitizeHTML": func(input string) template.HTML {
            return template.HTML(p.Sanitize(input))
        },
        "escapeJS": func(input string) template.JS {
            return template.JS(template.JSEscapeString(input))
        },
    })
    
    parsedTemplate, err := t.ParseFiles("templates/" + tmpl + ".html")
    if err != nil {
        http.Error(w, "Template error", http.StatusInternalServerError)
        return
    }
    
    // Set security headers
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "DENY")
    w.Header().Set("X-XSS-Protection", "1; mode=block")
    
    err = parsedTemplate.Execute(w, data)
    if err != nil {
        log.Errorf("Template execution error: %v", err)
    }
}

// API responses with proper escaping
func (h *APIHandler) ServeJSON(w http.ResponseWriter, data interface{}) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    
    encoder := json.NewEncoder(w)
    encoder.SetEscapeHTML(true) // Escape HTML in JSON strings
    
    if err := encoder.Encode(data); err != nil {
        log.Errorf("JSON encoding error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
```

## Rate Limiting & DDoS Protection

### Rate Limiting Implementation
```go
type RateLimiter struct {
    store    cache.Store
    rules    []RateRule
    enforcer RateEnforcer
}

type RateRule struct {
    Path     string        `yaml:"path"`
    Method   string        `yaml:"method"`
    Limit    int          `yaml:"limit"`
    Window   time.Duration `yaml:"window"`
    Scope    string        `yaml:"scope"` // ip, user, api_key
    Burst    int          `yaml:"burst"`
}

func (rl *RateLimiter) CheckLimit(req *http.Request, user *User) error {
    rule := rl.findMatchingRule(req)
    if rule == nil {
        return nil // No rate limit for this endpoint
    }
    
    key := rl.buildKey(req, user, rule)
    
    // Check current usage
    current, err := rl.store.GetCounter(key)
    if err != nil {
        log.Errorf("Failed to get rate limit counter: %v", err)
        return nil // Fail open
    }
    
    if current >= rule.Limit {
        return &RateLimitError{
            Rule:      rule,
            Current:   current,
            ResetTime: rl.getResetTime(key),
        }
    }
    
    // Increment counter
    err = rl.store.IncrementCounter(key, rule.Window)
    if err != nil {
        log.Errorf("Failed to increment rate limit counter: %v", err)
    }
    
    return nil
}

func (rl *RateLimiter) buildKey(req *http.Request, user *User, rule *RateRule) string {
    var identifier string
    
    switch rule.Scope {
    case "ip":
        identifier = getClientIP(req)
    case "user":
        if user != nil {
            identifier = user.ID
        } else {
            identifier = getClientIP(req) // Fallback to IP
        }
    case "api_key":
        if apiKey := req.Header.Get("X-API-Key"); apiKey != "" {
            identifier = hashString(apiKey)
        } else {
            identifier = getClientIP(req) // Fallback to IP
        }
    default:
        identifier = getClientIP(req)
    }
    
    return fmt.Sprintf("ratelimit:%s:%s:%s:%s", 
        rule.Scope, rule.Path, rule.Method, identifier)
}
```

### DDoS Protection
```yaml
ddos_protection:
  enabled: true
  
  # Connection limits
  max_connections_per_ip: 100
  max_new_connections_per_second: 10
  
  # Request size limits
  max_request_size: 1048576  # 1MB
  max_header_size: 8192      # 8KB
  
  # Timeout settings
  read_timeout: 30s
  write_timeout: 30s
  idle_timeout: 120s
  
  # Blacklist/Whitelist
  ip_whitelist:
    - "10.0.0.0/8"
    - "172.16.0.0/12"
    - "192.168.0.0/16"
  
  ip_blacklist:
    - "1.2.3.4"
    - "5.6.7.0/24"
  
  # Geographic blocking
  geo_blocking:
    enabled: false
    allowed_countries: ["US", "CA", "GB", "DE", "FR"]
    blocked_countries: ["CN", "RU"]
```

## Data Encryption

### Encryption at Rest
```go
type EncryptionService struct {
    keyManager KeyManager
    cipher     cipher.AEAD
}

func NewEncryptionService(keyManager KeyManager) (*EncryptionService, error) {
    // Get current encryption key
    key, err := keyManager.GetCurrentKey()
    if err != nil {
        return nil, fmt.Errorf("failed to get encryption key: %w", err)
    }
    
    // Create AES-GCM cipher
    block, err := aes.NewCipher(key)
    if err != nil {
        return nil, fmt.Errorf("failed to create cipher: %w", err)
    }
    
    gcm, err := cipher.NewGCM(block)
    if err != nil {
        return nil, fmt.Errorf("failed to create GCM: %w", err)
    }
    
    return &EncryptionService{
        keyManager: keyManager,
        cipher:     gcm,
    }, nil
}

func (es *EncryptionService) Encrypt(plaintext []byte) ([]byte, error) {
    // Generate random nonce
    nonce := make([]byte, es.cipher.NonceSize())
    if _, err := rand.Read(nonce); err != nil {
        return nil, fmt.Errorf("failed to generate nonce: %w", err)
    }
    
    // Encrypt data
    ciphertext := es.cipher.Seal(nil, nonce, plaintext, nil)
    
    // Prepend nonce to ciphertext
    result := make([]byte, len(nonce)+len(ciphertext))
    copy(result, nonce)
    copy(result[len(nonce):], ciphertext)
    
    return result, nil
}

func (es *EncryptionService) Decrypt(data []byte) ([]byte, error) {
    if len(data) < es.cipher.NonceSize() {
        return nil, fmt.Errorf("invalid encrypted data")
    }
    
    // Extract nonce and ciphertext
    nonce := data[:es.cipher.NonceSize()]
    ciphertext := data[es.cipher.NonceSize():]
    
    // Decrypt data
    plaintext, err := es.cipher.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, fmt.Errorf("decryption failed: %w", err)
    }
    
    return plaintext, nil
}

// Encrypt sensitive fields before storing in database
func (db *Database) StoreUser(user *User) error {
    // Encrypt sensitive fields
    if user.MFASecret != "" {
        encrypted, err := db.encryptionService.Encrypt([]byte(user.MFASecret))
        if err != nil {
            return fmt.Errorf("failed to encrypt MFA secret: %w", err)
        }
        user.MFASecret = base64.StdEncoding.EncodeToString(encrypted)
    }
    
    // Hash password
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("failed to hash password: %w", err)
    }
    user.PasswordHash = string(hashedPassword)
    user.Password = "" // Clear plaintext password
    
    // Store in database
    _, err = db.users.InsertOne(context.Background(), user)
    return err
}
```

### Key Management
```go
type KeyManager interface {
    GetCurrentKey() ([]byte, error)
    GetKey(keyID string) ([]byte, error)
    RotateKey() error
    ListKeys() ([]KeyInfo, error)
}

type VaultKeyManager struct {
    client   *api.Client
    keyPath  string
    keyCache map[string]*CachedKey
    mu       sync.RWMutex
}

type CachedKey struct {
    Key       []byte
    ExpiresAt time.Time
}

func NewVaultKeyManager(config *VaultConfig) (*VaultKeyManager, error) {
    client, err := api.NewClient(&api.Config{
        Address: config.Address,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create Vault client: %w", err)
    }
    
    // Authenticate with Vault
    err = client.SetToken(config.Token)
    if err != nil {
        return nil, fmt.Errorf("failed to set Vault token: %w", err)
    }
    
    return &VaultKeyManager{
        client:   client,
        keyPath:  config.KeyPath,
        keyCache: make(map[string]*CachedKey),
    }, nil
}

func (vkm *VaultKeyManager) GetCurrentKey() ([]byte, error) {
    // Check cache first
    vkm.mu.RLock()
    cached, exists := vkm.keyCache["current"]
    vkm.mu.RUnlock()
    
    if exists && time.Now().Before(cached.ExpiresAt) {
        return cached.Key, nil
    }
    
    // Fetch from Vault
    secret, err := vkm.client.Logical().Read(vkm.keyPath + "/current")
    if err != nil {
        return nil, fmt.Errorf("failed to read key from Vault: %w", err)
    }
    
    if secret == nil || secret.Data == nil {
        return nil, fmt.Errorf("key not found in Vault")
    }
    
    keyBase64, ok := secret.Data["key"].(string)
    if !ok {
        return nil, fmt.Errorf("invalid key format in Vault")
    }
    
    key, err := base64.StdEncoding.DecodeString(keyBase64)
    if err != nil {
        return nil, fmt.Errorf("failed to decode key: %w", err)
    }
    
    // Cache the key
    vkm.mu.Lock()
    vkm.keyCache["current"] = &CachedKey{
        Key:       key,
        ExpiresAt: time.Now().Add(5 * time.Minute),
    }
    vkm.mu.Unlock()
    
    return key, nil
}
```

## Audit Logging

### Security Event Logging
```go
type SecurityLogger struct {
    logger       *logrus.Logger
    db          *Database
    alertManager AlertManager
}

type SecurityEvent struct {
    ID          string                 `json:"id"`
    Timestamp   time.Time             `json:"timestamp"`
    EventType   string                `json:"event_type"`
    Severity    string                `json:"severity"`
    UserID      string                `json:"user_id,omitempty"`
    DeviceID    string                `json:"device_id,omitempty"`
    IPAddress   string                `json:"ip_address"`
    UserAgent   string                `json:"user_agent,omitempty"`
    Details     map[string]interface{} `json:"details"`
    Risk        int                   `json:"risk_score"`
}

func (sl *SecurityLogger) LogEvent(eventType string, details map[string]interface{}) {
    event := &SecurityEvent{
        ID:        generateEventID(),
        Timestamp: time.Now(),
        EventType: eventType,
        Severity:  sl.calculateSeverity(eventType, details),
        Details:   details,
        Risk:      sl.calculateRiskScore(eventType, details),
    }
    
    // Extract context information
    if userID, ok := details["user_id"].(string); ok {
        event.UserID = userID
    }
    
    if deviceID, ok := details["device_id"].(string); ok {
        event.DeviceID = deviceID
    }
    
    if ipAddress, ok := details["ip_address"].(string); ok {
        event.IPAddress = ipAddress
    }
    
    // Log to structured logger
    sl.logger.WithFields(logrus.Fields{
        "event_id":   event.ID,
        "event_type": event.EventType,
        "severity":   event.Severity,
        "user_id":    event.UserID,
        "device_id":  event.DeviceID,
        "ip_address": event.IPAddress,
        "risk_score": event.Risk,
    }).Info("Security event")
    
    // Store in database for analysis
    go func() {
        if err := sl.db.StoreSecurityEvent(event); err != nil {
            sl.logger.Errorf("Failed to store security event: %v", err)
        }
    }()
    
    // Send alerts for high-risk events
    if event.Risk >= 80 {
        go sl.alertManager.SendAlert(event)
    }
}

// Security event types
const (
    EventLoginSuccess           = "login_success"
    EventLoginFailure           = "login_failure"
    EventMultipleFailedLogins   = "multiple_failed_logins"
    EventPasswordChange         = "password_change"
    EventMFAEnabled            = "mfa_enabled"
    EventMFADisabled           = "mfa_disabled"
    EventAccountLocked         = "account_locked"
    EventPrivilegeEscalation   = "privilege_escalation"
    EventUnauthorizedAccess    = "unauthorized_access"
    EventSuspiciousActivity    = "suspicious_activity"
    EventDataExfiltration      = "data_exfiltration"
    EventSystemCompromise      = "system_compromise"
    EventDeviceCompromise      = "device_compromise"
    EventCertificateExpired    = "certificate_expired"
    EventCertificateRevoked    = "certificate_revoked"
    EventRateLimitExceeded     = "rate_limit_exceeded"
    EventDDoSDetected         = "ddos_detected"
)
```

### Compliance Reporting
```go
type ComplianceReporter struct {
    db       *Database
    logger   *SecurityLogger
    config   *ComplianceConfig
    storage  ReportStorage
}

type ComplianceReport struct {
    ID          string                 `json:"id"`
    Type        string                 `json:"type"`        // SOX, GDPR, HIPAA, PCI-DSS
    Period      Period                 `json:"period"`
    GeneratedAt time.Time             `json:"generated_at"`
    GeneratedBy string                 `json:"generated_by"`
    Data        map[string]interface{} `json:"data"`
    Status      string                 `json:"status"`     // draft, final, submitted
}

func (cr *ComplianceReporter) GenerateGDPRReport(startDate, endDate time.Time) (*ComplianceReport, error) {
    report := &ComplianceReport{
        ID:          generateReportID(),
        Type:        "GDPR",
        Period:      Period{Start: startDate, End: endDate},
        GeneratedAt: time.Now(),
        GeneratedBy: "system",
        Status:      "draft",
    }
    
    // Data processing activities
    activities, err := cr.db.GetDataProcessingActivities(startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get processing activities: %w", err)
    }
    
    // Data subject requests
    requests, err := cr.db.GetDataSubjectRequests(startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get data subject requests: %w", err)
    }
    
    // Security incidents
    incidents, err := cr.db.GetSecurityIncidents(startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get security incidents: %w", err)
    }
    
    // Data breaches
    breaches, err := cr.db.GetDataBreaches(startDate, endDate)
    if err != nil {
        return nil, fmt.Errorf("failed to get data breaches: %w", err)
    }
    
    report.Data = map[string]interface{}{
        "processing_activities": activities,
        "data_subject_requests": requests,
        "security_incidents":   incidents,
        "data_breaches":        breaches,
        "privacy_impact_assessments": cr.getPIAStatus(),
        "data_protection_measures":   cr.getDataProtectionMeasures(),
    }
    
    // Store report
    err = cr.storage.StoreReport(report)
    if err != nil {
        return nil, fmt.Errorf("failed to store report: %w", err)
    }
    
    return report, nil
}
```

This security guide provides comprehensive information for implementing and maintaining security controls in the OpenUSP platform.