package openshift

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/remotecommand"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/types"

	"github.com/integr8ly/cloud-resource-operator/pkg/resources"

	errorUtil "github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	controllerruntime "sigs.k8s.io/controller-runtime"

	k8serr "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/integr8ly/cloud-resource-operator/pkg/apis/integreatly/v1alpha1"
	"github.com/integr8ly/cloud-resource-operator/pkg/providers"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	postgresProviderName = "openshift-postgres-template"
	// default openshift create paramaters
	defaultPostgresPort        = 5432
	defaultPostgresUser        = "keycloak"
	defaultPostgresPassword    = "password"
	defaultPostgresUserKey     = "user"
	defaultPostgresDatabase    = "keycloak"
	defaultPostgresPasswordKey = "password"
	defaultPostgresDatabaseKey = "database"
	defaultCredentialsSec      = "postgres-credentials"
)

// PostgresStrat to be used to unmarshal strat map
type PostgresStrat struct {
	_ struct{} `type:"structure"`

	PostgresDeploymentSpec *appsv1.DeploymentSpec        `json:"deploymentSpec"`
	PostgresServiceSpec    *v1.ServiceSpec               `json:"serviceSpec"`
	PostgresPVCSpec        *v1.PersistentVolumeClaimSpec `json:"pvcSpec"`
	PostgresSecretData     map[string]string             `json:"secretData"`
}

type OpenShiftPostgresProvider struct {
	Client        client.Client
	Logger        *logrus.Entry
	ConfigManager ConfigManager
}

func NewOpenShiftPostgresProvider(client client.Client, logger *logrus.Entry) *OpenShiftPostgresProvider {
	return &OpenShiftPostgresProvider{
		Client:        client,
		Logger:        logger.WithFields(logrus.Fields{"provider": postgresProviderName}),
		ConfigManager: NewDefaultConfigManager(client),
	}
}

func (p *OpenShiftPostgresProvider) GetName() string {
	return postgresProviderName
}

func (p *OpenShiftPostgresProvider) SupportsStrategy(d string) bool {
	return d == providers.OpenShiftDeploymentStrategy
}

func (p *OpenShiftPostgresProvider) CreatePostgres(ctx context.Context, ps *v1alpha1.Postgres) (*providers.PostgresInstance, v1alpha1.StatusMessage, error) {
	// handle provider-specific finalizer
	if err := resources.CreateFinalizer(ctx, p.Client, ps, DefaultFinalizer); err != nil {
		errMsg := "failed to set finalizer"
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}

	// get postgres config
	postgresCfg, _, err := p.getPostgresConfig(ctx, ps)
	if err != nil {
		errMsg := fmt.Sprintf("failed to retrieve openshift postgres config for instance %s", ps.Name)
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrapf(err, errMsg)
	}

	// deploy pvc
	if err := p.CreatePVC(ctx, buildDefaultPostgresPVC(ps), postgresCfg); err != nil {
		errMsg := fmt.Sprintf("failed to create or update postgres PVC for instance %s", ps.Name)
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}
	// deploy credentials secret
	if err := p.CreateSecret(ctx, buildDefaultPostgresSecret(ps), postgresCfg); err != nil {
		errMsg := fmt.Sprintf("failed to create or update postgres secret for instance %s", ps.Name)
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}
	// deploy deployment
	if err := p.CreateDeployment(ctx, buildDefaultPostgresDeployment(ps), postgresCfg); err != nil {
		errMsg := fmt.Sprintf("failed to create or update postgres deployment for instance %s", ps.Name)
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}
	// deploy service
	if err := p.CreateService(ctx, buildDefaultPostgresService(ps), postgresCfg); err != nil {
		errMsg := fmt.Sprintf("failed to create or update postgres service for instance %s", ps.Name)
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}

	// check deployment status
	dpl := &appsv1.Deployment{}
	err = p.Client.Get(ctx, types.NamespacedName{Name: ps.Name, Namespace: ps.Namespace}, dpl)
	if err != nil {
		errMsg := "failed to get postgres deployment"
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}

	// get the cred secret
	sec := &v1.Secret{}
	err = p.Client.Get(ctx, types.NamespacedName{Name: defaultCredentialsSec, Namespace: ps.Namespace}, sec)
	if err != nil {
		errMsg := "failed to get postgres creds"
		return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}

	// check if deployment is ready and return connection details
	for _, s := range dpl.Status.Conditions {
		if s.Type == appsv1.DeploymentAvailable && s.Status == "True" {
			// create user, db and set privilages
			pgCommand := getPostgresProvisionCommand()
			podToExec, err := k8sclient.getDeploymentPod(ps.Name, ps.Namespace)
			if err != nil {
				errMsg := "failed to get postgres pod"
				return nil, v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}
			updatedPod := execIntoPod(podToExec, pgCommand, ps.Namespace)

			if updatedPod {
				p.Logger.Info("Found postgres deployment")
				return &providers.PostgresInstance{
					DeploymentDetails: &providers.PostgresDeploymentDetails{
						Username: string(sec.Data["user"]),
						Password: string(sec.Data["password"]),
						Database: string(sec.Data["database"]),
						Host:     fmt.Sprintf("%s.%s.svc.cluster.local", ps.Name, ps.Namespace),
						Port:     defaultPostgresPort,
					},
				}, "postgres deployment is complete", nil
			}
		}
	}

	// deployment is in progress
	p.Logger.Info("postgres deployment is not ready")
	return nil, "creation in progress", nil
}

