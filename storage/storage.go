package storage

// DataCell  想象为 MySQL 中的一行数据。
type DataCell struct {
	Data map[string]interface{}
}

func (d *DataCell) GetTaskName() string {
	return d.Data["Task"].(string)
}

func (d *DataCell) GetTableName() string {
	return d.Data["Task"].(string)
}

type Storage interface {
	Save(datas ...*DataCell) error
}
