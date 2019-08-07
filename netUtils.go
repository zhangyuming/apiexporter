package main

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type resultCallBack func(result *Result)

func httpRun(worker *httpWorker, back resultCallBack)  {
	log.Debug("http worker invoke",worker)
	result := Result{}
	result.Name = worker.Name
	if worker.Method == "" {
		worker.Method = "GET"
	}
	req,err := http.NewRequest(worker.Method,worker.Url,bytes.NewReader([]byte(worker.Body)))
	if err != nil {
		log.Warn("new http request error",worker, err)
		return
	}
	for _,h := range worker.Header {
		header := &req.Header
		header.Add(h.Key,h.Value)
	}

	start := time.Now()
	client := &http.Client{}
	resp,err := client.Do(req)
	rest := time.Since(start).Seconds()
	if err != nil {
		log.Warn("send http error",worker, err)
		return
	}
	result.Status = resp.StatusCode

	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil{
		log.Warn("read response filed",worker, err)
		result.Body = err.Error()
	}else{
		result.Body = string(bt)
	}
	result.ResponseTime = rest
	worker.checkResult(&result)
	back(&result)

}

func tcpRun(worker *tcpWorker, back resultCallBack){

	result := Result{}
	result.Name = worker.Name

	start := time.Now()
	ad :=  fmt.Sprint(worker.Host,":",worker.Port)
	conn, err := net.DialTimeout("tcp", ad,time.Second*20)

	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	responseTime := time.Since(start).Seconds()

	if err != nil {
		result.Status = 0
		result.Body = err.Error()
		result.ResponseTime = responseTime
	}else{
		result.Status = 1
		result.ResponseTime = responseTime
	}
	worker.checkResult(&result)
	back(&result)
	//fmt.Println(result)
}