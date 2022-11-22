package commConsts

import "os"

const (
	App        = "ztf"
	AppServer  = "server"
	AppAgent   = "agent"
	AppCommand = "cmd"

	Ip   = "127.0.0.1"
	Port = 8085

	JobTimeoutTime = 60 * 3
	JobRetryTime   = 3

	ConfigVersion      = "3.0"
	ConfigDir          = "conf"
	ConfigFile         = "ztf.conf"
	LogDirName         = "log"
	ExtNameSuite       = "cs"
	LogText            = "log.txt"
	ResultText         = "result.txt"
	ResultJson         = "result.json"
	ResultZip          = "result.zip"
	ExecZip            = "exec.zip"
	ExecZipPath        = "uploadTmp"
	DownloadServerPath = "serverTmp"
	DownloadPath       = "downloadTmp"
	ExecProxyPath      = "proxyExecDir"

	ExpectResultPass = "pass"

	PathInfo = "PATH_INFO"
	Get      = "GET"
	PthSep   = string(os.PathSeparator)
)

var (
	UnitBuildToolMap = map[string]BuildTool{
		"mvn": Maven,
	}

	AutoTestTypeSelenium = "selenium"
	AutoTestTypeAppium   = "appium"
	AutoTestTypes        = []string{AutoTestTypeSelenium, AutoTestTypeAppium}

	UnitTestTypeAllure     = "allure"
	UnitTestTypeJunit      = "junit"
	UnitTestTypeTestNG     = "testng"
	UnitTestTypeRobot      = "robot"
	UnitTestTypeCypress    = "cypress"
	UnitTestTypePlaywright = "playwright"
	UnitTestTypePuppeteer  = "puppeteer"
	UnitTestPhpUnit        = "phpunit"
	UnitTestTypePyTest     = "pytest"
	UnitTestTypeJest       = "jest"
	UnitTestTypeCppUnit    = "cppunit"
	UnitTestTypeGTest      = "gtest"
	UnitTestTypeQTest      = "qtest"
	UnitTestTypes          = []string{
		UnitTestTypeAllure, UnitTestTypeJunit, UnitTestTypeTestNG,
		UnitTestTypeRobot, UnitTestTypeCypress, UnitTestTypePlaywright, UnitTestTypePuppeteer,
		UnitTestPhpUnit, UnitTestTypePyTest, UnitTestTypeJest, UnitTestTypeCppUnit, UnitTestTypeGTest, UnitTestTypeQTest,
	}

	UnitTestToolMvn   = "mvn"
	UnitTestToolMocha = "mocha"
	UnitTestToolRobot = "robot"
)
