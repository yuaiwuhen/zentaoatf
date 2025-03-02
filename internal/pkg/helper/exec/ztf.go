package execHelper

import (
	"path"
	"strconv"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	analysisHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/analysis"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	zentaoHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/zentao"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	channelUtils "github.com/easysoft/zentaoatf/pkg/lib/channel"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12/websocket"
	"github.com/mattn/go-runewidth"
)

func ExecCases(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (report commDomain.ZtfReport, pathMaxWidth int, err error) {
	return RunZtf(ch, testSet.WorkspacePath, 0, 0, commConsts.Case, testSet.Cases, msg)
}

//func execCmd(cmd string, workspacePath string) (err error) {
//	pth := filepath.Join(workspacePath, ".cmd.tmp")
//	fileUtils.WriteFile(pth, cmd)
//
//	conf := configHelper.LoadByWorkspacePath(workspacePath)
//
//	stdOutput, errOutput := RunFile(pth, workspacePath, conf, nil, nil)
//	if errOutput != "" {
//		logUtils.Infof("failed to exec command '%s' without output %s, err %v.", pth, stdOutput, errOutput)
//	} else {
//		logUtils.Infof("exec command '%s' with output %v.", pth, stdOutput)
//	}
//
//	return
//}

func ExecModule(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	cases, err := zentaoHelper.GetCasesByModuleInDir(testSet.ProductId, testSet.ModuleId,
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)
	if err != nil {
		return
	}

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath, testSet.ProductId, testSet.ModuleId, commConsts.Module, cases, msg)
}

func ExecSuite(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases, err := zentaoHelper.GetCasesBySuiteInDir(testSet.ProductId, testSet.SuiteId,
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath,
		testSet.ProductId, testSet.SuiteId, commConsts.Suite, cases, msg)
}

func ExecTask(ch chan int, testSet serverDomain.TestSet, msg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {
	cases, err := zentaoHelper.GetCasesByTaskInDir(testSet.ProductId, testSet.TaskId,
		testSet.WorkspacePath, testSet.ScriptDirParamFromCmdLine)
	if err != nil {
		return
	}

	if testSet.Seq != "" {
		cases = analysisHelper.FilterCaseByResult(cases, testSet)
	}

	return RunZtf(ch, testSet.WorkspacePath,
		testSet.ProductId, testSet.TaskId, commConsts.Task, cases, msg)
}

func RunZtf(ch chan int,
	workspacePath string, productId, id int, by commConsts.ExecBy, cases []string, wsMsg *websocket.Message) (
	report commDomain.ZtfReport, pathMaxWidth int, err error) {

	conf := configHelper.LoadByWorkspacePath(workspacePath)

	casesToRun, casesToIgnore := FilterCases(cases, &conf)

	numbMaxWidth := 0
	numbMaxWidth, pathMaxWidth, titleMaxWidth := getNumbMaxWidth(casesToRun)
	report = genReport(productId, id, by)

	params := commDomain.ExecParams{
		CasesToRun:    casesToRun,
		CasesToIgnore: casesToIgnore,
		WorkspacePath: workspacePath,
		Conf:          conf,
		NumbMaxWidth:  numbMaxWidth,
		PathMaxWidth:  pathMaxWidth,
		TitleMaxWidth: titleMaxWidth,
		Report:        &report,
	}

	// exec scripts
	ExeScripts(params, ch, wsMsg)

	// gen report
	if len(casesToRun) > 0 {
		GenZTFTestReport(report, pathMaxWidth, workspacePath, ch, wsMsg)
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		websocketHelper.SendExecMsg("", "false", commConsts.Run, nil, wsMsg)
	}
	if ch != nil {
		if !channelUtils.IsChanClose(ch) {
			close(ch)
		}
	}

	return
}

func ExeScripts(execParams commDomain.ExecParams, ch chan int, wsMsg *websocket.Message) {

	now := time.Now()
	startTime := now.Unix()
	execParams.Report.StartTime = startTime

	workDir := commConsts.WorkDir
	if commConsts.ExecFrom == commConsts.FromClient {
		workDir = execParams.WorkspacePath
	}

	msg := ""
	if commConsts.Language == commConsts.LanguageZh {
		msg = i118Utils.Sprintf("found_scripts", workDir, len(execParams.CasesToRun), commConsts.ZtfDir)

	} else {
		msg = i118Utils.Sprintf("found_scripts", len(execParams.CasesToRun), workDir, commConsts.ZtfDir)
	}

	if commConsts.ExecFrom == commConsts.FromClient {
		msg = i118Utils.Sprintf("found_scripts_no_ztf_dir", len(execParams.CasesToRun), workDir)
		websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
	}
	logUtils.ExecConsolef(-1, msg)
	logUtils.ExecResult(msg)

	if len(execParams.CasesToIgnore) > 0 {
		temp := i118Utils.Sprintf("ignore_scripts", strconv.Itoa(len(execParams.CasesToIgnore)))
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendExecMsg(temp, "", commConsts.Run, nil, wsMsg)
		}

		logUtils.ExecConsolef(color.FgCyan, temp)
		logUtils.ExecResult(temp)
	}

	for idx, file := range execParams.CasesToRun {
		if fileUtils.IsDir(file) {
			continue
		}

		execParams.ScriptIdx = idx
		execParams.ScriptFile = file
		ExecScript(execParams, ch, wsMsg)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_all")
			if commConsts.ExecFrom == commConsts.FromClient {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run, nil, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitAllCase
		default:
		}
	}

ExitAllCase:
	endTime := time.Now().Unix()
	execParams.Report.EndTime = endTime
	execParams.Report.Duration = endTime - startTime
}

