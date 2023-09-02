package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"regexp"
	"strings"
)

type SourceFund struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var cex = []string{
	"Binance",
	"Huobi",
	"Coinbase",
	"Kraken",
	"Bitfinex",
	"Bitstamp",
	"Bithumb",
	"Bittrex",
	"Gemini",
	"OKEx",
	"KuCoin",
	"FTX",
	"BitMart",
	"Gate.io",
	"Bitforex",
	"Poloniex",
	"HitBTC",
	"ZB.com",
	"Bitrue",
	"Upbit",
	"FixedFloat",
	"Mexc",
}

func getSourceFundCallApi(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	addr := vars["address"]

	result, err := getSourceFund(addr)
	if err != nil {
		http.Error(w, "Error fetching data", http.StatusInternalServerError)
		return
	}

	// Send the result as a JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func getSourceFund(addr string) ([]SourceFund, error) {
	url := fmt.Sprintf("https://etherscan.io/address/%s", addr)

	// Make an HTTP GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Convert the body to a string
	bodyStr := string(body)

	// Use regular expressions to extract data between markers
	re := regexp.MustCompile(`data-bs-toggle="tooltip" data-bs-trigger="hover" data-bs-placement="top" title="(.*?)<\/a><a class="js-clipboar`)
	matches := re.FindAllStringSubmatch(bodyStr, -1)

	var result []map[string]string
	var finalData []SourceFund
	for _, match := range matches {
		data := strings.Split(match[1], "<br/>(")
		if len(data) >= 2 {
			name := data[0]
			address := strings.Split(data[1], ")")[0]
			result = append(result, map[string]string{
				"name":    name,
				"address": address,
			})
		}
	}

	// Filter name result using cex list
	for _, data := range result {
		for _, cexName := range cex {
			if strings.Contains(data["name"], cexName) {
				finalData = append(finalData, SourceFund{
					Name:    data["name"],
					Address: data["address"],
				})
			}
		}
	}

	// Remove duplicates from the result
	finalData = removeDuplicates(finalData)

	return finalData, nil
}

func removeDuplicates(elements []SourceFund) []SourceFund {
	encountered := map[string]bool{}
	result := []SourceFund{}

	for _, v := range elements {
		if !encountered[v.Address] {
			encountered[v.Address] = true
			result = append(result, v)
		}
	}

	return result
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/source-fund/address={address}", getSourceFundCallApi).Methods("GET")

	http.Handle("/", r)

	// Start the HTTP server on port 8080
	fmt.Println("Server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
