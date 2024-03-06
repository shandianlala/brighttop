package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sdll/brighttop/glog"
	"time"
)

type (
	StrangeTokenResp struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data string `json:"data"`
	}
)

var downloadTransport = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}
var strangeHttpClient = http.Client{
	Transport: downloadTransport,
	Timeout:   10 * time.Second,
}

func mainStange() {

	var urlPrefix = ""

	var getTokenUrl = urlPrefix + "/account/verify"
	var contentType = "application/json"
	var strangeAppSecurity = ""
	var applicationName = "sagittarius"

	var getTokenJsonStr = []byte(`{ "name":"` + applicationName + `","security":"` + strangeAppSecurity + `"}`)
	req, _ := http.NewRequest("POST", getTokenUrl, bytes.NewBuffer(getTokenJsonStr))
	req.Header.Set("Content-Type", contentType)
	resp, err := strangeHttpClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	getTokenResult := string(body)
	fmt.Println("response Body:", getTokenResult)
	var tokenResp StrangeTokenResp
	err = json.Unmarshal(body, &tokenResp)
	glog.Infof("响应code=%d, result=%s, ", tokenResp.Code, getTokenResult)
	if err != nil || tokenResp.Code != 200 {
		fmt.Println("response Body:", getTokenResult)
		return
	}

	var env = "default"
	var fetchAllUrl = urlPrefix + "/resource/fetchAllByToken?env=" + env
	fetchResourceReq, _ := http.NewRequest("GET", fetchAllUrl, nil)
	req.Header.Set("token", tokenResp.Data)
	response, err := strangeHttpClient.Do(fetchResourceReq)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	resourceBody, _ := io.ReadAll(resp.Body)
	fmt.Println("response Body:", string(resourceBody))

}
