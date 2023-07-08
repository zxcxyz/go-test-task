package main

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"io/ioutil"
	"log"
	"strings"
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"strconv"
	"github.com/prometheus/client_golang/prometheus"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
		totalRequests.WithLabelValues(path).Inc()

		timer.ObserveDuration()
	})
}

type Message struct {
	Username string
    Name string
    PublicRepositories int
    Followers int
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("\n got / request\n")

	_, err := io.WriteString(w, "This is my website!\n")

	if err != nil {
        fmt.Println(err.Error())
    }
}

func getGithub(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path[len("/github/"):]

	parcedURL := strings.Split(urlPath, "/")

	if parcedURL[1] == "repositories" {
		url := "https://api.github.com/users/" + parcedURL[0]

		resp, err := http.Get(url)

		// We Read the response body on the line below.
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("\n request body %s \n", body)

		var result map[string]interface{}
		err = json.Unmarshal(body, &result)

		if err != nil {
			// Prints the error if not nil
			fmt.Println("Error while decoding the data", err.Error())
		}

		m := Message{
			parcedURL[0],
			result["name"].(string), 
			int(result["public_repos"].(float64)),
			int(result["followers"].(float64)),
		}

		if err != nil {
		log.Fatalln(err)
		}

		w.Header().Set("Content-Type", "application/json")
		
		err = json.NewEncoder(w).Encode(m)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		_, err := io.WriteString(w, "Unsupported path:" + parcedURL[1] + "\n" )
		
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}

func main() {
	serverPort := ":" + os.Getenv("SERVER_PORT")
	metricsEndpoint := os.Getenv("METRICS_ENDPOINT")
	fmt.Printf(serverPort + "\n")
	fmt.Printf(metricsEndpoint + "\n")

	router := mux.NewRouter()
	router.Use(prometheusMiddleware)
	router.Path(metricsEndpoint).Handler(promhttp.Handler())
	router.HandleFunc("/github/{username}/{else}", getGithub)
	router.HandleFunc("/", getRoot)

	http.Handle("/", router)

	err := http.ListenAndServe(serverPort, router)

	log.Fatal(err)
}
