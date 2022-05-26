package test_cases

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	appsv1Typed "k8s.io/client-go/kubernetes/typed/apps/v1"
	"k8s.io/client-go/util/retry"
	"log"
	"strings"
)

type DeploymentResource struct {
	Client   appsv1Typed.DeploymentInterface
	Resource *appsv1.Deployment
	Error    error
}

func (r *DeploymentResource) GetObject() runtime.Object {
	//fmt.Printf("%#v\n\n", r.Resource)
	return r.Resource
}

func (r *DeploymentResource) GetError() error {
	return r.Error
}

func (r *DeploymentResource) GetResourceName() string {
	return r.Resource.Name
}

func (r *DeploymentResource) GetResourceKind() string {
	kind := strings.Split(fmt.Sprintf("%T", r.Resource), ".")
	return kind[len(kind)-1 : len(kind)][0]
}

func (r *DeploymentResource) IsReady() bool {
	for _, v := range r.Resource.Status.Conditions {
		if v.Reason == "MinimumReplicasAvailable" && v.Status != "True" {
			return false
		}
	}
	if r.Resource.Status.UnavailableReplicas != 0 {
		return false
	}
	return true
}

func (r *DeploymentResource) GetClient(namespace string) {
	r.Client = clientset.AppsV1().Deployments(namespace)
}

func (r *DeploymentResource) Get() {
	resource, err := r.Client.Get(context.TODO(), r.Resource.Name, metav1.GetOptions{})
	if getHandler(r.Resource.Kind, r.Resource.Name, err) {
		r.Resource = resource
		return
	}
	r.Error = err
}
func (r *DeploymentResource) Create() {
	result, err := r.Client.Create(context.TODO(), r.Resource, metav1.CreateOptions{})
	r.Error = err
	r.Resource = result
}
func (r *DeploymentResource) Update() {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := r.Client.Update(context.TODO(), r.Resource, metav1.UpdateOptions{})
		return updateErr
	})
	if retryErr != nil {
		log.Printf("Update failed for %s:%s: %v\n", r.Resource.Kind, r.Resource.Name, retryErr)
	}
}
func (r *DeploymentResource) Delete() {
	deletePolicy := metav1.DeletePropagationForeground
	if err := r.Client.Delete(context.TODO(), "demo-deployment", metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Println(err)
	}
}

