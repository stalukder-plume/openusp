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
	"encoding/xml"
	"time"
)

// SOAP Envelope structure for TR-069 CWMP
type SOAPEnvelope struct {
	XMLName xml.Name `xml:"soap:Envelope"`
	SoapNS  string   `xml:"xmlns:soap,attr"`
	CwmpNS  string   `xml:"xmlns:cwmp,attr"`
	XsiNS   string   `xml:"xmlns:xsi,attr"`
	XsdNS   string   `xml:"xmlns:xsd,attr"`
	Header  *SOAPHeader `xml:"soap:Header,omitempty"`
	Body    SOAPBody    `xml:"soap:Body"`
}

type SOAPHeader struct {
	ID                string `xml:"cwmp:ID,omitempty"`
	HoldRequests      bool   `xml:"cwmp:HoldRequests,omitempty"`
	NoMoreRequests    bool   `xml:"cwmp:NoMoreRequests,omitempty"`
	SessionTimeout    uint32 `xml:"cwmp:SessionTimeout,omitempty"`
}

type SOAPBody struct {
	XMLName   xml.Name    `xml:"soap:Body"`
	Content   interface{} `xml:",omitempty"`
	Fault     *SOAPFault  `xml:"soap:Fault,omitempty"`
}

type SOAPFault struct {
	FaultCode   string      `xml:"faultcode"`
	FaultString string      `xml:"faultstring"`
	Detail      *FaultDetail `xml:"detail,omitempty"`
}

type FaultDetail struct {
	CWMPFault *CWMPFault `xml:"cwmp:Fault,omitempty"`
}

type CWMPFault struct {
	FaultCode   uint32 `xml:"FaultCode"`
	FaultString string `xml:"FaultString"`
}

// TR-069 CWMP Method structures

// Inform method
type Inform struct {
	XMLName      xml.Name           `xml:"cwmp:Inform"`
	DeviceId     DeviceIdStruct     `xml:"DeviceId"`
	Event        []EventStruct      `xml:"Event>EventStruct"`
	MaxEnvelopes uint32            `xml:"MaxEnvelopes"`
	CurrentTime  time.Time         `xml:"CurrentTime"`
	RetryCount   uint32            `xml:"RetryCount"`
	ParameterList []ParameterValueStruct `xml:"ParameterList>ParameterValueStruct"`
}

type InformResponse struct {
	XMLName      xml.Name `xml:"cwmp:InformResponse"`
	MaxEnvelopes uint32   `xml:"MaxEnvelopes"`
}

// GetParameterValues method
type GetParameterValues struct {
	XMLName       xml.Name `xml:"cwmp:GetParameterValues"`
	ParameterNames []string `xml:"ParameterNames>string"`
}

type GetParameterValuesResponse struct {
	XMLName       xml.Name               `xml:"cwmp:GetParameterValuesResponse"`
	ParameterList []ParameterValueStruct `xml:"ParameterList>ParameterValueStruct"`
}

// SetParameterValues method
type SetParameterValues struct {
	XMLName       xml.Name               `xml:"cwmp:SetParameterValues"`
	ParameterList []ParameterValueStruct `xml:"ParameterList>ParameterValueStruct"`
	ParameterKey  string                 `xml:"ParameterKey"`
}

type SetParameterValuesResponse struct {
	XMLName xml.Name `xml:"cwmp:SetParameterValuesResponse"`
	Status  uint32   `xml:"Status"`
}

// GetParameterNames method
type GetParameterNames struct {
	XMLName       xml.Name `xml:"cwmp:GetParameterNames"`
	ParameterPath string   `xml:"ParameterPath"`
	NextLevel     bool     `xml:"NextLevel"`
}

type GetParameterNamesResponse struct {
	XMLName       xml.Name             `xml:"cwmp:GetParameterNamesResponse"`
	ParameterList []ParameterInfoStruct `xml:"ParameterList>ParameterInfoStruct"`
}

// AddObject method
type AddObject struct {
	XMLName      xml.Name `xml:"cwmp:AddObject"`
	ObjectName   string   `xml:"ObjectName"`
	ParameterKey string   `xml:"ParameterKey"`
}

type AddObjectResponse struct {
	XMLName        xml.Name `xml:"cwmp:AddObjectResponse"`
	InstanceNumber uint32   `xml:"InstanceNumber"`
	Status         uint32   `xml:"Status"`
}

// DeleteObject method
type DeleteObject struct {
	XMLName      xml.Name `xml:"cwmp:DeleteObject"`
	ObjectName   string   `xml:"ObjectName"`
	ParameterKey string   `xml:"ParameterKey"`
}

type DeleteObjectResponse struct {
	XMLName xml.Name `xml:"cwmp:DeleteObjectResponse"`
	Status  uint32   `xml:"Status"`
}

