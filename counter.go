package counter

// go-sql-driver features an automatic connection-pooling
import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"

	_ "github.com/go-sql-driver/mysql"
)

// sequentialExample gets and increments the counter of visits in
// a sequential fashion without concurrency
// returns the counter
func sequentialExample(db *sql.DB) int {
	// get the number of visitors
	counterValue := getCounterValue(db)
	// increment the counter of visitors
	incrementCounter(db)
	log.Printf("sequentialExample / Counter value: %v", counterValue)
	return counterValue
}

// concurrentExample gets and increments the counter of visits by
// using concurrent programming
// returns the counter
func concurrentExample(db *sql.DB) int {
	// we'll need to wait for the two goroutines to finnish
	var wg sync.WaitGroup
	wg.Add(2)

	// get the number of visitors
	var counterValue int
	go func() {
		counterValue = getCounterValue(db)
		wg.Done()
	}()
	// increment the counter of visitors
	go func() {
		incrementCounter(db)
		wg.Done()
	}()

	wg.Wait()
	log.Printf("concurrentExample / Counter value: %v", counterValue)
	return counterValue
}

// getCounterValue gets the value of the counter
// db is a reference to an opened database
// It returns the counter value
func getCounterValue(db *sql.DB) int {
	counterValue := 0
	err := db.QueryRow("SELECT value FROM counter").Scan(&counterValue)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Counter value: %v", counterValue)
	return counterValue
}

// incrementCounter increments the value of the visits counter
// db is a reference to an opened database
func incrementCounter(db *sql.DB) {
	update, err := db.Query("UPDATE counter SET value = value + 1")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer update.Close()
	log.Println("The counter has been incremented")
}

// openDB opens a connection to a database by using
// the environment variables
// returns a reference to a database connection
func openDB() *sql.DB {
	userDB := os.Getenv("MYSQL_USER")
	pwdDB := os.Getenv("MYSQL_PWD")
	nameDB := os.Getenv("MYSQL_DBNAME")
	addrDB := os.Getenv("MYSQL_ADDRESS")
	portDB := os.Getenv("MYSQL_PORT")
	connStrDB := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", userDB, pwdDB, addrDB, portDB, nameDB)
	//logging password is a bad practise
	//log.Printf("Connection string:%v\n", connStrDB)

	db, err := sql.Open("mysql", connStrDB)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Connected to the MySQL instance")
	}
	return db
}
