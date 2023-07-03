/*
 * @Author: easonchiu
 * @Date: 2023-07-03 15:01:22
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-03 17:17:40
 * @Description:
 */
package parser

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var IOS_APP_ID = "1142110895"

func TestGetAppStoreName(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	name := getAppStoreName(doc)
	if name == "" {
		t.Error("name 为空")
	}
}

func TestGetAppStorePackageSize(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	size := getAppStorePackageSize(doc)
	if size == "" {
		t.Error("size 为空")
	}
}

func TestGetAppStoreSupplier(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	supplier := getAppStoreSupplier(doc)
	if supplier == "" {
		t.Error("supplier 为空")
	}
}

func TestGetAppStoreCategory(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	category := getAppStoreCategory(doc)
	if category == "" {
		t.Error("category 为空")
	}

	if category != "娱乐" {
		t.Error("category 取错了")
	}
}

func TestGetAppStoreDesc(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	desc := getAppStoreDesc(doc)

	if desc == "" {
		t.Error("desc 为空")
	}
}

func TestGetAppStoreLanguage(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	language := getAppStoreLanguage(doc)
	if language == "" {
		t.Error("language 为空")
	}
	fmt.Println(language)

	if !strings.HasPrefix(language, "简体中文") {
		t.Error("language 取错了")
	}
}

func TestGetAppStoreRate(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	rate := getAppStoreRate(doc)
	if rate == "" {
		t.Error("rate 为空")
	}
	fmt.Println(rate)
	if !strings.HasPrefix(rate, "4.") {
		t.Error("rate 取错了")
	}
}

func TestGetAppStoreRateCount(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	rateCount := getAppStoreRateCount(doc)
	if rateCount == "" {
		t.Error("rateCount 为空")
	}

	reg := regexp.MustCompile("^[0-9].*万$") // todo: 其他形式的数字时
	if !reg.MatchString(rateCount) {
		t.Error("rateCount 取错了")
	}
}

func TestGetAppStoreLastUpdate(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	version, update := getAppStoreLastUpdate(doc)
	if version == "" {
		t.Error("version 为空")
	}
	if update == "" {
		t.Error("update 为空")
	}
}

func TestGetAppStorePrivacyPolicyUrl(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	url := getAppStorePrivacyPolicyUrl(doc)
	if url == "" {
		t.Error("url 为空")
	}

	if !strings.HasPrefix(url, "http") {
		t.Error("url 取错了")
	}
}

func TestGetAppStoreDeveloperOtherApps(t *testing.T) {
	doc, err := getAppStoreOtherAppsDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	list := getAppStoreDeveloperOtherApps(doc)
	if len(list) == 0 {
		t.Error("len = 0")
	}

	for _, a := range list {
		if a.Name == "" {
			t.Error("没找到 Name 字段")
		}
		if a.Category == "" {
			t.Error("没找到 Category 字段")
		}
		if a.ID == "" {
			t.Error("没找到 ID 字段")
		}
		if a.Icon == "" {
			t.Error("没找到 Icon 字段")
		}
	}
}

func TestGetAppStoreIAPList(t *testing.T) {
	doc, err := getAppStoreDoc(IOS_APP_ID)

	if err != nil {
		t.Error(err)
	}

	list := getAppStoreIAPList(doc)
	if len(list) == 0 {
		t.Error("len = 0")
	}

	for _, a := range list {
		if a.Name == "" {
			t.Error("没找到 Name 字段")
		}
		if a.Price == "" {
			t.Error("没找到 Price 字段")
		}
	}
}