// Reboot method
type Reboot struct {
	XMLName    xml.Name `xml:"cwmp:Reboot"`
	CommandKey string   `xml:"CommandKey"`
}

type RebootResponse struct {
	XMLName xml.Name `xml:"cwmp:RebootResponse"`
}

// FactoryReset method
type FactoryReset struct {
	XMLName xml.Name `xml:"cwmp:FactoryReset"`
}

type FactoryResetResponse struct {
	XMLName xml.Name `xml:"cwmp:FactoryResetResponse"`
}

// Download method
type Download struct {
	XMLName        xml.Name `xml:"cwmp:Download"`
	CommandKey     string   `xml:"CommandKey"`
	FileType       string   `xml:"FileType"`
	URL            string   `xml:"URL"`
	Username       string   `xml:"Username"`
	Password       string   `xml:"Password"`
	FileSize       uint32   `xml:"FileSize"`
	TargetFileName string   `xml:"TargetFileName"`
	DelaySeconds   uint32   `xml:"DelaySeconds"`
	SuccessURL     string   `xml:"SuccessURL"`
	FailureURL     string   `xml:"FailureURL"`
}

type DownloadResponse struct {
	XMLName      xml.Name  `xml:"cwmp:DownloadResponse"`
	Status       uint32    `xml:"Status"`
	StartTime    time.Time `xml:"StartTime"`
	CompleteTime time.Time `xml:"CompleteTime"`
}

// Upload method
type Upload struct {
	XMLName      xml.Name `xml:"cwmp:Upload"`
	CommandKey   string   `xml:"CommandKey"`
	FileType     string   `xml:"FileType"`
	URL          string   `xml:"URL"`
	Username     string   `xml:"Username"`
	Password     string   `xml:"Password"`
	DelaySeconds uint32   `xml:"DelaySeconds"`
}

type UploadResponse struct {
	XMLName      xml.Name  `xml:"cwmp:UploadResponse"`
	Status       uint32    `xml:"Status"`
	StartTime    time.Time `xml:"StartTime"`
	CompleteTime time.Time `xml:"CompleteTime"`
}

// Common structures
type DeviceIdStruct struct {
	Manufacturer  string `xml:"Manufacturer"`
	OUI           string `xml:"OUI"`
	ProductClass  string `xml:"ProductClass"`
	SerialNumber  string `xml:"SerialNumber"`
}

type EventStruct struct {
	EventCode  string `xml:"EventCode"`
	CommandKey string `xml:"CommandKey"`
}

type ParameterValueStruct struct {
	Name  string `xml:"Name"`
	Value string `xml:"Value"`
	Type  string `xml:"Type,attr,omitempty"`
}

type ParameterInfoStruct struct {
	Name     string `xml:"Name"`
	Writable bool   `xml:"Writable"`
}

// TR-069 Event codes
const (
	EventBootstrap        = "0 BOOTSTRAP"
	EventBoot            = "1 BOOT"
	EventPeriodic        = "2 PERIODIC"
	EventScheduled       = "3 SCHEDULED"
	EventValueChange     = "4 VALUE CHANGE"
	EventKicked          = "5 KICKED"
	EventConnectionRequest = "6 CONNECTION REQUEST"
	EventTransferComplete = "7 TRANSFER COMPLETE"
	EventDiagnosticsComplete = "8 DIAGNOSTICS COMPLETE"
	EventRequestDownload = "9 REQUEST DOWNLOAD"
	EventAutonomousTransferComplete = "10 AUTONOMOUS TRANSFER COMPLETE"
	EventDUStateChangeComplete = "11 DU STATE CHANGE COMPLETE"
	EventAutonomousDUStateChangeComplete = "12 AUTONOMOUS DU STATE CHANGE COMPLETE"
	EventWakeUp          = "13 WAKEUP"
)

// TR-069 CWMP Fault codes
const (
	FaultMethodNotSupported     = 9000
	FaultRequestDenied         = 9001
	FaultInternalError         = 9002
	FaultInvalidArguments      = 9003
	FaultResourcesExceeded     = 9004
	FaultInvalidParameterName  = 9005
	FaultInvalidParameterType  = 9006
	FaultInvalidParameterValue = 9007
	FaultAttemptToSetNonWritableParameter = 9008
	FaultNotificationRequestRejected = 9009
	FaultDownloadFailure       = 9010
	FaultUploadFailure         = 9011
	FaultFileTransferServerAuthenticationFailure = 9012
	FaultUnsupportedProtocolForFileTransfer = 9013
	FaultFileTransferFailure   = 9014
	FaultFileTransferFailureContactServer = 9015
	FaultFileTransferFailureAccessFile = 9016
	FaultFileTransferFailureCompleteDownload = 9017
	FaultFileTransferFailureFileCorrupted = 9018
	FaultFileTransferFailureFileAuthentication = 9019
)