package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/opensourceways/community-robot-lib/interrupts"
	"github.com/sirupsen/logrus"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"github.com/opensourceways/xihe-training-center/app"
	"github.com/opensourceways/xihe-training-center/controller"
)

type Service struct {
	Log *logrus.Entry

	Port    int
	Timeout time.Duration

	Training app.TrainingService
}

func StartWebServer(spec *swag.Spec, service *Service) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logRequest())

	setRouter(r, spec, service)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", service.Port),
		Handler: r,
	}

	defer interrupts.WaitForGracefulShutdown()

	interrupts.ListenAndServe(srv, service.Timeout)
}

//setRouter init router
func setRouter(engine *gin.Engine, spec *swag.Spec, service *Service) {
	spec.BasePath = "/api"
	spec.Title = "xihe-training-center"
	spec.Description = "APIs of xihe training center"

	v1 := engine.Group(spec.BasePath)
	{
		controller.AddRouterForTrainingController(
			v1,
			service.Training,
		)
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}

func logRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		endTime := time.Now()

		logrus.Infof(
			"| %d | %d | %s | %s |",
			c.Writer.Status(),
			endTime.Sub(startTime),
			c.Request.Method,
			c.Request.RequestURI,
		)
	}
}
