package common

type Headers struct { // 上传录像http请求中，需要放入http头部
	Authorization string `json:"Authorization"`
	Date          string `json:"Date"`
	ContentType   string `json:"Content-Type"`
}
