/*
 * @Author: easonchiu
 * @Date: 2023-07-03 17:09:07
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-07 17:31:53
 * @Description:
 */
package parser

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

const UA = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36"

// 华为市场
type HWData struct {
	HWID               string `bson:"hw_id"`                 // hw id
	HWPackageID        string `bson:"hw_package_id"`         // hw package id
	HWName             string `bson:"hw_name"`               // hw 名称
	HWSupplier         string `bson:"hw_supplier"`           // hw 供应商名称
	HWRate             string `bson:"hw_rate"`               // hw 评分
	HWRateCount        string `bson:"hw_rate_count"`         // hw 评价数
	HWLastVersion      string `bson:"hw_last_version"`       // hw 最新版本
	HWLastUpdate       string `bson:"hw_last_update"`        // hw 最新版本时间
	HWPackageSize      string `bson:"hw_package_size"`       // hw 包大小
	HWPrivacyPolicyUrl string `bson:"hw_privacy_policy_url"` // hw 隐私政策地址
	HWTargetSDK        string `bson:"hw_target_sdk"`         // hw 不知道是啥，感觉像第三方sdk的数量
	HWOtherApps        []*App `bson:"hw_other_apps"`         // hw 全部同主体的app
}

// 根据应用名获取华为id
func GetHWIdByName(name string) string {
	return getHWAppId(name)
}

// 获取华为市场数据
func ParseHWData(hwId string) (*HWData, error) {
	// 创建 hw data 结构体
	hwData := new(HWData)

	if strings.TrimSpace(hwId) == "" {
		return hwData, errors.New("hwId 不能为空")
	}

	json, err := getHWAppData(hwId)
	if err != nil {
		return hwData, err
	}

	hwData.HWID = hwId
	hwData.HWPackageID = getHWPackageID(json)
	hwData.HWName = getHWName(json)
	hwData.HWSupplier = getHWSupplier(json)
	hwData.HWRate = getHWRate(json)
	hwData.HWRateCount = getHWRateCount(json)
	hwData.HWLastVersion = getHWLastVersion(json)
	hwData.HWLastUpdate = getHWLastUpdate(json)
	hwData.HWPackageSize = getHWPackageSize(json)
	hwData.HWTargetSDK = getHWTargetSDK(json)
	hwData.HWPrivacyPolicyUrl = getHWPrivacyPolicyUrl(json)
	hwData.HWOtherApps = getHWOtherApps(json, hwData.HWID)

	return hwData, nil
}

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
	request.Header.Set("User-Agent", UA)

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
	findName := data.Get("name").String()

	// 判断名称是否匹配
	nameContains := strings.Contains(name, findName)
	if !nameContains {
		return ""
	}

	return data.Get("appid").String()
}

func getHWAppData(appid string) (*gjson.Result, error) {
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
		return nil, err
	}

	now := time.Now()
	request.Header.Set("Interface-Code", id+"_"+strconv.Itoa(int(now.UnixMicro())))
	request.Header.Set("User-Agent", UA)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json := gjson.ParseBytes(bytes)
	return &json, nil
}

func getHWOtherAppsData(appid string) (*gjson.Result, error) {
	id := getHWInterfaceCode()
	u := "https://web-drcn.hispace.dbankcloud.cn/uowap/index"

	params := url.Values{}
	params.Add("method", "internal.getTabDetail")
	params.Add("serviceType", "20")
	params.Add("reqPageNum", "1")
	params.Add("uri", "appdetailCommon|"+appid+"|automore|doublecolumncardwithstar|903428")
	params.Add("maxResults", "25")
	params.Add("locale", "zh")

	client := &http.Client{}
	request, err := http.NewRequest("GET", u+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	request.Header.Set("Interface-Code", id+"_"+strconv.Itoa(int(now.UnixMicro())))
	request.Header.Set("User-Agent", UA)

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	json := gjson.ParseBytes(bytes)
	return &json, nil
}

// 获取名称
func getHWName(json *gjson.Result) string {
	name := json.Get("layoutData.0.dataList.0.name")
	return strings.TrimSpace(name.String())
}

// 获取 package id
func getHWPackageID(json *gjson.Result) string {
	pkg := json.Get("layoutData.1.dataList.0.package")
	return strings.TrimSpace(pkg.String())
}

func getHWSupplier(json *gjson.Result) string {
	supplier := json.Get("layoutData.3.dataList.0.developer")
	return strings.TrimSpace(supplier.String())
}

// 获取评分
func getHWRate(json *gjson.Result) string {
	rate := json.Get("layoutData.0.dataList.0.starDesc")
	return strings.TrimSpace(rate.String())
}

// 获取评价数量
func getHWRateCount(json *gjson.Result) string {
	rate := json.Get("layoutData.0.dataList.0.gradeCount")
	return strings.TrimSpace(rate.String())
}

// 获取版本信息
func getHWLastVersion(json *gjson.Result) string {
	version := json.Get("layoutData.1.dataList.0.versionName")
	return strings.TrimSpace(version.String())
}

// 获取版本更新时间
func getHWLastUpdate(json *gjson.Result) string {
	update := json.Get("layoutData.3.dataList.0.releaseDate")
	time := strings.TrimSpace(update.String())
	time = strings.ReplaceAll(time, "/", "-")
	return time
}

// 获取包大小
func getHWPackageSize(json *gjson.Result) string {
	size := json.Get("layoutData.3.dataList.0.size")
	return strings.TrimSpace(size.String())
}

// 获取 target sdk
func getHWTargetSDK(json *gjson.Result) string {
	size := json.Get("layoutData.1.dataList.0.targetSDK")
	return strings.TrimSpace(size.String())
}

// 获取隐私政策网址
func getHWPrivacyPolicyUrl(json *gjson.Result) string {
	url := json.Get("layoutData.8.dataList.0.conceal.text")
	return strings.TrimSpace(url.String())
}

// 获取 同开发者的其他应用
func getHWOtherApps(json *gjson.Result, appid string) []*App {
	list := json.Get("layoutData.11.dataList.0.list")
	if !list.IsArray() {
		return nil
	}

	apps := make([]*App, 0)

	// 从应用详情页面获取，当数量超过5个时，拉另一个接口
	if len(list.Array()) <= 5 {
		list.ForEach(func(_, value gjson.Result) bool {
			apps = append(apps, &App{
				ID:       value.Get("appid").String(),
				Name:     value.Get("name").String(),
				Category: value.Get("kindName").String() + "-" + value.Get("tagName").String(),
				Icon:     value.Get("icon").String(),
			})
			return true
		})
	} else {
		json, err := getHWOtherAppsData(appid)
		if err == nil {
			list = json.Get("layoutData.0.dataList")

			if !list.IsArray() {
				return nil
			}

			list.ForEach(func(_, value gjson.Result) bool {
				apps = append(apps, &App{
					ID:       value.Get("appid").String(),
					Name:     value.Get("name").String(),
					Category: value.Get("kindName").String() + "-" + value.Get("tagName").String(),
					Icon:     value.Get("icon").String(),
				})
				return true
			})
		}
	}

	return apps
}
