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

	corev1 "k8s.io/api/core/v1"
)

var SpringBoot = Opinions{
	{
		Id: "spring-web-port",
		Applicable: func(applied AppliedOpinions, imageMetadata map[string]string) bool {
			// TODO apply if the metadata indicates a webapp
			return true
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata map[string]string) error {
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
		Id: "spring-boot-actuator-port",
		Applicable: func(applied AppliedOpinions, imageMetadata map[string]string) bool {
			// TODO apply if the metadata indicates a spring-boot-actuator is installed
			return true
		},
		Apply: func(ctx context.Context, podSpec *corev1.PodTemplateSpec, imageMetadata map[string]string) error {
			applicationProperties := SpringApplicationProperties(ctx)
			// TODO check for an existing port before clobbering
			applicationProperties["management.server.port"] = "8081"
			return nil
		},
	},
	// TODO add a whole lot more opinions
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
