// Package helloworld provides a set of Cloud Functions samples.
package counter

import (
	"fmt"
	"net/http"
)

// GetCounter is an HTTP Cloud Function with a request parameter.
func GetCounter(w http.ResponseWriter, r *http.Request) {
	db := openDB()
	defer db.Close()
	counterValue := concurrentExample(db)
	fmt.Fprintf(w, "%v", counterValue)
}
