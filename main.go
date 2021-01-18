package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/inexio/go-monitoringplugin"
	"github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"os"
)

var opts struct{
	URL 		string `short:"u" long:"url" description:"The url for requesting the federation check" required:"true"`
	ServerName 	string `short:"n" long:"server-name" description:"The Server name that you specified in the synapse configuration" required:"true"`
}

type CheckFederation struct {
	FederationOK bool `json:"FederationOK"`
}

type ErrorAndCode struct {
	ExitCode int
	Error    error
}

func CheckFederationSuccess(url string, serverName string) []ErrorAndCode {
	var errSlice []ErrorAndCode
	if url == "" || serverName == ""{
		errSlice = append(errSlice, ErrorAndCode{2, errors.New("URL and server name is required to send GET request")})
		return errSlice
	}
	url = url + serverName
	client := resty.New()
	request := client.SetDebugBodyLimit(1000).R()
	response, err := request.Get(url)
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.Wrap(err, "error during http request")})
		return errSlice
	}
	var resp CheckFederation
	err = json.Unmarshal(response.Body(), &resp)
	if resp.FederationOK != true {
		errSlice = append(errSlice, ErrorAndCode{2, errors.New("The federation check is not successful, please check what´s wrong")})
		return errSlice
	}
	errSlice = append(errSlice, ErrorAndCode{0, errors.New("Federation Check succeeded")})
	return errSlice

}

func OutputMonitoring(errSlice []ErrorAndCode, defaultMessage string) {
	response := monitoringplugin.NewResponse(defaultMessage)
	for i := 0; i < len(errSlice); i++ {
		response.UpdateStatus(errSlice[i].ExitCode, errSlice[i].Error.Error())
	}
	response.OutputAndExit()
}

func main() {
	var errSlice []ErrorAndCode
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		fmt.Println("error parsing flags")
		os.Exit(3)
	}
	errSlice = CheckFederationSuccess(opts.URL, opts.ServerName)
	OutputMonitoring(errSlice, "successfully checked: " + opts.ServerName)
}
