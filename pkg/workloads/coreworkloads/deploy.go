package coreworkloads

import (
	"context"
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	"log"
	"strings"
)

type Deployment struct {
	Client   *kubernetes.Clientset
	Resource *v1.Deployment
}

// Generate a base Deployment definition.
func (d *Deployment) Generate(data map[string]string) {
	d.Resource = &v1.Deployment{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
		Spec: v1.DeploymentSpec{
			Replicas: helpers.IntPtr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": data["name"],
				},
			},
			Strategy: v1.DeploymentStrategy{
				RollingUpdate: &v1.RollingUpdateDeployment{
					MaxUnavailable: &intstr.IntOrString{
						Type:   1,
						StrVal: "25%",
					},
					MaxSurge: &intstr.IntOrString{
						Type:   1,
						StrVal: "25%",
					},
				},
				Type: "RollingUpdate",
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": data["name"],
					},
				},
				Spec: apiv1.PodSpec{
					Affinity: &apiv1.Affinity{
						PodAffinity: &apiv1.PodAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []apiv1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: apiv1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchLabels: map[string]string{
												"app": data["affinityWith"],
											},
										},
										TopologyKey: "topology.kubernetes.io/zone",
									},
								},
							},
						},
						PodAntiAffinity: &apiv1.PodAntiAffinity{
							PreferredDuringSchedulingIgnoredDuringExecution: []apiv1.WeightedPodAffinityTerm{
								{
									Weight: 100,
									PodAffinityTerm: apiv1.PodAffinityTerm{
										LabelSelector: &metav1.LabelSelector{
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "app",
													Operator: metav1.LabelSelectorOpIn,
													Values:   []string{data["name"]},
												},
											},
										},
										TopologyKey: "topology.kubernetes.io/zone",
									},
								},
							},
						},
					},
					DNSPolicy:          apiv1.DNSClusterFirst,
					RestartPolicy:      apiv1.RestartPolicyAlways,
					ServiceAccountName: data["saName"],
				},
			},
		},
	}
}

// Create creates a Deployment on the Kubernetes cluster.
func (d *Deployment) Create() error {
	log.Printf("creating Deployment:%s...\n", d.Resource.Name)
	r := d.Client.AppsV1().Deployments(d.Resource.Namespace)
	res, err := r.Create(context.Background(), d.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	d.Resource = res
	log.Printf("Deployment:%s created.\n", d.Resource.Name)
	return nil
}

// Validate validates a Deployment on the Kubernetes cluster.
func (d *Deployment) Validate() error {
	var err error
	r := d.Client.AppsV1().Deployments(d.Resource.Namespace)
	d.Resource, err = r.Get(context.Background(), d.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a Deployment in the Kubernetes cluster.
func (d *Deployment) Update() error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		var err error
		r := d.Client.AppsV1().Deployments(d.Resource.Namespace)
		d.Resource, err = r.Update(context.TODO(), d.Resource, metav1.UpdateOptions{})
		return err
	})
	if retryErr != nil {
		return retryErr
	}
	return nil
}

// Delete deletes a Deployment from the Kubernetes cluster.
func (d *Deployment) Delete() error {
	name := d.Resource.Name
	log.Printf("deleting Deployment:%s...\n", name)
	r := d.Client.AppsV1().Deployments(d.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("Deployment:%s deleted.\n", name)
	return nil
}

func (d *Deployment) GetResourceName() string {
	return d.Resource.Name
}

func (d *Deployment) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", d.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (d *Deployment) IsReady() bool {
	if err := d.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if d.Resource.Status.UnavailableReplicas != 0 {
		return false
	}
	return true
}
