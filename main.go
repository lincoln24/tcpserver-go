package main

import (	
	"fmt"
	"log"
	"time"
    "database/sql"
	"tcpserver/globa"
	"tcpserver/tcpserver"
)

func main() {
	initSql()
	go tcpserver.ServerMain()

	for{
		time.Sleep(1e11)
	}
}


func initSql() {
    dbTmp, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/gangji?charset=utf8")
    if err != nil {
        log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
 	global.Db = dbTmp
    global.Db.SetMaxIdleConns(1000)
 
    err = global.Db.Ping() // This DOES open a connection if necessary. This makes sure the database is accessible
    if err != nil {
        log.Fatalf("Error on opening database connection: %s", err.Error())
    }else{
    	fmt.Println("database open success")
    }
}