/*
 * @Author: easonchiu
 * @Date: 2023-07-03 11:07:23
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-05 10:44:07
 * @Description:
 */
package parser

import (
	"errors"
	"strings"
)

type App struct {
	ID       string `bson:"id"`
	Name     string `bson:"name"`
	Category string `bson:"category"`
	Icon     string `bson:"icon"`
}

type IOSIAP struct {
	Name  string `bson:"name"`
	Price string `bson:"price"`
}

type APPData struct {
	IOSID               string    `bson:"ios_id"`                 // ios id
	IOSFullName         string    `bson:"ios_full_name"`          // 应用名称(会包含 - 后面的内容)
	IOSName             string    `bson:"ios_name"`               // 应用名称
	IOSIcon             string    `bson:"ios_icon"`               // ios 图标地址
	IOSBundleID         string    `bson:"ios_bundle_id"`          // ios bundle id
	IOSDesc             string    `bson:"ios_desc"`               // ios 描述
	IOSPrivacyPolicyUrl string    `bson:"ios_privacy_policy_url"` // ios 隐私政策地址
	IOSOtherApps        []*App    `bson:"ios_other_apps"`         // ios 全部同主体的app
	IOSCategory         string    `bson:"ios_category"`           // ios 分类
	IOSPackageSize      string    `bson:"ios_package_size"`       // ios 包大小
	IOSLanguage         string    `bson:"ios_language"`           // ios 支持语言
	IOSSupplier         string    `bson:"ios_supplier"`           // ios 供应商名称
	IOSRate             string    `bson:"ios_rate"`               // ios 评分
	IOSRateCount        string    `bson:"ios_rate_count"`         // ios 评价数
	IOSLastVersion      string    `bson:"ios_last_version"`       // ios 最新版本
	IOSLastUpdate       string    `bson:"ios_last_update"`        // ios 最新版本时间
	IOSIAPList          []*IOSIAP `bson:"ios_ipa_list"`           // ios 内购列表
	// 华为市场
	HWID               string `bson:"hw_id"`                 // hw id
	HWPackageID        string `bson:"hw_package_id"`         // hw package id
	HWSupplier         string `bson:"hw_supplier"`           // hw 供应商名称
	HWRate             string `bson:"hw_rate"`               // hw 评分
	HWRateCount        string `bson:"hw_rate_count"`         // hw 评价数
	HWLastVersion      string `bson:"hw_last_version"`       // hw 最新版本
	HWLastUpdate       string `bson:"hw_last_update"`        // hw 最新版本时间
	HWPackageSize      string `bson:"hw_package_size"`       // hw 包大小
	HWPrivacyPolicyUrl string `bson:"hw_privacy_policy_url"` // hw 隐私政策地址
	HWTargetSDK        string `bson:"hw_target_sdk"`         // hw 不知道是啥，感觉像第三方sdk的数量
	HWOtherApps        []*App `bson:"hw_other_apps"`         // hw 全部同主体的app
	// 小米市场
	MIExist       bool   `bson:"mi_exist"`        // mi 是否有
	MIPackageID   string `bson:"mi_package_id"`   // mi package id
	MIRateCount   string `bson:"mi_rate_count"`   // mi 评价数
	MILastVersion string `bson:"mi_last_version"` // mi 最新版本
	MILastUpdate  string `bson:"mi_last_update"`  // mi 最新版本时间
	// 应用宝
	QQExist       bool   `bson:"qq_exist"`        // QQ 是否有
	QQPackageID   string `bson:"qq_package_id"`   // qq package id
	QQLastVersion string `bson:"qq_last_version"` // QQ 最新版本
	QQLastUpdate  string `bson:"qq_last_update"`  // QQ 最新版本时间
}

func ParseAPPData(iosId string) (*APPData, error) {
	if strings.TrimSpace(iosId) == "" {
		return nil, errors.New("iosId 不能为空")
	}

	appStoreDoc, err := getAppStoreDoc(iosId)
	if err != nil {
		return nil, err
	}

	appStoreOtherAppsDoc, err := getAppStoreOtherAppsDoc(iosId)
	if err != nil {
		return nil, err
	}

	// 创建 app data 结构体
	appData := new(APPData)
	appData.IOSID = iosId

	// ios数据
	appData.IOSFullName = getAppStoreName(appStoreDoc)
	nameSp := strings.Split(appData.IOSFullName, "-")
	if len(nameSp) > 1 {
		appData.IOSName = strings.TrimSpace(nameSp[0])
	} else {
		appData.IOSName = appData.IOSFullName
	}
	appData.IOSIcon = getAppStoreIcon(appStoreDoc)
	appData.IOSBundleID = getAppStoreBundleID(iosId)
	appData.IOSSupplier = getAppStoreSupplier(appStoreDoc)
	appData.IOSCategory = getAppStoreCategory(appStoreDoc)
	appData.IOSDesc = getAppStoreDesc(appStoreDoc)
	appData.IOSLanguage = getAppStoreLanguage(appStoreDoc)
	appData.IOSRate = getAppStoreRate(appStoreDoc)
	appData.IOSRateCount = getAppStoreRateCount(appStoreDoc)
	appData.IOSPackageSize = getAppStorePackageSize(appStoreDoc)
	appData.IOSPrivacyPolicyUrl = getAppStorePrivacyPolicyUrl(appStoreDoc)
	appData.IOSOtherApps = getAppStoreDeveloperOtherApps(appStoreOtherAppsDoc)
	appData.IOSIAPList = getAppStoreIAPList(appStoreDoc)
	lastVersion, lastUpdate := getAppStoreLastUpdate(appStoreDoc)
	appData.IOSLastVersion = lastVersion
	appData.IOSLastUpdate = lastUpdate

	// 华为数据
	appData.HWID = getHWAppId(appData.IOSName)

	if appData.HWID != "" {
		json := getHWAppData(appData.HWID)
		appData.HWPackageID = getHWPackageID(json)
		appData.HWSupplier = getHWSupplier(json)
		appData.HWRate = getHWRate(json)
		appData.HWRateCount = getHWRateCount(json)
		appData.HWLastVersion = getHWLastVersion(json)
		appData.HWLastUpdate = getHWLastUpdate(json)
		appData.HWPackageSize = getHWPackageSize(json)
		appData.HWTargetSDK = getHWTargetSDK(json)
		appData.HWPrivacyPolicyUrl = getHWPrivacyPolicyUrl(json)
		appData.HWOtherApps = getHWOtherApps(json, appData.HWID)
	}

	// 小米数据
	if appData.HWPackageID != "" {
		doc, err := getMIDoc(appData.HWPackageID)
		if err == nil {
			appData.MIExist = getMIExist(doc)
			if appData.MIExist {
				appData.MIPackageID = appData.HWPackageID
				appData.MIRateCount = getMIRateCount(doc)
				appData.MILastVersion = getMILastVersion(doc)
				appData.MILastUpdate = getMILastUpdate(doc)
			}
		}
	}

	// 应用宝数据
	if appData.HWPackageID != "" {
		doc, err := getQQDoc(appData.HWPackageID)
		if err == nil {
			appData.QQExist = getQQExist(doc)
			if appData.QQExist {
				appData.QQPackageID = appData.HWPackageID
				appData.QQLastVersion = getQQLastVersion(doc)
				appData.QQLastUpdate = getQQLastUpdate(doc)
			}
		}
	}

	return appData, nil
}
