package main

import (
	"encoding/json"
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

type DeepLRequest struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"` // default is auto
	TragetLang string `json:"target_lang"` // ZH
}

type DeelLResponse struct {
	Code int32  `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
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

func HcfyRequestToDeepLRequest(hcfyRequest HcfyRequest) DeepLRequest {
	var originTragetLan string
	if len(hcfyRequest.Destination) > 1 && hcfyRequest.Destination[0] != "中文(简体)" {
		originTragetLan = hcfyRequest.Destination[0]
	} else {
		originTragetLan = hcfyRequest.Destination[0]
	}
	return DeepLRequest{
		Text:       hcfyRequest.Text,
		SourceLang: "auto",
		TragetLang: hcfyToDeeplMap[originTragetLan],
	}
}

func DeeplResponseToHcfyResponse(deeplResponse DeelLResponse, request HcfyRequest) HcfyResponse {
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

func HcfyToDeeplHanlder(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var hcfyRequest HcfyRequest
		err := decoder.Decode(&hcfyRequest)
		log.Println("Receive request: ", len(hcfyRequest.Text))
		if err != nil {
			log.Println("decode failed {}", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		deeplRequest := HcfyRequestToDeepLRequest(hcfyRequest)
		deeplRequestJson, _ := json.Marshal(deeplRequest)
		// log.Println(string(deeplRequestJson))
		resp, err := http.Post(targetEndpoint, "application/json", strings.NewReader(string(deeplRequestJson)))
		if err != nil {
			log.Println("request deepl api failed {}", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()
		decoder = json.NewDecoder(resp.Body)
		var deeplResponse DeelLResponse
		err = decoder.Decode(&deeplResponse)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hcfyResponse := DeeplResponseToHcfyResponse(deeplResponse, hcfyRequest)
		hcfyResponseJson, _ := json.Marshal(hcfyResponse)
		log.Println("Return the transfomated text with len {}", len(deeplResponse.Data))
		// log.Println(string(hcfyResponseJson))
		w.Header().Set("Content-Type", "application/json")
		w.Write(hcfyResponseJson)
		return
	} else {
		http.Error(w, "Invalid request", http.StatusInternalServerError)
		return
	}
}

func main() {
	for k, v := range deeplToHcfyMap {
		hcfyToDeeplMap[v] = k
	}
	if os.Getenv("DEEPL_ENDPOINT") != "" {
		targetEndpoint = os.Getenv("DEEPL_ENDPOINT")
	}
	log.Printf("targetEndpoint: %s", targetEndpoint)
	log.Println("start server.....")

	http.HandleFunc("/translate", HcfyToDeeplHanlder)

	http.ListenAndServe("0.0.0.0:8080", nil)
}
