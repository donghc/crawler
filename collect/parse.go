package collect

type RuleTree struct {
	// RuleTree.Root 是一个函数，用于生成爬虫的种子网站
	Root func() ([]*Request, error)
	// 规则哈希表，用于存储当前任务所有的规则，哈希表的 Key 为规则名，Value 为具体的规则,每一个规则就是一个 ParseFunc 解析函数。
	// 参数 RuleContext 为自定义结构体，用于传递上下文信息，也就是当前的请求参数以及要解析的内容字节数组。
	// 后续还会添加请求中的临时数据等上下文数据
	Trunk map[string]*Rule
}

// Rule 采集规则节点
type Rule struct {
	ParseFunc func(ctx *RuleContext) (ParseResult, error) // 内容解析函数
}
