package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sdll/brighttop/common"
	glog "sdll/brighttop/golog"
	"sdll/brighttop/util/cipherutil"
	"strconv"
	"strings"
	"time"
)

var NotFound = []byte("Page Not Found")

func main4() {
	//server := http.Server{
	//	Handler: http.HandlerFunc(handle),
	//	Addr:    ":9000",
	//}
	//log.Println("starting server")
	//log.Println(server.ListenAndServe())

	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/v1/api/file/encrypt", handleFileEncrypt)
	keyHttpServer := http.Server{
		Addr:         ":9000",
		Handler:      serveMux,
		ReadTimeout:  time.Second * 60,
		WriteTimeout: time.Second * 60,
	}
	log.Println("starting server")
	err := keyHttpServer.ListenAndServe()
	if err != nil {
		log.Println("start key server listen and serve failed.", err)
	}

}

func handleFileEncrypt(w http.ResponseWriter, req *http.Request) {
	// 获取文件流参数，加密秘钥，加密等级，上传uploadSignUrl地址
	uploadSignUrl := req.FormValue("uploadSignUrl")
	var headers common.Headers
	err := json.Unmarshal([]byte(req.FormValue("headers")), &headers)
	if err != nil {
		responseJson(w, common.GenerateResponseFailed(common.RspParamError))
		return
	}
	secretLevel := req.FormValue("secretLevel")
	key := req.FormValue("key")
	//glog.Infof("request params,secretLevel=%v, uploadSignUrl=%v ", secretLevel, uploadSignUrl)
	file, _, err := req.FormFile("file")
	defer file.Close()
	if err != nil || uploadSignUrl == "" || secretLevel == "" || key == "" {
		responseJson(w, common.GenerateResponseFailed(common.RspParamError))
		return
	}
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return
	}
	fileData := buf.Bytes()
	var encryptIV [16]byte // 初始化加密向量
	for i := 0; i < 16; i++ {
		encryptIV[i] = byte(rand.Intn(0xFF))
	}
	secretKey := []byte(key)

	// 开始加密文件流
	sl, err := strconv.ParseInt(secretLevel, 10, 0)
	if err != nil {
		glog.Infof("Error during conversion,secretLevel=%v", secretLevel)
		return
	}
	if sl == 2 {
		//return cipherutil.Chacha20Decrypt(fileBytes[64:], securityKey, iv)
		responseJson(w, common.GenerateResponseFailed(common.RspParamError))
		return
	}
	if sl == 4 {
		//return cipherutil.AesGCMDecryptNew(fileBytes[64:], securityKey, iv)
		responseJson(w, common.GenerateResponseFailed(common.RspParamError))
		return
	}
	if sl == 1 || sl == 5 {
		//return nil, errors.New(fmt.Sprintf("securityLevel not support,securityLevel=%v", sl))
		responseJson(w, common.GenerateResponseFailed(common.RspParamError))
		return
	}

	encryptData, err := cipherutil.AesCBCEncrypt(fileData, secretKey, encryptIV[0:])

	wbuf := &bytes.Buffer{}
	// version(4 bytes) + iv向量(16 bytes) + size(4 bytes) + reserve(40 bytes)
	reserve := make([]byte, 39)
	binary.Write(wbuf, binary.LittleEndian, uint32(0))                // 写入版本号
	binary.Write(wbuf, binary.LittleEndian, encryptIV[0:16])          // 写入加密向量
	binary.Write(wbuf, binary.LittleEndian, uint32(len(encryptData))) // 写入数据长度
	binary.Write(wbuf, binary.LittleEndian, byte(sl))                 // 写入安全等级security_level
	binary.Write(wbuf, binary.LittleEndian, reserve[:39])             // 写入39个字节保留字段
	binary.Write(wbuf, binary.LittleEndian, encryptData)              // 写入数据
	// 开始上传加密文件流到对象存储服务器
	request, err := http.NewRequest("PUT", uploadSignUrl, wbuf)
	if err != nil {
		return
	}
	// 添加头部，在添加头部之前先判断字段是否为空
	// 如果字段为空，说明该字段没有传递下来，不需要，就不用设置
	if headers.Authorization != "" {
		request.Header.Add("Authorization", headers.Authorization)
	}
	if headers.Date != "" {
		request.Header.Add("Date", headers.Date)
	}
	if headers.ContentType != "" {
		request.Header.Add("Content-Type", headers.ContentType)
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		responseJson(w, common.GenerateResponseSuccess("upload file fail, http do request fail"))
		return
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		responseJson(w, common.GenerateResponseSuccess("upload file fail, status code != 200"))
		//w.WriteHeader(response.StatusCode)
		return
	}
	// 返回上传成功或失败
	responseJson(w, common.GenerateResponseSuccess("upload file ok"))
	return
}

