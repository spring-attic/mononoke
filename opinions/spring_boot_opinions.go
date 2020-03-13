/*
Copyright 2020 the original author or authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package opinions

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/spring-cloud-incubator/mononoke/cnb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var SpringBoot = Opinions{
	{
		Id: "spring-boot-graceful-shutdown",
		Applicable: func(applied AppliedOpinions, imageMetadata cnb.BuildMetadata) bool {
			bootMetadata := NewSpringBootBOMMetadata(imageMetadata)
			// TODO(scothis) only apply to Boot 2.3+
			return bootMetadata.HasDependency("spring-boot-starter-tomcat") ||
				bootMetadata.HasDependency("spring-boot-starter-jetty") ||
				bootMetadata.HasDependency("spring-boot-starter-reactor-netty") ||
				bootMetadata.HasDependency("spring-boot-starter-undertow")
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata cnb.BuildMetadata) error {
			applicationProperties := SpringApplicationProperties(ctx)
			if _, ok := applicationProperties["server.shutdown.grace-period"]; ok {
				// boot grace period already defined, skipping
				return nil
			}
			var k8sGracePeriodSeconds int64 = 30 // default k8s grace period is 30 seconds
			if podSpec.Spec.TerminationGracePeriodSeconds != nil {
				k8sGracePeriodSeconds = *podSpec.Spec.TerminationGracePeriodSeconds
			}
			podSpec.Spec.TerminationGracePeriodSeconds = &k8sGracePeriodSeconds
			// allocate 80% of the k8s grace period to boot
			bootGracePeriodSeconds := int(math.Floor(0.8 * float64(k8sGracePeriodSeconds)))
			applicationProperties["server.shutdown.grace-period"] = fmt.Sprintf("%ds", bootGracePeriodSeconds)
			return nil
		},
	},
	{
		Id: "spring-web-port",
		Applicable: func(applied AppliedOpinions, imageMetadata cnb.BuildMetadata) bool {
			bootMetadata := NewSpringBootBOMMetadata(imageMetadata)
			return bootMetadata.HasDependency("spring-web")
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata cnb.BuildMetadata) error {
			applicationProperties := SpringApplicationProperties(ctx)
			// TODO be smarter about resolving the correct container
			c := &podSpec.Spec.Containers[0]
			// TODO check for an existing port before clobbering
			c.Ports = append(c.Ports, corev1.ContainerPort{
				ContainerPort: 8080,
				Protocol:      corev1.ProtocolTCP,
			})
			applicationProperties["server.port"] = "8080"
			return nil
		},
	},
	{
		Id: "spring-boot-actuator-probes",
		Applicable: func(applied AppliedOpinions, imageMetadata cnb.BuildMetadata) bool {
			bootMetadata := NewSpringBootBOMMetadata(imageMetadata)
			return bootMetadata.HasDependency("spring-boot-actuator")
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata cnb.BuildMetadata) error {
			applicationProperties := SpringApplicationProperties(ctx)

			// TODO check for an existing value before clobbering
			managementPort := 9001

			applicationProperties["management.server.port"] = strconv.Itoa(managementPort)
			applicationProperties["management.server.ssl.enabled"] = "false"
			applicationProperties["management.endpoint.health.enabled"] = "true"
			applicationProperties["management.endpoint.info.enabled"] = "true"

			// TODO check for an existing value before clobbering
			managementBasePath := "/actuator"
			applicationProperties["management.endpoints.web.base-path"] = managementBasePath

			// TODO be smarter about resolving the correct container
			c := &podSpec.Spec.Containers[0]

			// define probes
			if c.StartupProbe == nil {
				// requires k8s 1.16+
				// TODO(scothis) add if k8s can handle it
			}
			if c.LivenessProbe == nil {
				c.LivenessProbe = &corev1.Probe{
					InitialDelaySeconds: 30,
					PeriodSeconds:       5,
					TimeoutSeconds:      5,
				}
			}
			if c.LivenessProbe.Handler == (corev1.Handler{}) {
				c.LivenessProbe.Handler = corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: managementBasePath + "/info",
						Port: intstr.FromInt(managementPort),
					},
				}
			}
			if c.ReadinessProbe == nil {
				c.ReadinessProbe = &corev1.Probe{
					InitialDelaySeconds: 5,
					PeriodSeconds:       1,
					TimeoutSeconds:      5,
				}
			}
			if c.ReadinessProbe.Handler == (corev1.Handler{}) {
				c.ReadinessProbe.Handler = corev1.Handler{
					HTTPGet: &corev1.HTTPGetAction{
						Path: managementBasePath + "/health",
						Port: intstr.FromInt(managementPort),
					},
				}
			}

			return nil
		},
	},
	{
		// fallback if spring-boot-actuator-probes is not applied
		Id: "spring-web-probes",
		Applicable: func(applied AppliedOpinions, imageMetadata cnb.BuildMetadata) bool {
			return !applied.Has("spring-boot-actuator-probes") && applied.Has("spring-web-port")
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata cnb.BuildMetadata) error {
			applicationProperties := SpringApplicationProperties(ctx)

			if _, ok := applicationProperties["server.port"]; !ok {
				// no port, so we can't provide probes
				return nil
			}

			port, err := strconv.Atoi(applicationProperties["server.port"])
			if err != nil {
				return err
			}

			// TODO be smarter about resolving the correct container
			c := &podSpec.Spec.Containers[0]

			// define probes
			if c.StartupProbe == nil {
				// requires k8s 1.16+
				// TODO(scothis) add if k8s can handle it
			}
			if c.LivenessProbe == nil {
				c.LivenessProbe = &corev1.Probe{
					InitialDelaySeconds: 30,
					PeriodSeconds:       5,
					TimeoutSeconds:      5,
				}
			}
			if c.LivenessProbe.Handler == (corev1.Handler{}) {
				c.LivenessProbe.Handler = corev1.Handler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.FromInt(port),
					},
				}
			}
			if c.ReadinessProbe == nil {
				c.ReadinessProbe = &corev1.Probe{
					InitialDelaySeconds: 5,
					PeriodSeconds:       1,
					TimeoutSeconds:      5,
				}
			}
			if c.ReadinessProbe.Handler == (corev1.Handler{}) {
				c.ReadinessProbe.Handler = corev1.Handler{
					TCPSocket: &corev1.TCPSocketAction{
						Port: intstr.FromInt(port),
					},
				}
			}

			return nil
		},
	},
	// TODO add a whole lot more opinions
}

func NewSpringBootBOMMetadata(imageMetadata cnb.BuildMetadata) SpringBootBOMMetadata {
	// TODO(scothis) find a better way to convert map[string]interface{} to SpringBootBOMMetadata{}
	bom := imageMetadata.FindBOM("spring-boot")
	bootMetadata := SpringBootBOMMetadata{}
	bytes, err := json.Marshal(bom.Metadata)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(bytes, &bootMetadata)
	return bootMetadata
}

type SpringBootBOMMetadata struct {
	Classes      string                            `json:"classes"`
	ClassPath    []string                          `json:"classpath"`
	Dependencies []SpringBootBOMMetadataDependency `json:"dependencies"`
}

func (m *SpringBootBOMMetadata) HasDependency(name string) bool {
	for _, dep := range m.Dependencies {
		if dep.Name == name {
			return true
		}
	}
	return false
}

type SpringBootBOMMetadataDependency struct {
	Name    string `json:"name"`
	Sha256  string `json:"sha256"`
	Version string `json:"version"`
}

type springApplicationPropertiesKey struct{}

func StashSpringApplicationProperties(ctx context.Context, props map[string]string) context.Context {
	return context.WithValue(ctx, springApplicationPropertiesKey{}, props)
}

func SpringApplicationProperties(ctx context.Context) map[string]string {
	value := ctx.Value(springApplicationPropertiesKey{})
	if props, ok := value.(map[string]string); ok {
		return props
	}
	return nil
}
