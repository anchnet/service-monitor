package http

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/51idc/service-monitor/smartAPI-proxy/g"
)

func configSmartAPIRoutes() {
	http.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		var (
			Ip         string
			agent      string
			method     string
			requestURL string
			proxyURL   *url.URL
		)
		if r.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}
		if g.Config().Debug {
			requestAddr := r.RemoteAddr
			requestIpPort := strings.Split(requestAddr, ":")
			Ip = requestIpPort[0]
			agent = r.UserAgent()
			method = r.Method
			requestURL = r.RequestURI
		}
		smartAPI_url := g.Config().SmartAPI.Url
		director, err := url.Parse(smartAPI_url)
		if err != nil {
			log.Println("url Parse Error: ", err)
			return
		}
		proxy := httputil.NewSingleHostReverseProxy(director)
		proxy.ServeHTTP(w, r)
		if g.Config().Debug {
			proxyURL = r.URL
			log.Println(Ip, agent, method, requestURL, proxyURL)
		}

	})

}
