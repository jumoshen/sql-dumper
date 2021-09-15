package handler

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	model "dumper/dao"
	"dumper/svc"
	"dumper/utils"
)

func SaveVisitLog(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//todo 使用百度地图获取位置信息
		fmt.Println(r.RequestURI)
		fmt.Printf("header:%#v\n", r.Header)

		//todo https://api.map.baidu.com/location/ip?ak=vaXwBb4abtGvukEOOLA19Qltx3Ndua4c&ip=39.155.169.2&coor=bd09ll
		_, err := model.ManagerDao.SaveVisitLog(r.Context(), &model.VisitLogs{
			Url:         "https://www.baidu.com",
			Ip:          utils.ClientIP(r),
			Address:     "beijing",
			Point:       "sss",
			CreatedTime: time.Now(),
		})

		if err != nil {
			utils.OkJson(w, err)
		} else {
			_ = proxy(w, r)
		}
	}
}

func proxy(rw http.ResponseWriter, req *http.Request) error {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	transport := http.DefaultTransport

	outReq := new(http.Request)
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
		return err
	}

	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}

	rw.WriteHeader(res.StatusCode)
	_, err = io.Copy(rw, res.Body)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return err
	}
}
