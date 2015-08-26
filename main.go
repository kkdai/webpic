package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/cobra"
)

var (
	baseDir string
	imageId = regexp.MustCompile(`([^\/]+)\.(png|jpg)`)
)

var configSetting Parser
var targetSiteSetting WebSite

func init() {
	file, _ := ioutil.ReadFile("./parser.json")
	json.Unmarshal(file, &configSetting)
}

func worker(destDir string, linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		resp, err := http.Get(target)
		if err != nil {
			fmt.Printf("Http.Get\nerror: %s\ntarget: %s", err, target)
			continue
		}
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			fmt.Printf("image.Decode\nerror: %s\ntarget: %s", err, target)
			continue
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
			imgInfo := imageId.FindStringSubmatch(target)
			out, err := os.Create(destDir + "/" + imgInfo[1] + "." + imgInfo[2])
			if err != nil {
				fmt.Printf("os.Create\nerror: %s", err)
				continue
			}
			defer out.Close()
			switch imgInfo[2] {
			case "jpg":
				jpeg.Encode(out, m, nil)
			case "png":
				png.Encode(out, m)
			}
		}
	}
}

func findDomainByURL(url string) WebSite {
	//using matching first, should goes to regexp
	var targetDomain string
	urlRegex := regexp.MustCompile(`^((http[s]?|ftp):\/)?\/?([^:\/\s]+)((\/\w+)*\/)([\w\-\.]+[^#?\s]+)(.*)?(#[\w\-]+)?$`)
	tokenStrings := urlRegex.FindStringSubmatch(url)
	if len(tokenStrings) > 0 {
		targetDomain = tokenStrings[3]
	}

	for index, webside := range configSetting.SupportSites {
		if strings.Contains(targetDomain, webside.WebSite) {
			return configSetting.SupportSites[index]
		}
	}

	return configSetting.SupportSites[0]
}

func crawler(target string, workerNum int) {
	doc, err := goquery.NewDocument(target)
	if err != nil {
		panic(err)
	}

	//find web site from URL.
	targetSiteSetting := findDomainByURL(target)
	title := doc.Find(targetSiteSetting.TitlePattern).Text()
	fmt.Println("[", targetSiteSetting.WebSite, "]:", title, " starting downloading...")

	dir := fmt.Sprintf("%v/%v - %v", baseDir, targetSiteSetting.WebSite, title)

	os.MkdirAll(dir, 0755)

	linkChan := make(chan string)
	wg := new(sync.WaitGroup)
	for i := 0; i < workerNum; i++ {
		wg.Add(1)
		go worker(dir, linkChan, wg)
	}

	doc.Find(targetSiteSetting.ImgPattern).Each(func(i int, img *goquery.Selection) {
		imgUrl, _ := img.Attr(targetSiteSetting.ImgAttrPattern)
		linkChan <- imgUrl
	})

	close(linkChan)
	wg.Wait()
}

func main() {
	usr, _ := user.Current()
	baseDir = fmt.Sprintf("%v/Pictures/ilovedlimg", usr.HomeDir)

	var postUrl string
	var workerNum int

	rootCmd := &cobra.Command{
		Use:   "ilovedlimg",
		Short: "Download all the images in given post url",
		Run: func(cmd *cobra.Command, args []string) {
			crawler(postUrl, workerNum)
		},
	}
	rootCmd.Flags().StringVarP(&postUrl, "url", "u", "http://ck101.com/thread-2876990-1-1.html", "Url of post")
	rootCmd.Flags().IntVarP(&workerNum, "worker", "w", 25, "Number of workers")
	rootCmd.Execute()
}
