package main

import (
	"net/http"
	"strconv"
	"strings"

	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// type userInfo struct {
// 	ID    int    `json:"id"`
// 	Name  string `json:"name"`
// 	Phone string `json:"phone"`
// }

var addressBook map[int][]string

func init() {
	addressBook = make(map[int][]string)
	addressBook[1000] = []string{"1001", "1002"}
	addressBook[1001] = []string{"1002"}
	addressBook[1002] = []string{"1000"}
}

func main() {
	router := gin.Default()

	authRouter := router.Group("/auth")
	authRouter.Use(AuthMiddleware())
	authRouter.GET("/showaddressbook", showAddressBook)

	nologinRouter := router.Group("/nologin")
	nologinRouter.GET("/userinfo", queryUserInfoByName)

	router.Run(":8080")
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query("userid")
		password := ctx.Query("password")

		if password == "" || userId == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error":"please specify userid and username"})
			return
		}
		if ok := checkLogin(userId, password); !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "username or password is wrong"})
			return
		}
		ctx.Next()
	}
}

func checkLogin(userId, password string) bool {
	//建立数据库连接
	db, err := sql.Open("mysql", "root:yiqiu666...@tcp(127.0.0.1:3306)/sctf_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	stmtout, err := db.Prepare("SELECT password FROM users WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmtout.Close()

	var storedPwd string
	err = stmtout.QueryRow(userId).Scan(&storedPwd)
	if err != nil {
		log.Fatal(err)
	}

	if password == storedPwd {
		return true
	} else {
		return false
	}
}

func queryUserInfoByName(ctx *gin.Context) {
	userName := ctx.Query("username")
	if userName == "" {
		log.Print("username not set")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username not set"})
		return
	}

	//打开数据库连接
	db, err := sql.Open("mysql", "root:yiqiu666...@tcp(127.0.0.1:3306)/sctf_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("select id, name, phone from info where name='" + userName + "'")
	if err != nil {
		log.Print(err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userId int
		var userName string
		var userPhone string
		err = rows.Scan(&userId, &userName, &userPhone)
		if err != nil {
			log.Print(err)
			return
		}

		log.Print(userId, userName, userPhone)

		extraMsg := "can you help me to find alice, the flag is sctf{md5(find-path)}, find-path's format is like tom->jack->alice"

		if userName == "tom" {
			ctx.JSON(http.StatusOK, gin.H{
				"id":userId,
				"name":userName,
				"phone":userPhone,
				"message":extraMsg,
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"id":userId,
				"name":userName,
				"phone":userPhone,
			})
		}
	}
}

func showAddressBook(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("userid"))
	numbers := addressBook[userId]
	ctx.String(http.StatusOK, strings.Join(numbers, ","))
}
