package main

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-resty/resty"
	"github.com/pkg/errors"
	"github.com/inexio/go-monitoringplugin"
)

var Client = resty.New()

type CheckFederation struct {
	FederationOK bool `json:"FederationOK"`
}

type ErrorAndCode struct {
	ExitCode  int
	Error  error
}

func CheckFederationSuccess () []ErrorAndCode{
	var errSlice []ErrorAndCode
	request := Client.SetDebugBodyLimit(1000).R()
	response, err := request.Get("https://matrix.org/federationtester/api/report?server_name=thola.io")
	if err != nil {
		errSlice = append(errSlice, ErrorAndCode{3, errors.Wrap(err, "error during http request")})
		return errSlice
	}
	var resp CheckFederation
	err = json.Unmarshal(response.Body(), &resp)
	if resp.FederationOK != true {
		errSlice = append(errSlice, ErrorAndCode{2, errors.New("The federation check is not successfull, please check what´s wrong")})
		return errSlice
	} else {
		errSlice = append(errSlice, ErrorAndCode{1, errors.New("Federation Check succeeded")})
		return errSlice
	}

}

func OutputMonitoring(errSlice []ErrorAndCode, defaultMessage string, performanceDataSlice []monitoringplugin.PerformanceDataPoint) {
	response := monitoringplugin.NewResponse(defaultMessage)
	for i := 0; i < len(errSlice); i++ {
		response.UpdateStatus(errSlice[i].ExitCode, errSlice[i].Error.Error())
	}
	for i := 0; i < len(performanceDataSlice); i++ {
		err := response.AddPerformanceDataPoint(&performanceDataSlice[i])
		if err != nil {
			spew.Dump(err)
		}
	}
	response.OutputAndExit()
}


func main () {
	var errSlice []ErrorAndCode
	errSlice = CheckFederationSuccess()
	fmt.Println(errSlice)
}
