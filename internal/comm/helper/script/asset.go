package scriptHelper

import (
	"encoding/json"
	commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/comm/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/comm/helper/lang"
	"github.com/easysoft/zentaoatf/internal/pkg/consts"
	commonUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/file"
	logUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/internal/pkg/lib/string"
	serverDomain "github.com/easysoft/zentaoatf/internal/server/modules/v1/domain"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func LoadScriptTreeByDir(workspace model.Workspace, scriptIdsFromZentao map[int]string) (asset serverDomain.TestAsset, err error) {
	workspaceId := int(workspace.ID)
	workspaceDir := workspace.Path

	if !fileUtils.FileExist(workspaceDir) {
		logUtils.Infof("workspaceDir %s not exist", workspaceDir)
		return
	}

	asset = serverDomain.TestAsset{
		Type:          commConsts.Workspace,
		WorkspaceId:   workspaceId,
		WorkspaceType: workspace.Type,
		Path:          workspaceDir,
		Title:         fileUtils.GetDirName(workspaceDir),
		Slots:         iris.Map{"icon": "icon"},

		Checkable: true,
		IsLeaf:    false,
	}

	loadScriptNodesInDir(workspaceDir, &asset, 0, scriptIdsFromZentao)

	return
}

func LoadScriptByWorkspace(workspacePath string) (scriptFiles []string) {
	LoadScriptListInDir(workspacePath, &scriptFiles, 0)

	return
}

func GetScriptContent(pth string, workspaceId int) (script serverDomain.TestScript, err error) {
	script.Path = pth
	script.Code = fileUtils.ReadFile(pth)
	script.Lang = getScriptLang(pth)
	script.WorkspaceId = workspaceId

	return
}
func getScriptLang(pth string) (lang string) {
	extName := strings.TrimLeft(fileUtils.GetExtName(pth), ".")
	if extName != "" {
		return commConsts.EditorExtToLangMap[extName]
	}

	fileName := strings.ToLower(fileUtils.GetFileName(pth))
	return commConsts.EditorExtToLangMap[fileName]
}

func loadScriptNodesInDir(folder string, parent *serverDomain.TestAsset, level int, scriptIdsFromZentao map[int]string) (err error) {
	folder = fileUtils.AddFilePathSepIfNeeded(fileUtils.AbsolutePath(folder))

	list, err := ioutil.ReadDir(folder)
	if err != nil {
		return err
	}

	for _, grandson := range list {
		name := grandson.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		childPath := folder + name
		if grandson.IsDir() && level < 3 { // 目录, 递归遍历
			dirNode := AddDir(childPath, "", parent)

			loadScriptNodesInDir(childPath, dirNode, level+1, scriptIdsFromZentao)
		} else {
			content := fileUtils.ReadFile(childPath)
			caseIdStr := ReadCaseId(content)
			caseId, _ := strconv.Atoi(caseIdStr)

			if scriptIdsFromZentao == nil || caseId < 1 { // not to filter
				AddScript(childPath, caseId, parent)
				continue
			}

			_, ok := scriptIdsFromZentao[caseId]
			if ok {
				AddScript(childPath, caseId, parent)
			}
		}
	}

	return
}

func LoadScriptListInDir(path string, files *[]string, level int) error {
	regx := langHelper.GetSupportLanguageExtRegx()

	if !fileUtils.IsDir(path) { // first call, param is a file
		pass, _ := regexp.MatchString(`.*\.`+regx+`$`, path)
		if pass {
			pass = CheckFileIsScript(path)
			if pass {
				*files = append(*files, path)
			}
		}

		return nil
	}

	path = fileUtils.AbsolutePath(path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		if fi.IsDir() && level < 3 { // 目录, 递归遍历
			LoadScriptListInDir(path+name+consts.FilePthSep, files, level+1)
		} else {
			path := path + name
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", path)

			if pass {
				pass = CheckFileIsScript(path)
				if pass {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func AddScript(pth string, caseId int, parent *serverDomain.TestAsset) {
	regx := langHelper.GetSupportLanguageExtRegx()
	pass, _ := regexp.MatchString("^*.\\."+regx+"$", pth)

	if pass {
		pass = CheckFileIsScript(pth)
		if pass {
			childScript := &serverDomain.TestAsset{
				Type:   commConsts.File,
				CaseId: caseId,

				WorkspaceId:   parent.WorkspaceId,
				WorkspaceType: parent.WorkspaceType,
				Path:          pth,
				Title:         fileUtils.GetFileName(pth),
				Slots:         iris.Map{"icon": "icon"},

				Checkable: true,
				IsLeaf:    true,
			}

			parent.Children = append(parent.Children, childScript)
			parent.ScriptCount += 1
		}
	}
}
func AddDir(pth string, moduleName string, parent *serverDomain.TestAsset) (dirNode *serverDomain.TestAsset) {
	nodeType := commConsts.Dir
	if moduleName == "" {
		moduleName = fileUtils.GetDirName(pth)
	} else {
		nodeType = commConsts.ZentaoModule
	}

	dirNode = &serverDomain.TestAsset{
		Type:          nodeType,
		WorkspaceId:   parent.WorkspaceId,
		WorkspaceType: parent.WorkspaceType,
		Path:          pth,
		Title:         moduleName,
		Slots:         iris.Map{"icon": "icon"},

		Checkable: true,
		IsLeaf:    false,
	}
	parent.Children = append(parent.Children, dirNode)

	return
}

func GetScriptType(scripts []string) []string {
	exts := make([]string, 0)
	for _, script := range scripts {
		ext := path.Ext(script)
		if ext != "" {
			ext = ext[1:]
			name := commConsts.ScriptExtToNameMap[ext]

			if !stringUtils.FindInArr(name, exts) {
				exts = append(exts, name)
			}
		}
	}

	return exts
}

func GetFailedCasesDirectlyFromTestResult(resultFile string) []string {
	cases := make([]string, 0)

	extName := path.Ext(resultFile)

	if extName == "."+consts.ExtNameResult {
		resultFile = strings.Replace(resultFile, extName, "."+consts.ExtNameJson, -1)
	}

	//if vari.ServerWorkspaceDir != "" {
	//	resultFile = vari.ServerWorkspaceDir + resultFile
	//}

	content := fileUtils.ReadFile(resultFile)

	var report commDomain.ZtfReport
	json.Unmarshal([]byte(content), &report)

	for _, cs := range report.FuncResult {
		if cs.Status != "pass" {
			cases = append(cases, cs.Path)
		}
	}

	return cases
}

func GetCaseByDirAndFile(files []string) []string {
	cases := make([]string, 0)

	for _, file := range files {
		GetAllScriptsInDir(file, &cases)
	}

	return cases
}

func GetAllScriptsInDir(path string, files *[]string) error {
	if !fileUtils.IsDir(path) { // first call, param is file
		regx := langHelper.GetSupportLanguageExtRegx()

		pass, _ := regexp.MatchString(`.*\.`+regx+`$`, path)

		if pass {
			pass := CheckFileIsScript(path)
			if pass {
				*files = append(*files, path)
			}
		}

		return nil
	}

	path = fileUtils.AbsolutePath(path)

	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, fi := range dir {
		name := fi.Name()
		if commonUtils.IgnoreFile(name) {
			continue
		}

		if fi.IsDir() { // 目录, 递归遍历
			GetAllScriptsInDir(path+name+consts.FilePthSep, files)
		} else {
			path := path + name
			regx := langHelper.GetSupportLanguageExtRegx()
			pass, _ := regexp.MatchString("^*.\\."+regx+"$", path)

			if pass {
				pass = CheckFileIsScript(path)
				if pass {
					*files = append(*files, path)
				}
			}
		}
	}

	return nil
}

func GetCaseByListInMap(caseIds []int, mp map[int]string) (cases []string) {
	for _, id := range caseIds {
		pth, ok := mp[id]
		if ok && pth != "" {
			cases = append(cases, pth)
		}
	}

	return
}
