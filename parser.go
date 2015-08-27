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
