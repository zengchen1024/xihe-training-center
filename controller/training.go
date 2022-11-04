package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/opensourceways/xihe-training-center/app"
)

func AddRouterForTrainingController(
	rg *gin.RouterGroup,
	ts app.TrainingService,
) {
	ctl := TrainingController{ts: ts}

	rg.POST("/v1/training", ctl.Create)
	rg.DELETE("/v1/training/:id", ctl.Delete)
	rg.PUT("/v1/training/:id", ctl.Terminate)
	rg.GET("/v1/training/:id/log", ctl.GetLog)
	rg.GET("/v1/training/:id/result/:file", ctl.GetDownloadURL)
}

type TrainingController struct {
	baseController

	ts app.TrainingService
}

// @Summary Create
// @Description create training
// @Tags  Training
// @Param	body	body 	TrainingCreateRequest	true	"body of creating training"
// @Accept json
// @Success 201 {object} app.TrainingInfoDTO
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

	v, err := ctl.ts.Create(&cmd)
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusCreated, newResponseData(v))
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
	if err := ctl.ts.Delete(jobId); err != nil {
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
	if err := ctl.ts.Terminate(jobId); err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusAccepted, newResponseData("success"))
}

// @Summary GetLog
// @Description get log url of training for downloading
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Accept json
// @Success 200 {object} TrainingResultResp
// @Failure 500 system_error        system error
// @Router /v1/training/{id}/log [get]
func (ctl *TrainingController) GetLog(ctx *gin.Context) {
	v, err := ctl.ts.GetLogDownloadURL(ctx.Param("id"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(TrainingResultResp{v}))
}

// @Summary GetDownloadURL
// @Description get download url of training result such as log or output.
// @Tags  Training
// @Param	id	path	string	true	"id of training"
// @Param	file	path	string	true	"obs file path to download"
// @Accept json
// @Success 200 {object} TrainingResultResp
// @Failure 500 system_error        system error
// @Router /v1/training/{id}/result/{file} [get]
func (ctl *TrainingController) GetDownloadURL(ctx *gin.Context) {
	v, err := ctl.ts.GenFileDownloadURL(ctx.Param("file"))
	if err != nil {
		ctl.sendRespWithInternalError(ctx, newResponseError(err))

		return
	}

	ctx.JSON(http.StatusOK, newResponseData(TrainingResultResp{v}))
}
