package masque

import (
	"fmt"
	"net/http"

	"github.com/dunglas/httpsfv"
)

const (
	flowIDHeader     = "Datagram-Flow-Id"
	methodConnectUDP = "CONNECT-UDP"
)

func HandleMASQUE(h http.Handler) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != methodConnectUDP {
			h.ServeHTTP(w, r)
			return
		}
		// TODO: check for the masque scheme
		// TODO(#4): This header is only used when using datagram-based MASQUE.
		flowIDRaw, ok := r.Header[flowIDHeader]
		if !ok {
			w.WriteHeader(400)
			return
		}
		flowIDItem, err := httpsfv.UnmarshalItem(flowIDRaw)
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
