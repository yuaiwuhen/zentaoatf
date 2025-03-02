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

func CreateSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")
	webpage.Click("text=新建站点")
	locator = webpage.Locator("#siteFormModal input")
	locator.FillNth(0, "单元测试站点")
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl)
	locator.FillNth(2, "admin")
	locator.FillNth(3, "Test123456.")
	webpage.Click("text=确定")
	webpage.WaitForSelector(".list-item-content span:has-text('单元测试站点')")
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: "单元测试站点"})
}
func EditSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")
	plwConf.DisableErr()
	locator = webpage.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	c := locator.Count()
	if c == 0 {
		CreateSite(t)
		EditSite(t)
		plwConf.EnableErr()
		return
	}
	plwConf.EnableErr()
	locator = webpage.Locator(".list-item", playwright.PageLocatorOptions{HasText: "单元测试站点"})
	webpage.Click("text=编辑")
	locator = webpage.Locator("#siteFormModal input")
	locator.FillNth(0, "单元测试站点-update")
	locator.FillNth(1, constTestHelper.ZentaoSiteUrl)
	locator.FillNth(2, "admin")
	locator.FillNth(3, "Test123456.")
	webpage.Click("#siteFormModal>>.modal-action>>span:has-text(\"确定\")")
	webpage.WaitForSelector(".list-item-content span:has-text('单元测试站点-update')")
	locator = webpage.Locator(".list-item-content span", playwright.PageLocatorOptions{HasText: "单元测试站点-update"})
}
func DeleteSite(t provider.T) {
	t.ID("5466")
	t.AddParentSuite("配置禅道站点")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	locator := webpage.Locator("#siteMenuToggle")
	locator.Click()
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("text=禅道站点管理")
	locator = webpage.Locator(".list-item:has-text('单元测试站点')")
	webpage.Click("text=删除")
	webpage.WaitForTimeout(1000)
	webpage.Click(":nth-match(.modal-action > button, 1)")
	webpage.WaitForSelector(".list-item-content span:has-text('单元测试站点')", playwright.PageWaitForSelectorOptions{State: playwright.WaitForSelectorStateDetached})
	plwConf.DisableErr()
	defer plwConf.EnableErr()
	locator = webpage.Locator(".list-item-content:has-text('单元测试站点')")
	c := locator.Count()
	if c > 0 {
		t.Errorf("Delete site fail")
		t.FailNow()
	}
}

func TestUiSite(t *testing.T) {
	runner.Run(t, "客户端-编辑禅道站点", EditSite)
	runner.Run(t, "客户端-删除禅道站点", DeleteSite)
	runner.Run(t, "客户端-创建禅道站点", CreateSite)
}
