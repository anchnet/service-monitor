package funcs

var Mertics = map[string]string{
	"Total Bytes Received":   "COUNTER",
	"Total Bytes Sent":       "COUNTER",
	"Total Delete Requests":  "COUNTER",
	"Total Get Requests":     "COUNTER",
	"Total Post Requests":    "COUNTER",
	"Total Put Requests":     "COUNTER",
	"Total Not Found Errors": "COUNTER",
	"Maximum Connections":    "GUAGE",
	"Current Connections":    "GUAGE",
	"Service Uptime":         "GUAGE",
}
