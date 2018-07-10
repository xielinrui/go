package util

import (
	"database/sql"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"time"
)


func SaveYongJin(money float64,ownerId int) (bool,error) {
	db, err := sql.Open("mysql", "username:password@tcp(127.0.0.1:3306)/gamelog?parseTime=true")

	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()
	now :=time.Now()
	ye := now.Format("2006-01-02 03:04:05")
	fmt.Println(money)
	fmt.Println(ye)
	fmt.Println(ownerId)
	zhangMoney :=money*100
	que,err :=db.Prepare("update commision set commission =commission+?,charge_time=? where uid=?")
	if err != nil {
		fmt.Println(err)
		return false,err
	}
	d,erro := que.Exec(zhangMoney,ye,ownerId)
	if erro != nil{
		return false,erro
	}
	affect,err :=d.RowsAffected()
	if err != nil {
		//fmt.Println(err)
		return false,err
	}
	fmt.Println(affect)
	return true,err
}
