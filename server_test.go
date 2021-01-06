package masque

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("uses the original http.Handler for non-CONNECT-UDP requests", func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(418)
		})
		h := HandleMASQUE(mux)

		req, err := http.NewRequest(http.MethodGet, "/", nil)
		Expect(err).ToNot(HaveOccurred())

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(418))
	})

	It("errors if the flow ID header is missing", func() {
		h := HandleMASQUE(http.DefaultServeMux)

		req, err := http.NewRequest(methodConnectUDP, "/", nil)
		Expect(err).ToNot(HaveOccurred())

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(400))
	})

	It("errors if the flow ID header can't be parsed", func() {
		h := HandleMASQUE(http.DefaultServeMux)

		req, err := http.NewRequest(methodConnectUDP, "/", nil)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Add(flowIDHeader, "\"")

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(400))
	})

	It("errors if the flow ID header is not a number", func() {
		h := HandleMASQUE(http.DefaultServeMux)

		req, err := http.NewRequest(methodConnectUDP, "/", nil)
		Expect(err).ToNot(HaveOccurred())
		req.Header.Add(flowIDHeader, "foobar")

		rr := httptest.NewRecorder()
		h.ServeHTTP(rr, req)
		Expect(rr.Code).To(Equal(400))
	})
})
