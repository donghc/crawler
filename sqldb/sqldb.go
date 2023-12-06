package sqldb

type DBer interface {
	CreateTable(t TableData) error
	Insert(t TableData) error
}

type Field struct {
	Title string // 字段名
	Typ   string // 字段属性
}

// TableData 表的元数据
type TableData struct {
	TableName   string
	ColumnNames []Field       // 标题字段
	Args        []interface{} // 插入的数据
	DataCount   int           // 插入数量的数量
	AutoKey     bool          // 是否自增
}
