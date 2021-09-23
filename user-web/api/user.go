package api

import (
	"context"
	"fmt"
	"mxshop-api/user-web/global/response"
	"mxshop-api/user-web/proto"
	"net/http"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func HandleGrpcErrorToHttp(err error, c *gin.Context) {
	// 将grpc的code 转换成http的状态码
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其它错误" + e.Message(),
				})
			}
			return
		}
	}
}

func GetUserList(c *gin.Context) {
	zap.S().Debug("获取用户列表页")
	ip := "127.0.0.1"
	port := 50051
	// 拨号连接用户grpc服务器
	userConn, err := grpc.Dial(fmt.Sprintf("%s:%d", ip, port), grpc.WithInsecure())
	if err != nil {
		zap.S().Errorw("[GetUserList] 连接 【服务失败】",
			"msg", err.Error(),
		)
	}

	// 生成grpc的client并调用接口
	userSrvClient := proto.NewUserClient(userConn)
	rsp, err := userSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    0,
		PSize: 0,
	})
	if err != nil {
		zap.S().Errorw("[GetUserList] 查询用户列表失败")
		HandleGrpcErrorToHttp(err, c)
		return
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.Data {
		//data := make(map[string]interface{})
		user := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			//Birthday: time.Time(time.Unix(int64(value.BirthDay), 0)).Format("2006-01-02"),
			Birthday: response.JsonTime(time.Unix(int64(value.BirthDay), 0)),
			Gender:   value.Gender,
			Mobile:   value.Mobile,
		}
		//data["id"] = value.Id
		//data["name"] = value.NickName
		//data["birthday"] = value.BirthDay
		//data["gender"] = value.Gender
		//data["mobile"] = value.Mobile

		result = append(result, user)

	}
	c.JSON(http.StatusOK, result)

}
