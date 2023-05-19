package main

import (
	"crypto/md5"
	"crypto/rc4"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)



var addressBook  = map[int][]string{
	0: {"15501901520", "10289996646", "18439933215", "10191038147", "16190421735", "12143338624", },
	1: {"12305388723", "19733995153", "18439933215", "19614047237", "15199568833", },
	2: {"13797758294", "16700516379", "10289996646", "11750221901", "14388635671", "10217737207", },
	3: {"13341480707", "14388635671", "12310521340", },
	4: {"10289996646", "11750221901", "13341480707", "12681953709", "14388635671", "14453732679", },
	5: {"12681953709", "16115770066", "14036782017", "10217737207", "17929447200", "19614047237", },
	6: {"18439933215", "17275693007", "14900139871", "13797758294", "10651024818", "10191038147", },
	7: {"18336810545", "13731127217", "19733995153", "15199568833", "14388635671", "13797758294", },
	8: {"13731127217", "18336810545", "12594905378", "13941873216", "14453732679", "14900139871", },
	9: {"18336810545", "14882230442", "14388635671", "10268697568", "13252411646", "14076297572", },
	10: {"16411322012", "10217737207", "14453732679", "10339076498", "12681953709", "17436289978", },
	11: {"17322562212", "10554289923", "14076297572", "16700516379", "15501901520", "12948378676", },
	12: {"10289996646", "15080757553", "10191038147", "14882230442", "10268697568", "11843465773", },
	13: {"13252411646", "12594905378", "16029695919", "13252411646", "12594905378", },
	14: {"13252411646", "16700516379", },
	15: {"16029695919", "17322562212", "15080757553", "10268697568", },
	16: {"18336810545", "13797758294", "13731127217", },
	17: {"17929447200", },
	18: {"16411322012", "12143338624", "10191038147", "16190421735", "14900139871", "17076184312", },
	19: {"18336810545", "12305388723", "19733995153", "11454924698", "10330455483", "10651024818", },
	20: {"15501901520", "10554289923", "12948378676", "11454924698", "15199568833", "11843465773", },
	21: {"13252411646", "13731127217", "11454924698", "10554289923", },
	22: {"18439933215", },
	23: {"15501901520", "10798900163", "17436289978", "17076184312", "17436289978", "10289996646", },
	24: {"17124359578", "14036782017", "13739507274", "15080757553", "10289996646", },
	25: {"17322562212", "19614047237", "18988933231", },
	26: {"12305388723", "16115770066", "17275693007", "12305388723", "10651024818", "17322562212", },
	27: {"15501901520", "18988933231", "12681953709", "16029695919", "13731127217", "15199568833", },
	28: {"10268697568", "10798900163", "11454924698", "14036782017", },
	29: {"18439933215", "19614047237", "10651024818", "13941873216", "19614047237", "10554289923", },
	30: {"19733995153", "13731127217", "12143338624", },
	31: {"16029695919", "17124359578", "15080757553", "14882230442", "12143338624", "18988933231", },
	32: {"14882230442", "10268697568", "12594905378", "13739507274", "15199568833", },
	33: {"18336810545", "13252411646", "16411322012", "12310521340", "12948378676", "14036782017", },
	34: {"12681953709", "16029695919", "13341480707", },
	35: {"17275693007", "14076297572", "14453732679", "13797758294", },
	36: {"13797758294", "16411322012", "10191038147", "12310521340", "11750221901", "10651024818", },
	37: {"19733995153", "17124359578", "15080757553", "18988933231", "19614047237", "13341480707", },
	38: {"16411322012", "17275693007", "16115770066", "14453732679", },
	39: {"17275693007", "14388635671", "14453732679", "10217737207", "11454924698", "14076297572", },
	40: {"13341480707", "12310521340", "15501901520", "14036782017", "10554289923", },
	41: {"13941873216", "17929447200", "14076297572", "11750221901", },
	42: {"12305388723", "13739507274", "12143338624", },
	43: {"16190421735", "10268697568", "15199568833", "19733995153", "17124359578", "10191038147", },
	44: {"16411322012", "17929447200", "10339076498", "17322562212", "14882230442", "16029695919", },
	45: {"14036782017", "12310521340", "12305388723", "12948378676", "17076184312", },
	46: {"14882230442", "11843465773", "12948378676", "17929447200", "17275693007", "10339076498", },
	47: {"18439933215", "16700516379", "15080757553", "12310521340", "17929447200", "17076184312", },
	48: {"17436289978", "17436289978", "12143338624", "11750221901", },
	49: {"12681953709", "18853090573", "13739507274", "13341480707", "17322562212", "11454924698", },
}

