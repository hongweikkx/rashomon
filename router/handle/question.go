package handle

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hongweikkx/rashomon/log"
	dgd "github.com/hongweikkx/rashomon/router/middleware/degrade"
	"github.com/hongweikkx/rashomon/service"
)

func AnswerList(c *gin.Context) {
	dgd.Fuse(c, nil, func(ctx context.Context) (int, interface{}) {
		ids, err := service.RedisClient.SRandMemberN(ctx, "answers", 10).Result()
		if err != nil {
			log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
			return 500, ""
		}
		ress := []map[string]string{}
		for _, id := range ids {
			m, err := service.RedisClient.HGetAll(ctx, AnswerKey(id)).Result()
			if err != nil {
				log.SugarLogger.Errorf("err:%+v", err)
				continue
			}
			ress = append(ress, m)
		}
		return 200, map[string]interface{}{
			"code": 20000,
			"data": map[string]interface{}{
				"answers": ress,
			},
		}
	})
}

func Like(c *gin.Context) {
	dgd.Fuse(c, nil, func(ctx context.Context) (int, interface{}) {
		return 200, map[string]interface{}{
			"code": 20000,
			"data": map[string]interface{}{},
		}
	})
}

func Search(c *gin.Context) {
	dgd.Fuse(c, nil, func(ctx context.Context) (int, interface{}) {
		ids, err := service.RedisClient.SRandMemberN(ctx, "answers", 10).Result()
		if err != nil {
			log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
			return 500, ""
		}
		ress := []map[string]string{}
		for _, id := range ids {
			m, err := service.RedisClient.HGetAll(ctx, AnswerKey(id)).Result()
			if err != nil {
				continue
			}
			ress = append(ress, m)
		}
		return 200, map[string]interface{}{
			"code": 20000,
			"data": map[string]interface{}{
				"likelist": ress,
			},
		}
	})
}

func LikeList(c *gin.Context) {
	dgd.Fuse(c, nil, func(ctx context.Context) (int, interface{}) {
		ids, err := service.RedisClient.SRandMemberN(ctx, "answers", 10).Result()
		if err != nil {
			log.SugarLogger.Errorf("answerlist err:%+v", err.Error())
			return 500, ""
		}
		ress := []map[string]string{}
		for _, id := range ids {
			m, err := service.RedisClient.HGetAll(ctx, AnswerKey(id)).Result()
			if err != nil {
				continue
			}
			m["isLike"] = "1"
			ress = append(ress, m)
		}
		return 200, map[string]interface{}{
			"code": 20000,
			"data": map[string]interface{}{
				"likelist": ress,
			},
		}
	})
}

func AnswerKey(id string) string {
	return fmt.Sprintf("answer:zhihu:%s", id)
}
