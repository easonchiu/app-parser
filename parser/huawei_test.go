/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:36:20
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-06 15:28:58
 * @Description:
 */
package parser

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

var HW_APP_ID = "C34075"

func TestGetHWInterfaceCode(t *testing.T) {
	code := getHWInterfaceCode()
	if code == "" {
		t.Error("code 为空")
	}
}

func TestGetHWId(t *testing.T) {
	id := getHWAppId("抖音")
	if id == "" {
		t.Error("id 为空")
	}

	fmt.Println(id)
}

func TestGetHWAppData(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	if json.String() == "" {
		t.Error("appData 为空")
	}
}

func TestGetHWPackageID(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	pkg := getHWPackageID(json)

	if pkg == "" {
		t.Error("package 为空")
	}
}

func TestGetHWSupplier(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	supplier := getHWSupplier(json)

	if supplier == "" {
		t.Error("supplier 为空")
	}
}

func TestGetHWRate(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	rate := getHWRate(json)

	if rate == "" {
		t.Error("rate 为空")
	}

	if !strings.HasSuffix(rate, "分") {
		t.Error("rate 取错了")
	}
}

func TestGetHWRateCount(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	rateCount := getHWRateCount(json)

	if rateCount == "" {
		t.Error("rateCount 为空")
	}

	reg := regexp.MustCompile("^[0-9]+$")
	if !reg.MatchString(rateCount) {
		t.Error("rateCount 取错了")
	}
}

func TestGetHWLastVersion(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	version := getHWLastVersion(json)

	if version == "" {
		t.Error("version 为空")
	}

	reg := regexp.MustCompile("^[0-9.]+$")
	if !reg.MatchString(version) {
		t.Error("version 取错了")
	}
}

func TestGetHWLastUpdate(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	update := getHWLastUpdate(json)

	if update == "" {
		t.Error("update 为空")
	}

	reg := regexp.MustCompile("^[0-9-]+$")
	if !reg.MatchString(update) {
		t.Error("update 取错了")
	}
}

func TestGetHWPackageSize(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	size := getHWPackageSize(json)

	if size == "" {
		t.Error("size 为空")
	}

	reg := regexp.MustCompile("^[0-9]+$")
	if !reg.MatchString(size) {
		t.Error("size 取错了")
	}
}

func TestGetHWTargetSDK(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	sdk := getHWTargetSDK(json)

	if sdk == "" {
		t.Error("sdk 为空")
	}

	reg := regexp.MustCompile("^[0-9]+$")
	if !reg.MatchString(sdk) {
		t.Error("sdk 取错了")
	}
}

func TestGetHWPrivacyPolicyUrl(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	url := getHWPrivacyPolicyUrl(json)

	if url == "" {
		t.Error("url 为空")
	}

	reg := regexp.MustCompile("^http")
	if !reg.MatchString(url) {
		t.Error("url 取错了")
	}
}

func TestGetHWOtherApps(t *testing.T) {
	json, err := getHWAppData(HW_APP_ID)

	if err != nil {
		t.Error(err)
	}

	apps := getHWOtherApps(json, HW_APP_ID)

	for _, a := range apps {
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
