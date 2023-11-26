package collect

type TaskModel struct {
	Property
	Root  string      `json:"root"`  // 初始化种子节点的JS脚本
	Rules []RuleModel `json:"rules"` //具体爬虫任务的规则树
}

type RuleModel struct {
	Name      string `json:"name,omitempty"`
	ParseFunc string `json:"parse_func,omitempty"`
}
