package models

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

type Tenant struct {
	ID        int
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// GenerateSchemaNameByTime creates a unique, time-based schema name.
// Example output: tenant_20251031_142530_a3f1
func (t *Tenant) GenerateSchemaNameByTime() error {
	// Get current UTC time
	timestamp := time.Now().UTC().Format("20060102_150405")

	// Generate 2 random bytes (4 hex chars)
	buf := make([]byte, 2)
	if _, err := rand.Read(buf); err != nil {
		return err
	}
	randomSuffix := hex.EncodeToString(buf)

	t.Name = fmt.Sprintf("tenant_%s_%s", timestamp, randomSuffix)
	return nil
}

func (t *Tenant) GetTenantByName(name string) (tenant Tenant, err error) {
	conn, err := createConn()
	if err != nil {
		return
	}

	rows, _ := conn.Query("SELECT * FROM services WHERE name=`?`", name)
	defer rows.Close()
	for rows.Next() {
		_ = rows.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, tenant.UpdatedAt)
	}

	return
}
