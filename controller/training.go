package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-training-center/app"
	"github.com/opensourceways/xihe-training-center/domain/training"
)

func AddRouterForTrainingController(
	rg *gin.RouterGroup,
	ts training.Training,
) {
	ctl := TrainingController{
		s: app.NewTrainingService(ts),
	}

	rg.POST("/v1/training", ctl.Create)
	rg.DELETE("/v1/training/:id", ctl.Delete)
	rg.GET("/v1/training/:id", ctl.Get)
	rg.PUT("/v1/training/:id", ctl.Terminate)
	rg.GET("/v1/training/:id/log", ctl.GetLog)
}

type TrainingController struct {
	baseController

	s app.TrainingService
}

// @Summary Create
// @Description create training
// @Tags  Training
// @Param	body	body 	TrainingCreateRequest	true	"body of creating training"
// @Accept json
// @Success 201 {object} TrainingCreateResp
// @Failure 400 bad_request_body    can't parse request body
// @Failure 401 bad_request_param   some parameter of body is invalid
// @Failure 500 system_error        system error
// @Router /v1/training [post]
func (ctl *TrainingController) Create(ctx *gin.Context) {
	req := TrainingCreateRequest{}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, respBadRequestBody)

		return
	}

	cmd, err := req.toCmd()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, newResponseCodeError(
			errorBadRequestParam, err,
		))
		return
	}

	jobId, err := ctl.s.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(TrainingCreateResp{jobId}))
}

// @Summary Delete
// @Description delete training
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Accept json
// @Success 204
// @Failure 500 system_error        system error
// @Router /v1/training/{id} [delete]
func (ctl *TrainingController) Delete(ctx *gin.Context) {
	jobId := ctx.Param("id")
	if err := ctl.s.Delete(jobId); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusNoContent, newResponseData("success"))
}

// @Summary Terminate
// @Description terminate training
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Accept json
// @Success 202
// @Failure 500 system_error        system error
// @Router /v1/training/{id} [put]
func (ctl *TrainingController) Terminate(ctx *gin.Context) {
	jobId := ctx.Param("id")
	if err := ctl.s.Terminate(jobId); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary Get
// @Description get training info
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Accept json
// @Success 200 {object} app.TrainingDTO
// @Failure 500 system_error        system error
// @Router /v1/training/{id} [get]
func (ctl *TrainingController) Get(ctx *gin.Context) {
	v, err := ctl.s.Get(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(v))
}

// @Summary GetLog
// @Description get log url of training for downloading
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Accept json
// @Success 200 {object} TrainingLogResp
// @Failure 500 system_error        system error
// @Router /v1/training/{id}/log [get]
func (ctl *TrainingController) GetLog(ctx *gin.Context) {
	v, err := ctl.s.GetLogURL(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData(TrainingLogResp{v}))
}
