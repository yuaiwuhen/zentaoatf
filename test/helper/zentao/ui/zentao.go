package uiTest

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	"github.com/easysoft/zentaoatf/test/ui/conf"
	playwright "github.com/playwright-community/playwright-go"
)

var page playwright.Page
var pw *playwright.Playwright
var zentaoVersion = ""

func GetStatus(url string) bool {
	if _, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return false
	}
	return true
}

func Login(url string) (err error) {
	if _, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	title, _ := page.Title()
	if !strings.Contains(title, "用户登录") {
		return
	}
	err = page.Fill(`input[name="account"]`, "admin")
	if err != nil {
		return
	}
	err = page.Fill(`input[name="password"]`, "Test123456.")
	if err != nil {
		return
	}
	err = page.Click(`button:has-text("登录")`)
	if err != nil {
		return
	}
	_, err = page.WaitForSelector("#login", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	if err != nil {
		return
	}
	title, err = page.Title()
	if err != nil {
		return
	}
	if title == "流程 - 禅道" || title == "地盘-个性化设置 - 禅道" {
		err = page.Click(`button:has-text("保存")`)
		if err != nil {
			return
		}
	}
	page.WaitForTimeout(1000)
	for {
		page.WaitForTimeout(100)
		isVisible, err := page.IsVisible("#triggerModal")
		if err != nil {
			return err
		}
		if !isVisible {
			break
		}
		isVisible, _ = page.IsVisible("#iframe-triggerModal")
		if !isVisible {
			continue
		}
		iframeName := "iframe-triggerModal"
		iframe := page.Frame(playwright.PageFrameOptions{Name: &iframeName})
		isVisible, _ = iframe.IsVisible("footer>>text=下一步")
		if isVisible {
			err = iframe.Click("footer>>text=下一步")
			continue
		}
		isVisible, _ = iframe.IsVisible("footer>>text=关闭")
		if isVisible {
			err = iframe.Click("footer>>text=关闭")
			continue
		}
		return errors.New("Find close features fail")
	}
	page.WaitForTimeout(1000)
	return
}

func createModule() (err error) {
	if _, err = page.Goto(constTestHelper.ZentaoSiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	page.Click(".nav>>li>>text=测试")
	iframeName := "app-qa"
	iframe := page.Frame(playwright.PageFrameOptions{Name: &iframeName})
	if iframe != nil {
		iframe.Click(".nav>>li>>text=用例")
		iframe.Click("#mainContent>>a>>text=维护模块")
		err = iframe.Fill(`input[name="modules\[\]"]>>nth=0`, "module1")
		if err != nil {
			return
		}
		err = iframe.Fill(`input[name="modules\[\]"]>>nth=1`, "module2")
		if err != nil {
			return
		}
		err = iframe.Click(`#submit`)
		if err != nil {
			return
		}
	} else {
		page.Click(".nav>>li>>text=用例")
		page.Click("#currentItem")
		page.WaitForSelector("#dropMenu>>.list-group", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateVisible})
		page.Click("#dropMenu>>a:has-text('企业网站建设')")
		page.Click("#mainContent>>a>>text=维护模块")
		err = page.Fill(`input[name="modules\[\]"]>>nth=0`, "module1")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="modules\[\]"]>>nth=1`, "module2")
		if err != nil {
			return
		}
		err = page.Click(`#submit`)
		if err != nil {
			return
		}
	}

	page.WaitForTimeout(1000)
	return
}

