package main

import (
	"testing"

	constTestHelper "github.com/easysoft/zentaoatf/test/helper/conf"
	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
	playwright "github.com/playwright-community/playwright-go"
)

func SaveScript(t provider.T) {
	t.ID("5470")
	t.AddParentSuite("禅道站点脚本")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click(".tree-node-title:has-text('1_string_match.php')")
	locator := webpage.Locator(".view-line>>text=title=check string matches pattern")
	var positionX, positionY float64 = 400, 0
	force := true
	locator.ClickWithOption(playwright.PageClickOptions{Force: &force, Position: &playwright.PageClickOptionsPosition{X: &positionX, Y: &positionY}})
	locator.Type("-test")
	webpage.Click(".tabs-nav-toolbar>>[title=\"保存\"]")
	webpage.WaitForSelector(".toast-notification-close")
	locator = webpage.Locator(".toast-notification-container", playwright.PageLocatorOptions{HasText: "保存成功"})
	locator.Click()
	webpage.Click(".tree-node-title>>text=1_string_match.php")

	locator = webpage.Locator(".view-line>>:has-text('title=check string matches pattern')")
	locator.ClickWithOption(playwright.PageClickOptions{Force: &force, Position: &playwright.PageClickOptionsPosition{X: &positionX, Y: &positionY}})
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	locator.Press("Backspace")
	webpage.Click(".tabs-nav-toolbar>>[title=\"保存\"]")
}

func ViewScript(t provider.T) {
	t.ID("5469")
	t.AddParentSuite("禅道站点脚本")
	webpage, _ := plwHelper.OpenUrl(constTestHelper.ZtfUrl, t)
	defer webpage.Close()
	ztfTestHelper.SelectSite(webpage)
	ztfTestHelper.ExpandWorspace(webpage)
	scriptLocator := webpage.Locator(".tree-node-title:has-text('1_string_match.php')")
	scriptLocator.Click()
	webpage.Locator(".view-line>>text=title=check string matches pattern")
}

func TestUiScript(t *testing.T) {
	runner.Run(t, "客户端-编辑保存禅道站点脚本", SaveScript)
	runner.Run(t, "客户端-显示禅道站点脚本", ViewScript)
}
