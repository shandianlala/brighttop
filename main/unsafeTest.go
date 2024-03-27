package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"unsafe"
)

type H264Header struct { /* 24字节 */
	FrameType int32  /* 4字节 */
	Size      uint32 /* 4字节 */
	Timestamp uint64 /* 8字节 */
	Pts       uint64 /* 8字节 */
}

func main() {

	filePath := "/Users/caoti/Downloads/1711417746_0000.media" // 本地文件路径
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("无法读取文件:", err)
		return
	}
	//fmt.Printf("文件内容: %s\n", string(data))
	var index int = 0
	bodyLength := len(data)
	fmt.Println(bodyLength)
	num := int(binary.LittleEndian.Uint32(data[0:4]))
	fmt.Printf("num value is %d\n", num)
	for {
		//size := int(unsafe.Sizeof(H264Header{}))
		//fmt.Println(fmt.Sprintf("size is %d", size))

		//先判断剩余的长度是否够一个header的长度
		if index+int(unsafe.Sizeof(H264Header{})) > bodyLength {
			break
		}
		header := *(*H264Header)(unsafe.Pointer(&data[index:][0]))
		index += int(unsafe.Sizeof(H264Header{}))

		//判断剩余长度是否大于一帧的长度
		if index+int(header.Size) > bodyLength {
			break
		}
		//body := data[index : index+int(header.Size)]
		jsonData, _ := json.Marshal(header)
		fmt.Println(string(jsonData))
		index += int(header.Size)

	}

}
