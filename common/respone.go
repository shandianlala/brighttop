package common

//noinspection GoSnakeCaseUsage

type ErrorInfo interface {
	GetCode() int
	GetMsg() string
}

type DefaultErrorInfo struct {
	Code int
	Msg  string
}

func (info *DefaultErrorInfo) GetCode() int {
	return info.Code
}

func (info *DefaultErrorInfo) GetMsg() string {
	return info.Msg
}

//noinspection ALL
var (
	RspSuccess      = DefaultErrorInfo{0, "ok"}
	RspUnknownError = DefaultErrorInfo{1, "unknown error"}
	RspParamError   = DefaultErrorInfo{11, "invalid parameter"}
	RspServerError  = DefaultErrorInfo{20, "system error"}

	////////////////////// 具体的业务错误 ///////////////////////
	// 图片获取
	RspPicDncryptError  = DefaultErrorInfo{101, "failed to decrypt image"}
	RspPicDownloadError = DefaultErrorInfo{102, "failed to download image"}
	RspPicResizeError   = DefaultErrorInfo{104, "failed to resize image"}

	// 获取 key 失败
	RspKeyError        = DefaultErrorInfo{200, "failed to get key"}
	RspKeyDncryptError = DefaultErrorInfo{201, "invalid key"}

	//限流超时
	RspLimiterTimeout = DefaultErrorInfo{300, "requests are too frequent"}
)

// 返回 DTO
type ResponseDTO struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Rs   interface{} `json:"rs"`
}

func GenerateResponse(errInfo DefaultErrorInfo, rs interface{}) *ResponseDTO {
	return &ResponseDTO{Code: errInfo.GetCode(), Msg: errInfo.GetMsg(), Rs: rs}
}

func GenerateResponseSuccess(rs interface{}) *ResponseDTO {
	return GenerateResponse(RspSuccess, rs)
}

func GenerateResponseFailed(errInfo DefaultErrorInfo) *ResponseDTO {
	return GenerateResponse(errInfo, nil)
}
