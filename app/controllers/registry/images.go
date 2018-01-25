package registry

import (
	"net/url"

	"github.com/astaxie/beego"
	"github.com/zhujingfa/docker-registry-manager/app/models"
)

// ImagesController controls access to any meta information surrounding a registry image
type ImagesController struct {
	beego.Controller
}

// GetImages returns the template for the images page
func (c *ImagesController) GetImages() {

	registryName := FormatRegistryName(c.Ctx.Input.Param(":registryName"))
	repositoryName, _ := url.QueryUnescape(c.Ctx.Input.Param(":splat"))
	repositoryNameEncode := url.QueryEscape(repositoryName)
	c.Data["tagName"] = c.Ctx.Input.Param(":tagName")

	registry, _ := manager.AllRegistries.Registries[registryName]
	c.Data["registry"] = registry

	tag, _ := registry.Repositories[repositoryName].Tags[c.Ctx.Input.Param(":tagName")]
	c.Data["tag"] = tag

	labels := make(map[string]manager.KeywordInfo)
	for _, h := range tag.History {
		// run each command through the keyword lookup
		for _, cmd := range h.Commands {
			for _, keyword := range cmd.Keywords {
				labels[keyword] = manager.KeywordMapping[keyword]
			}
		}
	}
	c.Data["labels"] = labels

	// build the js chart dataset
	type segmentInfo struct {
		Stage    int      `json:"stage"`
		Cmd      string   `json:"cmd"`
		Keywords []string `json:"keywords"`
		Size     string
	}

	type dataset struct {
		Label            string   `json:"label"`
		Data             []int64  `json:"data"`
		BackgroundColor  []string `json:"backgroundColor"`
		BorderColor      []string `json:"borderColor"`
		BorderWidth      int64    `json:"borderWidth"`
		CutoutPercentage int64    `json:"cutoutPercentage"`

		// Custom data fields
		Info []segmentInfo `json:"info"`
	}

	var chart []dataset
	colors := []string{"#43A19E", "#7B43A1", "#F2317A", "#FF9824", "#58CF6C"}
	for i, history := range tag.History {
		ds := dataset{}
		for _, cmd := range history.Commands {
			ds.Data = append(ds.Data, 10)
			color := colors[0]
			colors = append(colors[:0], colors[1:]...)
			ds.BackgroundColor = append(ds.BackgroundColor, color)
			// add stage tooltip info
			ds.Info = append(ds.Info, segmentInfo{Stage: i + 1, Keywords: cmd.Keywords, Cmd: cmd.Cmd})
			colors = append(colors, color)
			if len(history.Commands) != 1 {
				ds.BorderColor = []string{"#FFF"}
				ds.BorderWidth = 1
			}
		}
		chart = append(chart, ds)
	}

	for i := len(chart)/2 - 1; i >= 0; i-- {
		opp := len(chart) - 1 - i
		chart[i], chart[opp] = chart[opp], chart[i]
	}

	c.Data["chart"] = chart
	c.Data["registryName"] = registryName
	c.Data["repositoryNameEncode"] = repositoryNameEncode
	c.Data["repositoryName"] = repositoryName

	// Compare the two manifest layers
	c.TplName = "images.tpl"
	// remove dockerhub info
}
