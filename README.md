webpic: [網頁圖片小幫手] Download all images from web site.
======================
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/kkdai/webpic/master/LICENSE) [![Build Status](https://travis-ci.org/kkdai/webpic.svg)](https://travis-ci.org/kkdai/iloveptt)

中文說明: (Chinese Description)
--------------

[網頁圖片小幫手] 是一個可以幫你抓取所有網頁上所有圖片的． 這個小幫手主要有以下功能：

- 下載網頁圖片: 自動下載所有網頁上的圖片到個人圖片資料夾．
- 背景剪貼簿監控: 可以自動查看系統剪貼簿，發現有相關網址則會自動下載．
- 自動更新解析引擎: 不需要重新更新程式，只要下一個指令就可以更新解析引擎．


Description
--------------

[WebPic]: A website image downloader to parse whole content on website and download and store all images.

### Features:

It shpport following features.

#### Download all images from specific web URL

WebPic can download all images automatically and store in your PIC directory.

#### Daemon mode to monitor clipboard and download automatically.

WebPic support daemon mode as option in `-d`, it will monitor your clipboard.

#### Update parser file without rebuild binary. 

WebPic support you to update your parser setting directly without rebuild this application. Please use command `webpic update` to update your parser directly without rebuild process.



How to Install
--------------

    go get -u -x github.com/kkdai/webpic

Usage
---------------------

    webpic -u URL

All the photos will download to `USERS/Pictures/iloveptt` and it will separate folder by article name.

For Windows user, it will store in your personal pictures folder.


Options
---------------

- `-w` number of workers. (concurrency), default workers is "25"
- `-u` input URL to download directly.
- `-d` Daemon mode.
- `update` Update parser file to latest one from Github.

Examples
---------------

Download all photos from Scottie Pippen facebook pages with 10 workers.

        //Run app to download url=http://ck101.com/thread-2876990-1-1.html
        webpic -u http://ck101.com/thread-2876990-1-1.html
        
        //Download image with URL with 10 running thread.
        webpic -u http://ck101.com/thread-2876990-1-1.html -w 10
        
        //Enable daemon mode 
        webpic -d
        >> Start watching clipboard.... (press ctrl+c to exit)
        
        //Update parser 
        webpic update
        >> Trying load local parser file
        >> Download new parser from github and update.
            

Support Websites
---------------

Currently WebPic support following website, welcome to file your request in issue, or submit your PR.

- [卡提諾](http://ck101.com/)
- [提姆家族](http://www.timliao.com/bbs/)
- [GigaCircle](gigacircle.com)
- [痞客邦](pixnet.net)

Please check our latest update [parser file](https://raw.githubusercontent.com/kkdai/webpic/master/parser.json) for more detail.

Known Issues on Go 1.4.2
---------------

You might get such error if you use `go 1.4.2`.
    
        jpgimage.Decode 
        error: unsupported JPEG feature: SOF has wrong length

It is known [issue](https://github.com/golang/go/issues/4500) in Go 1.4.2 for CMYK image decode. Please upgrade to `Go 1.5` to fixed this issue.


TODOs
---------------

Welcome to file your suggestion in issues.


Contribute
---------------

Please open up an issue on GitHub before you put a lot efforts on pull request.
The code submitting to PR must be filtered with `gofmt`

Related Project
---------------

An Instagram photo downloader also here. [https://github.com/kkdai/goInstagramDownloader](https://github.com/kkdai/goInstagramDownloader)

An Facebook Album downloader also here. [https://github.com/kkdai/goFBPages](https://github.com/kkdai/goFBPages)

A Ptt web site crawler here. [https://github.com/kkdai/iloveptt](https://github.com/kkdai/iloveptt)


Advertising
---------------

If you want to browse facebook page on your iPhone, why not check my App here :p [粉絲相簿](https://itunes.apple.com/tw/app/fen-si-xiang-bu/id839324997?l=zh&mt=8)

License
---------------

This package is licensed under MIT license. See LICENSE for details.
