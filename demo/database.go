package demo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql" // _ 空白标识符导入（blank import）
)

type DB struct {
	*sql.DB // DB 结构体里匿名嵌入了一个 *sql.DB 字段
}

func RunDatabase() {
	db := connect()
	defer db.Close()

	// 创建表
	if err := db.CreateUsersTable(); err != nil {
		log.Fatal("create users table failed:", err)
	}

	// 插入一条测试数据，再按插入后的 id 查询
	email := fmt.Sprintf("alice_%d@example.com", time.Now().UnixNano())
	userID, err := db.InsertUser("Alice", email, "123456")
	if err != nil {
		log.Fatal("insert user failed:", err)
	}

	db.TestQuery(userID)
    db.DropUsersTable()
}

func connect() *DB {
	// 格式: 用户:密码@tcp(主机:端口)/数据库名?参数
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("open db failed:", err)
	}
	//defer db.Close()
	// 连接池配置（建议）
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 用超时上下文做连通性检查
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		log.Fatal("ping db failed:", err)
	}
	fmt.Println("MySQL 连接成功")
	return &DB{DB: sqlDB}
}

func (db *DB) TestQuery(id int64) {
	var name string
	err := db.QueryRow("SELECT name FROM users WHERE id = ?", id).Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}

func (db *DB) CreateUsersTable() error {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    query := `
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (id),
    UNIQUE KEY uk_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

    _, err := db.ExecContext(ctx, query)
    return err
}

func (db *DB) DropUsersTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DROP TABLE IF EXISTS users;`
	_, err := db.ExecContext(ctx, query)
	return err
}

func (db *DB) InsertUser(name, email, password string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `INSERT INTO users (name, email, password) VALUES (?, ?, ?)`
	result, err := db.ExecContext(ctx, query, name, email, password)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
