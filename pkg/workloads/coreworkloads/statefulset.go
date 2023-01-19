package coreworkloads

import (
	"context"
	"fmt"
	"github.com/drew-viles/k8s-e2e-tester/pkg/helpers"
	v1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type StatefulSet struct {
	Client   *kubernetes.Clientset
	Resource *v1.StatefulSet
}

// Generate the base StatefulSet.
func (s *StatefulSet) Generate(data map[string]string) {
	volumeMode := apiv1.PersistentVolumeFilesystem
	storageClassName := data["storageClassName"]
	s.Resource = &v1.StatefulSet{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["label"]),
		Spec: v1.StatefulSetSpec{
			Replicas: helpers.IntPtr(3),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": data["name"],
				},
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
											MatchExpressions: []metav1.LabelSelectorRequirement{
												{
													Key:      "app",
													Operator: metav1.LabelSelectorOpIn,
													Values:   []string{data["affinityWith"]},
												},
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
											MatchLabels: map[string]string{
												"app": data["name"],
											},
										},
										TopologyKey: "topology.kubernetes.io/zone",
									},
								},
							},
						},
					},
					ServiceAccountName: data["saName"],
				},
			},
			VolumeClaimTemplates: []apiv1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "data",
					},
					Spec: apiv1.PersistentVolumeClaimSpec{
						AccessModes:      []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
						StorageClassName: &storageClassName,
						VolumeMode:       &volumeMode,
						Resources: apiv1.ResourceRequirements{
							Requests: apiv1.ResourceList{
								"storage": resource.MustParse(data["storageSize"]),
							},
						},
					},
				},
			},
		},
	}
}

// Create creates a StatefulSet on the Kubernetes cluster.
func (s *StatefulSet) Create() error {
	log.Printf("creating StatefulSet:%s...\n", s.Resource.Name)
	r := s.Client.AppsV1().StatefulSets(s.Resource.Namespace)
	res, err := r.Create(context.Background(), s.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	s.Resource = res
	log.Printf("StatefulSet:%s created.\n", s.Resource.Name)
	return nil
}

// Validate validates a StatefulSet on the Kubernetes cluster.
func (s *StatefulSet) Validate() error {
	var err error
	r := s.Client.AppsV1().StatefulSets(s.Resource.Namespace)
	s.Resource, err = r.Get(context.Background(), s.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Update modifies a StatefulSet in the Kubernetes cluster.
func (s *StatefulSet) Update() error {
	return nil
}

// Delete deletes a StatefulSet from the Kubernetes cluster.
func (s *StatefulSet) Delete() error {
	name := s.Resource.Name
	log.Printf("deleting StatefulSet:%s...\n", name)
	r := s.Client.AppsV1().StatefulSets(s.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}
	log.Printf("StatefulSet:%s deleted.\n", name)
	return nil
}

func (s *StatefulSet) GetResourceName() string {
	return s.Resource.Name
}

func (s *StatefulSet) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", s.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (s *StatefulSet) IsReady() bool {
	if err := s.Validate(); err != nil {
		log.Println(err)
		return false
	}
	if s.Resource.Status.AvailableReplicas != s.Resource.Status.Replicas {
		return false
	}
	return true
}
