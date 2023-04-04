package non_dependant

import (
	_ "embed"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
	"github.com/maistra/maistra-test-tool/pkg/util/env"
	"github.com/maistra/maistra-test-tool/pkg/util/log"
)

type SMCP struct {
	Name      string `default:"basic"`
	Namespace string `default:"istio-system"`
}

var (
	//go:embed yaml/subscription-jaeger.yaml
	jaegerSubscription string

	//go:embed yaml/subscription-kiali.yaml
	kialiSubscription string

	//go:embed yaml/subscription-ossm.yaml
	ossmSubscription string
)

var (
	smcpName      = env.Getenv("SMCPNAME", "basic")
	meshNamespace = env.Getenv("MESHNAMESPACE", "istio-system")
	smcp          = SMCP{smcpName, meshNamespace}
)

func createNamespaces() {
	log.Log.Info("creating namespaces")
	util.ShellSilent(`oc new-project bookinfo`)
	util.ShellSilent(`oc new-project foo`)
	util.ShellSilent(`oc new-project bar`)
	util.ShellSilent(`oc new-project legacy`)
	util.ShellSilent(`oc new-project mesh-external`)
	util.ShellSilent(`oc new-project %s`, meshNamespace)
}

// Install nightly build operators from quay.io. This is used in Jenkins daily build pipeline.
func installNightlyOperators() {
	util.KubeApplyContents("openshift-operators", jaegerSubscription)
	util.KubeApplyContents("openshift-operators", kialiSubscription)
	util.KubeApplyContents("openshift-operators", ossmSubscription)
	time.Sleep(time.Duration(60) * time.Second)
	util.CheckPodRunning("openshift-operators", "name=istio-operator")
	time.Sleep(time.Duration(30) * time.Second)
}

// Initialize a Namespace for the Mesh
func SetupNamespacesAndControlPlane() {
	log.Log.Info("Setting up namespaces and OSSM control plane")
	createNamespaces()
	if env.Getenv("NIGHTLY", "false") == "true" {
		installNightlyOperators()
	}
	//TODO: set more setup steps if needed
}
