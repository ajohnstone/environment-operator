package cluster

import (
	"testing"

	"github.com/pearsontechnology/environment-operator/pkg/bitesize"
	"github.com/pearsontechnology/environment-operator/pkg/diff"
	"github.com/pearsontechnology/environment-operator/pkg/util"

	log "github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/pkg/api/resource"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/extensions/v1beta1"
)

func TestKubernetesClusterClient(t *testing.T) {
	// t.Run("service count", testServiceCount)
	// t.Run("volumes", testVolumes)
	t.Run("full bitesize construct", testFullBitesizeEnvironment)
	t.Run("test service ports", testServicePorts)
	// t.Run("a/b deployment service", testABSingleService)
}

func TestApplyEnvironment(t *testing.T) {

	log.SetLevel(log.FatalLevel)
	client := fake.NewSimpleClientset(
		&v1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name: "environment-dev",
				Labels: map[string]string{
					"environment": "environment2",
				},
			},
		},
	)
	cluster := Cluster{Interface: client, TestMode: true}

	e1, err := bitesize.LoadEnvironment("../../test/assets/environments.bitesize", "environment2")
	if err != nil {
		t.Fatalf("Unexpected err: %s", err.Error())
	}

	cluster.ApplyEnvironment(e1)

	e2, err := cluster.LoadEnvironment("environment-dev")
	if err != nil {
		t.Fatalf("Unexpected err: %s", err.Error())
	}

	if d := diff.Compare(*e1, *e2); d != "" {
		t.Errorf("Expected loaded environments to be equal, yet diff is: %s", d)
	}
}

func newDeployment(namespace, name string) *v1beta1.Deployment {
	d := v1beta1.Deployment{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1beta1.DeploymentSpec{
			Template: v1.PodTemplateSpec{},
		},
	}
	return &d
}

func newService(namespace, serviceName, portName string, portNumber int32) *v1.Service {
	labels := map[string]string{"creator": "pipeline"}

	service := v1.Service{
		ObjectMeta: v1.ObjectMeta{
			Name:      serviceName,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.ServiceSpec{
			ClusterIP: "1.1.1.1",
			Ports: []v1.ServicePort{
				{Port: portNumber, Name: portName, Protocol: "TCP"},
			},
		},
	}
	return &service
}

func validMeta(namespace, name string) v1.ObjectMeta {
	validLabels := map[string]string{"creator": "pipeline"}

	if namespace != "" {
		return v1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    validLabels,
		}
	}
	return v1.ObjectMeta{
		Name:   name,
		Labels: validLabels,
	}
}

func loadTestEnvironment() *fake.Clientset {
	capacity, _ := resource.ParseQuantity("59G")
	validLabels := map[string]string{"creator": "pipeline"}
	nsLabels := map[string]string{"environment": "Development"}
	replicaCount := int32(1)

	return fake.NewSimpleClientset(
		&v1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				Name:   "test",
				Labels: nsLabels,
			},
		},
		&v1.PersistentVolumeClaim{
			ObjectMeta: v1.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&v1.PersistentVolumeClaim{
			ObjectMeta: validMeta("test", "test"),
			Spec: v1.PersistentVolumeClaimSpec{
				AccessModes: []v1.PersistentVolumeAccessMode{
					v1.ReadWriteOnce,
				},
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceStorage: capacity,
					},
				},
			},
		},
		&v1.PersistentVolume{
			ObjectMeta: v1.ObjectMeta{
				Name:   "test",
				Labels: validLabels,
			},
		},
		&v1beta1.Deployment{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "test",
				Labels:    validLabels,
			},
			Spec: v1beta1.DeploymentSpec{
				Replicas: &replicaCount,
				Template: v1.PodTemplateSpec{
					ObjectMeta: v1.ObjectMeta{
						Name:      "test",
						Namespace: "test",
						Labels:    validLabels,
					},
					Spec: v1.PodSpec{
						Containers: []v1.Container{
							{
								Env: []v1.EnvVar{
									{
										Name:  "test",
										Value: "1",
									},
									{
										Name:  "test2",
										Value: "2",
									},
									{
										Name:  "test3",
										Value: "3",
									},
								},

								VolumeMounts: []v1.VolumeMount{
									{
										Name:      "test",
										MountPath: "/tmp/blah",
										ReadOnly:  true,
									},
								},
							},
						},
						Volumes: []v1.Volume{
							{
								Name: "test",
								VolumeSource: v1.VolumeSource{
									PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
										ClaimName: "test",
									},
								},
							},
						},
					},
				},
			},
		},
		&v1.Service{
			ObjectMeta: v1.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&v1.Service{
			ObjectMeta: validMeta("test", "test"),
			Spec: v1.ServiceSpec{
				Ports: []v1.ServicePort{
					{
						Name:     "whatevs",
						Protocol: "TCP",
						Port:     80,
					},
					{
						Name:     "whatevs2",
						Protocol: "TCP",
						Port:     8081,
					},
				},
			},
		},
		&v1.Service{
			ObjectMeta: validMeta("test", "test2"),
		},
		&v1beta1.Ingress{
			ObjectMeta: v1.ObjectMeta{
				Name:      "ts",
				Namespace: "test",
			},
		},
		&v1beta1.Ingress{
			ObjectMeta: v1.ObjectMeta{
				Name:      "test",
				Namespace: "test",
				Labels: map[string]string{
					"creator": "pipeline",
				},
			},
			Spec: v1beta1.IngressSpec{
				Rules: []v1beta1.IngressRule{
					{
						Host: "www.test.com",
					},
				},
			},
		},
	)
}

func testServicePorts(t *testing.T) {
	client := loadTestEnvironment()
	cluster := Cluster{Interface: client, TestMode: true}
	environment, err := cluster.LoadEnvironment("test")
	if err != nil {
		t.Error(err)
	}

	svc := environment.Services.FindByName("test")
	if !util.EqualArrays(svc.Ports, []int{80, 8081}) {
		t.Errorf("Ports not equal. Expected: [80 8081], got: %v", svc.Ports)
	}
}

func testFullBitesizeEnvironment(t *testing.T) {

	client := loadTestEnvironment()
	cluster := Cluster{Interface: client, TestMode: true}
	environment, err := cluster.LoadEnvironment("test")
	if err != nil {
		t.Error(err)
	}
	if environment == nil {
		t.Error("Bitesize object is nil")
	}

	if environment.Name != "Development" {
		t.Errorf("Unexpected environment name: %s", environment.Name)
	}

	if len(environment.Services) != 2 {
		t.Errorf("Unexpected service count: %d, expected: 2", len(environment.Services))
	}

	svc := environment.Services[0]
	if svc.Name != "test" {
		t.Errorf("Unexpected service name: %s, expected: test", svc.Name)
	}
	// TODO: test ingresses, env variables, replica count

	if svc.ExternalURL != "www.test.com" {
		t.Errorf("Unexpected external URL: %s, expected: www.test.com", svc.ExternalURL)
	}

	if len(svc.EnvVars) != 3 {
		t.Errorf("Unexpected environment variable count: %d, expected: 3", len(svc.EnvVars))
	}
}