package tools

import (
	"errors"
	"github.com/onebids/onecommon/base"
	"github.com/onebids/onecommon/consts/errno"
)

func BuildBaseResp(err error) *base.BaseResponse {
	if err == nil {
		return baseResp(errno.Success)
	}
	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(&e)
	}
	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(&s)
}

func baseResp(err *errno.ErrNo) *base.BaseResponse {
	return &base.BaseResponse{
		Code: err.ErrCode,
		Msg:  err.ErrMsg,
	}
}

func BuildBaseRespSuccess(msg string) *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.Success.ErrCode,
		Msg:  msg,
	}
}

func BuildBaseRespSuccessNoParams() *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.Success.ErrCode,
		Msg:  "Success",
	}
}

func BuildBaseRespFailNoParams() *base.BaseResponse {
	return &base.BaseResponse{
		Code: errno.BadRequest.ErrCode,
		Msg:  "Fail",
	}
}

//func ParseBaseResp(resp *base.BaseResponse) error {
//	if resp.StatusCode == errno.Success.ErrCode {
//		return nil
//	}
//	return errno.NewErrNo(resp.StatusCode, resp.StatusMsg)
//}
