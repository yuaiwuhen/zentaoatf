package main

import (
	"testing"

	ztfTestHelper "github.com/easysoft/zentaoatf/test/helper/ztf"
	plwHelper "github.com/easysoft/zentaoatf/test/ui/helper"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/runner"
)

func SwitchProduct(t provider.T) {
	t.ID("5496")
	t.AddParentSuite("切换禅道产品")
	webpage, _ := plwHelper.OpenUrl("http://127.0.0.1:8000/", t)
	defer webpage.Close()
	ztfTestHelper.ExpandWorspace(webpage)
	webpage.Click("#productMenuToggle")
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("#navbar .list-item>>text=企业内部工时管理系统")
	webpage.WaitForTimeout(100)
	productName := webpage.InnerText("#productMenuToggle>>span")
	if productName != "企业内部工时管理系统" {
		webpage.ScreenShot()
		t.Error("Switch product fail")
		t.FailNow()
	}
	webpage.Click("#productMenuToggle")
	webpage.WaitForSelector("#navbar .list-item")
	webpage.Click("#navbar .list-item>>text=公司企业网站建设")
}

func TestUiProduct(t *testing.T) {
	runner.Run(t, "客户端-切换禅道产品", SwitchProduct)
}