func (p *OpenShiftPostgresProvider) DeletePostgres(ctx context.Context, ps *v1alpha1.Postgres) (v1alpha1.StatusMessage, error) {
	// check deployment status
	dpl := &appsv1.Deployment{}
	err := p.Client.Get(ctx, types.NamespacedName{Name: ps.Name, Namespace: ps.Namespace}, dpl)
	if err != nil {
		if k8serr.IsNotFound(err) {
			return v1alpha1.StatusEmpty, nil
		}
		errMsg := "failed to get postgres deployment"
		return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
	}

	for _, s := range dpl.Status.Conditions {
		if s.Type == appsv1.DeploymentAvailable && s.Status == "True" {
			// delete service
			p.Logger.Info("Deleting postgres service")
			svc := &v1.Service{
				ObjectMeta: controllerruntime.ObjectMeta{
					Name:      ps.Name,
					Namespace: ps.Namespace,
				},
			}
			err = p.Client.Delete(ctx, svc)
			if err != nil && !k8serr.IsNotFound(err) {
				errMsg := "failed to delete postgres service"
				return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}

			// delete pvc
			p.Logger.Info("Deleting postgres persistent volume claim")
			pvc := &v1.PersistentVolumeClaim{
				ObjectMeta: controllerruntime.ObjectMeta{
					Name:      ps.Name,
					Namespace: ps.Namespace,
				},
			}
			err = p.Client.Delete(ctx, pvc)
			if err != nil && !k8serr.IsNotFound(err) {
				errMsg := "failed to delete postgres persistent volume claim"
				return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}

			// delete secret
			p.Logger.Info("Deleting postgres secret")
			sec := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      defaultCredentialsSec,
					Namespace: ps.Namespace,
				},
			}
			err = p.Client.Delete(ctx, sec)
			if err != nil && !k8serr.IsNotFound(err) {
				errMsg := "failed to deleted postgres secrets"
				return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}

			// clean up objects
			p.Logger.Info("Deleting postgres deployment")
			err = p.Client.Delete(ctx, dpl)
			if err != nil && !k8serr.IsNotFound(err) {
				errMsg := "failed to delete postgres deployment"
				return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}

			// remove the finalizer added by the provider
			p.Logger.Info("Removing postgres finalizer")
			resources.RemoveFinalizer(&ps.ObjectMeta, DefaultFinalizer)
			if err := p.Client.Update(ctx, ps); err != nil {
				errMsg := "failed to update instance as part of the postgres finalizer reconcile"
				return v1alpha1.StatusMessage(errMsg), errorUtil.Wrap(err, errMsg)
			}

			p.Logger.Infof("deletion handler for postgres %s in namespace %s finished successfully", ps.Name, ps.Namespace)
		}
	}

	return "deletion in progress", nil
}

