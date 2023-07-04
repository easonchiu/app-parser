/*
 * @Author: easonchiu
 * @Date: 2023-07-03 11:07:23
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-04 14:36:16
 * @Description:
 */
package parser

import (
	"errors"
	"strings"
)

type App struct {
	ID       string
	Name     string
	Category string
	Icon     string
}

type IOSIAP struct {
	Name  string
	Price string
}

type APPData struct {
	IOSID               string    // ios id
	IOSName             string    // 应用名称
	IOSIcon             string    // ios 图标地址
	IOSDesc             string    // ios 描述
	IOSPrivacyPolicyUrl string    // ios 隐私政策地址
	IOSOtherApps        []*App    // ios 全部同主体的app
	IOSCategory         string    // ios 分类
	IOSPackageSize      string    // ios 包大小
	IOSLanguage         string    // ios 支持语言
	IOSSupplier         string    // ios 供应商名称
	IOSRate             string    // ios 评分
	IOSRateCount        string    // ios 评价数
	IOSLastVersion      string    // ios 最新版本
	IOSLastUpdate       string    // ios 最新版本时间
	IOSIAPList          []*IOSIAP // ios 内购列表
	// 华为市场
	HWID               string // hw id
	HWPackageID        string // hw package id
	HWSupplier         string // hw 供应商名称
	HWRate             string // hw 评分
	HWRateCount        string // hw 评价数
	HWLastVersion      string // hw 最新版本
	HWLastUpdate       string // hw 最新版本时间
	HWPackageSize      string // hw 包大小
	HWPrivacyPolicyUrl string // hw 隐私政策地址
	HWTargetSDK        string // hw 不知道是啥，感觉像第三方sdk的数量
	HWOtherApps        []*App // hw 全部同主体的app
	// 小米市场
	MIExist       bool   // mi 是否有
	MIRateCount   string // mi 评价数
	MILastVersion string // mi 最新版本
	MILastUpdate  string // mi 最新版本时间
	// 应用宝
	QQExist       bool   // QQ 是否有
	QQLastVersion string // QQ 最新版本
	QQLastUpdate  string // QQ 最新版本时间
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
	appData.IOSName = getAppStoreName(appStoreDoc)
	appData.IOSSupplier = getAppStoreSupplier(appStoreDoc)
	appData.IOSCategory = getAppStoreCategory(appStoreDoc)
	appData.IOSDesc = getAppStoreDesc(appStoreDoc)
	appData.IOSLanguage = getAppStoreLanguage(appStoreDoc)
	appData.IOSRate = getAppStoreRate(appStoreDoc)
	appData.IOSRateCount = getAppStoreRateCount(appStoreDoc)
	appData.IOSPackageSize = getAppStorePackageSize(appStoreDoc)
	appData.IOSPrivacyPolicyUrl = getAppStorePrivacyPolicyUrl(appStoreDoc)
	appData.IOSOtherApps = getAppStoreDeveloperOtherApps(appStoreOtherAppsDoc)
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
				appData.QQLastVersion = getQQLastVersion(doc)
				appData.QQLastUpdate = getQQLastUpdate(doc)
			}
		}
	}

	return appData, nil
}
