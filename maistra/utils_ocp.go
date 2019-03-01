// Copyright 2019 Red Hat, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package maistra


import (
	"fmt"
	"regexp"
	"strings"
	"time"
	
	"istio.io/istio/pkg/log"
	"istio.io/istio/tests/util"
)


// TBD
// OcLogin runs oc login command to log into the OCP CLI

func ocCommand(subCommand, namespace, yamlFileName string, kubeconfig string) string {
	if namespace == "" {
		return fmt.Sprintf("oc %s -f %s --kubeconfig=%s", subCommand, yamlFileName, kubeconfig)
	}
	return fmt.Sprintf("oc %s -n %s -f %s --kubeconfig=%s", subCommand, namespace, yamlFileName, kubeconfig)
}

// OcGrantPermission OCP cluster specific requirements for deploying an application with sidecar.
// This is a temporary permission config
func OcGrantPermission(namespace, kubeconfig string) {
	util.Shell("oc adm policy add-scc-to-user privileged -z default -n %s --kubeconfig=%s", namespace, kubeconfig)
	util.Shell("oc adm policy add-scc-to-user anyuid -z default -n %s --kubeconfig=%s", namespace, kubeconfig)
}

// OcApply oc apply from file
func OcApply(namespace, yamlFileName string, kubeconfig string) error {
	_, err := util.Shell(ocCommand("apply", namespace, yamlFileName, kubeconfig))
	return err
}

// OcDelete kubectl delete from file
func OcDelete(namespace, yamlFileName string, kubeconfig string) error {
	_, err := util.Shell(ocCommand("delete", namespace, yamlFileName, kubeconfig))
	return err
}

// GetOCPIngressgateway returns the OCP cluster ingressgateway host URL.
func GetOCPIngressgateway(podLabel, namespace, kubeconfig string) (string, error) {
	ingress, err := util.Shell("kubectl get routes -l %s -n %s -o jsonpath='{.items[0].spec.host}' --kubeconfig=%s",
							podLabel, namespace, kubeconfig)
	
	for i := 0; i < testRetryTimes; i++ {
		if err == nil {
			break
		}
		time.Sleep(time.Duration(5) * time.Second)
		ingress, err = util.Shell("kubectl get routes -l %s -n %s -o jsonpath='{.items[0].spec.host}' --kubeconfig=%s",
							podLabel, namespace, kubeconfig)
	}
	if err != nil {
		return "", err
	}
	return ingress, nil
}

// GetSecureIngressPort returns the https ingressgateway port
// "$(${OC_COMMAND} -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].port}')"
func GetSecureIngressPort(namespace, serviceName, kubeconfig string) (string, error) {
	port, err := util.Shell(
		"kubectl -n %s get service %s -o jsonpath='{.spec.ports[?(@.name==\"https\")].port}' --kubeconfig=%s",
		namespace, serviceName, kubeconfig)
	if err != nil {
		return "", err
	}
	port = strings.Trim(port, "'")
	rp := regexp.MustCompile(`^[0-9]{1,5}$`)
	if rp.FindString(port) == "" {
		err = fmt.Errorf("unable to find the port of %s", serviceName)
		log.Warna(err)
		return "", err
	}
	return port, nil
}

// GetTCPIngressPort returns the tcp ingressgateway port
// kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="tcp")].port}'
func GetTCPIngressPort(namespace, serviceName, kubeconfig string) (string, error) {
	port, err := util.Shell(
		"kubectl -n %s get service %s -o jsonpath='{.spec.ports[?(@.name==\"tcp\")].port}' --kubeconfig=%s",
		namespace, serviceName, kubeconfig)
	if err != nil {
		return "", err
	}
	port = strings.Trim(port, "'")
	rp := regexp.MustCompile(`^[0-9]{1,5}$`)
	if rp.FindString(port) == "" {
		err = fmt.Errorf("unable to find the port of %s", serviceName)
		log.Warna(err)
		return "", err
	}
	return port, nil
}

// GetIngressHostIP returns the OCP ingressgateway Host IP address from the OCP router endpoint
func GetIngressHostIP(kubeconfig string) (string, error) {
	ip, err := util.Shell("kubectl get endpoints -n default -l router -o jsonpath='{.items[0].subsets[0].addresses[0].ip}' --kubeconfig=%s", kubeconfig)
	if err != nil {
		return "", err
	}
	return ip, nil
}