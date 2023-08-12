package middleware

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Limit(engine *gin.Engine) {
	// 务必先进行初始化
	err := sentinel.InitDefault()
	if err != nil {
		log.Println(err)
	}

	//配置限流规则
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "some-test", //熔断器规则生效的埋点资源的名称
			TokenCalculateStrategy: flow.Direct, //Direct表示直接使用字段
			ControlBehavior:        flow.Reject, //匀速通过 Reject表示超过阈值直接拒绝
			Threshold:              300,         //QPS流控阈值30个
			StatIntervalInMs:       10000,       //QPS时间长度1s
		},
	})

	if err != nil {
		log.Println(err)
	}

	engine.Use(func(c *gin.Context) {
		e, err := sentinel.Entry("some-test", sentinel.WithTrafficType(base.Inbound))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{

				"msg": "访问量过多！",
			})
		} else {
			e.Exit()
			c.Next()
		}
	})

}
