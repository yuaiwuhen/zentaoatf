package service

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"

	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/model"
	"github.com/easysoft/zentaoatf/internal/server/modules/v1/repo"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	fileUtils "github.com/easysoft/zentaoatf/pkg/lib/file"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	shellUtils "github.com/easysoft/zentaoatf/pkg/lib/shell"
)

type InterpreterService struct {
	InterpreterRepo *repo.InterpreterRepo `inject:""`
}

func NewInterpreterService() *InterpreterService {
	return &InterpreterService{}
}

func (s *InterpreterService) List() (ret []model.Interpreter, err error) {
	ret, err = s.InterpreterRepo.List()
	return
}

func (s *InterpreterService) Get(id uint) (interpreter model.Interpreter, err error) {
	return s.InterpreterRepo.Get(id)
}

func (s *InterpreterService) Create(interpreter model.Interpreter) (id uint, err error) {
	if !fileUtils.FileExist(interpreter.Path) {
		err = errors.New(i118Utils.Sprintf("wrong_interpreter_format", interpreter.Path))
		return
	}

	id, err = s.InterpreterRepo.Create(interpreter)
	return
}

func (s *InterpreterService) Update(interpreter model.Interpreter) (err error) {
	if !fileUtils.FileExist(interpreter.Path) {
		err = errors.New(fmt.Sprintf("可执行文件%s不存在", interpreter.Path))
		return
	}

	err = s.InterpreterRepo.Update(interpreter)
	return
}

func (s *InterpreterService) Delete(id uint) error {
	return s.InterpreterRepo.Delete(id)
}

func (s *InterpreterService) GetMap(pos []model.Interpreter) (mp map[string]string, err error) {
	mp = map[string]string{}

	for _, item := range pos {
		mp[item.Lang] = item.Path
	}

	return
}

func (s *InterpreterService) GetLangSettings() (mp map[string]interface{}, err error) {
	allLangs := langHelper.GetSupportLanguageArrSort()

	langs := []string{}
	mpData := map[string]map[string]string{}
	for _, lang := range allLangs {
		mp := commConsts.LangMap[lang]
		if mp["interpreter"] == "" {
			continue
		}

		subMap := map[string]string{
			"name":        mp["name"],
			"interpreter": mp["interpreter"],
			"versionCmd":  mp["versionCmd"],
		}
		mpData[lang] = subMap
		langs = append(langs, lang)
	}

	mp = map[string]interface{}{}
	mp["languages"] = langs
	mp["languageMap"] = mpData

	return
}

func (s *InterpreterService) GetLangInterpreter(language string) (list []map[string]interface{}, err error) {
	if commonUtils.IsWin() {
		return s.GetLangInterpreterWin(language)
	} else {
		return s.GetLangInterpreterUnix(language)
	}
}

func (s *InterpreterService) GetLangInterpreterUnix(language string) (list []map[string]interface{}, err error) {
	langSettings := commConsts.LangMap[language]
	whereCmd := strings.TrimSpace(langSettings["linuxWhereCmd"])
	versionCmd := strings.TrimSpace(langSettings["versionCmd"])

	output, _ := shellUtils.ExeSysCmd(whereCmd)
	pathArr := strings.Split(output, "\n")

	for _, path := range pathArr {
		path = strings.TrimSpace(path)

		if path == "" {
			continue
		}

		var vcmd string
		if language == "tcl" {
			vcmd = versionCmd + " | " + path
		} else {
			vcmd = path + " " + versionCmd + " 2>&1"
		}

		versionInfo, err1 := shellUtils.ExeSysCmd(vcmd)
		if err1 != nil {
			continue
		}

		mp := map[string]interface{}{}
		mp["path"] = path
		mp["info"] = versionInfo
		list = append(list, mp)
	}

	return
}

func (s *InterpreterService) GetLangInterpreterWin(language string) (list []map[string]interface{}, err error) {
	langSettings := commConsts.LangMap[language]
	whereCmd := strings.TrimSpace(langSettings["whereCmd"])
	versionCmd := strings.TrimSpace(langSettings["versionCmd"])

	path := langSettings["interpreter"]
	info := ""

	if language == "autoit" {
		if fileUtils.IsDir(filepath.Dir(path)) {
			mp := map[string]interface{}{}
			mp["path"] = path
			mp["info"] = "AutoIt V3"

			list = append(list, mp)
		}

		return
	}

	if !commonUtils.IsWin() || whereCmd == "" {
		return
	}

	output, _ := shellUtils.ExeSysCmd(whereCmd)
	pathArr := s.GetNoEmptyLines(strings.TrimSpace(output), ".exe", false)

	for _, path := range pathArr {
		if strings.Index(path, ".exe") != len(path)-4 {
			continue
		}
		if language == "lua" && strings.Index(path, "luac") > -1 { // compile exec file
			continue
		}

		var cmd *exec.Cmd
		if language == "tcl" {
			cmd = exec.Command("cmd", "/C", versionCmd, "|", path)
		} else {
			cmd = exec.Command("cmd", "/C", path, versionCmd)
		}

		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr

		err = cmd.Run()
		if err != nil {
			err = nil
			continue
		}

		infoArr := s.GetNoEmptyLines(out.String(), "", true)
		if len(infoArr) > 0 {
			info = infoArr[0]
		} else {
			infoArr = s.GetNoEmptyLines(stderr.String(), "", true)
			if len(infoArr) > 0 {
				info = infoArr[0]
			}
		}

		mp := map[string]interface{}{}
		mp["path"] = path
		mp["info"] = info
		list = append(list, mp)
	}

	return
}

func (s *InterpreterService) GetNoEmptyLines(text, find string, getOne bool) (ret []string) {
	arr := regexp.MustCompile("\r?\n").Split(text, -1)
	for _, item := range arr {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}

		if find == "" || (find != "" && strings.Contains(item, find)) {
			ret = append(ret, item)

			if getOne {
				break
			}
		}
	}

	return
}