// getPostgresConfig retrieves the postgres config from the cloud-resources-openshift-strategies configmap
func (p *OpenShiftPostgresProvider) getPostgresConfig(ctx context.Context, ps *v1alpha1.Postgres) (*PostgresStrat, *StrategyConfig, error) {
	stratCfg, err := p.ConfigManager.ReadStorageStrategy(ctx, providers.PostgresResourceType, ps.Spec.Tier)
	if err != nil {
		return nil, nil, errorUtil.Wrap(err, "failed to read openshift strategy config")
	}

	// unmarshal the postgres config
	postgresCfg := &PostgresStrat{}
	if err := json.Unmarshal(stratCfg.RawStrategy, postgresCfg); err != nil {
		return nil, nil, errorUtil.Wrap(err, "failed to unmarshal openshift postgres configuration")
	}

	return postgresCfg, stratCfg, nil
}

func (p *OpenShiftPostgresProvider) CreateDeployment(ctx context.Context, d *appsv1.Deployment, postgresCfg *PostgresStrat) error {
	or, err := immutableCreateOrUpdate(ctx, p.Client, d, func(existing runtime.Object) error {
		e := existing.(*appsv1.Deployment)

		if postgresCfg.PostgresDeploymentSpec == nil {
			e.Spec = d.Spec
			return nil
		}

		e.Spec = *postgresCfg.PostgresDeploymentSpec
		return nil
	})
	if err != nil {
		return errorUtil.Wrapf(err, "failed to create or update deployment %s, action was %s", d.Name, or)
	}
	return nil
}

func (p *OpenShiftPostgresProvider) CreateService(ctx context.Context, s *v1.Service, postgresCfg *PostgresStrat) error {
	or, err := immutableCreateOrUpdate(ctx, p.Client, s, func(existing runtime.Object) error {
		e := existing.(*v1.Service)

		if postgresCfg.PostgresServiceSpec == nil {
			clusterIP := e.Spec.ClusterIP
			e.Spec = s.Spec
			e.Spec.ClusterIP = clusterIP
			return nil
		}

		e.Spec = *postgresCfg.PostgresServiceSpec
		return nil
	})
	if err != nil {
		return errorUtil.Wrapf(err, "failed to create or update service %s, action was %s", s.Name, or)
	}
	return nil
}

func (p *OpenShiftPostgresProvider) CreateSecret(ctx context.Context, s *v1.Secret, postgresCfg *PostgresStrat) error {
	or, err := controllerutil.CreateOrUpdate(ctx, p.Client, s, func(existing runtime.Object) error {
		e := existing.(*v1.Secret)

		if postgresCfg.PostgresSecretData == nil {
			e.Data = s.Data
			return nil
		}

		e.StringData = postgresCfg.PostgresSecretData
		return nil
	})
	if err != nil {
		return errorUtil.Wrapf(err, "failed to create or update secret %s, action was %s", s.Name, or)
	}
	return nil
}

func (p *OpenShiftPostgresProvider) CreatePVC(ctx context.Context, pvc *v1.PersistentVolumeClaim, postgresCfg *PostgresStrat) error {
	or, err := immutableCreateOrUpdate(ctx, p.Client, pvc, func(existing runtime.Object) error {
		e := existing.(*v1.PersistentVolumeClaim)

		if strings.ToLower(string(e.Status.Phase)) != "bound" {
			return nil
		}
		if postgresCfg.PostgresPVCSpec == nil {
			e.Spec.Resources.Requests = pvc.Spec.Resources.Requests
			return nil
		}

		e.Spec.Resources.Requests = postgresCfg.PostgresPVCSpec.Resources.Requests
		return nil
	})
	if err != nil {
		return errorUtil.Wrapf(err, "failed to create or update persistent volume claim %s, action was %s", pvc.Name, or)
	}
	return nil
}

