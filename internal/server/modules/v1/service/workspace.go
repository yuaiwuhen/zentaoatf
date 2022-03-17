package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	configUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/config"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	commonUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/common"
	fileUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/file"
	logUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
)

type WorkspaceService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
}

func NewWorkspaceService() *WorkspaceService {
	return &WorkspaceService{}
}

func (s *WorkspaceService) Paginate(req serverDomain.ReqPaginate) (ret domain.PageData, err error) {
	ret, err = s.WorkspaceRepo.Paginate(req)
	return
}

func (s *WorkspaceService) Create(workspace model.Workspace) (id uint, err error) {
	id, err = s.WorkspaceRepo.Create(workspace)

	return
}

func (s *WorkspaceService) FindById(id uint) (model.Workspace, error) {
	return s.WorkspaceRepo.FindById(id)
}
func (s *WorkspaceService) FindByPath(workspacePath string) (po model.Workspace, err error) {
	return s.WorkspaceRepo.FindByPath(workspacePath)
}

func (s *WorkspaceService) Update(workspace model.Workspace) (err error) {
	err = s.WorkspaceRepo.Update(workspace)
	return
}

func (s *WorkspaceService) DeleteByPath(pth string) (err error) {
	err = s.WorkspaceRepo.DeleteByPath(pth)
	if err != nil {
		return
	}

	err = s.WorkspaceRepo.SetCurrWorkspace("")

	return
}

func (s *WorkspaceService) ListWorkspaceByUser() (workspaces []model.Workspace, err error) {
	workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()

	return
}

func (s *WorkspaceService) GetByUser(currWorkspacePath string) (
	workspaces []model.Workspace, currWorkspace model.Workspace, currWorkspaceConfig commDomain.WorkspaceConf, scriptTree serverDomain.TestAsset, err error) {
	workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()

	found := false
	for _, p := range workspaces {
		if p.Path == currWorkspacePath {
			found = true
			break
		}
	}

	if !found {
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		name := fileUtils.GetDirName(currWorkspacePath)
		newLocalWorkspace := model.Workspace{Path: currWorkspacePath, Name: name, Type: commConsts.TestFunc}

		_, err = s.WorkspaceRepo.Create(newLocalWorkspace)
		if err != nil {
			logUtils.Errorf("db operation error %s", err.Error())
			return
		}

		workspaces, err = s.WorkspaceRepo.ListWorkspaceByUser()
	}

	s.WorkspaceRepo.SetCurrWorkspace(currWorkspacePath)

	currWorkspace, err = s.WorkspaceRepo.GetCurrWorkspaceByUser()
	if err != nil {
		logUtils.Errorf("db operation error %s", err.Error())
		return
	}

	if currWorkspace.Type == commConsts.TestFunc {
		scriptTree, err = scriptUtils.LoadScriptTree(currWorkspace.Path)
	}

	currWorkspaceConfig = configUtils.ReadFromFile(currWorkspace.Path)
	currWorkspaceConfig.IsWin = commonUtils.IsWin()

	return
}