func FilterCases(cases []string, conf *commDomain.WorkspaceConf) (casesToRun, casesToIgnore []string) {
	for _, cs := range cases {
		ext := path.Ext(cs)
		if ext != "" {
			ext = ext[1:]
		}
		lang := commConsts.ScriptExtToNameMap[ext]
		if lang == "" {
			continue
		}

		if commonUtils.IsWin() {
			filterWinCases(cs, lang, conf, &casesToIgnore, &casesToRun)
			continue
		}

		if path.Ext(cs) == ".bat" {
			continue
		}
		casesToRun = append(casesToRun, cs)
	}

	return
}

func filterWinCases(cs, lang string, conf *commDomain.WorkspaceConf, casesToIgnore, casesToRun *[]string) {
	if path.Ext(cs) == ".sh" { // filter by os
		return
	}

	interpreter := configHelper.GetFieldVal(*conf, stringUtils.UcFirst(lang))

	if interpreter == "-" || interpreter == "" {
		interpreter = ""
		if lang != "bat" {
			ok := AddInterpreterIfExist(conf, lang)
			if !ok {
				*casesToIgnore = append(*casesToIgnore, cs)
			} else {
				interpreter = configHelper.GetFieldVal(*conf, stringUtils.UcFirst(lang))
			}
		}
	}

	if lang != "bat" && interpreter == "" { // ignore the ones with no interpreter set
		return
	}

	*casesToRun = append(*casesToRun, cs)
}

func genReport(productId, id int, by commConsts.ExecBy) (report commDomain.ZtfReport) {
	report = commDomain.ZtfReport{
		TestEnv: commonUtils.GetOs(), ExecBy: by, ExecById: id, ProductId: productId,
		Pass: 0, Fail: 0, Total: 0, FuncResult: make([]commDomain.FuncResult, 0)}

	report.TestType = commConsts.TestFunc
	report.TestTool = commConsts.AppServer

	return
}

func getNumbMaxWidth(casesToRun []string) (numbMaxWidth, pathMaxWidth, titleMaxWidth int) {
	for _, cs := range casesToRun {
		lent := runewidth.StringWidth(cs)
		if lent > pathMaxWidth {
			pathMaxWidth = lent
		}

		content := fileUtils.ReadFile(cs)
		caseId := scriptHelper.ReadCaseId(content)
		if len(caseId) > numbMaxWidth {
			numbMaxWidth = len(caseId)
		}

		_, _, _, title, _ := scriptHelper.GetCaseInfo(cs)
		titleLength := runewidth.StringWidth(title)
		if titleLength > titleMaxWidth {
			titleMaxWidth = titleLength
		}
	}

	return
}
