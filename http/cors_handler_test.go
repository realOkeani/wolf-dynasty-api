package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	. "github.com/realOkeani/wolf-dynasty-api/http"
)

var _ = Describe("CorsHandler", func() {
	handler := CorsHandler(ghttp.RespondWith(200, "ok", nil))
	var (
		req *http.Request
		res *httptest.ResponseRecorder
	)
	BeforeEach(func() {
		req = httptest.NewRequest("GET", "/", nil)
		res = httptest.NewRecorder()
	})
	JustBeforeEach(func() {
		handler.ServeHTTP(res, req)
	})

	Describe("access-control-allow-origin", func() {
		It("sets to * when no origin is given", func() {
			Expect(res.Result().Header.Get("access-control-allow-origin")).Should(
				Equal("*"))
		})
	})

	Describe("access-control-allow-headers", func() {
		BeforeEach(func() {
			req.Method = http.MethodOptions
		})
		It("sets to `authorization` for an OPTIONS request", func() {
			Expect(res.Result().Header.Get("access-control-allow-headers")).Should(
				Equal("authorization"))
		})
	})

	Describe("access-control-allow-methods", func() {
		BeforeEach(func() {
			req.Method = http.MethodOptions
		})
		It("sets to POST,PUT,GET,DELETE,OPTIONS for an OPTIONS request", func() {
			methods := strings.Split(
				strings.Replace(
					res.Result().Header.Get("access-control-allow-methods"),
					", ", ",", -1),
				",")
			Expect(methods).Should(ConsistOf("DELETE", "OPTIONS", "POST", "GET", "PUT", "PATCH"))
		})
	})
})

