package collect

type Request struct {
	Ulr       string                   // 表示要访问的网站
	ParseFunc func([]byte) ParseResult // ParseFunc 函数会解析从网站获取到的网站信息，并返回 Requests 数组用于进一步获取数据。而 Items 表示获取到的数据。
}

type ParseResult struct {
	Requests []*Request
	Items    []interface{} // Items 表示获取到的数据。
}
