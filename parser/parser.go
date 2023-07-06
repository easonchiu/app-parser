/*
 * @Author: easonchiu
 * @Date: 2023-07-03 11:07:23
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-06 17:09:14
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

// IOS市场
type IOSData struct {
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
}

// 华为市场
type HWData struct {
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
}

// 小米市场
type MIData struct {
	MIExist       bool   `bson:"mi_exist"`        // mi 是否有
	MIPackageID   string `bson:"mi_package_id"`   // mi package id
	MIRateCount   string `bson:"mi_rate_count"`   // mi 评价数
	MILastVersion string `bson:"mi_last_version"` // mi 最新版本
	MILastUpdate  string `bson:"mi_last_update"`  // mi 最新版本时间
}

// 应用宝市场
type QQData struct {
	QQExist       bool   `bson:"qq_exist"`        // QQ 是否有
	QQPackageID   string `bson:"qq_package_id"`   // qq package id
	QQLastVersion string `bson:"qq_last_version"` // QQ 最新版本
	QQLastUpdate  string `bson:"qq_last_update"`  // QQ 最新版本时间
}

type APPData struct {
	IOSData
	// 华为市场
	HWData
	// 小米市场
	MIData
	// 应用宝
	QQData
}

func ParseAPPData(iosId string) (*APPData, error) {
	var (
		iosData *IOSData
		hwData  *HWData
		miData  *MIData
		qqData  *QQData
	)

	// IOS数据
	iosData, err := ParseIOSData(iosId)
	if err != nil {
		return nil, err
	}

	// 华为数据
	hwId := getHWAppId(iosData.IOSName)
	if hwId != "" {
		hwData, _ = ParseHWData(hwId)
	}

	// 小米数据
	if hwId != "" {
		miData, _ = ParseMIData(hwId)
	}

	// 应用宝数据
	if hwId != "" {
		qqData, _ = ParseQQData(hwId)
	}

	return &APPData{
		*iosData,
		*hwData,
		*miData,
		*qqData,
	}, nil
}

// 获取ios数据
func ParseIOSData(iosId string) (*IOSData, error) {
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

	// 创建 ios data 结构体
	iosData := new(IOSData)
	iosData.IOSID = iosId

	// ios数据
	iosData.IOSFullName = getAppStoreName(appStoreDoc)
	nameSp := strings.Split(iosData.IOSFullName, "-")
	if len(nameSp) > 1 {
		iosData.IOSName = strings.TrimSpace(nameSp[0])
	} else {
		iosData.IOSName = iosData.IOSFullName
	}
	iosData.IOSIcon = getAppStoreIcon(appStoreDoc)
	iosData.IOSBundleID = getAppStoreBundleID(iosId)
	iosData.IOSSupplier = getAppStoreSupplier(appStoreDoc)
	iosData.IOSCategory = getAppStoreCategory(appStoreDoc)
	iosData.IOSDesc = getAppStoreDesc(appStoreDoc)
	iosData.IOSLanguage = getAppStoreLanguage(appStoreDoc)
	iosData.IOSRate = getAppStoreRate(appStoreDoc)
	iosData.IOSRateCount = getAppStoreRateCount(appStoreDoc)
	iosData.IOSPackageSize = getAppStorePackageSize(appStoreDoc)
	iosData.IOSPrivacyPolicyUrl = getAppStorePrivacyPolicyUrl(appStoreDoc)
	iosData.IOSOtherApps = getAppStoreDeveloperOtherApps(appStoreOtherAppsDoc)
	iosData.IOSIAPList = getAppStoreIAPList(appStoreDoc)
	lastVersion, lastUpdate := getAppStoreLastUpdate(appStoreDoc)
	iosData.IOSLastVersion = lastVersion
	iosData.IOSLastUpdate = lastUpdate

	return iosData, nil
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

// 获取小米市场数据
func ParseMIData(pkgId string) (*MIData, error) {
	// 创建 mi data 结构体
	miData := new(MIData)

	if strings.TrimSpace(pkgId) == "" {
		return miData, errors.New("pkgId 不能为空")
	}

	doc, err := getMIDoc(pkgId)
	if err != nil {
		return miData, err
	}

	miData.MIExist = getMIExist(doc)
	if miData.MIExist {
		miData.MIPackageID = pkgId
		miData.MIRateCount = getMIRateCount(doc)
		miData.MILastVersion = getMILastVersion(doc)
		miData.MILastUpdate = getMILastUpdate(doc)
	}

	return miData, nil
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
		qqData.QQLastVersion = getQQLastVersion(doc)
		qqData.QQLastUpdate = getQQLastUpdate(doc)
	}

	return qqData, nil
}
