package handler

import (
	"net/http"
	"time"

	model "dumper/dao"
	"dumper/svc"
	"dumper/utils"
)

func SaveVisitLog(ctx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//todo 使用百度地图获取位置信息
		result, err := model.ManagerDao.SaveVisitLog(r.Context(), &model.VisitLogs{
			Url:         "https://www.baidu.com",
			Ip:          utils.ClientIP(r),
			Address:     "beijing",
			Point:       "sss",
			CreatedTime: time.Now(),
		})

		if err != nil {
			utils.OkJson(w, err)
		} else {
			utils.OkJson(w, result)
		}
	}
}
