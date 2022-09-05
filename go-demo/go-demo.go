package main

import (
	"database/sql"
	"encoding/json"
	"example/go-demo-api/index"
	"example/go-demo-api/model"
	"example/go-demo-api/responsemodel"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

//var users = []user{}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:1234@tcp(localhost:3306)/godb")

	if err != nil {
		panic(err)
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	//defer db.Close()

	router := gin.Default()
	router.GET("/", index.GoHello)
	router.GET("/user", getUsers)
	router.POST("/api/insert", insertUser)
	router.POST("/api/update", updateUser)
	router.POST("/api/delete", deleteUser)

	router.Run("localhost:8080")
}

// getUser responds with the list of all users as JSON.
func getUsers(c *gin.Context) {
	var user = []model.User{}

	// query all data
	rows, err := db.Query("select id, firstname, lastname, email, tel from user")
	ErrorCheck(err)

	defer rows.Close()
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.Id, &u.Firstname, &u.Lastname, &u.Email, &u.Tel)
		ErrorCheck(err)
		user = append(user, u)
		//log.Println(u.Id, u.Firstname, u.Lastname, u.Email, u.Tel)

	}
	ErrorCheck(err)

	var res = &responsemodel.Responsemodel{Code: "S200", Desc: "Success", Data: user}
	log.Println(res)
	//c.IndentedJSON(http.StatusOK, req)
	c.JSON(http.StatusOK, res)
	/*
		   convert json

		pigeon := &Bird{
			Species:     "Pigeon",
			Description: "likes to eat seed",
		}

		data, _ := json.Marshal(pigeon)
		fstruct := &Bird{}
		var d string = string(data)
		log.Println(d)
		json.Unmarshal([]byte(d), fstruct)
		c.IndentedJSON(http.StatusOK, fstruct)

	*/
}

func insertUser(c *gin.Context) {
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)

	var statusCheckData = CheckDataUser(user)
	var resp = &responsemodel.Responsemodel{}
	if statusCheckData == true {
		// INSERT INTO DB
		// prepare
		stmt, e := db.Prepare("insert into user(ID, FIRSTNAME, LASTNAME, EMAIL, TEL) values (?, ?, ?, ?, ?)")
		ErrorCheck(e)

		//execute
		res, e := stmt.Exec(user.Id, user.Firstname, user.Lastname, user.Email, user.Tel)
		ErrorCheck(e)

		id, e := res.LastInsertId()
		ErrorCheck(e)

		fmt.Println("Insert id", id)

		lastId, err := res.LastInsertId()

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("The last inserted row id: %d\n", lastId)
		if lastId != 0 {
			resp = &responsemodel.Responsemodel{Code: "S200", Desc: "Insert Success", Data: []model.User{}}
		} else {
			resp = &responsemodel.Responsemodel{Code: "E200", Desc: "Insert fail", Data: []model.User{}}
		}

	} else {
		resp = &responsemodel.Responsemodel{Code: "S200", Desc: "Email นี้มีอยู่ในระบบแล้วไม่สามารถลงทะเบียนซ้ำ !", Data: []model.User{}}
	}

	//c.IndentedJSON(http.StatusOK, req)
	c.JSON(http.StatusOK, resp)

}

func updateUser(c *gin.Context) {
	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)

	var resp = &responsemodel.Responsemodel{}

	//Update db
	stmt, e := db.Prepare("update user set FIRSTNAME = ?, LASTNAME = ?, EMAIL = ?, TEL = ? where ID=?")
	ErrorCheck(e)

	// execute
	res, e := stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Tel, user.Id)
	ErrorCheck(e)

	a, e := res.RowsAffected()
	ErrorCheck(e)

	fmt.Println(a)
	if a != 0 {
		resp = &responsemodel.Responsemodel{Code: "S200", Desc: "Update Success", Data: []model.User{}}
	} else {
		resp = &responsemodel.Responsemodel{Code: "E200", Desc: "Update fail", Data: []model.User{}}
	}

	c.IndentedJSON(http.StatusOK, resp)
}

func deleteUser(c *gin.Context) {

	var resp = &responsemodel.Responsemodel{}

	reqBody, _ := ioutil.ReadAll(c.Request.Body)
	var user model.User
	json.Unmarshal(reqBody, &user)

	// delete data
	stmt, e := db.Prepare("delete from user where ID=?")
	ErrorCheck(e)

	// delete set ID
	res, e := stmt.Exec(user.Id)
	ErrorCheck(e)

	// affected rows
	a, e := res.RowsAffected()
	ErrorCheck(e)
	fmt.Println(a) // 1
	//if res.LastInsertId()
	if a != 0 {
		resp = &responsemodel.Responsemodel{Code: "S200", Desc: "Delete Success", Data: []model.User{}}
	} else {
		resp = &responsemodel.Responsemodel{Code: "E200", Desc: "Delete fail", Data: []model.User{}}
	}

	c.JSON(http.StatusOK, resp)
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func CheckDataUser(user model.User) bool {
	var result = true
	var userObject = findUserByEmail(user.Email)
	if (model.User{}) != userObject {
		result = false
	}
	return result
}

func findUserByEmail(email string) model.User {
	var result model.User
	var user = model.User{}

	// query by email
	rows, err := db.Query("select id, firstname, lastname, email, tel from user where EMAIL = ?", email)

	ErrorCheck(err)

	defer rows.Close()
	for rows.Next() {
		var u model.User
		err := rows.Scan(&u.Id, &u.Firstname, &u.Lastname, &u.Email, &u.Tel)
		ErrorCheck(err)
		user = u
		//log.Println("+++++++++++++++++++++++++++++++++++++++++++++")
		//log.Println(u.Id, u.Firstname, u.Lastname, u.Email, u.Tel)

	}
	ErrorCheck(err)
	result = user
	return result
}
