package invdata

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	invRtDataUrl = "/solar_api/v1/GetInverterRealtimeData.cgi"
)

type cumInvRtDataReq struct {
	url        string
	params     [6]string
	bodyString string
}

type cumInvRtDataReqResp struct {
	Body struct {
		Data struct {
			DayEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 float64 `json:"1"`
				} `json:"Values"`
			} `json:"DAY_ENERGY"`
			Pac struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"PAC"`
			TotalEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"TOTAL_ENERGY"`
			YearEnergy struct {
				Unit   string `json:"Unit"`
				Values struct {
					Num1 int `json:"1"`
				} `json:"Values"`
			} `json:"YEAR_ENERGY"`
		} `json:"Data"`
	} `json:"Body"`
	Head struct {
		RequestArguments struct {
			DeviceClass string `json:"DeviceClass"`
			Scope       string `json:"Scope"`
		} `json:"RequestArguments"`
		Status struct {
			Code        int    `json:"Code"`
			Reason      string `json:"Reason"`
			UserMessage string `json:"UserMessage"`
		} `json:"Status"`
		Timestamp time.Time `json:"Timestamp"`
	} `json:"Head"`
}

func FetchCumInvRtData(addr string) int64 {
	invRTDataReq := cumInvRtDataReq{"", [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1"}, "",
	}
	resp := invRTDataReq.httpGet(addr)
	defer resp.Body.Close()

	body := invRTDataReq.parseJsonData()
	dailyPower := int64(body.Body.Data.DayEnergy.Values.Num1)
	//fmt.Printf("Pwr: %.0f\n", dailyPower)
	return dailyPower
}

func (d *cumInvRtDataReq) parseJsonData() cumInvRtDataReqResp {
	var body cumInvRtDataReqResp
	json.Unmarshal([]byte(d.bodyString), &body)
	return body
}

func (d *cumInvRtDataReq) httpGet(invAddr string) *http.Response {
	d.url = invAddr + invRtDataUrl
	d.params = [6]string{
		"Scope", "System", "DeviceId", "1", "DataCollection", "CumulationInverterData1",
	}
	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// Add request info to the GET request
	q := req.URL.Query()
	q.Add(d.params[0], d.params[1])
	q.Add(d.params[2], d.params[3])
	q.Add(d.params[4], d.params[5])

	req.URL.RawQuery = q.Encode()

	resp, err := http.Get(req.URL.String())
	if err != nil {
		panic(err)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	d.bodyString = string(bodyBytes)

	//fmt.Printf(d.bodyString)
	//fmt.Println("Fronius HTTP GET Response status:", resp.Status)
	return resp
}
