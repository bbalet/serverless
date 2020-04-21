package counter

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

// Benchmark of the sequential access
func BenchmarkSequential(b *testing.B) {
	db := openDB()
	defer db.Close()

	// run the sequential example function b.N times
	for n := 0; n < b.N; n++ {
		sequentialExample(db)
	}
}

// Benchmark of the concurrent access
func BenchmarkConcurrent(b *testing.B) {
	db := openDB()
	defer db.Close()

	// run the concurrent example function b.N times
	for n := 0; n < b.N; n++ {
		concurrentExample(db)
	}
}

// Test the concurrent programing model
func TestConcurrent(t *testing.T) {
	db := openDB()
	defer db.Close()
	initialValue := getCounterValue(db)
	incrementCounter(db)
	counterValue := concurrentExample(db)
	if counterValue < initialValue {
		t.Errorf("TestConcurrent was incorrect, got: intial=%v and final=%v.", initialValue, counterValue)
	}
}

// Test the sequential programing model
func TestSequential(t *testing.T) {
	db := openDB()
	defer db.Close()
	initialValue := getCounterValue(db)
	incrementCounter(db)
	counterValue := sequentialExample(db)
	if counterValue < initialValue {
		t.Errorf("TestSequential was incorrect, got: intial=%v and final=%v.", initialValue, counterValue)
	}
}

// Test of the prinitives functions
func TestPrimitives(t *testing.T) {
	db := openDB()
	defer db.Close()
	initialValue := getCounterValue(db)
	incrementCounter(db)
	counterValue := getCounterValue(db)
	if counterValue != (initialValue + 1) {
		t.Errorf("Test of primitives was incorrect, got: intial=%v and final=%v.", initialValue, counterValue)
	}
}
