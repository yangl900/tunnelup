package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

var (
	resourceURI = "https://management.azure.com/subscriptions/1489e197-43c4-4a6d-9c72-7cdbf920cec6/resourceGroups/MC_aks2_yanglaks2_westus/providers/Microsoft.ContainerInstance/containerGroups/jumpbox-aci/containers/jumpbox-aci/exec?api-version=2018-04-01"
	reqBody     = "{\"command\": \"/ncssh\", \"terminalSize\":{\"rows\":80,\"cols\":180}}}"
)

// ExecResponse is the EXEC response
type ExecResponse struct {
	WebsocketURI string `json:"webSocketUri"`
	Passowrd     string `json:"password"`
}

func getSocketURI() (*ExecResponse, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", resourceURI, bytes.NewReader([]byte(reqBody)))

	token, err := acquireAuthTokenCurrentTenant()
	if err != nil {
		return nil, errors.New("Failed to acquire auth token: " + err.Error())
	}

	req.Header.Set("Authorization", token)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, errors.New("Request failed: " + err.Error())
	}

	defer response.Body.Close()
	buf, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, errors.New("Request failed: " + err.Error())
	}

	resp := ExecResponse{}
	json.Unmarshal(buf, &resp)

	return &resp, nil
}
