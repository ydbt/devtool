package dbpg

/**  ***注意***
  1. 所有的资源建议使用schema.[table|view|...]
*/

type DbPgI interface {
	// Query
	// DQL查询数据
	Query(_strSql string) (DbTable, error)
	// AutoCommitStmt
	// 自动提交的事务操作
	AutoCommitStmt(_strSql string) error
	// ExecuteDdl
	// DDL数据定义语句执行
	ExecuteDdl(_strSql string) error
	// AutoCommitBindStmt
	// DML自动提交的事务操作，绑定变量
	AutoCommitStmtBindByRow(_strSql string, _dbRow DbRow) error
	AutoCommitStmtBindByArgs(_strSql string, _argVals ...interface{}) error
	// BatchInsertBindStmt
	// DML批次插入操作，绑定变量，自动提交的事务操作
	// 返回：插入失败索引和异常
	BatchInsertStmtBindByTable(_strSql string, _dbTable DbTable) ([]int, error)
}

// 开启事务
//StartTranscation() error
// 提交事务
//Commit() error
// 回滚事务
//Rollback() error
// 手动提交的事务操作(绑定变量)
//HandCommitBindStmt(_strSql string) error
// 手动提交的事务操作
//HandCommitStmt(_strSql string) error
