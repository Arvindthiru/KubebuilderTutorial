package controllers

import (
	"context"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	cronjobv1 "tutorial.kubebuilder.io/project/api/v1"
)

const (
	CronjobName      = "test-cronjob"
	CronjobNamespace = "default"
	JobName          = "test-job"
)

func TestCronJobReconciler_Reconcile(t *testing.T) {
	r := &CronJobReconciler{
		Client: k8sClient,
		Scheme: scheme.Scheme,
	}
	//ctx := context.Background()

	ctlResult, _ := r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "cluster",
			Namespace: "namespace",
		},
	})

	if (ctlResult == ctrl.Result{}) {
		t.Fatalf("fail")
	}

}

func (r *CronJobReconciler) Get(ctx context.Context, name types.NamespacedName, cronjob *cronjobv1.CronJob) error {
	cronjob.TypeMeta = metav1.TypeMeta{
		APIVersion: "batch.tutorial.kubebuilder.io/v1",
		Kind:       "CronJob",
	}
	cronjob.ObjectMeta = metav1.ObjectMeta{
		Name:      CronjobName,
		Namespace: CronjobNamespace,
	}
	cronjob.Spec = cronjobv1.CronJobSpec{
		Schedule: "1 * * * *",
		JobTemplate: batchv1beta1.JobTemplateSpec{
			Spec: batchv1.JobSpec{
				// For simplicity, we only fill out the required fields.
				Template: v1.PodTemplateSpec{
					Spec: v1.PodSpec{
						// For simplicity, we only fill out the required fields.
						Containers: []v1.Container{
							{
								Name:  "test-container",
								Image: "test-image",
							},
						},
						RestartPolicy: v1.RestartPolicyOnFailure,
					},
				},
			},
		},
	}
	return nil
}

func (r *CronJobReconciler) List(ctx context.Context, childJobs *batchv1.JobList, namespace client.InNamespace, fields client.MatchingFields) error {
	testJob := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      JobName,
			Namespace: CronjobNamespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					// For simplicity, we only fill out the required fields.
					Containers: []v1.Container{
						{
							Name:  "test-container",
							Image: "test-image",
						},
					},
					RestartPolicy: v1.RestartPolicyOnFailure,
				},
			},
		},
		Status: batchv1.JobStatus{
			Active: 2,
		},
	}
	var Items []batchv1.Job = []batchv1.Job{*testJob}
	childJobs.Items = Items

	return nil
}
