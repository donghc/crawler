package sqlstorage

import (
	"encoding/json"

	"go.uber.org/zap"

	"github.com/donghc/crawler/engine"
	"github.com/donghc/crawler/sqldb"
	"github.com/donghc/crawler/sqldb/mysql"
	"github.com/donghc/crawler/storage"
)

// SqlStore 是对 Storage 接口的实现，SqlStore 实现了 option 模式，同时它的内部包含了操作数据库的 DBer 接口。
type SqlStore struct {
	dataDocker  []*storage.DataCell // 分批输出结果缓存
	columnNames []sqldb.Field       // 标题字段
	db          sqldb.DBer
	Table       map[string]struct{}
	options
}

func NewSqlStore(opts ...Option) (*SqlStore, error) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}
	s := &SqlStore{}
	s.options = o
	s.Table = make(map[string]struct{})
	var err error
	s.db, err = mysql.NewSqlDB(
		mysql.WithConnectionURL(s.sqlUrl),
		mysql.WithLogger(s.logger),
	)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// Save 循环遍历要存储的 DataCell，并判断当前 DataCell 对应的数据库表是否已经被创建。
// 1: 如果表格没有被创建，则调用 CreateTable 创建表格。
// 在存储数据时，getFields 用于获取当前数据的表字段与字段类型，这是从采集规则节点的 ItemFields 数组中获得的。
// 你可能想问，那我们为什么不直接用 DataCell 中 Data 对应的哈希表中的 Key 生成字段名呢？
// 这一方面是因为它的速度太慢，另外一方面是因为 Go 中的哈希表在遍历时的顺序是随机的，而生成的字段列表需要顺序固定。
// 2 : 如果当前的数据小于 s.BatchCount，则将数据放入到缓存中直接返回（使用缓冲区批量插入数据库可以提高程序的性能）。
// 3 : 如果缓冲区已经满了，则调用 SqlStore.Flush() 方法批量插入数据
func (s *SqlStore) Save(cells ...*storage.DataCell) error {
	for _, cell := range cells {
		tableName := cell.GetTableName()
		if _, ok := s.Table[tableName]; !ok {
			// 创建表
			columnNames := getFields(cell)

			err := s.db.CreateTable(sqldb.TableData{
				TableName:   tableName,
				ColumnNames: columnNames,
				AutoKey:     true,
			})

			if err != nil {
				s.logger.Error("create table failed", zap.Error(err))
			}
			s.Table[tableName] = struct{}{}
		}
		s.dataDocker = append(s.dataDocker, cell)
		if len(s.dataDocker) > s.BatchCount {
			s.Flush()
			s.dataDocker = nil
		}
	}
	return nil
}

func getFields(cell *storage.DataCell) []sqldb.Field {
	taskName := cell.Data["Task"].(string)
	ruleName := cell.Data["Rule"].(string)
	fields := engine.GetFields(taskName, ruleName)

	var columnNames []sqldb.Field
	for _, field := range fields {
		columnNames = append(columnNames, sqldb.Field{
			Title: field,
			Typ:   "MEDIUMTEXT",
		})
	}
	columnNames = append(columnNames,
		sqldb.Field{Title: "URL", Typ: "VARCHAR(255)"},
		sqldb.Field{Title: "Time", Typ: "VARCHAR(255)"},
	)
	return columnNames
}

func (s *SqlStore) Flush() error {
	if len(s.dataDocker) == 0 {
		return nil
	}
	args := make([]interface{}, 0)
	for _, datacell := range s.dataDocker {
		ruleName := datacell.Data["Rule"].(string)
		taskName := datacell.Data["Task"].(string)
		fields := engine.GetFields(taskName, ruleName)
		data := datacell.Data["Data"].(map[string]interface{})
		var value []string
		for _, field := range fields {
			v := data[field]
			switch v.(type) {
			case nil:
				value = append(value, "")
			case string:
				value = append(value, v.(string))
			default:
				j, err := json.Marshal(v)
				if err != nil {
					value = append(value, "")
				} else {
					value = append(value, string(j))
				}
			}
		}
		value = append(value, datacell.Data["URL"].(string), datacell.Data["Time"].(string))
		for _, v := range value {
			args = append(args, v)
		}
	}

	return s.db.Insert(sqldb.TableData{
		TableName:   s.dataDocker[0].GetTableName(),
		ColumnNames: getFields(s.dataDocker[0]),
		Args:        args,
		DataCount:   len(s.dataDocker),
	})
}
