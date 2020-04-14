package http_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/gorilla/mux"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	wolf "github.com/realOkeani/wolf-dynasty-api"
	. "github.com/realOkeani/wolf-dynasty-api/http"
)

var _ = Describe("Health", func() {
	var (
		svr    *httptest.Server
		uri    string
		router *mux.Router
	)

	BeforeEach(func() {
		router = mux.NewRouter().StrictSlash(true)

		AddRoutes(wolf.Services{}, router)
		svr = httptest.NewServer(router)
		uri = svr.URL + "/health"
	})

	AfterEach(func() {
		svr.Close()
	})

	Describe("GET", func() {
		var (
			req *http.Request
			res *http.Response
			err error
		)

		BeforeEach(func() {
			req, _ = http.NewRequest("GET", uri, nil)
		})

		JustBeforeEach(func() {
			res, err = http.DefaultClient.Do(req)
		})

		Context("when a health check is performed", func() {
			Context("and the application is down", func() {
				BeforeEach(func() {
					svr.Close()
				})

				It("returns an error", func() {
					Expect(err).To(HaveOccurred())
				})
			})

			Context("and the application is running", func() {
				It("returns status OK", func() {
					Expect(err).ToNot(HaveOccurred())
					Expect(res.StatusCode).To(Equal(http.StatusOK))
				})
			})
		})
	})
})

