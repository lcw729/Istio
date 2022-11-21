// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1beta1 "istio.io/client-go/pkg/apis/networking/v1beta1"
	networkingv1beta1 "istio.io/client-go/pkg/applyconfiguration/networking/v1beta1"
	scheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// DestinationRulesGetter has a method to return a DestinationRuleInterface.
// A group's client should implement this interface.
type DestinationRulesGetter interface {
	DestinationRules(namespace string) DestinationRuleInterface
}

// DestinationRuleInterface has methods to work with DestinationRule resources.
type DestinationRuleInterface interface {
	Create(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.CreateOptions) (*v1beta1.DestinationRule, error)
	Update(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.UpdateOptions) (*v1beta1.DestinationRule, error)
	UpdateStatus(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.UpdateOptions) (*v1beta1.DestinationRule, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.DestinationRule, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.DestinationRuleList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.DestinationRule, err error)
	Apply(ctx context.Context, destinationRule *networkingv1beta1.DestinationRuleApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.DestinationRule, err error)
	ApplyStatus(ctx context.Context, destinationRule *networkingv1beta1.DestinationRuleApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.DestinationRule, err error)
	DestinationRuleExpansion
}

// destinationRules implements DestinationRuleInterface
type destinationRules struct {
	client rest.Interface
	ns     string
}

// newDestinationRules returns a DestinationRules
func newDestinationRules(c *NetworkingV1beta1Client, namespace string) *destinationRules {
	return &destinationRules{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the destinationRule, and returns the corresponding destinationRule object, and an error if there is any.
func (c *destinationRules) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.DestinationRule, err error) {
	result = &v1beta1.DestinationRule{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("destinationrules").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of DestinationRules that match those selectors.
func (c *destinationRules) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.DestinationRuleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.DestinationRuleList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("destinationrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested destinationRules.
func (c *destinationRules) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("destinationrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a destinationRule and creates it.  Returns the server's representation of the destinationRule, and an error, if there is any.
func (c *destinationRules) Create(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.CreateOptions) (result *v1beta1.DestinationRule, err error) {
	result = &v1beta1.DestinationRule{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("destinationrules").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(destinationRule).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a destinationRule and updates it. Returns the server's representation of the destinationRule, and an error, if there is any.
func (c *destinationRules) Update(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.UpdateOptions) (result *v1beta1.DestinationRule, err error) {
	result = &v1beta1.DestinationRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("destinationrules").
		Name(destinationRule.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(destinationRule).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *destinationRules) UpdateStatus(ctx context.Context, destinationRule *v1beta1.DestinationRule, opts v1.UpdateOptions) (result *v1beta1.DestinationRule, err error) {
	result = &v1beta1.DestinationRule{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("destinationrules").
		Name(destinationRule.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(destinationRule).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the destinationRule and deletes it. Returns an error if one occurs.
func (c *destinationRules) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("destinationrules").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *destinationRules) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("destinationrules").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched destinationRule.
func (c *destinationRules) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.DestinationRule, err error) {
	result = &v1beta1.DestinationRule{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("destinationrules").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied destinationRule.
func (c *destinationRules) Apply(ctx context.Context, destinationRule *networkingv1beta1.DestinationRuleApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.DestinationRule, err error) {
	if destinationRule == nil {
		return nil, fmt.Errorf("destinationRule provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(destinationRule)
	if err != nil {
		return nil, err
	}
	name := destinationRule.Name
	if name == nil {
		return nil, fmt.Errorf("destinationRule.Name must be provided to Apply")
	}
	result = &v1beta1.DestinationRule{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("destinationrules").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *destinationRules) ApplyStatus(ctx context.Context, destinationRule *networkingv1beta1.DestinationRuleApplyConfiguration, opts v1.ApplyOptions) (result *v1beta1.DestinationRule, err error) {
	if destinationRule == nil {
		return nil, fmt.Errorf("destinationRule provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(destinationRule)
	if err != nil {
		return nil, err
	}

	name := destinationRule.Name
	if name == nil {
		return nil, fmt.Errorf("destinationRule.Name must be provided to Apply")
	}

	result = &v1beta1.DestinationRule{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("destinationrules").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
