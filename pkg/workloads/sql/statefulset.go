package sql

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/constants"
	"github.com/drew-viles/k8s-e2e-tester/pkg/workloads/coreworkloads"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// GeneratePostgresStatefulSet returns a StatefulSet that will be consumed by the postgres service
func GeneratePostgresStatefulSet(namespace, storageClassName string) *coreworkloads.StatefulSet {
	sts := &coreworkloads.StatefulSet{}

	sts.Generate(map[string]string{"namespace": namespace, "name": constants.PGSqlName, "label": constants.PGSqlName, "saName": constants.PGSqlSAName,
		"storageClassName": storageClassName, "storageSize": "5Gi", "affinityWith": constants.NginxName})

	sts.Resource.Spec.Template.Spec.Volumes = []apiv1.Volume{
		coreworkloads.GenerateVolumeFromConfigMap("init-data", constants.PGSqlConfName, 0755, map[string]string{
			"init-user-db.sh": "init-user-db.sh",
		}),
	}

	sts.Resource.Spec.Template.Spec.InitContainers = generatePostgresInitContainer()
	sts.Resource.Spec.Template.Spec.Containers = generatePostgresContainer()

	return sts
}

func generatePostgresInitContainer() []apiv1.Container {
	c := coreworkloads.GenerateContainer("postgres-clean", "busybox", "latest")
	c.Command = []string{"/bin/sh"}
	c.Args = []string{"-c", "rm -rf /var/lib/postgresql/data/*", "rm -rf /var/lib/postgresql/data/.*"}
	c.VolumeMounts = []apiv1.VolumeMount{
		{
			Name:      "data",
			MountPath: "/var/lib/postgresql/data",
		},
	}
	return []apiv1.Container{
		c,
	}
}

func generatePostgresContainer() []apiv1.Container {
	c := coreworkloads.GenerateContainer("postgres", "postgres", "15.1-alpine")
	c.Env = []apiv1.EnvVar{
		{
			Name: "POSTGRES_PASSWORD",
			ValueFrom: &apiv1.EnvVarSource{
				SecretKeyRef: &apiv1.SecretKeySelector{
					LocalObjectReference: apiv1.LocalObjectReference{Name: "pg-password"},
					Key:                  "passwd",
				},
			},
		},
		{
			Name:  "POSTGRES_USER",
			Value: constants.DBUser,
		},
		{
			Name:  "POSTGRES_DB",
			Value: constants.DBName,
		},
	}
	c.Ports = []apiv1.ContainerPort{
		{
			ContainerPort: 5432,
			Name:          "sql",
		},
	}
	c.ReadinessProbe = &apiv1.Probe{
		ProbeHandler: apiv1.ProbeHandler{
			TCPSocket: &apiv1.TCPSocketAction{
				Port: intstr.IntOrString{
					Type:   0,
					IntVal: 5432,
				},
			},
		},
		InitialDelaySeconds: 15,
		PeriodSeconds:       5,
	}
	c.VolumeMounts = []apiv1.VolumeMount{
		{
			Name:      "data",
			MountPath: "/var/lib/postgresql/data",
		},
		{
			Name:      "init-data",
			MountPath: "/docker-entrypoint-initdb.d",
		},
	}

	return []apiv1.Container{
		c,
	}
}
