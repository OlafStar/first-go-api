package advertisment

import (
	"database/sql"
	"fmt"
	"net/http"
)

func GetAdvertisment(db *sql.DB, w http.ResponseWriter, r*http.Request) {
	fmt.Fprint(w, "Get advertisment")
}

func CreateAdvertisment(db *sql.DB, w http.ResponseWriter, r*http.Request) {
	fmt.Fprint(w, "Create advertisment")
}