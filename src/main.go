package zetton

import (
	"net/http"
)

func init() {
	http.HandleFunc("/spaces", spacesHandler)
}
