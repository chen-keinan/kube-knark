package main

import (
	"fmt"
	"github.com/chen-keinan/kube-knark/pkg/model"
	"net/http"
	"strings"
)

//OnK8sFileConfigChangeHook this plugin method accept k8s file config change event
//(event include file change command chmod / chown) and it send the event command args to an external events service
func OnK8sFileConfigChangeHook(k8sConfigFileChangeEvent model.K8sConfigFileChangeEvent) error {
	if k8sConfigFileChangeEvent.Spec.Severity == "CRITICAL" {
		commArgs := fmt.Sprintf("%s", k8sConfigFileChangeEvent.Msg.Args)
		req, err := http.NewRequest("POST", "http://localhost:8090/events", strings.NewReader(commArgs))
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
