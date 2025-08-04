package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var sqlDb *sql.DB

func initMySQL() (err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/myData"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	err = db.Ping() //检测是否连接成功
	if err != nil {
		return err
	}
	db.SetMaxOpenConns(200)                 //最大连接数
	db.SetMaxIdleConns(10)                  //连接池里最大空闲连接数。必须要比maxOpenConns小
	db.SetConnMaxLifetime(time.Second * 10) //最大存活保持时间
	db.SetConnMaxIdleTime(time.Second * 10) //最大空闲保持时间
	sqlDb = db
	return nil
}

func ConnectSql() {
	if err := initMySQL(); err != nil {
		fmt.Printf("connect to db failed,err:%v\n", err)
	} else {
		fmt.Println("connect to db success")
	}
}

func Exec(sql string, args ...interface{}) (int64, error) {

	stmt, err := sqlDb.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println("spl prepare err =", err)
		return 0, err
	}
	result, err := stmt.Exec(args...)
	if err != nil {
		fmt.Println("spl Exec err =", err)
		return 0, err
	}
	code, _ := result.LastInsertId()
	fmt.Println("spl exec code =", code)
	return code, nil
}

func Query(sql string, args ...interface{}) (*sql.Rows, error) {

	stmt, err := sqlDb.Prepare(sql)
	defer stmt.Close()
	if err != nil {
		fmt.Println("spl prepare err =", err)
		return nil, err
	}
	rows, err := stmt.Query(args...)
	if err != nil {
		fmt.Println("spl query err =", err)
		return nil, err
	}
	return rows, nil
}
