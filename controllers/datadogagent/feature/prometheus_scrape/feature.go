// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package prometheusscrape

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"

	"github.com/DataDog/datadog-operator/apis/datadoghq/v1alpha1"
	"github.com/DataDog/datadog-operator/apis/datadoghq/v2alpha1"
	apiutils "github.com/DataDog/datadog-operator/apis/utils"

	apicommon "github.com/DataDog/datadog-operator/apis/datadoghq/common"
	apicommonv1 "github.com/DataDog/datadog-operator/apis/datadoghq/common/v1"
	"github.com/DataDog/datadog-operator/controllers/datadogagent/feature"
)

func init() {
	err := feature.Register(feature.PrometheusScrapeIDType, buildPrometheusScrapeFeature)
	if err != nil {
		panic(err)
	}
}

func buildPrometheusScrapeFeature(options *feature.Options) feature.Feature {
	prometheusScrapeFeat := &prometheusScrapeFeature{}

	return prometheusScrapeFeat
}

type prometheusScrapeFeature struct {
	enable                 bool
	enableServiceEndpoints bool
	additionalConfigs      string
}

// Configure is used to configure the feature from a v2alpha1.DatadogAgent instance.
func (f *prometheusScrapeFeature) Configure(dda *v2alpha1.DatadogAgent) (reqComp feature.RequiredComponents) {
	prometheusScrape := dda.Spec.Features.PrometheusScrape

	if prometheusScrape != nil && apiutils.BoolValue(prometheusScrape.Enabled) {
		f.enable = true
		f.enableServiceEndpoints = apiutils.BoolValue(prometheusScrape.EnableServiceEndpoints)
		if prometheusScrape.AdditionalConfigs != nil {
			f.additionalConfigs = *prometheusScrape.AdditionalConfigs
		}
		reqComp = feature.RequiredComponents{
			Agent: feature.RequiredComponent{
				IsRequired: &f.enable,
				Containers: []apicommonv1.AgentContainerName{
					apicommonv1.CoreAgentContainerName,
				},
			},
			ClusterAgent: feature.RequiredComponent{
				IsRequired: &f.enable,
				Containers: []apicommonv1.AgentContainerName{
					apicommonv1.ClusterAgentContainerName,
				},
			},
		}
	}

	return reqComp
}

// ConfigureV1 use to configure the feature from a v1alpha1.DatadogAgent instance.
func (f *prometheusScrapeFeature) ConfigureV1(dda *v1alpha1.DatadogAgent) (reqComp feature.RequiredComponents) {
	prometheusScrape := dda.Spec.Features.PrometheusScrape

	if apiutils.BoolValue(prometheusScrape.Enabled) {
		f.enable = true
		f.enableServiceEndpoints = apiutils.BoolValue(prometheusScrape.ServiceEndpoints)
		if prometheusScrape.AdditionalConfigs != nil {
			f.additionalConfigs = *prometheusScrape.AdditionalConfigs
		}
		reqComp = feature.RequiredComponents{
			Agent: feature.RequiredComponent{
				IsRequired: &f.enable,
				Containers: []apicommonv1.AgentContainerName{
					apicommonv1.CoreAgentContainerName,
				},
			},
			ClusterAgent: feature.RequiredComponent{
				IsRequired: &f.enable,
				Containers: []apicommonv1.AgentContainerName{
					apicommonv1.ClusterAgentContainerName,
				},
			},
		}
	}

	return reqComp
}

// ManageDependencies allows a feature to manage its dependencies.
// Feature's dependencies should be added in the store.
func (f *prometheusScrapeFeature) ManageDependencies(managers feature.ResourceManagers) error {
	return nil
}

// ManageClusterAgent allows a feature to configure the ClusterAgent's corev1.PodTemplateSpec
// It should do nothing if the feature doesn't need to configure it.
func (f *prometheusScrapeFeature) ManageClusterAgent(managers feature.PodTemplateManagers) error {
	managers.EnvVar().AddEnvVarToContainer(apicommonv1.ClusterAgentContainerName, &corev1.EnvVar{
		Name:  apicommon.DDPrometheusScrapeEnabled,
		Value: "true",
	})
	managers.EnvVar().AddEnvVarToContainer(apicommonv1.ClusterAgentContainerName, &corev1.EnvVar{
		Name:  apicommon.DDPrometheusScrapeServiceEndpoints,
		Value: strconv.FormatBool(f.enableServiceEndpoints),
	})
	if f.additionalConfigs != "" {
		managers.EnvVar().AddEnvVarToContainer(apicommonv1.ClusterAgentContainerName, &corev1.EnvVar{
			Name:  apicommon.DDPrometheusScrapeChecks,
			Value: apiutils.YAMLToJSONString(f.additionalConfigs),
		})
	}

	return nil
}

// ManageNodeAgent allows a feature to configure the Node Agent's corev1.PodTemplateSpec
// It should do nothing if the feature doesn't need to configure it.
func (f *prometheusScrapeFeature) ManageNodeAgent(managers feature.PodTemplateManagers) error {
	managers.EnvVar().AddEnvVarToContainer(apicommonv1.CoreAgentContainerName, &corev1.EnvVar{
		Name:  apicommon.DDPrometheusScrapeEnabled,
		Value: "true",
	})
	managers.EnvVar().AddEnvVarToContainer(apicommonv1.CoreAgentContainerName, &corev1.EnvVar{
		Name:  apicommon.DDPrometheusScrapeServiceEndpoints,
		Value: strconv.FormatBool(f.enableServiceEndpoints),
	})
	if f.additionalConfigs != "" {
		managers.EnvVar().AddEnvVarToContainer(apicommonv1.CoreAgentContainerName, &corev1.EnvVar{
			Name:  apicommon.DDPrometheusScrapeChecks,
			Value: apiutils.YAMLToJSONString(f.additionalConfigs),
		})
	}

	return nil
}

// ManageClusterChecksRunner allows a feature to configure the ClusterChecksRunner's corev1.PodTemplateSpec
// It should do nothing if the feature doesn't need to configure it.
func (f *prometheusScrapeFeature) ManageClusterChecksRunner(managers feature.PodTemplateManagers) error {
	return nil
}
