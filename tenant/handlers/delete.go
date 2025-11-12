package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/leedrum/centerlized-emails-management/tenant/config"
	"github.com/leedrum/centerlized-emails-management/tenant/models"
)

func Delete(c *gin.Context, cfg config.Config) {
	tenant := models.Tenant{}
	tenant, err := tenant.GetTenantByName(c.Query("name"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}

	kafkaProducer, err := createSyncProducer(cfg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}
	defer kafkaProducer.Close()

	err = sendMessage(kafkaProducer, tenant, "delete")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "something went wrong"})
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("delete success %s in a few mins", tenant.Name)})
}
