package main

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-redis/redis"
	"gopkg.in/resty.v1"
)

var (
	ZHIns    ZH
	redisCli *redis.Client
)

type ZH struct{}

func init() {
	opt := redis.Options{
		Addr:         "localhost:6379",
		DialTimeout:  5 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
		PoolSize:     runtime.NumCPU() / 2,
		MinIdleConns: 1,
	}
	redisCli = redis.NewClient(&opt)
}

const USER_AGENT = "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Mobile Safari/537.36"

func (zh *ZH) R() *resty.Request {
	return resty.R().SetHeaders(zh.getHeader())
}

func (zh *ZH) getHeader() map[string]string {
	return map[string]string{
		"Accept":     "*/*",
		"Connection": "keep-alive",
		"Host":       "www.zhihu.com",
		"Origin":     "http://www.zhihu.com",
		"Pragma":     "no-cache",
		"User-Agent": USER_AGENT,
	}
}

func (zh *ZH) InitData() error {
	fmt.Println("zhihu init data...")
	for _, url := range ZHIHU_URLS {
		fmt.Printf("resolve:%s\n", url)
		ans, err := zh.GetAnswer(url)
		if err != nil {
			continue
		}
		redisCli.HMSet(zh.AnswerKey(ans.Id), ans.Format())
	}
	return nil
}

func (zh *ZH) AnswerKey(id string) string {
	return fmt.Sprintf("answer:zhihu:%s", id)
}

func (zh *ZH) GetAnswer(url string) (*Answer, error) {
	res, err := zh.R().Get(url)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		return nil, err
	}
	id, err := zh.GetId(url)
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(res.Body()))
	if err != nil {
		fmt.Printf("err:%+v\n", err)
		return nil, err
	}
	title := doc.Find("title").Contents().Text()
	content := doc.Find("p").Contents().Text()
	if title == "" || content == "" {
		return nil, errors.New("title or content is empty")
	}
	return &Answer{
		Id:        id,
		Title:     title,
		Content:   content,
		OriginUrl: url,
	}, nil
}

func (zh *ZH) GetId(url string) (string, error) {
	if url == "" {
		return "", errors.New("url is empty")
	}
	ss := strings.Split(url, "/")
	id := ss[len(ss)-1]
	_, err := strconv.Atoi(id)
	return id, err
}

var ZHIHU_URLS = []string{
	"https://www.zhihu.com/question/28626263/answer/41992632",
}