var (
	dbHost string
	dbUser string
	dbPwd string
	dbName string
	dbPort string
)

func init() {
	flag.StringVar(&dbHost, "h", "", "set the host of database")
	flag.StringVar(&dbUser, "u", "", "set the username to connect database")
	flag.StringVar(&dbPwd, "p", "", "set the password of user")
	flag.StringVar(&dbName, "n", "", "set the database name to connect")
	flag.StringVar(&dbPort, "P", "", "set the port to connect database")
}

func main() {

	flag.Parse()
	checkFlag()

	router := gin.Default()

	authRouter := router.Group("/auth")
	authRouter.Use(AuthMiddleware())
	authRouter.GET("/showaddressbook", showAddressBook)

	nologinRouter := router.Group("/nologin")
	nologinRouter.GET("/userinfo", queryUserInfoByName)

	router.Run(":8080")
}

//检查命令行参数是否正确设置
func checkFlag() {
	if dbHost == "" {
		log.Println("please set the host to database")
		os.Exit(1)
	}

	if dbPort == "" {
		log.Println("please set the port to connect database")
		os.Exit(1)
	}

	if dbName == "" {
		log.Println("please set the database name")
		os.Exit(1)
	}

	if dbUser == "" || dbPwd == "" {
		log.Println("please set the user and password to connect database")
		os.Exit(1)
	}
}

//鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.Query("userid")
		password := ctx.Query("password")

		if password == "" || userId == "" {
			ctx.String(http.StatusUnauthorized, "please specify userid and username")
			ctx.Abort()
			return
		}
		if ok := checkLogin(userId, password); !ok {
			ctx.String(http.StatusUnauthorized, "username or password is wrong")
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

//检查用户名和参数是否正确
func checkLogin(userId, password string) bool {
	//建立数据库连接
	DBSrcName := getDBSourceName()
	db, err := sql.Open("mysql", DBSrcName)
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

	//计算password的md5值
	temp := md5.Sum([]byte(password))
	md5Pwd := hex.EncodeToString(temp[:])

	if md5Pwd == storedPwd {
		return true
	} else {
		return false
	}
}

type userInfoResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Msg   string `json:"message"`
}

func queryUserInfoByName(ctx *gin.Context) {
	userName := ctx.Query("username")
	if userName == "" {
		log.Print("username not set")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username not set"})
		return
	}

	//打开数据库连接
	DBSrcName := getDBSourceName()
	db, err := sql.Open("mysql", DBSrcName)
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

		extraMsg := "can you help me to contact Christopher, the flag is sctf{md5(find-path)}, find-path's format is like tom->jack->alice"

		userInfo := userInfoResp{
			ID: userId,
			Name: userName,
			Phone: userPhone,
		}
		if userId == 14 {
			userInfo.Msg = extraMsg
		}

		userInfoJSON, err := json.Marshal(userInfo)
		if err != nil {
			log.Println(err)
			ctx.String(http.StatusInternalServerError, "some internal error happend")
			ctx.Abort()
		}

		cryptoResp := cryptoJSONByte(userInfoJSON)
		ctx.Data(http.StatusOK, "application/octet-stream", cryptoResp)
	}
}

type contact struct {
	Phones []string `json:"phones"`
}

func showAddressBook(ctx *gin.Context) {
	userId, _ := strconv.Atoi(ctx.Query("userid"))
	phones := addressBook[userId]

	respContact := contact{
		Phones: phones,
	}
	
	respJSON, err := json.Marshal(respContact)
	if err != nil {
		log.Println(err)
		ctx.String(http.StatusInternalServerError, "some internal error happend")
		ctx.Abort()
	}

	cryptoData := cryptoJSONByte(respJSON)
	ctx.Data(http.StatusOK, "application/octet-stream", cryptoData)
}

func getDBSourceName() string {
	DBSrcName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPwd, dbHost, dbPort, dbName)
	return DBSrcName
}

func cryptoJSONByte(jsonByte []byte) []byte {
	key := []byte("sctf2023")
	rc4Enc, err := rc4.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil
	}

	crypto := make([]byte, len(jsonByte))
	rc4Enc.XORKeyStream(crypto, jsonByte)

	return crypto
}

