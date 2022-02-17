package otusdb

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-sql/sqlexp"
	"os"
)

var Db *sql.DB

var Quoter sqlexp.Quoter

func InitDb() {

	dbHost, _ := os.LookupEnv("DB_HOST")
	dbUser, _ := os.LookupEnv("DB_USER")
	dbPassword, _ := os.LookupEnv("DB_PASSWORD")
	dbName, _ := os.LookupEnv("DB_NAME")
	dbPort, _ := os.LookupEnv("DB_PORT")
	db, err := sql.Open("mysql", dbUser+":"+dbPassword+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	db.SetMaxOpenConns(10)
	if err != nil {
		panic(err)
	}
	Db = db
	Quoter, err = sqlexp.QuoterFromDriver(Db.Driver(), context.Background())
	if err != nil {
		Quoter = TmpMysqlQuoter{}
	}
}

func CloseDb() error {
	return Db.Close()
}

type TmpMysqlQuoter struct{}

func (p TmpMysqlQuoter) ID(name string) string {
	return name
}

func (p TmpMysqlQuoter) Value(v interface{}) string {
	return escape(v.(string))
}

func escape(source string) string {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}
