package uiTest

import (
	"errors"
	"os"
	"strconv"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	playwright "github.com/playwright-community/playwright-go"
)

var page playwright.Page
var zentaoVersion = ""

func Login() (err error) {
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
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
	title, err := page.Title()
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
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
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
	if _, err = page.Goto("http://127.0.0.1:8081/", playwright.PageGotoOptions{
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
	}
	return
}

func goToLastUnitTestResult() {

}

func checkUnitTestResult() {

}

func InstallExt(version, codeDir string) error {
	versions := strings.Split(version, ".")
	versionNumber, _ := strconv.Atoi(versions[0])
	if versionNumber < 17 {
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
	if _, err = page.Goto("http://127.0.0.1:8081", playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded}); err != nil {
		return
	}
	title, err := page.Title()
	if err != nil {
		return
	}
	if strings.Contains(title, "欢迎使用禅道") {
		err = page.Click("text=开始安装")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		err = page.Fill(`input[name="dbPassword"]`, "123456")
		if err != nil {
			return
		}
		err = page.Click(`input[name="clearDB\[\]"]`)
		if err != nil {
			return
		}
		err = page.Click("text=保存")
		if err != nil {
			return
		}
		err = page.Click("text=下一步")
		if err != nil {
			return
		}
		_, err = page.WaitForSelector(".modal-header>>:has-text('保存配置文件')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
		if err != nil {
			return
		}
		title, err = page.Title()
		if err != nil {
			return
		}
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
		err = Login()
		if err != nil {
			return
		}
		err = createModule()
		if err != nil {
			return
		}
		err = createSuite()
		if err != nil {
			return
		}
		err = InstallExt(version, codeDir)
		if err != nil {
			return
		}
	}
	return
}

func init() {
	pw, err := playwright.Run()
	if err != nil {
		return
	}
	headless := false
	var slowMo float64 = 100
	runBrowser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{Headless: &headless, SlowMo: &slowMo})
	if err != nil {
		return
	}
	page, err = runBrowser.NewPage()
	if err != nil {
		return
	}
}
