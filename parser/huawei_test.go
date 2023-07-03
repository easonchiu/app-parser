/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:36:20
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-03 19:16:21
 * @Description:
 */
package parser

import (
	"fmt"
	"strings"
	"testing"
)

var HW_APP_ID = "C10652857"

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
	json := getHWAppData(HW_APP_ID)
	if json.String() == "" {
		t.Error("appData 为空")
	}
}

func TestGetHWRate(t *testing.T) {
	json := getHWAppData(HW_APP_ID)
	rate := getHWRate(json)

	if rate == "" {
		t.Error("rate 为空")
	}

	if !strings.HasSuffix(rate, "分") {
		t.Error("rate 取错了")
	}
}
