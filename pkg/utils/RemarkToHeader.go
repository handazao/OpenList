package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// RemarkData 用于解析 Remark 字段
type RemarkData struct {
	Headers map[string]string `json:"headers"`
	Params  map[string]string `json:"params"`
}

// RemarkToHeader 将 Remark 转成 http.Header
func RemarkToHeader(remarkStr string) http.Header {
	if remarkStr == "" {
		return nil
	}
	var remark RemarkData
	if err := json.Unmarshal([]byte(remarkStr), &remark); err != nil {
		log.Printf("failed to parse remark: %v", err)
		return nil
	}

	h := http.Header{}

	// 添加 headers
	for k, v := range remark.Headers {
		h.Set(k, v)
	}

	fmt.Println("Driver Remark Header:", h)
	return h
}

// RemarkToParam 将 Remark 转成 map[string]string
func RemarkToParam(remarkStr string) map[string]string {
	if remarkStr == "" {
		return nil
	}
	var remark RemarkData
	if err := json.Unmarshal([]byte(remarkStr), &remark); err != nil {
		log.Printf("failed to parse remark: %v", err)
		return nil
	}

	params := make(map[string]string)
	// 添加 params
	for k, v := range remark.Params {
		params[k] = v
	}

	fmt.Println("Driver Remark Params:", params)
	return params
}
