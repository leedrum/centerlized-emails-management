package models

import (
	"database/sql"
	"fmt"
	"time"
)

type Service struct {
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func GetList() (services []Service, err error) {
	conn, err := createConn()
	if err != nil {
		return nil, err
	}

	rows, _ := conn.Query("SELECT * FROM services")
	defer rows.Close()
	for rows.Next() {
		var service Service
		_ = rows.Scan(&service.Name, &service.CreatedAt, service.UpdatedAt)
		services = append(services, service)
	}

	return
}

func (service Service) Create(name string) error {
	conn, err := createConn()
	if err != nil {
		return err
	}

	sql := "CREATE table services (foo INTEGER, bar TEXT) IF NOT EXISTS"
	_, _ = conn.Exec(sql)

	sql = "INSERT INTO services (name, create_at, updated_at) values (?, ?, ?)"
	stmt, err := conn.Prepare(sql)
	defer stmt.Close()
	_, _ = stmt.Exec(name, time.Now(), time.Now())
	rows, _ := conn.Query("SELECT * from services")
	defer rows.Close()
	for rows.Next() {
		var tenant Tenant
		_ = rows.Scan(&tenant.ID, &tenant.Name, &tenant.CreatedAt, tenant.UpdatedAt)
		tenantMessage := TenantMessage{}
		err := tenantMessage.Init(tenant)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service Service) Delete(name string) error {
	conn, err := createConn()

	sql := "DELETE services WHERE NAME='?'"
	_, err = conn.Exec(sql, name)
	defer conn.Close()

	return err
}

func createConn() (*sql.DB, error) {
	conn, err := sql.Open("turso", ":memory:")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return nil, err
	}

	return conn, err
}
