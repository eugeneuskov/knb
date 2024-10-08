package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"knb/app/repositories"
	"knb/app/services"
	"knb/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testEnvFilePath = "../../.env.test"
	testRoutesMode  = "debug"
)

type responseError struct {
	Message string `json:"message"`
}

type expectedError struct {
	code    int
	message string
}

type testRequestContext struct {
	key   string
	value interface{}
}

type testRequestHeader struct {
	key   string
	value string
}

type testAppLayers struct {
	bootstrap  *tests.BootstrapTest
	db         *gorm.DB
	repository *repositories.Repository
	service    *services.Service
	router     *gin.Engine
}

type requestData struct {
	router      *gin.Engine
	context     []*testRequestContext
	headers     []*testRequestHeader
	requestBody []byte
	method      string
	url         string
}

func preparationForTest(t *testing.T) testAppLayers {
	bootstrapTest := tests.NewBootstrapTest(testEnvFilePath)
	if err := bootstrapTest.SetupTestDB(); err != nil {
		t.Errorf("Failed to setup test DB, %s", err)
	}

	db := bootstrapTest.DB()
	repositoryMap := repositories.NewRepository(db)
	serviceMap := services.NewService(repositoryMap, bootstrapTest.Config())

	return testAppLayers{
		bootstrap:  bootstrapTest,
		db:         db,
		repository: repositoryMap,
		service:    serviceMap,
		router:     NewHandler(serviceMap).InitRoutes(testRoutesMode),
	}
}

func sendRequestAndGetResponse(data requestData) ([]byte, int) {
	req, _ := http.NewRequest(data.method, data.url, getRequestBody(data.requestBody))
	for _, header := range data.headers {
		req.Header.Set(header.key, header.value)
	}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	for _, ctx := range data.context {
		c.Set(ctx.key, ctx.value)
	}

	data.router.ServeHTTP(w, req)

	return w.Body.Bytes(), w.Code
}

func getRequestBody(requestBody []byte) *bytes.Buffer {
	switch requestBody {
	case nil:
		return bytes.NewBuffer(nil)
	default:
		return bytes.NewBuffer(requestBody)
	}
}
