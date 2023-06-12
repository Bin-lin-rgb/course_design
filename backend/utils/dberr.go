package utils

import (
	"github.com/go-sql-driver/mysql"
)

const (
	ErrMySQLDupEntry            = 1062
	ErrMySQLDupEntryWithKeyName = 1586
)

// IsUniqueConstraintError 判断是否主键冲突/唯一键冲突
func IsUniqueConstraintError(err error) bool {
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		if mysqlErr.Number == ErrMySQLDupEntry ||
			mysqlErr.Number == ErrMySQLDupEntryWithKeyName {
			return true
		}
	}
	return false
}
