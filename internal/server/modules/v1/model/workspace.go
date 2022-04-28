package model

import commConsts "github.com/easysoft/zentaoatf/internal/comm/consts"

type Workspace struct {
	BaseModel

	Path string              `json:"path"`
	Name string              `json:"name"`
	Desc string              `json:"desc" gorm:"column:descr"`
	Type commConsts.TestTool `json:"type" gorm:"default:ztf"`
	Lang string              `json:"lang"`
	Cmd  string              `json:"cmd"`

	SiteId    uint `json:"siteId"`
	ProductId uint `json:"productId"`

	IsDefault bool `json:"isDefault"`
}

func (Workspace) TableName() string {
	return "biz_workspace"
}
