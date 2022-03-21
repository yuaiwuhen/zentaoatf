package service

import (
	commConsts "github.com/aaronchen2k/deeptest/internal/comm/consts"
	scriptUtils "github.com/aaronchen2k/deeptest/internal/comm/helper/script"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"github.com/kataras/iris/v12"
)

type TestScriptService struct {
	WorkspaceRepo *repo.WorkspaceRepo `inject:""`
	SiteService   *SiteService        `inject:""`
}

func NewTestScriptService() *TestScriptService {
	return &TestScriptService{}
}

func (s *TestScriptService) LoadTestScriptsBySiteProduct(
	site serverDomain.ZentaoSite, product serverDomain.ZentaoProduct, workspaceId int) (root serverDomain.TestAsset, err error) {

	workspaces, _ := s.WorkspaceRepo.ListWorkspacesByProduct(site.Id, product.Id)

	root = serverDomain.TestAsset{Path: "", Title: "测试脚本", Type: commConsts.Root, Slots: iris.Map{"icon": "icon"}}
	for _, workspace := range workspaces {
		if workspace.Type == commConsts.ZTF {
			if workspaceId > 0 && uint(workspaceId) != workspace.ID {
				continue
			}

			scriptsInDir, _ := scriptUtils.LoadScriptTree(workspace.Path)

			root.Children = append(root.Children, &scriptsInDir)
		}
	}

	return
}