//TODO: Implement this method
//func rawDeployment(){
//	deploy := &appsv1.Deployment{
//		TypeMeta:metav1.TypeMeta{Kind:"", APIVersion:""},
//		ObjectMeta:metav1.ObjectMeta{
//			Name:        "nginx-e2e",
//			Namespace:   "cloudops-testing",
//			Labels:      map[string]string{"app": "nginx-e2e", "app.kubernetes.io/instance": "web-frontend", "app.kubernetes.io/managed-by": "kubectl", "app.kubernetes.io/name": "web-frontend"},
//			Annotations: map[string]string{"deployment.kubernetes.io/revision": "1"},
//		},
//		Spec: appsv1.DeploymentSpec{
//			Replicas:(*int32)(1),
//			Selector:(*metav1.LabelSelector)(0xc00056b120),
//			Template:v12.PodTemplateSpec{
//				ObjectMeta:metav1.ObjectMeta{
//					Labels:map[string]string{"app":"web-frontend"},
//				},
//				Spec:v12.PodSpec{
//					Volumes:[]v12.Volume{
//						v12.Volume{
//							Name:"index-html",
//							VolumeSource:v12.VolumeSource{
//								HostPath:(*v12.HostPathVolumeSource)(nil),
//								EmptyDir:(*v12.EmptyDirVolumeSource)(nil),
//								Secret:(*v12.SecretVolumeSource)(nil),
//								PersistentVolumeClaim:(*v12.PersistentVolumeClaimVolumeSource)(nil),
//								ConfigMap:(*v12.ConfigMapVolumeSource)(0xc0005e5f00),
//							},
//						},
//						v12.Volume{
//							Name:"conf",
//							VolumeSource:v12.VolumeSource{
//								HostPath:(*v12.HostPathVolumeSource)(nil),
//								EmptyDir:(*v12.EmptyDirVolumeSource)(nil),
//								Secret:(*v12.SecretVolumeSource)(nil),
//								PersistentVolumeClaim:(*v12.PersistentVolumeClaimVolumeSource)(nil),
//								ConfigMap:(*v12.ConfigMapVolumeSource)(0xc0005e5f40),
//							},
//						},
//					},
//					InitContainers:[]v12.Container(nil),
//					Containers:[]v12.Container{
//						v12.Container{
//							Name:"nginx",
//							Image:"nginx:1.21.5",
//							Command:[]string(nil),
//							Args:[]string(nil),
//							WorkingDir:"",
//							Ports:[]v12.ContainerPort{
//								v12.ContainerPort{
//									Name:"http",
//									HostPort:0,
//									ContainerPort:80,
//									Protocol:"TCP",
//									HostIP:""},
//							},
//							EnvFrom:[]v12.EnvFromSource(nil),
//							Env:[]v12.EnvVar(nil),
//							Resources:v12.ResourceRequirements{
//								Limits:v12.ResourceList{
//									"cpu":resource.Quantity{
//										i:resource.int64Amount{
//											value:1,
//											scale:0,
//										},
//										d:resource.infDecAmount{
//											Dec:(*inf.Dec)(nil),
//										},
//										s:"1",
//										Format:"DecimalSI",
//									},
//									"memory":resource.Quantity{
//										i:resource.int64Amount{
//											value:1073741824,
//											scale:0,
//										},
//										d:resource.infDecAmount{
//											Dec:(*inf.Dec)(nil),
//										},
//										s:"1Gi",
//										Format:"BinarySI",
//									},
//								},
//								Requests:v12.ResourceList{
//									"cpu":resource.Quantity{
//										i:resource.int64Amount{
//											value:500,
//											scale:-3,
//										},
//										d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"500m", Format:"DecimalSI"}, "memory":resource.Quantity{i:resource.int64Amount{value:524288000, scale:0}, d:resource.infDecAmount{Dec:(*inf.Dec)(nil)}, s:"500Mi", Format:"BinarySI"}}}, VolumeMounts:[]v12.VolumeMount{v12.VolumeMount{Name:"index-html", ReadOnly:false, MountPath:"/usr/share/nginx/html", SubPath:"", MountPropagation:(*v12.MountPropagationMode)(nil), SubPathExpr:""}, v12.VolumeMount{Name:"conf", ReadOnly:false, MountPath:"/etc/nginx/conf.d", SubPath:"", MountPropagation:(*v12.MountPropagationMode)(nil), SubPathExpr:""}}, VolumeDevices:[]v12.VolumeDevice(nil), LivenessProbe:(*v12.Probe)(nil), ReadinessProbe:(*v12.Probe)(nil), StartupProbe:(*v12.Probe)(nil), Lifecycle:(*v12.Lifecycle)(nil), TerminationMessagePath:"/dev/termination-log", TerminationMessagePolicy:"File", ImagePullPolicy:"IfNotPresent", SecurityContext:(*v12.SecurityContext)(nil), Stdin:false, StdinOnce:false, TTY:false}, v12.Container{Name:"nginx-prometheus", Image:"nginx/nginx-prometheus-exporter", Command:[]string(nil), Args:[]string(nil), WorkingDir:"", Ports:[]v12.ContainerPort{v12.ContainerPort{Name:"http-metrics", HostPort:0, ContainerPort:9113, Protocol:"TCP", HostIP:""}}, EnvFrom:[]v12.EnvFromSource(nil), Env:[]v12.EnvVar(nil), Resources:v12.ResourceRequirements{Limits:v12.ResourceList(nil), Requests:v12.ResourceList(nil)}, VolumeMounts:[]v12.VolumeMount(nil), VolumeDevices:[]v12.VolumeDevice(nil), LivenessProbe:(*v12.Probe)(nil), ReadinessProbe:(*v12.Probe)(nil), StartupProbe:(*v12.Probe)(nil), Lifecycle:(*v12.Lifecycle)(nil), TerminationMessagePath:"/dev/termination-log", TerminationMessagePolicy:"File", ImagePullPolicy:"IfNotPresent", SecurityContext:(*v12.SecurityContext)(nil), Stdin:false, StdinOnce:false, TTY:false}}, EphemeralContainers:[]v12.EphemeralContainer(nil), RestartPolicy:"Always", TerminationGracePeriodSeconds:(*int64)(0xc0005f1958), ActiveDeadlineSeconds:(*int64)(nil), DNSPolicy:"ClusterFirst", NodeSelector:map[string]string(nil), ServiceAccountName:"", DeprecatedServiceAccount:"", AutomountServiceAccountToken:(*bool)(nil), NodeName:"", HostNetwork:false, HostPID:false, HostIPC:false, ShareProcessNamespace:(*bool)(nil), SecurityContext:(*v12.PodSecurityContext)(0xc000219260), ImagePullSecrets:[]v12.LocalObjectReference(nil), Hostname:"", Subdomain:"", Affinity:(*v12.Affinity)(0xc0004bec30), SchedulerName:"default-scheduler", Tolerations:[]v12.Toleration(nil), HostAliases:[]v12.HostAlias(nil), PriorityClassName:"", Priority:(*int32)(nil), DNSConfig:(*v12.PodDNSConfig)(nil), ReadinessGates:[]v12.PodReadinessGate(nil), RuntimeClassName:(*string)(nil), EnableServiceLinks:(*bool)(nil), PreemptionPolicy:(*v12.PreemptionPolicy)(nil), Overhead:v12.ResourceList(nil), TopologySpreadConstraints:[]v12.TopologySpreadConstraint(nil), SetHostnameAsFQDN:(*bool)(nil), OS:(*v12.PodOS)(nil)}}, Strategy:appsv1.DeploymentStrategy{Type:"RollingUpdate", RollingUpdate:(*v12.RollingUpdateDeployment)(0xc000031aa0)}, MinReadySeconds:0, RevisionHistoryLimit:(*int32)(0xc0005f19b4), Paused:false, ProgressDeadlineSeconds:(*int32)(0xc0005f19bc)}, Status:appsv1.DeploymentStatus{ObservedGeneration:1, Replicas:3, UpdatedReplicas:3, ReadyReplicas:3, AvailableReplicas:3, UnavailableReplicas:0, Conditions:[]v12.DeploymentCondition{v12.DeploymentCondition{Type:"Available", Status:"True", LastUpdateTime:time.Date(2022, time.May, 6, 10, 20, 30, 0, time.Local), LastTransitionTime:time.Date(2022, time.May, 6, 10, 20, 30, 0, time.Local), Reason:"MinimumReplicasAvailable", Message:"Deployment has minimum availability."}, appsv1.DeploymentCondition{Type:"Progressing", Status:"True", LastUpdateTime:time.Date(2022, time.May, 6, 10, 20, 30, 0, time.Local), LastTransitionTime:time.Date(2022, time.May, 6, 10, 20, 27, 0, time.Local), Reason:"NewReplicaSetAvailable", Message:"ReplicaSet \"nginx-e2e-5fff85cfc6\" has successfully progressed."}}, CollisionCount:(*int32)(nil)}}
//
//}
