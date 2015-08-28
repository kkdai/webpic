package main

type WebSite struct {
	WebSite        string `json:"web_site"`
	TitlePattern   string `json:"title_pattern"`
	ImgPattern     string `json:"img_pattern"`
	ImgAttrPattern string `json:"img_attr_pattern"`
	ForceBig5      bool   `json:"forceBig5"`
}

type Parser struct {
	SupportSites []WebSite `json:"support_web_sites"`
}

var DefaultJson = []byte(`
{
	"support_web_sites" : [

		{
			"web_site": "pixnet.net",
			"title_pattern": "title",
			"img_pattern" : ".article-content-inner p img",
			"img_attr_pattern": "src"
		},
		{
			"web_site": "ck101.com",
			"title_pattern": "h1#thread_subject",
			"img_pattern" : "div[itemprop=articleBody] img",
			"img_attr_pattern": "file"
		},
		{
			"web_site": "timliao.com",
			"title_pattern": "title",
			"img_pattern" : ".bodycontent img",
			"img_attr_pattern": "src",
			"forceBig5" : true
		},
		{
			"web_site": "gigacircle.com",
			"title_pattern": "title",
			"img_pattern" : ".usercontent p img",
			"img_attr_pattern": "src"
		}
	]
}`)
