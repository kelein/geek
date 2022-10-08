package main

import (
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"

	"geek/glog"
	"geek/orm"
)

func init() {
	log.SetPrefix("[GeeORM] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmsgprefix)
}

func rawQuery() {
	db, err := sql.Open("sqlite3", "geek.db")
	if err != nil {
		log.Fatalf("open sqlite3 failed: %v", err)
	}
	defer func() {
		err := db.Close()
		log.Printf("close sqlite3 failed: %v", err)
	}()

	log.Printf("db info: %v", db)

	db.Exec("DROP TABLE IF EXISTS user;")
	db.Exec("CREATE TABLE user(name text);")
	res, err := db.Exec("INSERT INTO user(name) values (?), (?)", "kallen", "torres")
	if err == nil {
		n, _ := res.RowsAffected()
		log.Printf("%v row affected", n)
	}

	name := ""
	row := db.QueryRow("SELECT name FROM user LIMIT 1")
	if err := row.Scan(&name); err == nil {
		log.Printf("found user: %s", name)
	}

	names := []string{}
	rows, err := db.Query("SELECT name FROM user")
	if err != nil {
		log.Printf("query users failed: %v", err)
	}

	if err := rows.Scan(&names); err != nil {
		log.Printf("search users failed: %v", err)
	}
	log.Printf("search users: %v", names)
}

func ormQuery() {
	engine, err := orm.NewEngine("sqlite3", "geek.db")
	if err != nil {
		glog.Errorf("create orm engine failed: %v", err)
		return
	}
	defer engine.Close()

	s := engine.NewSession()
	s.Raw("DROP TABLE IF EXISTS user;").Exec()
	s.Raw("CREATE TABLE user(name text);").Exec()
	s.Raw("CREATE TABLE user(name text);").Exec()

	res, err := s.Raw("INSERT INTO user(name) values (?), (?)", "kallen", "torres").Exec()
	if err != nil {
		glog.Errorf("insert data err: %v", err)
	}
	n, _ := res.RowsAffected()
	glog.Infof("Exec success, %v affected", n)
}

func main() {
	// rawQuery()

	glog.Warn(glog.Yellow.Add(strings.Repeat("=", 20)))

	ormQuery()
}