func createSuite() (err error) {
	if _, err = page.Goto(constTestHelper.ZentaoSiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	page.Click(".nav>>li>>text=测试")
	iframeName := "app-qa"
	iframe := page.Frame(playwright.PageFrameOptions{Name: &iframeName})
	if iframe != nil {
		iframe.Click(".nav>>li>>text=套件")
		iframe.Click("#mainMenu>>text=建套件")
		err = iframe.Fill(`#name`, "test_suite")
		if err != nil {
			return
		}
		err = iframe.Click(`#submit`)
		if err != nil {
			return
		}
		_, err = iframe.WaitForSelector("#submit", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}
		iframe.Click("#mainContent>>a[title=\"关联用例\"]")
		iframe.Click(`input[name="cases\[\]"]>>nth=-1`)
		iframe.Click("#submit:has-text('保存')")
	} else {
		page.Click(".nav>>li>>text=套件")
		page.Click("#mainMenu>>text=建套件")
		err = page.Fill(`#name`, "test_suite")
		if err != nil {
			return
		}
		err = page.Click(`#submit`)
		if err != nil {
			return
		}
		_, err = page.WaitForSelector("#submit", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}
		page.Click("#mainContent>>a[title=\"关联用例\"]")
		page.Click(`input[name="cases\[\]"]>>nth=-1`)
		page.Click("#submit:has-text('保存')")
	}
	return
}

func getLastUnitTestResult() (results []map[string]string, err error) {
	if _, err = page.Goto(constTestHelper.ZentaoSiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	page.Click(".nav>>li>>text=测试")
	iframeName := "app-qa"
	iframe := page.Frame(playwright.PageFrameOptions{Name: &iframeName})
	results = []map[string]string{}
	if iframe != nil {
		iframe.Click(".nav>>li>>text=用例")
		navbarHtml, _ := iframe.InnerHTML("#navbar")
		if !strings.Contains(navbarHtml, "单元测试") {
			iframe.Click("#byTypeTab")
		}
		iframe.Click("#mainMenu>>a>>text=单元测试")
		iframe.Click("#taskList>>tr>>nth=1>>td>>nth=1>>a")
		iframe.WaitForSelector("#taskList", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		tds, err := iframe.QuerySelectorAll("table>>tr")
		if err != nil {
			return results, err
		}
		for index := 1; index < len(tds); index++ {
			titleNth := "2"
			statusNth := "5"
			if index == 1 {
				titleNth = "3"
				statusNth = "6"
			}
			titleSelector, err := iframe.QuerySelector("table>>tr>>nth=" + strconv.Itoa(index) + ">>td>>nth=" + titleNth)
			if err != nil || titleSelector == nil {
				continue
			}
			title, err := titleSelector.InnerText()
			if err != nil {
				continue
			}
			statusSelector, err := iframe.QuerySelector("table>>tr>>nth=" + strconv.Itoa(index) + ">>td>>nth=" + statusNth)
			if err != nil {
				continue
			}
			status, err := statusSelector.InnerText()
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"title":  title,
				"status": status,
			})
		}
	} else {
		page.Hover(".nav>>li>>text=用例")
		page.Click(".nav>>li>>text=单元测试")
		page.Click("#currentItem")
		page.Click("#dropMenu>>a>>text=公司企业网站建设")
		page.Click("#taskList>>tr>>nth=1>>td>>nth=1>>a")
		page.WaitForSelector("#taskList", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		tds, err := page.QuerySelectorAll("table>>tr")
		if err != nil {
			return results, err
		}
		for index := 1; index < len(tds); index++ {
			titleNth := "2"
			statusNth := "5"
			if index == 1 {
				titleNth = "3"
				statusNth = "6"
			}
			titleSelector, err := page.QuerySelector("table>>tr>>nth=" + strconv.Itoa(index) + ">>td>>nth=" + titleNth)
			if err != nil || titleSelector == nil {
				continue
			}
			title, err := titleSelector.InnerText()
			if err != nil {
				continue
			}
			statusSelector, err := page.QuerySelector("table>>tr>>nth=" + strconv.Itoa(index) + ">>td>>nth=" + statusNth)
			if err != nil {
				continue
			}
			status, err := statusSelector.InnerText()
			if err != nil {
				continue
			}
			results = append(results, map[string]string{
				"title":  title,
				"status": status,
			})
		}
	}
	return results, err
}

func CheckUnitTestResult() bool {
	results, err := getLastUnitTestResult()
	if err != nil {
		return false
	}
	titleExist := map[string]bool{}
	for _, result := range results {
		titleExist[result["title"]] = true
		if result["title"] == "loginFail" && result["status"] != "通过" {
			return false
		}
		if result["title"] == "loginSuccess" && result["status"] != "失败" {
			return false
		}
	}
	return titleExist["loginFail"] == true && titleExist["loginSuccess"] == true
}

func InstallExt(version, codeDir string) error {
	versions := strings.Split(version, ".")
	versionNumber, _ := strconv.Atoi(versions[0])
	if versionNumber < 17 && version != "latest" {
		return downloadExt(codeDir)
	}
	return nil
}

func downloadExt(codeDir string) (err error) {
	if _, err = page.Goto("https://www.zentao.net/extension-browseRelease-186-front.html", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	err = page.Click("#siteNav>>a:has-text('登录')")
	if err != nil {
		return
	}
	err = page.Click("#loginModal>>a>>text=密码登录")
	if err != nil {
		return
	}
	err = page.Fill("#loginModal>>#account", "wx_62ba567413304")
	if err != nil {
		return
	}
	err = page.Fill("#loginModal>>#password", "zhaoke@easycorp.ltd")
	if err != nil {
		return
	}
	err = page.Click("#loginModal>>.login-form>>#submit")
	if err != nil {
		return
	}
	downloadInfo, err := page.ExpectDownload(func() error {
		err = page.Click("td>>a>>text=下载")
		return err
	})

	if err != nil {
		return
	}
	filePath, err := downloadInfo.Path()
	if err != nil {
		return
	}
	_, err = fileUtils.CopyFile(filePath, "restful.zip")
	if err != nil {
		return
	}
	err = fileUtils.Unzip("restful.zip", "")
	if err != nil {
		return
	}
	err = fileUtils.CopyDir("restful"+commConsts.PthSep, codeDir)
	if err != nil {
		return
	}
	os.RemoveAll("restful")
	os.Remove("restful.zip")
	return
}

func InitZentaoData(version string, codeDir string) (err error) {
	zentaoVersion = version
	if _, err = page.Goto(constTestHelper.ZentaoSiteUrl, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	title, err := page.Title()
	fmt.Println(title)
	if err != nil {
		return
	}
	if strings.Contains(title, "欢迎使用禅道") {
		err = page.Click("text=开始安装")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Fill(`input[name="dbPassword"]`, "123456")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Click(`input[name="clearDB\[\]"]`)
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Click("text=保存")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		retryCount := 0
		for {
			title, err = page.Title()
			fmt.Println(title)
			retryCount++
			if retryCount > 20 {
				break
			}
			locator, _ := page.Locator("text=数据库连接失败")
			c, _ := locator.Count()
			if c > 0 {
				page.Click("text=返回")
			} else {
				break
			}
			page.WaitForTimeout(1000)
			err = page.Fill(`input[name="dbPassword"]`, "123456")
			if err != nil {
				return
			}
			err = page.Click("text=保存")
			if err != nil {
				return
			}
		}
		title, err = page.Title()
		fmt.Println(title)
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		title, err = page.Title()
		fmt.Println(title)
		_, err = page.WaitForSelector(".modal-header>>:has-text('保存配置文件')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}

		page.WaitForLoadState()
		title, err = page.Title()
		fmt.Println(title)
		if err != nil {
			return
		}

		if strings.Contains(title, "使用模式") {
			page.Click("text=使用全生命周期管理模式")
		}
		title, err = page.Title()
		fmt.Println(title)
		if strings.Contains(title, "功能介绍") {
			err = page.Click(`button:has-text("下一步")`)
			if err != nil {
				return
			}
		}

		err = page.Fill(`input[name="company"]`, "test")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="account"]`, "admin")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="password"]`, "Test123456.")
		if err != nil {
			return
		}
		err = page.Click(`input[name="importDemoData\[\]"]`)
		if err != nil {
			return
		}
		err = page.Click("text=保存")
		if err != nil {
			return
		}
		_, err = page.WaitForSelector("text=保存", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}
		err = Login(constTestHelper.ZentaoSiteUrl)
		if err != nil {
			return
		}
		err = createModule()
		if err != nil {
			return
		}
		err = createSuite()
		title, err = page.Title()
		fmt.Println(title)
		if err != nil {
			return
		}
		if codeDir == "" {
			return
		}
		err = InstallExt(version, codeDir)
		if err != nil {
			return
		}
	}
	return
}

func Close() {
	if page == nil {
		return
	}
	page.Close()
	pw.Stop()
}

func init() {
	if page != nil {
		return
	}
	var err error
	pw, err = playwright.Run(&playwright.RunOptions{
		SkipInstallBrowsers: true,
	})
	if err != nil {
		fmt.Println(err)
		return
	}
	headless := conf.Headless
	var slowMo float64 = 100
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		fmt.Println(err)
		return
	}
	page, err = browser.NewPage(playwright.BrowserNewContextOptions{Locale: &conf.Locale})
	if err != nil {
		fmt.Println(err)
		return
	}
	Login(constTestHelper.ZentaoSiteUrl)
}
