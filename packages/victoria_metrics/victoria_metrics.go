package victoria_metrics

import (
	"fmt"
	"time"
	"bytes"
	"errors"
	"net/url"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
)

type VictoriaMetrics struct {
	Config Config
	Client *http.Client
}

// New creates a new GoVictoria instance
func New(config Config) *VictoriaMetrics {
	return &VictoriaMetrics{
		Config: config,
		Client: &http.Client{},
	}
}

// SendMetrics sends the metrics to VictoriaMetrics
func (g *VictoriaMetrics) SendMetrics(requests []VictoriaMetricsRequest) error {
	if len(requests) == 0 {
		return errors.New("No requests to send")
	}

	// Loop through the request and build the body
	body := ""
	for i, requestBody := range requests {
		jsonRequest, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}

		body += string(jsonRequest)

		if i != len(requests) - 1 {
			body += "\n"
		}
	}

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("POST", g.Config.URL+"/api/v1/import", bytes.NewBuffer([]byte(body)))
	request.Header.Add("Authorization", "Basic "+BasicAuth(g.Config.Username, g.Config.Password))
	request.Header.Add("User-Agent", "Vonkje (github.com/GJSBRT/vonkje)")

	// Send the request to Victoria Metrics
	response, err := g.Client.Do(request)
	if err != nil {
		return err
	}

	// Close the response body
	defer response.Body.Close()

	// Check if the status code is not 204
	if response.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Victoria Metrics returned a non-200 status code: %d", response.StatusCode))
	}

	return nil
}

// QueryTimeRange queries Victoria Metrics for metrics in a time range
func (g *VictoriaMetrics) QueryTimeRange(promql string, startTime time.Time, endTime time.Time, step string) (VictoriaMetricsQueryResponse, error) {
	// Check if the start time is before the end time
	if startTime.After(endTime) {
		return VictoriaMetricsQueryResponse{}, errors.New("Start time must be before end time")
	}

	// Add the query parameters to the request
	params := url.Values{}
	params.Add("query", promql)
	params.Add("start", strconv.FormatInt(startTime.Unix(), 10))
	params.Add("end", strconv.FormatInt(endTime.Unix(), 10))
	params.Add("step", step)

	url := g.Config.URL + "/api/v1/query_range?" + params.Encode()

	// Create the request to Victoria Metrics
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Add the query parameters to the request
	request.Header.Add("Authorization", "Basic "+BasicAuth(g.Config.Username, g.Config.Password))
	request.Header.Add("User-Agent", "Vonkje (github.com/GJSBRT/vonkje)")

	// Send the request to Victoria Metrics
	response, err := g.Client.Do(request)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Close the response body
	defer response.Body.Close()

	// Check if the status code is not 200
	if response.StatusCode != http.StatusOK {
		return VictoriaMetricsQueryResponse{}, errors.New(fmt.Sprintf("Victoria Metrics returned a non-200 status code: %d", response.StatusCode))
	}

	// Read the response body
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	// Unmarshal the response
	var metrics VictoriaMetricsQueryResponse
	err = json.Unmarshal([]byte(body), &metrics)
	if err != nil {
		return VictoriaMetricsQueryResponse{}, err
	}

	return metrics, nil
}

// BasicAuth returns the base64 encoded string for basic auth
func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
