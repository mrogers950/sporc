package controllers

import (
	sporcv1alpha1 "github.com/mrogers950/sporc/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func newResponseConfig(name, namespace string) *sporcv1alpha1.ResponseConfig {
	return &sporcv1alpha1.ResponseConfig{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: sporcv1alpha1.ResponseConfigSpec{
			ResponseConfigMap: sporcv1alpha1.ResponseConfigMapRef{
				Name:      name,
				Namespace: namespace,
			},
		},
	}
}
