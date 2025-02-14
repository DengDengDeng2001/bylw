package modules

import (
	"strconv"
	"strings"

	"gopkg.in/yaml.v2"
)

// M is map
type M map[string]interface{}

// S is slice
type S []interface{}

// Prom ...
type Prom struct {
	ID  int64
	URL string
}

// Rule 定义了一个 Prometheus 规则，用于监控和告警的配置
type Rule struct {
	// ID 唯一标识规则的 ID
	ID int64 `json:"id"`
	// PromID 表示该规则所属的 Prometheus 实例的 ID
	PromID int64 `json:"prom_id"`
	// Expr 是 Prometheus 查询语言（PromQL）表达式，定义了数据查询的逻辑，
	Expr string `json:"expr"`
	// Op 是用于与 `Value` 字段进行比较的操作符。
	Op string `json:"op"`
	// Value 是与查询结果进行比较的值，通常用于定义阈值。
	Value string `json:"value"`
	// For 定义了一个持续时间，只有当表达式结果满足条件持续 `For` 指定的时间后，才会触发报警。这可以避免因瞬时波动导致的误报。
	For string `json:"for"`
	// Labels 是一个标签集合，用于给警报加上附加的元数据标签。
	Labels      map[string]string `json:"labels"`
	Summary     string            `json:"summary"`
	Description string            `json:"description"`
}

// Rules ...
type Rules []Rule

// PromRules ...
type PromRules struct {
	Prom  Prom
	Rules Rules
}

// RulesResp ...
type RulesResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Rules  `json:"data"`
}

// PromsResp ...
type PromsResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []Prom `json:"data"`
}

// Content get prom rules
// 将规则转化为 Prometheus 规则引擎可识别的 YAML 格式，这样 Prometheus 就能够通过配置文件加载和执行这些规则。
func (r Rules) Content() ([]byte, error) {
	rules := S{}
	for _, i := range r {
		rules = append(rules, M{
			"alert":  strconv.FormatInt(i.ID, 10),
			"expr":   strings.Join([]string{i.Expr, i.Op, i.Value}, " "),
			"for":    i.For,
			"labels": i.Labels,
			"annotations": M{
				"rule_id":     strconv.FormatInt(i.ID, 10),
				"prom_id":     strconv.FormatInt(i.PromID, 10),
				"summary":     i.Summary,
				"description": i.Description,
			},
		})
	}
	result := M{
		"groups": S{
			M{
				"name":  "ruleengine",
				"rules": rules,
			},
		},
	}

	return yaml.Marshal(result)
}

// PromRules cut prom rules
// 将所有规则按照 PromID 进行分组。这样可以确保每个 Prometheus 实例都有对应的规则集合，便于后续的管理和更新
func (r Rules) PromRules() []PromRules {
	tmp := map[int64]Rules{}

	for _, rule := range r {
		if v, ok := tmp[rule.PromID]; ok {
			tmp[rule.PromID] = append(v, rule)
		} else {
			tmp[rule.PromID] = Rules{rule}
		}
	}

	data := []PromRules{}
	for id, rules := range tmp {
		data = append(data, PromRules{
			Prom:  Prom{ID: id},
			Rules: rules,
		})
	}

	return data
}
