/*
 * @Author: easonchiu
 * @Date: 2023-07-03 14:35:38
 * @LastEditors: easonchiu
 * @LastEditTime: 2023-07-03 17:26:16
 * @Description:
 */
package parser

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// app store 的 doc 内容
func getAppStoreDoc(id string) (*goquery.Document, error) {
	url := "https://apps.apple.com/cn/app/id" + id + "/"
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("返回状态错误 %v", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// app store 更多此开发人员的 app 页面的 doc
func getAppStoreOtherAppsDoc(id string) (*goquery.Document, error) {
	url := "https://apps.apple.com/cn/app/id" + id + "?see-all=developer-other-apps"
	client := &http.Client{}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("返回状态错误 %v", resp.StatusCode)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// 获取 section，app store 详情页面划分了 n 块 section 放不同内容，没有特定样式名，所以用标题匹配得到不同的 secion 块
func getAppStoreSection(doc *goquery.Document, title string) *goquery.Selection {
	content := doc.Find("section.section")

	return content.FilterFunction(func(i int, s *goquery.Selection) bool {
		h := s.Find(".section__headline").Text()
		return strings.HasPrefix(strings.TrimSpace(h), title)
	})
}

// 获取应用名称
func getAppStoreName(doc *goquery.Document) string {
	content := doc.Find("meta[property='og:title']")

	return strings.TrimSpace(content.AttrOr("content", ""))
}

// 获取描述文案
func getAppStoreDesc(doc *goquery.Document) string {
	content := doc.Find("meta[property='og:description']")

	return strings.TrimSpace(content.AttrOr("content", ""))
}

// 获取供应商信息
func getAppStoreSupplier(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "信息")

	items := sel.Find(".information-list .information-list__item")
	return strings.TrimSpace(items.Eq(0).Find("dd").Text())
}

// 获取分类
func getAppStoreCategory(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "信息")

	items := sel.Find(".information-list .information-list__item")
	return strings.TrimSpace(items.Eq(2).Find("dd").Text())
}

// 获取语言
func getAppStoreLanguage(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "信息")

	items := sel.Find(".information-list .information-list__item")
	return strings.TrimSpace(items.Eq(4).Find("dd").Text())
}

// 获取评分
func getAppStoreRate(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "评分及评论")

	return sel.Find(".we-customer-ratings__averages .we-customer-ratings__averages__display").Text()
}

// 获取评价数量
func getAppStoreRateCount(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "评分及评论")

	count := sel.Find(".we-customer-ratings__stats .we-customer-ratings__count").Text()
	count = strings.ReplaceAll(count, "个评分", "")

	return strings.TrimSpace(count)
}

// 获取最新版本号与时间，version update
func getAppStoreLastUpdate(doc *goquery.Document) (string, string) {
	content := doc.Find("section.whats-new")

	version := content.Find(".whats-new__latest__version").Text()
	version = strings.ReplaceAll(version, "版本", "")

	update := content.Find("time").Text()

	return strings.TrimSpace(version), strings.TrimSpace(update)
}

// 获取包大小
func getAppStorePackageSize(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "信息")

	items := sel.Find(".information-list .information-list__item")
	return strings.TrimSpace(items.Eq(1).Find("dd").Text())
}

// 获取隐私政策链接地址
func getAppStorePrivacyPolicyUrl(doc *goquery.Document) string {
	sel := getAppStoreSection(doc, "信息")
	return sel.Find(".inline-list--app-extensions .inline-list__item").Eq(1).Find("a").AttrOr("href", "")
}

// 获取 更多来自此开发人员的 App
func getAppStoreDeveloperOtherApps(doc *goquery.Document) []*App {
	items := doc.Find(".section .l-row a")

	if items.Length() > 0 {
		list := make([]*App, 0)

		// 匹配 app store 链接的正则
		// 例：https://apps.apple.com/cn/app/%E6%8A%96%E9%9F%B3%E6%9E%81%E9%80%9F%E7%89%88/id1477031443
		// app 的 id 在链接的最后，以id开头 + 一串数字，这串数字就是其 id
		idReg, _ := regexp.Compile("id.*$")

		items.Each(func(i int, s *goquery.Selection) {
			app := new(App)

			srcset := s.Find("picture source[type='image/png']").AttrOr("srcset", "")
			// png 格式找不到的话，找 jpg 格式的
			if srcset == "" {
				srcset = s.Find("picture source[type='image/jpeg']").AttrOr("srcset", "")
			}

			// 找到 292w 尺寸的图留下来
			if srcset != "" {
				sp := strings.Split(srcset, ",")
				for _, v := range sp {
					vv := strings.TrimSpace(v)
					if strings.HasSuffix(vv, "292w") {
						vv = strings.TrimSpace(strings.Replace(vv, "292w", "", 1))
						app.Icon = vv
						break
					}
				}
			}

			app.Name = s.Find(".we-lockup__copy .we-lockup__text .we-lockup__title p").Text()
			app.Category = strings.TrimSpace(s.Find(".we-lockup__copy .we-lockup__text .we-lockup__subtitle").Text())

			href := s.AttrOr("href", "")
			app.ID = idReg.FindString(href)
			app.ID = strings.Replace(app.ID, "id", "", 1) // 把 id 去掉，只留数字部分

			list = append(list, app)
		})

		return list
	}

	return nil
}

// 获取内购信息
func getAppStoreIAPList(doc *goquery.Document) []*IOSIAP {
	sel := getAppStoreSection(doc, "信息")

	items := sel.Find(".list-with-numbers__item")
	if items.Length() > 0 {
		list := make([]*IOSIAP, 0)

		items.Each(func(i int, s *goquery.Selection) {
			iap := new(IOSIAP)

			iap.Name = s.Find(".list-with-numbers__item__title span").Text()
			iap.Name = strings.TrimSpace(iap.Name)

			iap.Price = s.Find(".list-with-numbers__item__price").Text()
			iap.Price = strings.TrimSpace(iap.Price)

			list = append(list, iap)
		})

		return list
	}

	return nil
}
