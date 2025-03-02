package main

import (
	"testing"

	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	plwConf "github.com/easysoft/zentaoatf/test/ui/conf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func CreateProxy(t provider.T) {
	t.ID("5740")
	t.AddParentSuite("设置界面语言")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelector("#settingModal .z-tbody-tr:has-text('本地节点')")
	plwConf.DisableErr()
	locator := webpage.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "测试执行节点"})
	c := locator.Count()
	if c > 0 {
		DeleteProxy(t)
	}
	plwConf.EnableErr()

	webpage.Click("#serverTable>>button:has-text('新建执行节点')")
	locator = webpage.Locator("#proxyFormModal input")

	locator.FillNth(0, "测试执行节点")
	webpage.WaitForTimeout(200)
	locator.FillNth(1, "http://127.0.0.1:8085")
	webpage.Click("#proxyFormModal>>text=确定")
	plwConf.DisableErr()
	err := webpage.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	if err != nil {
		webpage.Click("#proxyFormModal>>text=确定")
		webpage.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	}
	plwConf.EnableErr()
	webpage.WaitForSelector("#proxyTable .z-tbody-td >> :scope:has-text('测试执行节点')")
	webpage.Locator("#proxyTable .z-tbody-td >> :scope:has-text('测试执行节点')")
}
func EditProxy(t provider.T) {
	t.ID("5741")
	t.AddParentSuite("设置界面语言")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelectorTimeout("#proxyTable", 5000)
	webpage.WaitForSelectorTimeout("#proxyTable:has-text('测试执行节点')", 5000)
	locator := webpage.Locator("#proxyTable:has-text('测试执行节点')")
	locator = locator.Locator("text=编辑")
	locator.Click()
	locator = webpage.Locator("#proxyFormModal input")
	locator.FillNth(0, "测试执行节点-update")
	webpage.Click("#proxyFormModal>>text=确定")
	webpage.WaitForSelector("#proxyFormModal", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	webpage.WaitForTimeout(1000)
	webpage.Locator("#proxyTable .z-tbody-td >> :scope:has-text('测试执行节点')")
}
func DeleteProxy(t provider.T) {
	t.ID("5742")
	t.AddParentSuite("设置界面语言")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	webpage.Click("#navbar>>[title=\"设置\"]")
	webpage.WaitForSelectorTimeout("#proxyTable", 5000)
	webpage.WaitForSelectorTimeout("#proxyTable:has-text('测试执行节点')", 5000)
	locator := webpage.Locator("#proxyTable:has-text('测试执行节点')")
	locator = locator.Locator("text=删除")
	locator.Click()
	webpage.Click(":nth-match(.modal-action > button, 1)")
	webpage.WaitForTimeout(1000)
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	locator = webpage.Locator("#settingModal .z-tbody-tr", playwright.PageLocatorOptions{HasText: "测试执行节点-update"})
	c := locator.Count()
	if c > 0 {
		t.Errorf("Delete proxy fail")
		t.FailNow()
	}
}

func TestUiProxy(t *testing.T) {
	runner.Run(t, "客户端-创建执行节点", CreateProxy)
	runner.Run(t, "客户端-编辑执行节点", EditProxy)
	runner.Run(t, "客户端-删除执行节点", DeleteProxy)
}
