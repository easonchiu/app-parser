/*
 * @Author: easonchiu
 * @Date: 2023-07-04 11:47:45
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-04 14:16:23
 * @Description:
 */
package parser

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getMIDoc(id string) (*goquery.Document, error) {
	url := "https://app.mi.com/details?id=" + id + "&ref=search"
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("返回状态错误 %v", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// 判断是否在小米市场上架
func getMIExist(doc *goquery.Document) bool {
	node := doc.Find(".bigimg-scroll-title")
	return node.Length() > 0
}

// 获取评价数量
func getMIRateCount(doc *goquery.Document) string {
	node := doc.Find(".app-intro-comment")
	txt := node.Text()
	txt = strings.ReplaceAll(txt, "(", "")
	txt = strings.ReplaceAll(txt, ")", "")
	txt = strings.ReplaceAll(txt, "次评分", "")
	return strings.TrimSpace(txt)
}

// 获取最新版本号
func getMILastVersion(doc *goquery.Document) string {
	node := doc.Find(".main .container")
	sp := strings.Split(node.Text(), " ")

	index := 0
	flitered := make([]string, 0)
	for _, s := range sp {
		s = strings.TrimSpace(s)
		if s != "" {
			if s == "版本号" {
				index = len(flitered)
			}
			flitered = append(flitered, s)
		}
	}

	if index > 0 {
		return flitered[index+1]
	}

	return ""
}

// 获取最新更新时间
func getMILastUpdate(doc *goquery.Document) string {
	node := doc.Find(".main .container")
	sp := strings.Split(node.Text(), " ")

	index := 0
	flitered := make([]string, 0)
	for _, s := range sp {
		s = strings.TrimSpace(s)
		if s != "" {
			if s == "更新时间" {
				index = len(flitered)
			}
			flitered = append(flitered, s)
		}
	}

	if index > 0 {
		return flitered[index+1]
	}

	return ""
}
