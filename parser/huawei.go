/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:09:07
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-03 19:15:35
 * @Description:
 */
package parser

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// 获取华为市场的interface id
func getHWInterfaceCode() string {
	url := "https://web-drcn.hispace.dbankcloud.cn/webedge/getInterfaceCode"
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return ""
	}
	resp, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	return strings.ReplaceAll(string(bytes), "\"", "")
}

func getHWAppId(name string) string {
	id := getHWInterfaceCode()
	u := "https://web-drcn.hispace.dbankcloud.cn/uowap/index"

	params := url.Values{}
	params.Add("method", "internal.getTabDetail")
	params.Add("serviceType", "20")
	params.Add("reqPageNum", "1")
	params.Add("uri", "searchApp|"+name)
	params.Add("maxResults", "25")
	params.Add("version", "10.0.0")
	params.Add("locale", "zh")

	client := &http.Client{}
	request, err := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	if err != nil {
		return ""
	}

	now := time.Now()
	request.Header.Set("Interface-Code", id+"_"+strconv.Itoa(int(now.UnixMicro())))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		return ""
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	json := gjson.ParseBytes(bytes)
	data := json.Get("layoutData.0.dataList.0")

	// 判断名称是否匹配
	nameContains := strings.Contains(data.Get("name").String(), name)
	if !nameContains {
		return ""
	}

	return data.Get("appid").String()
}

func getHWAppData(appid string) *gjson.Result {
	id := getHWInterfaceCode()
	u := "https://web-drcn.hispace.dbankcloud.cn/uowap/index"

	params := url.Values{}
	params.Add("method", "internal.getTabDetail")
	params.Add("serviceType", "20")
	params.Add("reqPageNum", "1")
	params.Add("uri", "app|"+appid)
	params.Add("appid", appid)
	params.Add("locale", "zh")

	client := &http.Client{}
	request, err := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	if err != nil {
		return nil
	}

	now := time.Now()
	request.Header.Set("Interface-Code", id+"_"+strconv.Itoa(int(now.UnixMicro())))
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")

	resp, err := client.Do(request)
	if err != nil {
		return nil
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	json := gjson.ParseBytes(bytes)
	return &json
}

// 获取评分
func getHWRate(json *gjson.Result) string {
	rate := json.Get("layoutData.0.dataList.0.starDesc")
	return strings.TrimSpace(rate.String())
}
