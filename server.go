package masque

import (
	"fmt"
	"net/http"

	"github.com/dunglas/httpsfv"
)

const flowIDHeader = "Datagram-Flow-Id"

func HandleMASQUE(h http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT-UDP" {
			h.ServeHTTP(w, r)
			return
		}
		// TODO: check for the masque scheme
		flowIDItem, err := httpsfv.UnmarshalItem(r.Header[flowIDHeader])
		if err != nil {
			w.WriteHeader(400)
		}
		flowID, ok := flowIDItem.Value.(int64)
		if !ok {
			w.WriteHeader(400)
		}
		fmt.Println("Flow ID:", flowID)
		w.WriteHeader(200)
	})
	return mux
}
