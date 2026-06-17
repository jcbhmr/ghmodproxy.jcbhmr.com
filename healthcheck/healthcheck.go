package healthcheck

import "time"

const ContentType = "application/health+json"

type Response struct {
	Status      Status   `json:"status"`
	Version     string   `json:"version,omitempty"`
	ReleaseID   string   `json:"releaseId,omitempty"`
	Notes       []string `json:"notes,omitempty"`
	Output      string   `json:"output,omitempty"`
	Checks      Checks   `json:"checks,omitempty"`
	Links       Links    `json:"links,omitzero"`
	ServiceID   string   `json:"serviceId,omitempty"`
	Description string   `json:"description,omitempty"`
}

type Status string

const (
	StatusPass Status = "pass"
	StatusFail Status = "fail"
	StatusWarn Status = "warn"
)

type Checks map[string][]ComponentDetails

func (c Checks) Add(componentName string, measurementName string, value ComponentDetails) {
	key := componentName + ":" + measurementName
	c[key] = append(c[key], value)
}

func (c Checks) Set(componentName string, measurementName string, value []ComponentDetails) {
	key := componentName + ":" + measurementName
	c[key] = value
}

type ComponentDetails struct {
	ComponentID       string    `json:"componentId,omitempty"`
	ComponentType     string    `json:"componentType,omitempty"`
	ObservedValue     any       `json:"observedValue,omitzero"`
	ObservedUnit      string    `json:"observedUnit,omitempty"`
	Status            Status    `json:"status,omitempty"`
	AffectedEndpoints []string  `json:"affectedEndpoints,omitempty"`
	Time              time.Time `json:"omitzero"`
	Output            string    `json:"output,omitempty"`
	Links             Links     `json:"links,omitzero"`
}
