package coreworkloads

import (
	"context"
	"fmt"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"log"
	"strings"
)

type PodDisruptionBudget struct {
	Client   *kubernetes.Clientset
	Resource *policyv1.PodDisruptionBudget
}

// Generate the base PodDisruptionBudget.
func (p *PodDisruptionBudget) Generate(data map[string]string) {
	p.Resource = &policyv1.PodDisruptionBudget{
		ObjectMeta: GenerateMetadata(data["namespace"], data["name"], data["name"]),
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: &intstr.IntOrString{
				Type:   0,
				IntVal: 2,
			},
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": data["label"]},
			},
		},
	}
}

// Create creates a PodDisruptionBudget on the Kubernetes cluster.
func (p *PodDisruptionBudget) Create() error {
	log.Printf("creating PodDisruptionBudget:%s...\n", p.Resource.Name)
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	_, err := r.Create(context.Background(), p.Resource, metav1.CreateOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("PodDisruptionBudget:%s created.\n", p.Resource.Name)
	return nil
}

// Validate validates a PodDisruptionBudget on the Kubernetes cluster.
func (p *PodDisruptionBudget) Validate() error {
	var err error
	log.Printf("confirming PodDisruptionBudget:%s...\n", p.Resource.Name)
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	p.Resource, err = r.Get(context.Background(), p.Resource.Name, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
		return err
	}

	log.Printf("PodDisruptionBudget: %s exists\n", p.Resource.Name)
	return nil
}

// Delete deletes a PodDisruptionBudget from the Kubernetes cluster.
func (p *PodDisruptionBudget) Delete() error {
	name := p.Resource.Name
	log.Printf("deleting PodDisruptionBudget:%s...\n", name)
	r := p.Client.PolicyV1().PodDisruptionBudgets(p.Resource.Namespace)
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Delete(context.Background(), name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
		return err
	}

	log.Printf("PodDisruptionBudget:%s deleted.\n", name)
	return nil
}

func (p *PodDisruptionBudget) GetResourceName() string {
	return p.Resource.Name
}

func (p *PodDisruptionBudget) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", p.Resource), ".")
	return kind[len(kind)-1:][0]
}

func (p *PodDisruptionBudget) IsReady() bool {
	if p.Resource.Status.CurrentHealthy < p.Resource.Status.DesiredHealthy || p.Resource.Status.DisruptionsAllowed == 0 {
		return false
	}
	return true
}
