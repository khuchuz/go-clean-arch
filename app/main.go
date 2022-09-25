package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"

	_userHttpDelivery "github.com/khuchuz/go-clean-arch/user/delivery/http"
	_userHttpDeliveryMiddleware "github.com/khuchuz/go-clean-arch/user/delivery/http/middleware"
	_userRepo "github.com/khuchuz/go-clean-arch/user/repository/mysql"
	_userUcase "github.com/khuchuz/go-clean-arch/user/usecase"
)

func main() {
	dbHost := "localhost"
	dbPort := "3306"
	dbUser := "root"
	dbPass := ""
	dbName := "project1"
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())
	dbConn, err := sql.Open(`mysql`, dsn)

	if err != nil {
		log.Fatal(err)
	}
	err = dbConn.Ping()
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := dbConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	middL := _userHttpDeliveryMiddleware.InitMiddleware()
	e.Use(middL.CORS)
	ar := _userRepo.NewMysqlUserRepository(dbConn)

	timeoutContext := time.Duration(2) * time.Second
	au := _userUcase.NewUserUsecase(ar, timeoutContext)
	_userHttpDelivery.NewUserHandler(e, au)

	log.Fatal(e.Start(":9090"))
}
