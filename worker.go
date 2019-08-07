package main

import (
	"github.com/muesli/cache2go"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"regexp"
	"time"
)

var cache = cache2go.Cache("api_exporter_cache")

type workerBean struct {
	Http []httpWorker `yaml:http`
	Tcp []tcpWorker `yaml:tcp`

}

type Result struct {
	Name string
	Status int
	ResponseTime float64
	Success int
	Body string
}


type httpWorker struct {
	Name string `yaml:"name"`
	Interval int `yaml:"interval"`
	Url string `yaml:"url"`
	Method string `yaml:"method"`
	Header []struct{
		Key string `yaml:"key"`
		Value string `yaml:"value"`
	} `yaml:"header"`
	Body string `yaml:"body"`
	CheckStr string `yaml:"check_str"`
}

func (h *httpWorker) checkResult( result *Result)  {
	if result.Status < 400 {
		result.Success = 1
	}else{
		result.Success = 0
	}
	if h.CheckStr != "" {
		b,err := regexp.MatchString(h.CheckStr,result.Body)
		if err != nil{
			log.Error("regexp failed ",result.Body,h.CheckStr)
			result.Success = 0
		}else if b {
			result.Success = 1
		}else{
			result.Success = 0
		}
	}
}

type tcpWorker struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Interval int `yaml:"interval"`
}

func (t *tcpWorker) checkResult(result *Result)  {
	result.Success = result.Status
}

func CallBack(result *Result){

	if result != nil{
		cache.Add(result.Name,time.Minute*60,result)
	}else{
		log.Warn("save cache failed, result is null")
	}


}



func initWorker()(workerBean,error)  {
  w,err := loadConfig(configFile)
  if err != nil {
  	log.Fatal("parse config file error",err)
  	panic("parse config file error")
  }
  for _,h := range w.Http {
		if h.Interval == 0 {
			h.Interval = workerInterval
		}
		go func(worker httpWorker) {
			httpRun(&worker, CallBack)
			for range time.Tick(time.Duration(worker.Interval) * time.Second ) {
				httpRun(&worker, CallBack)
			}
		}(h)
  }
  for _,t := range w.Tcp{
	  if t.Interval == 0 {
		  t.Interval = workerInterval
	  }
  		go func(tcpw tcpWorker) {
			tcpRun(&tcpw, CallBack)
			for range time.Tick(time.Duration(tcpw.Interval) * time.Second ) {
				tcpRun(&tcpw, CallBack)
			}
		}(t)
  }
  return w,nil

}

func loadConfig(path string)(workerBean,error){
	w := workerBean{}

	data,err:= ioutil.ReadFile(path)
	//fmt.Print(string(data))
	if err != nil{
		log.Error("no found config file ",path,err)
		panic("config file no found")
	}

	err = yaml.Unmarshal(data,&w)
	if err != nil {
		log.Fatal("parse yml config file error", path,err)
		panic("parse yml config file error")
	}
	return w,nil

}