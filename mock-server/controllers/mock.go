package controllers

import (
	"encoding/json"
	"mayfly-go/base/biz"
	"mayfly-go/base/ctx"
	"mayfly-go/base/ginx"
	"mayfly-go/base/rediscli"
	"mayfly-go/base/utils"
	"mayfly-go/mock-server/controllers/form"
)

const key = "mock:data"

// @router /api/mock-datas/:method [get]
func GetMockData(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	method := g.Param("method")
	params := utils.MapBuilder("method", method).ToMap()
	// 该mock数据需要指定的生效用户才可访问
	username := g.Query("username")
	if username != "" {
		params["username"] = username
	}
	rc.ReqParam = params

	val := rediscli.HGet(key, method)
	mockData := &form.MockData{}
	json.Unmarshal([]byte(val), mockData)
	biz.IsTrue(mockData.Enable == 1, "无该mock数据")

	eu := mockData.EffectiveUser
	// 如果设置的生效用户为空，则表示所有用户都生效
	if len(eu) == 0 {
		rc.ResData = mockData.Data
		return
	}
	biz.IsTrue(utils.StrLen(username) != 0, "该用户无法访问该mock数据")
	for _, e := range eu {
		if username == e {
			rc.ResData = mockData.Data
			return
		}
	}
	panic(biz.NewBizErr("该用户无法访问该mock数据"))
}

// @router /api/mock-datas [put]
func UpdateMockData(rc *ctx.ReqCtx) {
	mockData := &form.MockData{}
	ginx.BindJsonAndValid(rc.GinCtx, mockData)
	rc.ReqParam = mockData.Method
	val, _ := json.Marshal(mockData)
	rediscli.HSet(key, mockData.Method, val)
}

// @router /api/mock-datas [post]
func CreateMockData(rc *ctx.ReqCtx) {
	mockData := &form.MockData{}
	ginx.BindJsonAndValid(rc.GinCtx, mockData)
	biz.IsTrue(!rediscli.HExist(key, mockData.Method), "该方法已存在")
	val, _ := json.Marshal(mockData)
	rediscli.HSet(key, mockData.Method, val)
}

// @router /api/mock-datas [get]
func GetAllData(rc *ctx.ReqCtx) {
	rc.ResData = rediscli.HGetAll(key)
}

// @router /api/mock-datas/:method [delete]
func DeleteMockData(rc *ctx.ReqCtx) {
	method := rc.GinCtx.Param("method")
	rc.ReqParam = method
	rediscli.HDel(key, method)
}