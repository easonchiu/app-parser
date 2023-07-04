/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:36:20
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-04 14:15:37
 * @Description:
 */
package parser

import (
	"regexp"
	"testing"
)

var MI_APP_ID = "com.ss.android.ugc.aweme"

func TestGetMiRateCount(t *testing.T) {
	doc, err := getMIDoc(MI_APP_ID)

	if err != nil {
		t.Error(err)
	}

	rateCount := getMIRateCount(doc)

	reg := regexp.MustCompile("^[0-9]+$")
	if !reg.MatchString(rateCount) {
		t.Error("rateCount 取错了")
	}
}

func TestGetMiExist(t *testing.T) {
	doc, err := getMIDoc(MI_APP_ID)

	if err != nil {
		t.Error(err)
	}

	exist := getMIExist(doc)
	if exist == false {
		t.Error("exist 取错了")
	}

	doc, err = getMIDoc(MI_APP_ID + "fake")

	if err != nil {
		t.Error(err)
	}

	exist = getMIExist(doc)
	if exist == true {
		t.Error("exist 取错了")
	}
}

func TestGetMiLastVersion(t *testing.T) {
	doc, err := getMIDoc(MI_APP_ID)

	if err != nil {
		t.Error(err)
	}

	version := getMILastVersion(doc)

	reg := regexp.MustCompile("^[0-9.]+$")
	if !reg.MatchString(version) {
		t.Error("version 取错了")
	}
}

func TestGetMiLastUpdate(t *testing.T) {
	doc, err := getMIDoc(MI_APP_ID)

	if err != nil {
		t.Error(err)
	}

	update := getMILastUpdate(doc)

	reg := regexp.MustCompile("^[0-9-]+$")
	if !reg.MatchString(update) {
		t.Error("update 取错了")
	}
}
