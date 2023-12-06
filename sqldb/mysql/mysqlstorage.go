package mysql

import (
	"database/sql"
	"errors"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"go.uber.org/zap"

	"github.com/donghc/crawler/sqldb"
)

type SqlDB struct {
	options
	db *sql.DB
}

func NewSqlDB(opts ...Option) (*SqlDB, error) {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	d := &SqlDB{}
	d.options = o
	err := d.OpenDB()
	return d, err
}

func (d *SqlDB) OpenDB() error {
	db, err := sql.Open("mysql", d.sqlDSN)
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(2048)
	db.SetMaxIdleConns(2048)
	if err = db.Ping(); err != nil {
		return err
	}
	d.db = db
	return nil
}

func (d *SqlDB) CreateTable(t sqldb.TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("column can not be empty")
	}
	sql := `CREATE TABLE IF NOT EXISTS ` + t.TableName + " ("
	if t.AutoKey {
		sql += `id INT(12) NOT NULL PRIMARY KEY AUTO_INCREMENT,`
	}
	for _, v := range t.ColumnNames {
		sql += v.Title + ` ` + v.Typ + `,`
	}
	sql = sql[:len(sql)-1] + `) ENGINE=MyISAM DEFAULT CHARSET=utf8;`

	d.logger.Debug("crate table", zap.String("sql", sql))

	_, err := d.db.Exec(sql)
	return err
}

func (d *SqlDB) Insert(t sqldb.TableData) error {
	if len(t.ColumnNames) == 0 {
		return errors.New("empty column")
	}
	sql := `INSERT INTO ` + t.TableName + `(`

	for _, v := range t.ColumnNames {
		sql += v.Title + ","
	}

	sql = sql[:len(sql)-1] + `) VALUES `

	blank := ",(" + strings.Repeat(",?", len(t.ColumnNames))[1:] + ")"
	sql += strings.Repeat(blank, t.DataCount)[1:] + `;`
	d.logger.Debug("insert table", zap.String("sql", sql))
	_, err := d.db.Exec(sql, t.Args...)
	return err
}
