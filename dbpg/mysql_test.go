package dbpg_test

// db_mysql "github.com/ydbt/devtool/database"
import (
	"testing"
	"github.com/ydbt/devtool/dbpg"
	logger "github.com/ydbt/devtool/logger"
)

var g_lg *logger.LoggerFake
var g_cfg *dbpg.MysqlCfg

const (
	gc_strHost     = "10.130.29.76"
	gc_iPort       = 3306
	gc_strAccount  = "ucds"
	gc_strPassword = "ucds"
	gc_strDatabase = "test"
	gc_strCharset  = "utf8"
)

func init() {
	g_lg = &logger.LoggerFake{
		Level:     "debug",
		IsConsole: true,
	}
	g_cfg = &dbpg.MysqlCfg{
		Host:        gc_strHost,
		Port:        gc_iPort,
		User:        gc_strAccount,
		Password:    gc_strPassword,
		Database:    gc_strDatabase,
		Charset:     gc_strCharset,
		MaxOpenConn: 5,
		MaxIdleConn: 2,
		MaxLifeTime: 60 * 60,
	}
}

// MysqlCfg 方法的基本测试
func TestMysqlCfg(t *testing.T) {
	yml := `
host: "127.0.0.1"
port: 3306
user: "ut-sdadmin"
password: "ut-admin"
database: "ut-sdadmin"
charset: "utf8"
maxopenconn: 5
maxidleconn: 2
maxlifetime: 3600
`
	ymlCfg, err := dbpg.Yaml2MysqlCfg(yml)
	if err != nil {
		t.Error(err)
		return
	}
	js := `
{
  "host": "127.0.0.1",
  "port": 3306,
  "user": "ut-sdadmin",
  "password": "ut-admin",
  "database": "ut-sdadmin",
  "charset": "utf8",
  "maxopenconn": 5,
  "maxidleconn": 2,
  "maxlifetime": 3600  
}
`
	jsCfg, err := dbpg.Yaml2MysqlCfg(js)
	if err != nil {
		t.Error(err)
		return
	}
	if (jsCfg.User != ymlCfg.User) && (jsCfg.User != "ut-sdadmin") {
		t.Errorf("\"%s\" != \"%s\" != \"ut-sdadmin\"", jsCfg.User, ymlCfg.User)
	}

}

// mysql 数据操作测试
func TestMysqlQueryInt(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}
	//	pMysql := new DbMysql()
	strSelect := "select 1 from dual"
	dbTable, err := db.Query(strSelect)
	if err != nil {
		_t.Error(err)
	}
	if len(dbTable) != 1 {
		_t.Error("query logic error, 'select 1 from dual' can't fetch data")
	}
	dbRow01 := dbTable[0]
	if len(dbRow01) != 1 {
		_t.Error("query logic error, 'select 1 from dual' row column wrong")
	}
	if dbRow01[0] != "1" {
		_t.Error("query logic error, 'select 1 from dual' select data(1) worng")
	}
}

func TestMysqlQueryStr(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}
	strSelect := "select 'a' from dual"
	dbTable, err := db.Query(strSelect)
	if err != nil {
		_t.Error(err)
	}
	if len(dbTable) != 1 {
		_t.Error("query logic error, 'select 1 from dual' can't fetch data")
	}
	dbRow01 := dbTable[0]
	if len(dbRow01) != 1 {
		_t.Error("query logic error, 'select 1 from dual' row column wrong")
	}
	if dbRow01[0] != "a" {
		_t.Error("query logic error, 'select 1 from dual' select data(1) worng")
	}
}

func TestMysqlQueryCreate(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}
	strTable := "test.ut_mysql_create_xxx"
	strCreate := "create table " + strTable + "(id int , name varchar(30))"
	err = db.ExecuteDdl(strCreate)
	if err != nil {
		_t.Error(err)
	}
	defer func() { // 测试清理
		strDrop := "drop table " + strTable
		err = db.ExecuteDdl(strDrop)
		if err != nil {
			_t.Error(err)
		}
	}()
}

func TestMysqlQueryInsert(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}

	strTable := "test.ut_mysql_create_xxx"
	strCreate := "create table " + strTable + "(id int , name varchar(30))"
	err = db.ExecuteDdl(strCreate)
	if err != nil {
		_t.Error(err)
	}
	defer func() { // 测试清理
		strDrop := "drop table " + strTable
		err = db.ExecuteDdl(strDrop)
		if err != nil {
			_t.Error(err)
		}
	}()
	strInsert := "insert into " + strTable + "(id,name) values(1,'a')"
	err = db.AutoCommitStmt(strInsert)
	if err != nil {
		_t.Error(err)
	}
}

func TestMysqlQueryBindInsert(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}
	strTable := "test.ut_mysql_create_xxx"
	strCreate := "create table " + strTable + "(id int , name varchar(30))"
	err = db.ExecuteDdl(strCreate)
	if err != nil {
		_t.Error(err)
	}
	defer func() { // 测试清理
		strDrop := "drop table " + strTable
		err = db.ExecuteDdl(strDrop)
		if err != nil {
			_t.Error(err)
		}
	}()
	strInsert := "insert into " + strTable + "(id,name) values(?,?)"
	var dbRow dbpg.DbRow
	dbRow = make([]string, 2)
	dbRow[0] = "1"
	dbRow[1] = "a"
	err = db.AutoCommitStmtBindByRow(strInsert, dbRow)
	if err != nil {
		_t.Error(err)
	}
}

func TestMysqlQuerySelect(_t *testing.T) {
	db, err := dbpg.NewDbMysql(g_cfg, g_lg)
	if err != nil {
		_t.Error(err)
		return
	}

	strTable := "test.ut_mysql_create_xxx"
	strCreate := "create table " + strTable + "(id int , name varchar(30))"

	err = db.ExecuteDdl(strCreate)
	if err != nil {
		_t.Error(err)
	}
	defer func() { // 测试清理
		strDrop := "drop table " + strTable
		err = db.ExecuteDdl(strDrop)
		if err != nil {
			_t.Error(err)
		}
	}()
	strInsert := "insert into " + strTable + "(id,name) values(1,'a')"
	err = db.AutoCommitStmt(strInsert)
	if err != nil {
		_t.Error(err)
	}

	strInsert = "insert into " + strTable + "(id,name) values(2,'b')"
	err = db.AutoCommitStmt(strInsert)
	if err != nil {
		_t.Error(err)
	}

	strSelect := "select id,name from " + strTable
	dbTable, err := db.Query(strSelect)
	if err != nil {
		_t.Error(err)
	}
	if len(dbTable) != 2 {
		_t.Error("query logic error, ", strSelect, " pre inserted 2 data")
	}
	dbRow01 := dbTable[0]
	if len(dbRow01) != 2 {
		_t.Error("query logic error, table actual 2 columen, but expect not")
	}
	if dbTable[0][1] != "a" || dbTable[1][1] != "b" {
		_t.Error("query logic error, select data worng")
	}
}
