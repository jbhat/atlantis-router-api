package client


import (
	"io"
	"io/ioutil"
	"strings"
	"net/http"

)


func BuildRequest(method, uri, data string) (*http.Request, error){
	var reader io.Reader 
	if data == "" {
		reader = nil
	} else {
		reader = strings.NewReader(data)
	}
	req, err := http.NewRequest(method, APIAddress + uri, reader) 
	if err != nil {
		return nil, err
	}
	//User and Secret are set in client.go
	req.Header.Add("User", User)
	req.Header.Add("Secret", Secret)
	if data != "" {
		req.Header.Add("Content-Type", "application/json")
	}	

	return req, nil
}

//returns statusCode, the response body, and an optional error
func BuildAndSendRequest(method, uri, data string) (int, string , error) {

	req, err := BuildRequest(method, uri, data)
	if err != nil {
		return 0, "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, "", err
	}

	statusCode := resp.StatusCode

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, "", err
	}

	return statusCode, string(body), nil

}
