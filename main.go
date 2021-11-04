package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	TokenParams GlovoTokenParams `json:"tokenParams"`
}

type GlovoTokenParams struct {
	GrantType string `json:"grantType"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

type glovoToken struct {
	AccessToken  string      `json:"accessToken"`
	TokenType    string      `json:"tokenType"`
	ExpiresIn    int         `json:"expiresIn"`
	RefreshToken string      `json:"refreshToken"`
	Scope        interface{} `json:"scope"`
	TwoFactor    interface{} `json:"twoFactorToken"`
}

type glovoData struct {
	Orders     []order     `json:"orders"`
	Pagination interface{} `json:"pagination"`
}

type order struct {
	OrderId    int        `json:"orderId"`
	OrderUrn   string     `json:"orderUrn"`
	Image      glovoImage `json:"image"`
	C          content    `json:"content"`
	F          footer     `json:"footer"`
	Style      string     `json:"style"`
	LayoutType string     `json:"layoutType"`
}

type glovoImage struct {
	LightImageId string `json:"lightImageId"`
	DarkImageId  string `json:"darkImageId"`
}

type content struct {
	Title string      `json:"title"`
	Body  []glovoBody `json:"body"`
}

type glovoBody struct {
	T string `json:"type"`
	D string `json:"data"`
}

type footer struct {
	L left  `json:"left"`
	R right `json:"right"`
}

type left struct {
	T string `json:"type"`
	D string `json:"data"`
}
type right struct {
	T string `json:"type"`
}

func main() {
	var body []byte
	var e error
	var gt glovoToken

	var file *os.File
	file, e = os.Open("config.json")

	var config Config

	decoder := json.NewDecoder(file)
	e = decoder.Decode(&config)
	if e != nil {
		log.Fatal(e)
	}

	gt, e = getGlovoToken(config.TokenParams)
	if e != nil {
		log.Fatal(e)
	}

	body, e = getGlovoOrders(gt)
	if e != nil {
		log.Fatal(e)
	}

	r := bytes.NewReader(body)

	// Code for getting glovo orders from json file
	/*
		f, e := os.Open("d:/docs/glovo_orders.json")
		if e != nil {
			log.Fatal(e)
		}
	*/

	b, _ := ioutil.ReadAll(r)

	var gData glovoData
	e = json.Unmarshal(b, &gData)
	if e != nil {
		log.Fatal(e)
	}

	var sum float64
	for i := range gData.Orders {
		var f1 float64
		f1, e = orderDataStringToFloat64(gData.Orders[i].F.L.D)
		if e != nil {
			log.Fatal(e, "\n number of iterator = ", i)
		}
		sum += f1
	}

	fmt.Println("you fat fuck spent for glovo", sum, "KZT. Go to fitness damn bloody lazy prick")
}

func orderDataStringToFloat64(s string) (float64, error) {
	s = strings.TrimSpace(s[4:])
	if _, e := strconv.ParseFloat(s[:1], 64); e != nil {
		return 0.0, nil
	}

	s = strings.Trim(s, " ")

	s = strings.Replace(s, ",", ".", -1)
	s = strings.Replace(s, string(rune(160)), "", -1)
	s = strings.Replace(s, string(rune(65533)), "", -1)

	f, e := strconv.ParseFloat(s, 64)
	if e != nil {
		fmt.Println("failed number", s)
	}
	return f, e
}

func getGlovoToken(tokenParams GlovoTokenParams) (glovoToken, error) {

	url := "https://api.glovoapp.com/oauth/token"

	payload := strings.NewReader("{\n\t\"grantType\": \"" + tokenParams.GrantType + "\",\n\t\"username\": \"" + tokenParams.Username + "\",\n\t\"password\": \"" + tokenParams.Password + "\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var gT glovoToken

	err := json.Unmarshal(body, &gT)
	return gT, err
}

func getGlovoOrders(gt glovoToken) (body []byte, err error) {
	url := "https://api.glovoapp.com/v3/customer/orders-list?offset=0&limit=2000"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("authorization", gt.AccessToken)
	req.Header.Add("path", "/v3/customer/orders-list?offset=0&limit=2000")
	req.Header.Add("scheme", "https")
	req.Header.Add("glovo-api-version", "14")
	req.Header.Add("glovo-app-development-state", "Production")
	req.Header.Add("glovo-app-platform", "web")
	req.Header.Add("glovo-app-type", "customer")
	req.Header.Add("glovo-app-version", "7")
	req.Header.Add("glovo-device-id", "560313773")
	req.Header.Add("glovo-language-code", "en")
	req.Header.Add("origin", "https://glovoapp.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, err = ioutil.ReadAll(res.Body)

	return
}
