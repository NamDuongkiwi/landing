package manager
import (
"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

func Connect() (db *sql.DB) {
	user := viper.Get("mysql.user")
	password := viper.Get("mysql.password")
	host := viper.Get("mysql.host")
	database := viper.Get("mysql.database")
	//dsn := user + ":" + password + "@tcp(" + host + ")/" + database
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v", user, password, host, database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic("Could not connect to database")
	}
	return db
}
