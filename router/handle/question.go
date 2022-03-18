package handle

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	"github.com/hongweikkx/rashomon/service"
)

func AnswerList(c *gin.Context) {
	ids, err := service.RedisClient.SRandMemberN("answers", 10).Result()
	if err != nil {
		log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
		c.JSON(500, err)
		return
	}
	ress := []map[string]string{}
	for _, id := range ids {
		m, err := service.RedisClient.HGetAll(AnswerKey(id)).Result()
		if err != nil {
			log.SugarLogger.Errorf("err:%+v", err)
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

func Like(c *gin.Context) {
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{},
	})
}

func Search(c *gin.Context) {
	ids, err := service.RedisClient.SRandMemberN("answers", 10).Result()
	if err != nil {
		log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
		c.JSON(500, err)
		return
	}
	ress := []map[string]string{}
	for _, id := range ids {
		m, err := service.RedisClient.HGetAll(AnswerKey(id)).Result()
		if err != nil {
			continue
		}
		ress = append(ress, m)
	}
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"likelist": ress,
		},
	})
}

func LikeList(c *gin.Context) {
	ids, err := service.RedisClient.SRandMemberN("answers", 10).Result()
	if err != nil {
		log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
		c.JSON(500, err)
		return
	}
	ress := []map[string]string{}
	for _, id := range ids {
		m, err := service.RedisClient.HGetAll(AnswerKey(id)).Result()
		if err != nil {
			continue
		}
		m["isLike"] = "1"
		ress = append(ress, m)
	}
	c.JSON(200, map[string]interface{}{
		"code": 20000,
		"data": map[string]interface{}{
			"likelist": ress,
		},
	})
}

func AnswerKey(id string) string {
	return fmt.Sprintf("answer:zhihu:%s", id)
}
