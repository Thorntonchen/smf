/*
 * Nsmf_PDUSession
 *
 * SMF PDU Session Service
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pdusession

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/free5gc/openapi"
	"github.com/free5gc/openapi/models"
	"github.com/free5gc/smf/internal/logger"
	"github.com/free5gc/smf/internal/sbi/producer"
	"github.com/free5gc/util/httpwrapper"
)

// HTTPPostSmContexts - Create SM Context
func HTTPPostSmContexts(c *gin.Context) {
	auth_err := authorizationCheck(c)
	if auth_err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": auth_err.Error()})
		return
	}
	
	logger.PduSessLog.Info("Receive Create SM Context Request")
	var request models.PostSmContextsRequest

	request.JsonData = new(models.SmContextCreateData)

	s := strings.Split(c.GetHeader("Content-Type"), ";")
	var err error
	switch s[0] {
	case "application/json":
		err = c.ShouldBindJSON(request.JsonData)
	case "multipart/related":
		err = c.ShouldBindWith(&request, openapi.MultipartRelatedBinding{})
	}

	if err != nil {
		problemDetail := "[Request Body] " + err.Error()
		rsp := models.ProblemDetails{
			Title:  "Malformed request syntax",
			Status: http.StatusBadRequest,
			Detail: problemDetail,
		}
		logger.PduSessLog.Errorln(problemDetail)
		c.JSON(http.StatusBadRequest, rsp)
		return
	}

	req := httpwrapper.NewRequest(c.Request, request)
	isDone := c.Done()
	HTTPResponse := producer.HandlePDUSessionSMContextCreate(isDone,
		req.Body.(models.PostSmContextsRequest))
	// Http Response to AMF
	for key, val := range HTTPResponse.Header {
		c.Header(key, val[0])
	}
	switch HTTPResponse.Status {
	case http.StatusCreated,
		http.StatusBadRequest,
		http.StatusForbidden,
		http.StatusNotFound,
		http.StatusInternalServerError,
		http.StatusServiceUnavailable,
		http.StatusGatewayTimeout:
		c.Render(HTTPResponse.Status, openapi.MultipartRelatedRender{Data: HTTPResponse.Body})
	default:
		c.JSON(HTTPResponse.Status, HTTPResponse.Body)
	}
}
