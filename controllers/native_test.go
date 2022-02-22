package controllers

import (
	"context"
	"errors"
	"fmt"
	batchv1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	client "sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
	cronjobv1 "tutorial.kubebuilder.io/project/api/v1"
)

const (
	CronjobName      = "test-cronjob"
	CronjobNamespace = "default"
	JobName          = "test-job"
)

func TestCronJobReconciler_Reconcile(t *testing.T) {
	mockClient := &mockClient1{
		Client:  k8sClient,
		errBool: false,
	}
	r := &CronJobReconciler{
		Client: mockClient,
		Scheme: scheme.Scheme,
		Clock:  mockClock{},
	}
	//ctx := context.Background()

	ctlResult, err := r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "cluster",
			Namespace: "namespace",
		},
	})

	if err != nil {
		t.Fatalf("fail")
	}
	if (ctlResult != ctrl.Result{}) {
		t.Fatalf("fail")
	}

}

func TestCronJobReconciler_Reconcile_WithGetError(t *testing.T) {
	mockClient := &mockClient1{
		Client:  k8sClient,
		errBool: true,
	}
	r := &CronJobReconciler{
		Client: mockClient,
		Scheme: scheme.Scheme,
		Clock:  mockClock{},
	}
	//ctx := context.Background()

	_, err := r.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "cluster",
			Namespace: "namespace",
		},
	})

	if err != nil {
		t.Fatalf("fail")
	}
}

func (m *mockClient1) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	if m.errBool {
		return errors.New("unable to fetch CronJob")
	}
	cronjob := obj.(*cronjobv1.CronJob)
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

func (m *mockClient1) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	childJobs := list.(*batchv1.JobList)
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

type mockClient1 struct {
	client.Client
	errBool bool
}

//func (m *mockClient1) Get(ctx context.Context, key client.ObjectKey, obj client.Object) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) Scheme() *runtime.Scheme {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (m *mockClient1) RESTMapper() meta.RESTMapper {
//	//TODO implement me
//	panic("implement me")
//}

type mockStatusWriter struct {
}

type mockClock struct{}

func (_ mockClock) Now() time.Time {
	fmt.Println("hello")
	return time.Now()
}

func (m *mockClient1) Status() client.StatusWriter {
	fmt.Println("here")
	return &mockStatusWriter{}
}

func (m *mockStatusWriter) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	fmt.Println("there")
	return nil
}

func (m *mockStatusWriter) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}
