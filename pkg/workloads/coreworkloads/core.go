package coreworkloads

import (
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GenerateContainer returns a base container definition.
func GenerateContainer(name, image, tag string) apiv1.Container {
	return apiv1.Container{
		Name:                     name,
		Image:                    fmt.Sprintf("%s:%s", image, tag),
		TerminationMessagePath:   "/dev/termination-log",
		TerminationMessagePolicy: apiv1.TerminationMessageReadFile,
		ImagePullPolicy:          apiv1.PullIfNotPresent,
	}
}

// GenerateMetadata creates basic metadata for any resources that require it.
func GenerateMetadata(namespace, name, label string) metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Namespace: namespace,
		Name:      name,
		Labels: map[string]string{
			"app":                        label,
			"app.kubernetes.io/instance": label,
			"app.kubernetes.io/name":     label,
		},
	}
}

// GenerateVolumeFromConfigMap generates a Volume from a ConfigMap to be used in []apiv1.Volume.
func GenerateVolumeFromConfigMap(volumeName, configMapName string, mode int32, items map[string]string) apiv1.Volume {
	var i []apiv1.KeyToPath

	for k, v := range items {
		item := apiv1.KeyToPath{
			Key:  k,
			Path: v,
		}
		i = append(i, item)
	}

	return apiv1.Volume{
		Name: volumeName,
		VolumeSource: apiv1.VolumeSource{
			ConfigMap: &apiv1.ConfigMapVolumeSource{
				LocalObjectReference: apiv1.LocalObjectReference{Name: configMapName},
				Items:                i,
				DefaultMode:          helpers.IntPtr(mode),
			},
		},
	}
}