// 响应 json 数据
func responseJson(w http.ResponseWriter, responseDTO *common.ResponseDTO) {
	if responseDTO.Code != 0 {
		glog.Warningf("请求出错,code=%v,msg=%v", responseDTO.Code, responseDTO.Msg)
	}

	jsonResult, err := json.Marshal(responseDTO)
	if err != nil {
		glog.Error("返回数据转json失败", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprint(w, string(jsonResult))
	if err != nil {
		glog.Error("响应数据失败", err)
		return
	}
}

func handle(w http.ResponseWriter, r *http.Request) {

	log.Println("path", r.URL.Path)

	switch r.URL.Path {
	case "/receive-file":
		upload(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write(NotFound)
	}

}

func upload(w http.ResponseWriter, r *http.Request) {

	//If its not multipart, We will expect file data in body.
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		handleFileInBody(w, r)
		return
	}

	handleFileInForm(w, r)

}

func handleFileInBody(w http.ResponseWriter, r *http.Request) {

	//Check if body if empty or not
	if r.ContentLength <= 0 {
		w.WriteHeader(400)
		w.Write([]byte("Got content length <= 0"))
		return
	}

	f, err := getFile("")
	if err != nil {
		somethingWentWrong(w)
		return
	}
	defer f.Close()

	written, err := io.Copy(f, r.Body)
	if err != nil {
		log.Println("copy error", err)
		somethingWentWrong(w)
		return
	}

	success(w)

	log.Println("Written", written)
}

func handleFileInForm(w http.ResponseWriter, r *http.Request) {
	uploadSignUrl := r.FormValue("uploadSignUrl")
	secretLevel := r.FormValue("secretLevel")
	key := r.FormValue("key")
	fmt.Println(uploadSignUrl, secretLevel, key)
	f, fh, err := r.FormFile("file")
	if err != nil {
		log.Println("formfile error", err)
		somethingWentWrong(w)
		return
	}

	if fh.Size <= 0 {
		w.WriteHeader(400)
		w.Write([]byte("Got File length <= 0"))
		return
	}

	outFile, err := getFile(fh.Filename)
	if err != nil {
		log.Println("getFile error", err)
		somethingWentWrong(w)
		return
	}

	written, err := io.Copy(outFile, f)
	if err != nil {
		log.Println("copy error", err)
		somethingWentWrong(w)
		return
	}

	success(w)
	log.Println("Written", written)
}

func getFile(fname string) (*os.File, error) {
	var fileName string

	now := time.Now()
	if fname != "" {
		fileName = strconv.Itoa(int(now.Unix())) + "_" + fname
	} else {
		fileName = "temp_" + strconv.Itoa(int(now.Unix())) + ".txt"
	}

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("create file error", err)
		return nil, err
	}

	return f, nil

}

func success(w http.ResponseWriter) {
	w.WriteHeader(200)
}

func somethingWentWrong(w http.ResponseWriter) {
	w.WriteHeader(500)
	w.Write([]byte("something went wrong"))
}
