/*
 * @Author: easonchiu
 * @Date: 2023-07-04 11:47:45
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-07 17:38:33
 * @Description:
 */
package parser

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 小米市场
type MIData struct {
	MIExist       bool   `bson:"mi_exist"`        // mi 是否有
	MIPackageID   string `bson:"mi_package_id"`   // mi package id
	MIName        string `bson:"mi_name"`         // mi 名称
	MIRateCount   string `bson:"mi_rate_count"`   // mi 评价数
	MILastVersion string `bson:"mi_last_version"` // mi 最新版本
	MILastUpdate  string `bson:"mi_last_update"`  // mi 最新版本时间
}

// 获取小米市场数据
func ParseMIData(pkgId string) (*MIData, error) {
	// 创建 mi data 结构体
	miData := new(MIData)

	if strings.TrimSpace(pkgId) == "" {
		return miData, errors.New("pkgId 不能为空")
	}

	doc, err := getMIDoc(pkgId)
	if err != nil {
		return miData, err
	}

	miData.MIExist = getMIExist(doc)
	if miData.MIExist {
		miData.MIPackageID = pkgId
		miData.MIName = getMIName(doc)
		miData.MIRateCount = getMIRateCount(doc)
		miData.MILastVersion = getMILastVersion(doc)
		miData.MILastUpdate = getMILastUpdate(doc)
	}

	return miData, nil
}

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

// 获取名称
func getMIName(doc *goquery.Document) string {
	node := doc.Find(".intro-titles h3")
	txt := node.Text()
	return strings.TrimSpace(txt)
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
