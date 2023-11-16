package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var deeplToHcfyMap = map[string]string{
	"ZH": "中文(简体)",
	"DE": "德语",
	"EN": "英语",
	"ES": "西班牙语",
	"FR": "法语",
	"IT": "意大利语",
	"JA": "日语",
	"NL": "荷兰语",
	"PL": "波兰语",
	"PT": "葡萄牙语",
	"RU": "俄语",
	"BG": "保加利亚语",
	"CS": "捷克语",
	"DA": "丹麦语",
	"EL": "希腊语",
	"ET": "爱沙尼亚语",
	"FI": "芬兰语",
	"HU": "匈牙利语",
	"LT": "立陶宛语",
	"LV": "拉脱维亚语",
	"RO": "罗马尼亚语",
	"SK": "斯洛伐克语",
	"SL": "斯洛文尼亚语",
	"SV": "瑞典语",
}

var HcfyToDeeplMap map[string]string

var hcfyToDeeplMap = map[string]string{}

var targetEndpoint string
var tagetName string

type DeepLXRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"` // default is auto
	TragetLang string `json:"target_lang"` // ZH
}

type DeepLXResponse struct {
	Code         int32    `json:"code"`
	Msg          string   `json:"msg"`
	Data         string   `json:"data"`
	SourceLang   string   `json:"source_lang"` // default is auto
	TragetLang   string   `json:"target_lang"` // ZH
	Alternatives []string `json:"alternatives"`
}

type HcfyRequest struct {
	Name        string   `json:"name"`
	Text        string   `json:"text"`
	Destination []string `json:"destination"` //["中文(简体)", "英语"]
	Source      string   `json:"source"`      // undefined -> auto
}

type HcfyResponse struct {
	Text   string   `json:"text"`
	From   string   `json:"from"`
	To     string   `json:"to"`
	Result []string `json:"result"`
}

func HcfyRequestToDeepLRequest(hcfyRequest HcfyRequest) DeepLXRequest {
	var originTragetLan string
	if len(hcfyRequest.Destination) > 1 && hcfyRequest.Destination[0] != "中文(简体)" {
		originTragetLan = hcfyRequest.Destination[0]
	} else {
		originTragetLan = hcfyRequest.Destination[0]
	}
	return DeepLXRequest{
		Text:       hcfyRequest.Text,
		SourceLang: "auto",
		TragetLang: hcfyToDeeplMap[originTragetLan],
	}
}

func DeeplResponseToHcfyResponse(deeplResponse DeepLXResponse, request HcfyRequest) HcfyResponse {
	var originTragetLan string
	if len(request.Destination) > 1 && request.Destination[0] != "中文(简体)" {
		originTragetLan = request.Destination[1]
	} else {
		originTragetLan = request.Destination[0]
	}
	if deeplResponse.Code != 200 {
		return HcfyResponse{
			Text:   request.Text,
			From:   request.Source,
			To:     originTragetLan,
			Result: []string{deeplResponse.Msg},
		}
	}
	return HcfyResponse{
		Text:   request.Text,
		From:   request.Source,
		To:     originTragetLan,
		Result: strings.Split(deeplResponse.Data, "\n"),
	}
}

func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HelloWorldHandler from: ", r.RemoteAddr)
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is accepted", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("Hello World"))
}

func HcfyToDeeplHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("HcfyToDeeplHandler from: ", r.RemoteAddr)
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is accepted", http.StatusMethodNotAllowed)
		return
	}

	var hcfyRequest HcfyRequest
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	if err := json.Unmarshal(requestBody, &hcfyRequest); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}

	if hcfyRequest.Name != tagetName {
		log.Printf("Wrong name: %s\n", hcfyRequest.Name)
		http.Error(w, "Wrong name", http.StatusBadRequest)
		return
	}
	// log.Printf("HcfyToDeeplHandler: %v\n", hcfyRequest)

	deeplRequest := HcfyRequestToDeepLRequest(hcfyRequest)
	log.Printf("deeplRequest: %v\n", deeplRequest)
	deeplRequestJson, err := json.Marshal(deeplRequest)
	if err != nil {
		log.Printf("Error marshalling DeepL request: %v", err)
		http.Error(w, "Error processing request", http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Post(targetEndpoint, "application/json", bytes.NewBuffer(deeplRequestJson))
	if err != nil {
		log.Printf("Error requesting DeepL API: %v", err)
		http.Error(w, "Error calling translation service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var deeplResponse DeepLXResponse
	requestBody, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading DeepL response: %v", err)
		http.Error(w, "Error reading response", http.StatusInternalServerError)
		return
	}

	if err = json.Unmarshal(requestBody, &deeplResponse); err != nil {
		log.Printf("Error decoding DeepL response: %v", err)
		http.Error(w, "Error decoding response", http.StatusInternalServerError)
		return
	}

	log.Printf("deeplResponse: %v\n", deeplResponse)

	hcfyResponse := DeeplResponseToHcfyResponse(deeplResponse, hcfyRequest)
	// log.Printf("hcfyResponse: %v\n", hcfyResponse)
	hcfyResponseJson, err := json.Marshal(hcfyResponse)
	if err != nil {
		log.Printf("Error marshalling Hcfy response: %v", err)
		http.Error(w, "Error processing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(hcfyResponseJson); err != nil {
		log.Printf("Error writing response: %v", err)
		// It's too late to change the HTTP status code here since the header has been sent
	}
}

func main() {
	for k, v := range deeplToHcfyMap {
		hcfyToDeeplMap[v] = k
	}
	if os.Getenv("DEEPLX_ENDPOINT") != "" {
		targetEndpoint = os.Getenv("DEEPLX_ENDPOINT")
	} else {
		log.Fatalln("DEEPLX_ENDPOINT is not set")
		os.Exit(1)
	}

	if os.Getenv("DEEPLX_NAME") != "" {
		tagetName = os.Getenv("DEEPLX_NAME")
	} else {
		log.Fatalln("DEEPLX_NAME is not set")
		os.Exit(1)
	}
	log.Printf("targetEndpoint: %s", targetEndpoint)
	log.Println("start server.....")

	http.HandleFunc("/", HcfyToDeeplHandler)
	// http.HandleFunc("/", HelloWorldHandler)

	http.ListenAndServe("0.0.0.0:9911", nil)
}
