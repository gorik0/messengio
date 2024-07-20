package error

import (
	"log"
	"testing"
)

func HandlerError(err error, msg string) {

	if err != nil {
		log.Fatalf("ERR :::: %s  ::: %v \n", msg, err)
	}
}
func TestHandlerError(t *testing.T, err error, msg string) {

	if err != nil {
		t.Fatalf("ERR :::: %s  ::: %v \n", msg, err)
	}
}
