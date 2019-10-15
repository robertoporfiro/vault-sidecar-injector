// Copyright © 2019 Talend
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

package config

import (
	"crypto/sha256"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"k8s.io/klog"
)

// Load : Load Vault Sidecar Injector's config
func Load(whSvrParams WhSvrParameters) (*InjectionConfig, error) {
	klog.Infof("annotationKeyPrefix=%s", whSvrParams.AnnotationKeyPrefix)
	klog.Infof("appLabelKey=%s", whSvrParams.AppLabelKey)
	klog.Infof("appServiceLabelKey=%s", whSvrParams.AppServiceLabelKey)

	// Load sidecar config
	var sidecarConfig SidecarConfig
	err := loadYaml(whSvrParams.SidecarCfgFile, &sidecarConfig)
	if err != nil {
		klog.Errorf("Failed to load sidecar configuration: %v", err)
		return nil, err
	}

	// Load Consul Template's block
	ctTemplateBlock, err := loadString(whSvrParams.ConsulTemplateTmplBlockFile)
	if err != nil {
		klog.Errorf("Failed to load Consul Template's template configuration: %v", err)
		return nil, err
	}

	// Load Consul Template's default template
	ctTemplateDefaultTmpl, err := loadString(whSvrParams.ConsulTemplateTmplDefaultFile)
	if err != nil {
		klog.Errorf("Failed to load Consul Template's default template: %v", err)
		return nil, err
	}

	// Load lifecycle hooks to inject into requesting pods
	var hooks LifecycleHooks
	err = loadYaml(whSvrParams.PodLifecycleHooksFile, &hooks)
	if err != nil {
		klog.Errorf("Failed to load pod's lifecycle hooks configuration: %v", err)
		return nil, err
	}

	return &InjectionConfig{
		VaultInjectorAnnotationKeyPrefix: whSvrParams.AnnotationKeyPrefix,
		ApplicationLabelKey:              whSvrParams.AppLabelKey,
		ApplicationServiceLabelKey:       whSvrParams.AppServiceLabelKey,
		SidecarConfig:                    &sidecarConfig,
		CtTemplateBlock:                  ctTemplateBlock,
		CtTemplateDefaultTmpl:            ctTemplateDefaultTmpl,
		PodslifecycleHooks:               &hooks,
	}, nil
}

func loadString(fileName string) (string, error) {
	data, err := loadRaw(fileName)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func loadYaml(fileName string, obj interface{}) error {
	data, err := loadRaw(fileName)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, obj)
}

func loadRaw(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	klog.Infof("Loading %s [sha256sum: %x]", fileName, sha256.Sum256(data))

	return data, nil
}