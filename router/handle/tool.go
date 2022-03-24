package handle

import (
	"context"

	"github.com/gin-gonic/gin"
	dgd "github.com/hongweikkx/rashomon/router/middleware/degrade"
)

type ToolS struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Desc string `json:"desc"`
}

var DownloadTOOLS = []ToolS{
	{
		Name: "抖音下载",
		Url:  "https://douyin.wtf",
		Desc: "抖音/TikTok无水印在线解析",
	},
}

func ToolList(c *gin.Context) {
	dgd.Fuse(c, nil, func(ctx context.Context) (int, interface{}) {
		return 200, map[string]interface{}{
			"code": 20000,
			"data": map[string]interface{}{
				"downloads": DownloadTOOLS,
			},
		}
	})
}
