package zentaoHelper

import (
	"encoding/json"
	commDomain "github.com/aaronchen2k/deeptest/internal/comm/domain"
	httpUtils "github.com/aaronchen2k/deeptest/internal/pkg/lib/http"
	"github.com/bitly/go-simplejson"
)

func GetProfile(config commDomain.WorkspaceConf) (profile commDomain.ZentaoUserProfile, err error) {
	err = Login(config)
	if err != nil {
		return
	}

	url := GenApiUrl("/user", nil, config.Url)
	bytes, err := httpUtils.Get(url)
	if err != nil {
		return
	}

	jsn, err := simplejson.NewJson(bytes)
	if err != nil {
		return
	}
	bytes, err = json.Marshal(jsn.Get("profile"))
	if err != nil {
		return
	}

	err = json.Unmarshal(bytes, &profile)
	profile.Avatar = config.Url + profile.Avatar

	return
}
