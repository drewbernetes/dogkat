package web

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GenerateNginxDeploy returns a Deployment that will be consumed by the nginx service
func GenerateNginxDeploy(namespace string) *coreworkloads.Deployment {
	d := &coreworkloads.Deployment{}
	d.Generate(map[string]string{"namespace": namespace, "name": constants.NginxName, "labels": constants.NginxName, "saName": constants.NginxSAName, "affinityWith": constants.PGSqlName})
	d.Resource.Spec.Template.Spec.Volumes = []apiv1.Volume{
		coreworkloads.GenerateVolumeFromConfigMap("index-html", constants.NginxPagesName, 0644, map[string]string{
			"index":   "index.php",
			"healthz": "healthz.php",
			"common":  "common.php",
		}),
		coreworkloads.GenerateVolumeFromConfigMap("conf", constants.NginxConfName, 0644, map[string]string{
			"default": "default.conf",
			"metrics": "metrics.conf",
		}),
	}
	d.Resource.Spec.Template.Spec.Containers = generateNginxContainers()

	return d
}

func generateNginxContainers() []apiv1.Container {
	// Nginx container
	n := coreworkloads.GenerateContainer("nginx", "nginx", "1.23.2-alpine")
	n.Env = []apiv1.EnvVar{
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{Name: constants.PGSqlPasswdName},
					Key:                  "passwd",
				},
			},
		},
	}
	n.Resources = apiv1.ResourceRequirements{}
	n.Ports = []apiv1.ContainerPort{
		{
			ContainerPort: 80,
			Name:          "http",
			Protocol:      apiv1.ProtocolTCP,
		},
	}
	n.ReadinessProbe = &apiv1.Probe{
		ProbeHandler: apiv1.ProbeHandler{
			HTTPGet: &apiv1.HTTPGetAction{
				Path: "/healthz.php",
				Port: intstr.IntOrString{
					Type:   0,
					IntVal: 80,
				},
			},
		},
		InitialDelaySeconds: 30,
		PeriodSeconds:       5,
	}
	n.VolumeMounts = []apiv1.VolumeMount{
		{
			Name:      "index-html",
			MountPath: "/usr/share/nginx/html",
		},
		{
			Name:      "conf",
			MountPath: "/etc/nginx/conf.d",
		},
	}

	// PHP container
	p := coreworkloads.GenerateContainer("php", "drewviles/php-pdo", "8.0.18-fpm")
	p.Env = []apiv1.EnvVar{
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{Name: constants.PGSqlPasswdName},
					Key:                  "passwd",
				},
			},
		},
	}
	p.Ports = []apiv1.ContainerPort{
		{
			ContainerPort: 9000,
			Protocol:      apiv1.ProtocolTCP,
			Name:          "php",
		},
	}
	p.VolumeMounts = []apiv1.VolumeMount{
		{
			Name:      "index-html",
			MountPath: "/usr/share/nginx/html",
		},
	}

	// Nginx-Prometheus container
	e := coreworkloads.GenerateContainer("nginx-prometheus", "nginx/nginx-prometheus-exporter", "latest")
	e.Ports = []apiv1.ContainerPort{
		{
			ContainerPort: 9113,
			Name:          "http-metrics",
		},
	}

	return []apiv1.Container{n, p, e}
}
