package dbpg

import (
	"database/sql"
	"fmt"
	"time"

	logger "github.com/ydbt/devtool/v3/logger"

	_ "github.com/go-sql-driver/mysql"
)

type DbMysql struct {
	handle *sql.DB
	log    logger.LogI
}

// NewDbMysql
func NewDbMysql(cfg *MysqlCfg, lg logger.LogI) (*DbMysql, error) {
	// username:password@tcp(ip:port)/database?charset=utf8
	strDns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset)
	db := new(DbMysql)
	var err error
	db.log = lg
	if db.handle, err = sql.Open("mysql", strDns); err != nil {
		lg.Errorf("%v", err)
		return nil, err
	}
	db.handle.SetMaxOpenConns(cfg.MaxOpenConn)
	db.handle.SetMaxIdleConns(cfg.MaxIdleConn)
	db.handle.SetConnMaxLifetime(time.Duration(cfg.MaxLifeTime) * time.Second)
	if err = db.handle.Ping(); err != nil {
		lg.Errorf("%v", err)
	}
	return db, nil
}

// UpdateCfg
// 热加载配置
func (db *DbMysql) UpdateCfg(cfg interface{}) {
}

// AutoCommitStmt sql语句由外部拼接完成
func (db *DbMysql) AutoCommitStmt(_strSql string) error {
	tx, err := db.handle.Begin()
	if nil != err {
		db.log.Errorf("\"AutoCommitStmt\":\"launch transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	_, err = tx.Exec(_strSql)
	if nil != err {
		tx.Rollback()
		db.log.Errorf("\"AutoCommitStmt\":\"exec transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	tx.Commit()
	return nil
}

// AutoCommitBindStmtByArray 变量值需要转化为字符串数组
func (db *DbMysql) AutoCommitStmtBindByRow(_strSql string, _dbRow DbRow) error {
	tx, err := db.handle.Begin()
	if nil != err {
		db.log.Errorf("\"AutoCommitStmtBindByRow \":\"launch transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	sqlStmt, err := tx.Prepare(_strSql)
	if nil != err {
		db.log.Errorf("\"AutoCommitStmtBindByRow\":\"prepare sql failed\",\"except\":\"%v\"", err)
		return err
	}
	colLen := len(_dbRow)
	arrVals := make([]interface{}, colLen)
	index := 0
	for index < colLen {
		arrVals[index] = _dbRow[index]
		index += 1
	}
	sqlResult, err := sqlStmt.Exec(arrVals...)
	if nil != err {
		tx.Rollback()
		db.log.Errorf("\"AutoCommitStmtBindByRow\":\"exec transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	if _, err := sqlResult.RowsAffected(); nil != err {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}
	return nil
}

// AutoCommitStmtBindByArgs
// 绑定变量方式的自动提交
func (db *DbMysql) AutoCommitStmtBindByArgs(_strSql string, _argVals ...interface{}) error {
	tx, err := db.handle.Begin()
	if nil != err {
		db.log.Errorf("\"AutoCommitStmtBindByArgs\":\"launch transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	sqlStmt, err := tx.Prepare(_strSql)
	if nil != err {
		db.log.Errorf("\"AutoCommitStmtBindByArgs\":\"prepare sql failed\",\"except\":\"%v\"", err)
		return err
	}
	sqlResult, err := sqlStmt.Exec(_argVals...)
	if nil != err {
		tx.Rollback()
		db.log.Errorf("\"AutoCommitStmtBindByArgs\":\"exec transcation failed\",\"except\":\"%v\"", err)
		return err
	}
	if _, err := sqlResult.RowsAffected(); nil != err {
		tx.Rollback()
		return err
	} else {
		tx.Commit()
	}
	return nil
}

// Query 数据库表数据查询
// 返回结果为二维字符串数组
func (db *DbMysql) Query(_strSql string) (DbTable, error) {
	var dbTable DbTable
	qryRows, err := db.handle.Query(_strSql)
	if err != nil {
		db.log.Errorf("\"Query\":\"sql mistake or network error\",\"except\":\"%v\"", err)
		return dbTable, err
	}
	defer qryRows.Close()
	colFields, err := qryRows.ColumnTypes()
	if err != nil {
		db.log.Errorf("\"Query\":\"fetch query column info\",\"except\":\"%v\"", err)
		return dbTable, err
	}
	colLen := len(colFields)
	arrStrs := make([]sql.NullString, colLen)
	arrVals := make([]interface{}, colLen)
	index := 0
	for index < colLen {
		arrVals[index] = &arrStrs[index]
		index += 1
	}
	iResultCnt := 0
	for qryRows.Next() {
		if err := qryRows.Scan(arrVals...); err != nil {
			db.log.Errorf("\"Query\":\"fetch fetch row failed\",\"except\":\"%v\"", err)
			return dbTable, err
		}
		rowStrs := make([]string, colLen)
		for i := 0; i < colLen; i++ {
			if arrStrs[i].Valid {
				rowStrs[i] = arrStrs[i].String
			} else {
				rowStrs[i] = ""
			}
		}
		dbTable = append(dbTable, rowStrs)
		iResultCnt += 1
	}
	return dbTable, nil
}

// ExecuteDdl 数据库执行ddl操作
func (db *DbMysql) ExecuteDdl(_strSql string) error {
	db.log.Tracef("\"ExecuteDdl\":\"successful\",\"sql\":\":%s\"", _strSql)
	sqlResult, err := db.handle.Exec(_strSql)
	if err != nil {
		db.log.Tracef("\"ExecuteDdl\":\"successful\",\"except\":\"%v\"", err)
		return err
	}
	if iRowAffect, err := sqlResult.RowsAffected(); nil != err {
		return err
	} else {
		db.log.Tracef("\"ExecuteDdl\":\"successful\",\"tip\":\"rowsaffected:%v\"", iRowAffect)
	}
	return nil
}

func (db *DbMysql) BatchInsertStmtBindByTable(_strSql string, _dbTable DbTable) ([]int, error) {
	var listFailedIndex []int
	lenRow := len(_dbTable)
	if lenRow == 0 {
		db.log.Warnf("\"BatchInsertBindStmt\":\"valid data is null\"")
		return listFailedIndex, nil
	}
	tx, err := db.handle.Begin()
	if nil != err {
		db.log.Errorf("%v", err)
		return listFailedIndex, err
	}
	sqlStmt, err := tx.Prepare(_strSql)
	if nil != err {
		db.log.Warnf("\"BatchInsertBindStmt\":\"prepare failed\",\"except\":\"%v\"", err)
		return listFailedIndex, err
	}
	lenCol := len(_dbTable[0])
	arrVals := make([]interface{}, lenCol)
	for indexRow := 0; indexRow < lenRow; indexRow++ {
		for indexCol := 0; indexCol < lenCol; indexCol++ {
			arrVals[indexCol] = _dbTable[indexRow][indexCol]
		}
		_, err := sqlStmt.Exec(arrVals...)
		if nil != err {
			db.log.Warnf("\"BatchInsertBindStmt\":\"failed\",\"except\":\"%v\"", err)
			listFailedIndex = append(listFailedIndex, indexRow)
		}
	}
	tx.Commit()
	return listFailedIndex, nil
}
