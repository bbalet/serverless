package counter

// go-sql-driver features an automatic connection-pooling
import (
	"fmt"
)

// Handle a serverless request
func GetCounterOpenFaas(req []byte) string {
	db := openDB()
	defer db.Close()
	counterValue := concurrentExample(db)
	return fmt.Sprintf("%s", string(counterValue))
}
