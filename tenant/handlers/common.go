package handlers

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/models"
	"github.com/leedrum/centerlized-emails-management/tenant/producer"
)

func sendMessage(kafkaProducer sarama.SyncProducer, tenant models.Tenant, action string) error {
	services, err := models.GetList()
	if err == nil {
		tenantMessage := models.TenantMessage{}
		err := tenantMessage.Init(tenant)
		if err != nil {
			log.Fatalf("Failed to init tenant message: %v", err)
			return err
		}

		if action == "create" {
			tenantMessage.CreateTenant()
		} else {
			tenantMessage.DeleteTenant()
		}

		for _, service := range services {
			topic := fmt.Sprintf("cem.migrate.%", service.Name)
			err = producer.PublishToKafka(tenantMessage, kafkaProducer, topic)
			if err != nil {
				log.Fatalf("Failed to create Kafka producer: %v", err)
				return err
			}
		}
	} else {
		log.Fatalf("Failed to get list services: %v", err)
		return err
	}

	return nil
}

func initTenantMessage(tenant models.Tenant) error {
	tenantMessage := models.TenantMessage{}
	err := tenantMessage.Init(tenant)
	if err != nil {
		log.Fatalf("Failed to init tenant message: %v", err)
		return err
	}

	return nil
}

func createSyncProducer(cfg config.Config) (sarama.SyncProducer, error) {
	kafkaProducer, err := producer.NewKafkaSyncProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
		return nil, err
	}
	return kafkaProducer, nil
}
