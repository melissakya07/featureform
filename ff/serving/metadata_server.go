package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MetaServer interface {
	ExposePort(port string)
	GetMetadata(c *gin.Context)
}

type MetadataServer struct {
	Server *FeatureServer
}

func NewMetadataServer(logger *zap.SugaredLogger, serv *FeatureServer) (*MetadataServer, error) {
	logger.Debug("Creating new metadata server")
	return &MetadataServer{
		Server: serv,
	}, nil
}

func (m MetadataServer) GetTrainingMetadata(c *gin.Context) {

	name := c.Param("name")
	version := c.Param("version")

	entry, err := m.Server.Metadata.TrainingSetMetadata(name, version)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Problem fetching metadata"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Metadata": entry})
}

func (m MetadataServer) GetFeatureMetadata(c *gin.Context) {

	name := c.Param("name")
	version := c.Param("version")

	entry, err := m.Server.Metadata.FeatureMetadata(name, version)
	if err != nil {
		c.JSON(500, gin.H{"Error": "Problem fetching metadata"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Metadata": entry})
}

func (m MetadataServer) ExposePort(port string) {
	router := gin.Default()
	router.GET("/training_set/:name/:version", m.GetTrainingMetadata)
	router.GET("/feature/:name/:version", m.GetFeatureMetadata)
	router.Run(port)
}