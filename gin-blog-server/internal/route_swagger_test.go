package ginblog

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

/*
每个注册的路由都要在 Swagger 文档中有对应的描述。

docs/ 由 swag init 生成, 新增接口忘记写注解时这里会失败, 提示需要
补注解并重新执行 ./swag_init.sh。
*/

type swaggerDoc struct {
	Paths map[string]map[string]struct {
		Summary string `json:"summary"`
	} `json:"paths"`
}

// gin 的 /:id 风格转成 swagger 的 /{id} 风格, 同时去掉 /api 前缀
func toSwaggerPath(path string) string {
	parts := strings.Split(strings.TrimPrefix(path, apiPrefix), "/")
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			parts[i] = "{" + part[1:] + "}"
		}
	}
	return strings.Join(parts, "/")
}

func TestAllRoutesHaveSwaggerDoc(t *testing.T) {
	raw, err := os.ReadFile("../docs/swagger.json")
	assert.Nil(t, err)

	var doc swaggerDoc
	assert.Nil(t, json.Unmarshal(raw, &doc))
	assert.NotEmpty(t, doc.Paths)

	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterHandlers(r)

	var missing []string
	for _, route := range r.Routes() {
		if strings.HasPrefix(route.Path, "/swagger") {
			continue
		}

		path := toSwaggerPath(route.Path)
		operations, ok := doc.Paths[path]
		if !ok {
			missing = append(missing, route.Method+" "+path)
			continue
		}
		operation, ok := operations[strings.ToLower(route.Method)]
		if !ok {
			missing = append(missing, route.Method+" "+path)
			continue
		}
		assert.NotEmpty(t, operation.Summary, "接口缺少 @Summary: "+route.Method+" "+path)
	}

	assert.Empty(t, missing, "以下接口没有出现在 docs/swagger.json 中, 需要补 swagger 注解并重新生成文档")
}

/*
反过来: 文档里不能有已经不存在的接口。

删接口时容易只删代码, docs/ 里的路径会一直留着, 在线文档就会给出
一个调不通的接口。
*/
func TestSwaggerDocHasNoExtraPath(t *testing.T) {
	raw, err := os.ReadFile("../docs/swagger.json")
	assert.Nil(t, err)

	var doc swaggerDoc
	assert.Nil(t, json.Unmarshal(raw, &doc))

	gin.SetMode(gin.TestMode)
	r := gin.New()
	RegisterHandlers(r)

	registered := make(map[string]bool)
	for _, route := range r.Routes() {
		registered[strings.ToLower(route.Method)+" "+toSwaggerPath(route.Path)] = true
	}

	var extra []string
	for path, operations := range doc.Paths {
		for method := range operations {
			if !registered[method+" "+path] {
				extra = append(extra, strings.ToUpper(method)+" "+path)
			}
		}
	}

	assert.Empty(t, extra, "以下接口只存在于 docs/swagger.json 中, 需要删掉注解并重新生成文档")
}
