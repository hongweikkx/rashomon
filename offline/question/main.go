package main

type Answer struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	OriginUrl string `json:"originUrl"`
}

func main() {
	InitData()
}

func InitData() {
	err := ZHIns.InitData()
	if err != nil {
		panic(err)
	}
}

func (ans *Answer) Format() map[string]interface{} {
	return map[string]interface{}{
		"id":        ans.Id,
		"title":     ans.Title,
		"content":   ans.Content,
		"originUrl": ans.OriginUrl,
	}
}
