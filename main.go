package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// ----------------------------------------------------------
// Created feiwo by 2019/8/10
// ----------------------------------------------------------
// @author feiwo
// ----------------------------------------------------------
// @version 1.0
// ----------------------------------------------------------
//  下载简书markdown 文档中的图片
// ----------------------------------------------------------

const rootPath = `/Users/feiwo/Desktop/ndk/docs`

func main() {
	_ = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if bytes, err := ioutil.ReadFile(path); err != nil {
			return err
		}else {
			compile := regexp.MustCompile(`!\[(.+)]\((http://upload-images.jianshu.io/upload_images/.+)\)`)
			allImageUrl := compile.FindAllSubmatch(bytes, -1)
			if len(allImageUrl)  == 0{
				return nil
			}
			dirPath := fmt.Sprintf("images/%s", strings.ReplaceAll(info.Name(),".md",""))
			if err := os.Mkdir(dirPath,0777); err != nil{
				return nil
			}
			fileContent := string(bytes)
			for _, originImageUrl := range allImageUrl {
				index := strings.Index(fileContent, string(originImageUrl[2]))
				endIndex := index + len(originImageUrl[2])
				sub := fileContent[index:endIndex]
				fileName := strings.Replace(fmt.Sprintf("%s",originImageUrl[1])," ","_",-1)
				fp := fmt.Sprintf("../assets/%s/%s.webp",dirPath,fileName)
				fileContent = strings.Replace(fileContent,sub, fp, endIndex)
			}
			for _, v := range allImageUrl {
				fileName := strings.Replace(fmt.Sprintf("%s",v[1])," ","_",-1)
				fp := fmt.Sprintf("%s/%s.webp",dirPath,fileName)
				fmt.Println(fmt.Sprintf("download....%s, url = %s",dirPath,fmt.Sprintf("%s/format/webp",v[2])))
				download(fp,fmt.Sprintf("%s/format/webp",v[2]))
			}
			err := ioutil.WriteFile(path, []byte(fileContent), 0666)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("写入文件成功....")
		}
		return nil
	})
}

// 下载图片
func download(fp string, url string) {
	request, e := http.NewRequest("GET", url, nil)
	if e != nil {
		log.Println(e)
	}
	response, e := http.DefaultClient.Do(request)
	if e != nil {
		log.Println(e)
	}
	defer response.Body.Close()


	file, e := os.Create(fp)

	if e != nil {
		log.Println(e)
	}

	defer file.Close()
	_, e = io.Copy(file, response.Body)
	if e != nil {
		log.Println(e)
	}
}
