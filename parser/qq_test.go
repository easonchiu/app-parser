/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:36:20
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-04 14:52:30
 * @Description:
 */
package parser

import (
	"regexp"
	"testing"
)

var QQ_APP_ID = "com.ss.android.ugc.aweme"

func TestGetQQExist(t *testing.T) {
	doc, err := getQQDoc(QQ_APP_ID)

	if err != nil {
		t.Error(err)
	}

	exist := getQQExist(doc)
	if exist == false {
		t.Error("exist 取错了")
	}

	doc, _ = getQQDoc(QQ_APP_ID + "fake")

	exist = getQQExist(doc)
	if exist == true {
		t.Error("exist 取错了")
	}
}

func TestGetQQLastVersion(t *testing.T) {
	doc, err := getQQDoc(QQ_APP_ID)

	if err != nil {
		t.Error(err)
	}

	version := getQQLastVersion(doc)

	reg := regexp.MustCompile("^[0-9.]+$")
	if !reg.MatchString(version) {
		t.Error("version 取错了")
	}
}

func TestGetQQLastUpdate(t *testing.T) {
	doc, err := getQQDoc(QQ_APP_ID)

	if err != nil {
		t.Error(err)
	}

	update := getQQLastUpdate(doc)

	reg := regexp.MustCompile("^[0-9.]+$")
	if !reg.MatchString(update) {
		t.Error("update 取错了")
	}
}
