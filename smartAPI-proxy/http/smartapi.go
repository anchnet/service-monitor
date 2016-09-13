package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/51idc/service-monitor/smartAPI-proxy/g"
)

func configSmartAPIRoutes() {
	http.HandleFunc("/api/service/version", func(w http.ResponseWriter, req *http.Request) {
		if req.ContentLength == 0 {
			http.Error(w, "body is blank", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(req.Body)
		var smartapi_data smartAPI_Data
		var smartapi_result smartAPI_Result
		var result smartAPI_Result

		err := decoder.Decode(&smartapi_data)
		if err != nil {
			http.Error(w, "connot decode body", http.StatusBadRequest)
			return
		}
		smartAPI_url := g.Config().SmartAPI.Url

		if g.Config().Debug {
			log.Println("Get smartAPI Request: ", smartapi_data)
		}
		body, res, err := sendData(smartAPI_url, smartapi_data)

		if err != nil {
			log.Println("smartAPI err ", err)
			result.Status = "error"
			result.Message = err.Error()
			js, _ := json.Marshal(result)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write(js)
			return
		}

		if res != 200 {
			log.Println("smartAPI error,statcode= ", res)
			result.Status = "error"
			result.Message = "smartAPI error,statcode=" + string(res)
			js, _ := json.Marshal(result)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write(js)
			return
		}

		err = json.Unmarshal(body, &smartapi_result)

		if err != nil {
			log.Println("json_Unmarshal_error", err)
			result.Status = "error"
			result.Message = "json_Unmarshal_error"
			js, _ := json.Marshal(result)
			w.Header().Set("Content-Type", "application/json; charset=UTF-8")
			w.Write(js)
			return
		}
		if smartapi_result.Status != "ok" {
			log.Println("SmartAPI return error: ", smartapi_result.Message)
		}
		if g.Config().Debug {
			log.Println("Get smartAPI Result: ", smartapi_result)
		}
		js, _ := json.Marshal(smartapi_result)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Write(js)

	})
}
