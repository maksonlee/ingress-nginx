/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package connection

import (
	"testing"

	api "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-nginx/internal/ingress/annotations/parser"
	"k8s.io/ingress-nginx/internal/ingress/resolver"
)

func TestParse(t *testing.T) {
	annotation := parser.GetAnnotationWithPrefix("connection-proxy-header")

	ap := NewParser(&resolver.Mock{})
	if ap == nil {
		t.Fatalf("expected a parser.IngressAnnotation but returned nil")
	}

	testCases := []struct {
		annotations map[string]string
		expected    *Config
		expectErr   bool
	}{
		{map[string]string{annotation: "keep-alive"}, &Config{Enabled: true, Header: "keep-alive"}, false},
		{map[string]string{annotation: "not-allowed-value"}, &Config{Enabled: false}, true},
		{map[string]string{}, &Config{Enabled: false}, true},
		{nil, &Config{Enabled: false}, true},
	}

	ing := &networking.Ingress{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: networking.IngressSpec{},
	}

	for _, testCase := range testCases {
		ing.SetAnnotations(testCase.annotations)
		i, err := ap.Parse(ing)
		if (err != nil) != testCase.expectErr {
			t.Fatalf("expected error: %t got error: %t err value: %s. %+v", testCase.expectErr, err != nil, err, testCase.annotations)
		}
		p, ok := i.(*Config)
		if !ok {
			t.Fatalf("expected a Config type")
		}
		if !p.Equal(testCase.expected) {
			t.Errorf("expected %v but returned %v, annotations: %s", testCase.expected, p, testCase.annotations)
		}
	}
}
