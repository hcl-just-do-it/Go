package controller

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
)

type Users struct {
	Username string
	Password string
}
type Student struct {
	Id   string
	Name string
}
type TestDBResponse struct { //Response（里面有StatusCode,StatusMsg）
	Response
	User          Users      `json:"user"` //左边是变量名字，右边是变量类型，如  a  int。下面返回值里的逻辑是:b=1;	a : b
	Users         []Users    `json:"Users"`
	TestStructHcl TestStruct `json:"test-struct"`
	Student       Student    `json:"student"`
}
type TestResponse struct { //Response（里面有StatusCode,StatusMsg）
	Response
	CommentList   []Comment  `json:"comment_list,omitempty"`
	UserId        int64      `json:"user_id,omitempty"`
	TestStructHcl TestStruct `json:"test-struct"`
}

func TestGorm(c *gin.Context) { //"gorm.io/driver/mysql"
	fmt.Println("TestGorm函数：")
	db, err := gorm.Open(
		mysql.Open("root:root@tcp(127.0.0.1:3306)/mydb"),
	)
	if err != nil {
		fmt.Println(err)
	}
	//db.SingularTable(true)
	var student Student
	db.First(&student, 1) // 根据整型主键查找
	//err = db.Select("id", "name").Find(&student, "1").Error
	c.JSON(http.StatusOK, TestDBResponse{
		Response: Response{StatusCode: 0, StatusMsg: "TestGorm成功"},
		Student:  student,
	})
}

func Test(c *gin.Context) {

	fmt.Println("func Test!")

	c.JSON(http.StatusOK, TestResponse{
		Response:      Response{StatusCode: 0, StatusMsg: "Test成功"},
		CommentList:   DemoComments,
		UserId:        123456,
		TestStructHcl: DemoTest,
	})
}
func TestDB2(c *gin.Context) { // github.com/jmoiron/sqlx
	db, err := sqlx.Connect("mysql", "root:root@(localhost:3306)/mydb")
	if err != nil {
		fmt.Printf("connect DataBase failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	fmt.Printf("connect DataBase success\n")
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
	fmt.Printf("username=%v,password=%v\r\n", username, password)
	sql := "SELECT * FROM user where username='" + username + "' and password ='" + password + "'"
	rows, err := db.Query(sql)
	defer rows.Close()
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.Username, &user.Password)
		if err != nil {
			return
		}
		fmt.Printf("%v", user)
		c.JSON(http.StatusOK, TestDBResponse{
			Response: Response{StatusCode: 0, StatusMsg: "TestDB成功"},
			User:     user,
		})
	}
}

func TestPOST(c *gin.Context) { //"database/sql"
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/mydb")
	//username := "hcl"
	username, _ := c.GetPostForm("username")
	password, _ := c.GetPostForm("password")
	rows, err := db.Query("select * from user where username=? and password = ?", username, password)
	if err != nil {
		fmt.Println(err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(rows)
	var users []Users
	for rows.Next() {
		var user Users
		err := rows.Scan(&user.Username, &user.Password)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}
	fmt.Printf("%v", users)
	if rows.Err() != nil {
		fmt.Println(rows.Err())
	}
	c.JSON(http.StatusOK, TestDBResponse{
		Response: Response{StatusCode: 0, StatusMsg: "TestDB成功"},
		Users:    users,
	})
}
