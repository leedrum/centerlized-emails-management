package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/models"
)

func Create(c *gin.Context, cfg config.Config) {
	tenant := models.Tenant{}
	err := tenant.GenerateSchemaNameByTime()
	// retry 1 time if exist
	if err != nil && strings.ContainsAny(err.Error(), "exist") {
		err = tenant.GenerateSchemaNameByTime()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "tenant already existed"})
		}
	}

	kafkaProducer, err := createSyncProducer(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}
	defer kafkaProducer.Close()

	err = sendMessage(kafkaProducer, tenant, "create")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("create success %s in a few mins", tenant.Name)})
}
