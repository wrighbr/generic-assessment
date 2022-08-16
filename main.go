package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	file := flag.String("file", "./test.json", "json file to be read")
	url := flag.String("url", "https://my.testapi.com", "url to post json data")
	flag.Parse()

	jsonData, err := readJsonFile(*file)
	if err != nil {
		log.Println(err)
	}

	_, err = validateJson(jsonData)
	if err != nil {
		log.Fatalln(err, file)
	}
	fmt.Println(string(jsonData))

	respBody, err := sendPostRequest(*url, jsonData)
	if err != nil {
		log.Fatalln(err)
	}

	var respJson ResponseBody

	err = json.Unmarshal([]byte(respBody), &respJson)
	if err != nil {
		log.Fatalln(err)
	}

	for _, v := range respJson {
		if v.Valid {
			fmt.Println(v)
		}
	}
}

func readJsonFile(filepath string) (jsonData []byte, err error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func validateJson(jsonData []byte) (bool, error) {
	var j UserDetails

	var prefixErrMsg string

	err := json.Unmarshal(jsonData, &j)
	if err != nil {
		return false, err
	}

	if j.FirstName == nil {
		prefixErrMsg = "First Name"
	}

	if j.Surname == nil {
		prefixErrMsg = "Surname"
	}

	if j.ContactDetails.Address == nil {
		prefixErrMsg = "Address"
	}

	if j.ContactDetails.Mobile == nil {
		prefixErrMsg = "Mobile number"
	}

	if j.ContactDetails.Email == nil {
		prefixErrMsg = "email address"
	}
	if j.ContactDetails.City == nil {
		prefixErrMsg = "City"
	}
	if j.ContactDetails.Country == nil {
		prefixErrMsg = "Country"
	}
	if j.ContactDetails.PostCode == nil {
		prefixErrMsg = "Postcode"
	}

	if prefixErrMsg != "" {
		err := errors.New(prefixErrMsg + " is missing from Json file")
		return false, err
	}
	return true, nil
}

func sendPostRequest(url string, payload []byte) (respBody string, err error) {

	data := bytes.NewReader(payload)
	req, err := http.NewRequest("POST", url, data)
	if err != nil {
		log.Println(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	if res.StatusCode < 200 || res.StatusCode > 202 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("request failed, expected status: 2xx got: %d, error message %s", res.StatusCode, string(body))
	}

	return string(body), nil
}

type UserDetails struct {
	FirstName      *string `json:"first_name"`
	Surname        *string `json:"surname"`
	ContactDetails struct {
		Email    *string `json:"email"`
		Mobile   *int    `json:"mobile"`
		Address  *string `json:"address"`
		City     *string `json:"city"`
		Country  *string `json:"country"`
		PostCode *int    `json:"post_code"`
	} `json:"contact_details"`
}

type ResponseBody []struct {
	Valid bool `json:"valid"`
}
