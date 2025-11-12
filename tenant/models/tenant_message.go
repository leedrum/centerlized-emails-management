package models

import (
	"encoding/json"
	"time"
)

type TenantMessage struct {
	Name        string
	AggregateID string
	Payload     string
	EventType   string
	CreatedAt   time.Time
}

func (t *TenantMessage) Init(tenant Tenant) error {
	err := t.setPayload(tenant)
	if err != nil {
		return err
	}

	return nil
}

func (t *TenantMessage) CreateTenant() {
	t.Name = "migrate new tenant"
	t.EventType = "create_tenant"
	t.CreatedAt = time.Now()
}

func (t *TenantMessage) DeleteTenant() {
	t.Name = "delete tenant"
	t.EventType = "delete_tenant"
	t.CreatedAt = time.Now()
}

func (t *TenantMessage) setPayload(tenant Tenant) error {
	eventPayload := map[string]interface{}{
		"id":   tenant.ID,
		"name": tenant.Name,
	}

	payloadBytes, err := json.Marshal(eventPayload)
	if err != nil {
		return err
	}

	t.Payload = string(payloadBytes)

	return nil
}
