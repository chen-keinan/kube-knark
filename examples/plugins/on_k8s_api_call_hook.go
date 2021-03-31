package main

import (
	"github.com/chen-keinan/kube-knark/pkg/model"
	"net/http"
	"strings"
)

//OnK8sAPICallHook this plugin method accept k8s api call event
//(event include http request/response, and match spec data) and it send the event URI to an external event service
func OnK8sAPICallHook(k8sAPICallEvent model.K8sAPICallEvent) error {
	if k8sAPICallEvent.Spec.Severity == "MAJOR" {
		req, err := http.NewRequest("POST", "http://localhost:8090/events", strings.NewReader(k8sAPICallEvent.Msg.HTTPRequestData.RequestURI))
		if err != nil {
			return err
		}
		client := http.Client{}
		_, err = client.Do(req)
		if err != nil {
			return err
		}
	}
	return nil
}
