package tools

import (
	"errors"
	"github.com/onebids/onecommon/errno"
	"github.com/onebids/onecommon/model"
)

func BuildBaseResp(err error) *model.BaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *model.BaseResponse {
	return &model.BaseResponse{
		StatusCode: err.ErrCode,
		StatusMsg:  err.ErrMsg,
	}
}

func BuildBaseRespSuccess(msg string) *model.BaseResponse {
	return &model.BaseResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  msg,
	}
}

func BuildBaseRespSuccessNoParams() *model.BaseResponse {
	return &model.BaseResponse{
		StatusCode: errno.Success.ErrCode,
		StatusMsg:  "Success",
	}
}

func BuildBaseRespFailNoParams() *model.BaseResponse {
	return &model.BaseResponse{
		StatusCode: errno.BadRequest.ErrCode,
		StatusMsg:  "Fail",
	}
}

func ParseBaseResp(resp *model.BaseResponse) error {
	if resp.StatusCode == errno.Success.ErrCode {
		return nil
	}
	return errno.NewErrNo(resp.StatusCode, resp.StatusMsg)
}