func buildDefaultPostgresService(ps *v1alpha1.Postgres) *v1.Service {
	return &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ps.Name,
			Namespace: ps.Namespace,
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name:       "postgresql",
					Protocol:   v1.ProtocolTCP,
					Port:       int32(defaultPostgresPort),
					TargetPort: intstr.FromInt(defaultPostgresPort),
				},
			},
			Selector: map[string]string{"deployment": ps.Name},
		},
	}
}

func buildDefaultPostgresPVC(ps *v1alpha1.Postgres) *v1.PersistentVolumeClaim {
	return &v1.PersistentVolumeClaim{
		TypeMeta: metav1.TypeMeta{
			Kind:       "PersistentVolumeClaim",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ps.Name,
			Namespace: ps.Namespace,
		},
		Spec: v1.PersistentVolumeClaimSpec{
			AccessModes: []v1.PersistentVolumeAccessMode{"ReadWriteOnce"},
			Resources: v1.ResourceRequirements{
				Requests: v1.ResourceList{
					"storage": resource.MustParse("1Gi"),
				},
			},
		},
	}
}

func buildDefaultPostgresDeployment(ps *v1alpha1.Postgres) *appsv1.Deployment {
	return &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      ps.Name,
			Namespace: ps.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Strategy: appsv1.DeploymentStrategy{
				Type: appsv1.RecreateDeploymentStrategyType,
			},
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"deployment": ps.Name,
				},
			},
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: ps.Name,
							VolumeSource: v1.VolumeSource{
								PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
									ClaimName: ps.Name,
								},
							},
						},
					},
					Containers: buildDefaultPostgresPodContainers(ps),
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"deployment": ps.Name,
					},
				},
			},
		},
	}
}

func buildDefaultPostgresPodContainers(ps *v1alpha1.Postgres) []v1.Container {
	return []v1.Container{
		{
			Name:  ps.Name,
			Image: "registry.redhat.io/rhscl/postgresql-96-rhel7",
			Ports: []v1.ContainerPort{
				{
					ContainerPort: int32(defaultPostgresPort),
					Protocol:      v1.ProtocolTCP,
				},
			},
			Env: []v1.EnvVar{
				envVarFromSecret("POSTGRESQL_USER", defaultCredentialsSec, defaultPostgresUserKey),
				envVarFromSecret("POSTGRESQL_PASSWORD", defaultCredentialsSec, defaultPostgresPasswordKey),
				envVarFromSecret("POSTGRESQL_DATABASE", defaultCredentialsSec, defaultPostgresDatabaseKey),
			},
			Resources: v1.ResourceRequirements{
				Limits: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("250m"),
					v1.ResourceMemory: resource.MustParse("2Gi"),
				},
				Requests: v1.ResourceList{
					v1.ResourceCPU:    resource.MustParse("50m"),
					v1.ResourceMemory: resource.MustParse("512Mi"),
				},
			},
			VolumeMounts: []v1.VolumeMount{
				{
					Name:      ps.Name,
					MountPath: "/var/lib/pgsql/data",
				},
			},
			LivenessProbe: &v1.Probe{
				Handler: v1.Handler{
					TCPSocket: &v1.TCPSocketAction{
						Port: intstr.IntOrString{
							Type:   intstr.Int,
							IntVal: int32(defaultPostgresPort),
						},
					},
				},
				InitialDelaySeconds: 30,
				PeriodSeconds:       10,
				TimeoutSeconds:      0,
				SuccessThreshold:    0,
				FailureThreshold:    0,
			},
			ReadinessProbe: &v1.Probe{
				Handler: v1.Handler{
					Exec: &v1.ExecAction{
						Command: []string{"/bin/sh", "-i", "-c", "psql -h 127.0.0.1 -U $POSTGRESQL_USER -q -d $POSTGRESQL_DATABASE -c 'SELECT 1'"}},
				},
				InitialDelaySeconds: 10,
				PeriodSeconds:       30,
				TimeoutSeconds:      5,
				SuccessThreshold:    0,
				FailureThreshold:    0,
			},
			ImagePullPolicy: v1.PullIfNotPresent,
		},
	}
}

