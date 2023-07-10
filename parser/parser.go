/*
 * @Author: easonchiu
 * @Date: 2023-07-03 11:07:23
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-07 17:20:36
 * @Description:
 */
package parser

type App struct {
	ID       string `bson:"id"`
	Name     string `bson:"name"`
	Category string `bson:"category"`
	Icon     string `bson:"icon"`
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
		iosData *IOSData = new(IOSData)
		hwData  *HWData  = new(HWData)
		miData  *MIData  = new(MIData)
		qqData  *QQData  = new(QQData)
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
