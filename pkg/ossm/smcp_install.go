// Copyright 2021 Red Hat, Inc.
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

package ossm

import (
	"strings"
	"testing"
	"time"

	"github.com/maistra/maistra-test-tool/pkg/util"
	"github.com/maistra/maistra-test-tool/pkg/util/env"
	"github.com/maistra/maistra-test-tool/pkg/util/log"
)

func installDefaultSMCP23() {
	log.Log.Info("Create SMCP v2.3 in ", meshNamespace)
	util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
	util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV23_template, smcp))
	util.KubeApplyContents(meshNamespace, smmr)

	// patch SMCP identity if it's on a ROSA cluster
	if env.Getenv("ROSA", "false") == "true" {
		util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
	}
	log.Log.Info("Waiting for mesh installation to complete")
	util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)
	time.Sleep(time.Duration(20) * time.Second)
}

func installDefaultSMCP22() {
	log.Log.Info("Create SMCP v2.2 in ", meshNamespace)
	util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
	util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV22_template, smcp))
	util.KubeApplyContents(meshNamespace, smmr)

	// patch SMCP identity if it's on a ROSA cluster
	if env.Getenv("ROSA", "false") == "true" {
		util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
	}
	log.Log.Info("Waiting for mesh installation to complete")
	util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)
	time.Sleep(time.Duration(20) * time.Second)
}

func TestSMCPInstall(t *testing.T) {
	defer installDefaultSMCP23()

	t.Run("smcp_test_install_2.3", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Create SMCP v2.3 in ", meshNamespace)
		util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV23_template, smcp))
		util.KubeApplyContents(meshNamespace, smmr)

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		msg, _ := util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		if !strings.Contains(msg, "ComponentsReady") {
			log.Log.Error("SMCP not Ready")
			t.Error("SMCP not Ready")
		}
		util.Shell(`oc get -n %s pods`, meshNamespace)
	})

	t.Run("smcp_test_uninstall_2.3", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Delete SMCP v2.3 in ", meshNamespace)
		util.KubeDeleteContents(meshNamespace, smmr)
		util.KubeDeleteContents(meshNamespace, util.RunTemplate(smcpV23_template, smcp))
		time.Sleep(time.Duration(40) * time.Second)
	})

	t.Run("smcp_test_install_2.2", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Create SMCP v2.2 in ", meshNamespace)
		util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV22_template, smcp))
		util.KubeApplyContents(meshNamespace, smmr)

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		msg, _ := util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		if !strings.Contains(msg, "ComponentsReady") {
			log.Log.Error("SMCP not Ready")
			t.Error("SMCP not Ready")
		}
		util.Shell(`oc get -n %s pods`, meshNamespace)
	})

	t.Run("smcp_test_uninstall_2.2", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Delete SMCP v2.2 in ", meshNamespace)
		util.KubeDeleteContents(meshNamespace, smmr)
		util.KubeDeleteContents(meshNamespace, util.RunTemplate(smcpV22_template, smcp))
		time.Sleep(time.Duration(40) * time.Second)
	})

	t.Run("smcp_test_install_2.1", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Create SMCP v2.1 in namespace ", meshNamespace)
		util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV21_template, smcp))
		util.KubeApplyContents(meshNamespace, smmr)

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		msg, _ := util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		if !strings.Contains(msg, "ComponentsReady") {
			log.Log.Error("SMCP not Ready")
			t.Error("SMCP not Ready")
		}
		util.Shell(`oc get -n %s pods`, meshNamespace)
	})

	t.Run("smcp_test_uninstall_2.1", func(t *testing.T) {
		defer util.RecoverPanic(t)
		log.Log.Info("Delete SMCP v2.1 in ", meshNamespace)
		util.KubeDeleteContents(meshNamespace, smmr)
		util.KubeDeleteContents(meshNamespace, util.RunTemplate(smcpV21_template, smcp))
		time.Sleep(time.Duration(60) * time.Second)
	})

	t.Run("smcp_test_upgrade_2.1_to_2.2", func(t *testing.T) {
		defer util.RecoverPanic(t)

		log.Log.Info("Create SMCP v2.1 in namespace ", meshNamespace)
		util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV21_template, smcp))
		util.KubeApplyContents(meshNamespace, smmr)

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		util.Shell(`oc get -n %s pods`, meshNamespace)

		log.Log.Info("Upgrade SMCP to v2.2")
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV22_template, smcp))

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		time.Sleep(time.Duration(10) * time.Second)
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 360s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		msg, _ := util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		if !strings.Contains(msg, "ComponentsReady") {
			log.Log.Error("SMCP not Ready")
			t.Error("SMCP not Ready")
		}
		util.Shell(`oc get -n %s pods`, meshNamespace)
	})

	t.Run("smcp_test_upgrade_2.2_to_2.3", func(t *testing.T) {
		defer util.RecoverPanic(t)

		log.Log.Info("Create SMCP v2.2 in namespace ", meshNamespace)
		util.ShellMuteOutputError(`oc new-project %s`, meshNamespace)
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV22_template, smcp))
		util.KubeApplyContents(meshNamespace, smmr)

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 300s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		util.Shell(`oc get -n %s pods`, meshNamespace)

		log.Log.Info("Upgrade SMCP to v2.2")
		util.KubeApplyContents(meshNamespace, util.RunTemplate(smcpV23_template, smcp))

		// patch SMCP identity if it's on a ROSA cluster
		if env.Getenv("ROSA", "false") == "true" {
			util.Shell(`oc patch -n %s smcp/%s --type merge -p '{"spec":{"security":{"identity":{"type":"ThirdParty"}}}}'`, meshNamespace, smcpName)
		}
		log.Log.Info("Waiting for mesh installation to complete")
		time.Sleep(time.Duration(10) * time.Second)
		util.Shell(`oc wait --for condition=Ready -n %s smcp/%s --timeout 360s`, meshNamespace, smcpName)

		log.Log.Info("Verify SMCP status and pods")
		msg, _ := util.Shell(`oc get -n %s smcp/%s -o wide`, meshNamespace, smcpName)
		if !strings.Contains(msg, "ComponentsReady") {
			log.Log.Error("SMCP not Ready")
			t.Error("SMCP not Ready")
		}
		util.Shell(`oc get -n %s pods`, meshNamespace)
	})
}
