package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"
	"strings"
	"sync"
	"time"

	"code.google.com/p/go.text/encoding/traditionalchinese"
	"code.google.com/p/go.text/transform"

	"github.com/PuerkitoBio/goquery"
	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

var (
	baseDir        string
	imageId        = regexp.MustCompile(`([^\/]+)\.(png|jpg)`)
	urlRegex       = regexp.MustCompile(`^((http[s]?|ftp):\/)?\/?([^:\/\s]+)((\/\w+)*\/)([\w\-\.]+[^#?\s]+)(.*)?(#[\w\-]+)?$`)
	jsonFileGithub = "https://raw.githubusercontent.com/kkdai/webpic/master/parser.json"
)

var configSetting Parser
var targetSiteSetting WebSite

func worker(destDir string, linkChan chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for target := range linkChan {
		imgInfo := imageId.FindStringSubmatch(target)
		if (len(imgInfo) > 0 && strings.EqualFold(imgInfo[2], "gif")) || strings.Contains(target, ".gif") {
			//GIF not support for now, skip
			continue
		}

		resp, err := http.Get(target)
		if err != nil {
			fmt.Printf("Http.Get\nerror: %s\ntarget: %s\n", err, target)
			continue
		}
		defer resp.Body.Close()

		m, _, err := image.Decode(resp.Body)
		if err != nil {
			fmt.Printf("image.Decode\nerror: %s\ntarget: %s\n", err, target)
			continue
		}

		// Ignore small images
		bounds := m.Bounds()
		if bounds.Size().X > 300 && bounds.Size().Y > 300 {
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

func BigDecodeUTF8(s []byte) ([]byte, error) {
	I := bytes.NewReader(s)
	O := transform.NewReader(I, traditionalchinese.Big5.NewDecoder())
	d, e := ioutil.ReadAll(O)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func findCharacterSet(targetUrl string) string {
	resp, _ := http.Get(targetUrl)
	return resp.Header.Get("Content-Type")
}

func findDomainByURL(url string) *WebSite {
	//using matching first, should goes to regexp
	var targetDomain string
	tokenStrings := urlRegex.FindStringSubmatch(url)
	if len(tokenStrings) > 0 {
		targetDomain = tokenStrings[3]
	}

	for index, webside := range configSetting.SupportSites {
		if strings.Contains(targetDomain, webside.WebSite) {
			fmt.Println("Use", configSetting.SupportSites[index].WebSite, " parser.")
			return &configSetting.SupportSites[index]
		}
	}

	//Cannot find using first as default parser.
	fmt.Println("Not in our support list, skip.")
	return nil
}

func crawler(target string, workerNum int) {
	doc, err := goquery.NewDocument(target)
	if err != nil {
		panic(err)
	}

	//find web site from URL.
	targetSiteSetting := findDomainByURL(target)
	if targetSiteSetting == nil {
		return
	}

	title := doc.Find(targetSiteSetting.TitlePattern).Text()

	if targetSiteSetting.ForceBig5 {
		byteStr, _ := BigDecodeUTF8([]byte(title))
		title = string(byteStr)
	}

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

func checkFileExist(fileAdd string) bool {
	if _, err := os.Stat(fileAdd); err == nil {
		return true
	}
	return false
}

func reloadParser(jsonFileAddr string) {
	//Load parser if exist.
	var file []byte

	if checkFileExist(jsonFileAddr) {
		file, _ = ioutil.ReadFile(jsonFileAddr)
		fmt.Println("Trying load local parser file")
	} else {
		//Download latest one from github.
		updateParser(jsonFileAddr)
		file, _ = ioutil.ReadFile(jsonFileAddr)
		fmt.Println("Trying load local parser file")
	}

	if len(file) == 0 {
		//file not exist, download new one.
		fmt.Println("Local parser is empty, load default parser in app.")
		file = DefaultJson
	}
	json.Unmarshal(file, &configSetting)
}

func updateParser(jsonAddr string) {
	fmt.Println("Download new parser from github and update.")
	out, err := os.Create(jsonAddr)
	if err != nil {
		return
	}
	defer out.Close()
	resp, err := http.Get(jsonFileGithub)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
}

func main() {
	usr, _ := user.Current()
	baseDir = fmt.Sprintf("%v/Pictures/webpic", usr.HomeDir)
	os.MkdirAll(baseDir, 0755)
	jsonFileAddr := fmt.Sprintf("%s/parser.json", baseDir)

	reloadParser(jsonFileAddr)

	var postUrl string
	var workerNum int
	var useDaemon bool

	rootCmd := &cobra.Command{
		Use:   "webpic",
		Short: "Download all the images in given post url",
		Run: func(cmd *cobra.Command, args []string) {
			if useDaemon {
				//Check clipboard
				var previousString string
				fmt.Println("Start watching clipboard.... (press ctrl+c to exit)")
				for {
					text, err := clipboard.ReadAll()
					if previousString != text {
						//Found something new in clipboard, check if it is URL.
						if err == nil && len(text) > 0 {
							urlInfo := urlRegex.FindStringSubmatch(text)
							if len(urlInfo) > 0 {
								go crawler(text, workerNum)
							}
						}
						previousString = text
					}

					time.Sleep(time.Second)
				}
			} else {
				if postUrl == "" {
					fmt.Println("Please use 'webpic -u URL'.")
					return
				}
				crawler(postUrl, workerNum)
			}
		},
	}
	rootCmd.Flags().StringVarP(&postUrl, "url", "u", "", "Url of post")
	rootCmd.Flags().IntVarP(&workerNum, "worker", "w", 25, "Number of workers")
	rootCmd.Flags().BoolVarP(&useDaemon, "daemon", "d", false, "Enable daemon mode to watch the clipboard.")

	updateCmd := &cobra.Command{
		Use:   "update",
		Short: "Download new parser from github and update local.",
		Run: func(cmd *cobra.Command, args []string) {
			updateParser(jsonFileAddr)
		},
	}

	rootCmd.AddCommand(updateCmd)
	rootCmd.Execute()
}
