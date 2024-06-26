package controlapi

import (
	"log/slog"

	cloudevents "github.com/cloudevents/sdk-go"
)

const (
	APIPrefix = "$NEX"
)

const (
	AuctionResponseType  = "io.nats.nex.v1.auction_response"
	InfoResponseType     = "io.nats.nex.v1.info_response"
	PingResponseType     = "io.nats.nex.v1.ping_response"
	RunResponseType      = "io.nats.nex.v1.run_response"
	StopResponseType     = "io.nats.nex.v1.stop_response"
	LameDuckResponseType = "io.nats.nex.v1.lameduck_response"

	TagOS       = "nex.os"
	TagArch     = "nex.arch"
	TagCPUs     = "nex.cpucount"
	TagUnsafe   = "nex.unsafe"
	TagLameDuck = "nex.lameduck"
)

type RunResponse struct {
	Started bool   `json:"started"`
	ID      string `json:"id"`
	Issuer  string `json:"issuer"`
	Name    string `json:"name"`
}

type NexWorkload string

const (
	NexWorkloadNative NexWorkload = "native"
	NexWorkloadV8     NexWorkload = "v8"
	NexWorkloadOCI    NexWorkload = "oci"
	NexWorkloadWasm   NexWorkload = "wasm"

	// cloud events can't have - in extensions
	EventExtensionNamespace = "namespace"
)

type NodeCapabilities struct {
	Sandboxable        bool              `json:"sandboxable"`
	SupportedProviders []NexWorkload     `json:"supported_providers"`
	NodeTags           map[string]string `json:"node_tags"`
}

type AuctionRequest struct {
	Arch          *string           `json:"arch,omitempty"`
	OS            *string           `json:"os,omitempty"`
	Sandboxed     *bool             `json:"sandboxed,omitempty"`
	Tags          map[string]string `json:"tags,omitempty"`
	WorkloadTypes []NexWorkload     `json:"workload_types,omitempty"`
}

type AuctionResponse PingResponse

// TODO: remove omitempty in next version bump
type PingResponse struct {
	NodeId          string            `json:"node_id"`
	Nexus           string            `json:"nexus,omitempty"`
	Version         string            `json:"version"`
	Uptime          string            `json:"uptime"`
	TargetXkey      string            `json:"target_xkey"`
	Tags            map[string]string `json:"tags,omitempty"`
	RunningMachines int               `json:"running_machines"`
}

type WorkloadPingResponse struct {
	NodeId          string                       `json:"node_id"`
	TargetXkey      string                       `json:"target_xkey"`
	Version         string                       `json:"version"`
	Tags            map[string]string            `json:"tags,omitempty"`
	Uptime          string                       `json:"uptime"`
	RunningMachines []WorkloadPingMachineSummary `json:"running_machines"`
}

type WorkloadPingMachineSummary struct {
	Id           string      `json:"id"`
	Namespace    string      `json:"namespace"`
	Name         string      `json:"name"`
	WorkloadType NexWorkload `json:"type"`
}

type LameDuckResponse struct {
	NodeId  string `json:"node_id"`
	Success bool   `json:"success"`
}

type MemoryStat struct {
	MemTotal     int `json:"total"`
	MemFree      int `json:"free"`
	MemAvailable int `json:"available"`
}

type InfoResponse struct {
	Version                string            `json:"version"`
	Uptime                 string            `json:"uptime"`
	PublicXKey             string            `json:"public_xkey"`
	Tags                   map[string]string `json:"tags,omitempty"`
	Memory                 *MemoryStat       `json:"memory,omitempty"`
	Machines               []MachineSummary  `json:"machines"`
	SupportedWorkloadTypes []NexWorkload     `json:"supported_workload_types,omitempty"`
}

type MachineSummary struct {
	Id        string          `json:"id"`
	Healthy   bool            `json:"healthy"`
	Uptime    string          `json:"uptime"`
	Namespace string          `json:"namespace,omitempty"`
	Workload  WorkloadSummary `json:"workload,omitempty"`
}

type WorkloadSummary struct {
	Name         string      `json:"name"`
	Description  string      `json:"description,omitempty"`
	Runtime      string      `json:"runtime"`
	WorkloadType NexWorkload `json:"type"`
	Hash         string      `json:"hash"`
}

type Envelope struct {
	PayloadType string      `json:"type"`
	Data        interface{} `json:"data,omitempty"`
	Error       interface{} `json:"error,omitempty"`
}

// Wrapper for what goes across the wire
type EmittedLog struct {
	Namespace string `json:"namespace"`
	NodeId    string `json:"node_id"`
	Workload  string `json:"workload_id"`
	Timestamp string `json:"timestamp"`
	RawLog
}

type RawLog struct {
	Text  string     `json:"text"`
	Level slog.Level `json:"level"`
	ID    string     `json:"id"`
}

// Note this a wrapper to add context to a cloud event
type EmittedEvent struct {
	cloudevents.Event
	Namespace string `json:"namespace"`
	EventType string `json:"event_type"`
}

func NewEnvelope(dataType string, data interface{}, err *string) Envelope {
	var e interface{}
	if err != nil {
		e = *err
	}
	return Envelope{
		PayloadType: dataType,
		Data:        data,
		Error:       e,
	}
}