func buildDefaultPostgresSecret(ps *v1alpha1.Postgres) *v1.Secret {
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      defaultCredentialsSec,
			Namespace: ps.Namespace,
		},
		StringData: map[string]string{
			"user":     defaultPostgresUser,
			"password": defaultPostgresPassword,
			"database": defaultPostgresDatabase,
		},
		Type: v1.SecretTypeOpaque,
	}
}

// create an environment variable referencing a secret
func envVarFromSecret(envVarName string, secretName, secretKey string) v1.EnvVar {
	return v1.EnvVar{
		Name: envVarName,
		ValueFrom: &v1.EnvVarSource{
			SecretKeyRef: &v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: secretName,
				},
				Key: secretKey,
			},
		},
	}
}

type k8s struct {
	clientset kubernetes.Interface
}

var (
	k8sclient = getK8Client()
)

// return the postgres command to create user, database and grant privileges,
func getPostgresProvisionCommand() (command string) {
	command = "OUT=$(psql postgres -tAc \"SELECT 1 FROM pg_roles WHERE rolname='keycloak'\"); " +
		"if [ $OUT -eq 1 ]; then echo \"DB exists\"; exit 0; fi " +
		"&& psql -c \"CREATE USER keycloak WITH PASSWORD '" + defaultPostgresPassword + "'\" " +
		"&& psql -c \"CREATE DATABASE " + defaultPostgresDatabase + "\" " +
		"&& psql -c \"GRANT ALL PRIVILEGES ON DATABASE " + defaultPostgresDatabase + " TO keycloak\" " +
		"&& psql -c \"ALTER USER " + defaultPostgresUser + " WITH SUPERUSER\""

	return command
}

// exec into a pod
func execIntoPod(podName string, provisionCommand string, ns string) (provisioned bool) {
	command := []string{"/bin/bash", "-c", provisionCommand}
	logrus.Infof("Running exec to %s in pod", podName)
	_, stderr, err := k8sclient.runExec(command, podName, ns)
	if err != nil {
		logrus.Errorf("Error exec into pod: %v: , command: %s", err, command)
		logrus.Errorf(stderr)
		return false
	}
	logrus.Info("Exec successfully completed")
	return true
}

// run exec command on pod
func (cl *k8s) runExec(command []string, podName, namespace string) (string, string, error) {

	req := cl.clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	req.VersionedParams(&v1.PodExecOptions{
		Command: command,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}, scheme.ParameterCodec)

	cfg, _ := config.GetConfig()
	exec, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	var stdin io.Reader
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return stdout.String(), stderr.String(), err
	}

	return stdout.String(), stderr.String(), nil
}

// return the k8s client
func getK8Client() *k8s {
	cfg, err := config.GetConfig()
	if err != nil {
		logrus.Errorf(err.Error())
	}
	client := k8s{}
	client.clientset, err = kubernetes.NewForConfig(cfg)
	if err != nil {
		logrus.Errorf(err.Error())
		return nil
	}
	return &client
}

//getDeploymentPod queries all pods is a selected namespace by LabelSelector = deployment
func (cl *k8s) getDeploymentPod(name string, ns string) (podName string, err error) {
	api := cl.clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "deployment=" + name,
	}
	podList, _ := api.Pods(ns).List(listOptions)
	podListItems := podList.Items
	if len(podListItems) == 0 {
		logrus.Errorf("Failed to find pod to exec into. List of pods: %v", podListItems)
		return "", err
	}
	podName = podListItems[0].Name
	return podName, nil
}
