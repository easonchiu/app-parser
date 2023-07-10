/*
 * @Author: easonchiu
 * @Date: 2023-07-04 11:47:45
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-07 17:41:50
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

// 应用宝市场
type QQData struct {
	QQExist       bool   `bson:"qq_exist"`        // QQ 是否有
	QQName        string `bson:"qq_name"`         // QQ 名称
	QQPackageID   string `bson:"qq_package_id"`   // qq package id
	QQLastVersion string `bson:"qq_last_version"` // QQ 最新版本
	QQLastUpdate  string `bson:"qq_last_update"`  // QQ 最新版本时间
}

// 获取应用宝数据
func ParseQQData(pkgId string) (*QQData, error) {
	// 创建 qq data 结构体
	qqData := new(QQData)

	if strings.TrimSpace(pkgId) == "" {
		return qqData, errors.New("pkgId 不能为空")
	}

	doc, err := getQQDoc(pkgId)
	if err != nil {
		return qqData, err
	}

	qqData.QQExist = getQQExist(doc)
	if qqData.QQExist {
		qqData.QQPackageID = pkgId
		qqData.QQName = getQQName(doc)
		qqData.QQLastVersion = getQQLastVersion(doc)
		qqData.QQLastUpdate = getQQLastUpdate(doc)
	}

	return qqData, nil
}

func getQQDoc(id string) (*goquery.Document, error) {
	url := "https://sj.qq.com/appdetail/" + id
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

// 判断是否在QQ市场上架
func getQQExist(doc *goquery.Document) bool {
	if doc == nil {
		return false
	}
	node := doc.Find(".Error_error__IPDOG")
	return node.Length() == 0
}

// 获取名称
func getQQName(doc *goquery.Document) string {
	node := doc.Find(".GameCard_title__Q4wFJ h1")

	name := node.Text()

	return strings.TrimSpace(name)
}

// 获取最新版本号
func getQQLastVersion(doc *goquery.Document) string {
	node := doc.Find(".GameDetail_detailItem__Lza1O")

	version := ""

	node.Map(func(i int, s *goquery.Selection) string {
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "版本号") {
			version = strings.ReplaceAll(text, "版本号", "")
		}
		return ""
	})

	return version
}

// 获取最新更新时间
func getQQLastUpdate(doc *goquery.Document) string {
	node := doc.Find(".GameDetail_detailItem__Lza1O")

	update := ""

	node.Map(func(i int, s *goquery.Selection) string {
		text := strings.TrimSpace(s.Text())
		if strings.HasPrefix(text, "更新时间") {
			update = strings.ReplaceAll(text, "更新时间", "")
		}
		return ""
	})

	return update
}
