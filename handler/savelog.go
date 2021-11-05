package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
	"net/url"

	model "dumper/dao"
	"dumper/svc"
	"dumper/utils"
)

func SaveVisitLog(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//todo 使用百度地图获取位置信息
		fmt.Println(r.RequestURI)
		fmt.Printf("header:%#v\n", r.Header)
	}
}

type Pxy struct {}

func (p *Pxy)ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	//todo https://api.map.baidu.com/location/ip?ak=vaXwBb4abtGvukEOOLA19Qltx3Ndua4c&ip=39.155.169.2&coor=bd09ll
	_, err := model.ManagerDao.SaveVisitLog(req.Context(), &model.VisitLogs{
		Url:         "https://www.baidu.com",
		Ip:          utils.ClientIP(req),
		Address:     "beijing",
		Point:       "sss",
		CreatedTime: time.Now(),
	})

	transport := http.DefaultTransport

	outReq := new(http.Request)
outReq.URL = &url.URL{
		Host: "https:www.jumoshen.cn",
		Path: req.RequestURI,
	}
	// this only does shallow copies of maps
	*outReq = *req

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
	}

	res, err := transport.RoundTrip(outReq)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, res.Body)

	res.Body.Close()
}
