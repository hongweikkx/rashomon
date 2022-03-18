package main

type Answer struct {
	Id        string
	Title     string
	Content   string
	OriginUrl string
	VoteUp    string
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
		"voteUp":    ans.VoteUp,
		"isLike":    "0",
	}
}
