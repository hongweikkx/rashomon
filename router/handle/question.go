package handle

import (
	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/service"
)

func AnswerList(c *gin.Context) {
	ids, err := service.RedisClient.LRange("answers", 0, 10).Result()
	if err != nil {
		log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
		c.JSON(500, err)
		return
	}
	ress := []map[string]string{}
	for _, id := range ids {
		m, err := service.RedisClient.HGetAll(id).Result()
		if err != nil {
			continue
		}
		ress = append(ress, m)
	}
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"answers": ress,
		},
	})
}
