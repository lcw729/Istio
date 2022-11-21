package controller

import (
	// "Hybrid_Cloud/Hybrid_Cloud/util/clientset"
	// "context"
	"fmt"
	// resourcev1alpha1 "hcp-pkg/apis/resource/v1alpha1"

	// Informer "istio.io/client-go/pkg/informers/externalversions/networking/v1alpha3"
	istiov1alpha3scheme "istio.io/client-go/pkg/clientset/versioned/scheme"
	// lister "istio.io/client-go/pkg/listers/networking/v1alpha3"
	istiov1alpha3clientset "istio.io/client-go/pkg/clientset/versioned"
	Informer "istio.io/client-go/pkg/informers/externalversions/networking/v1alpha3"
	lister "istio.io/client-go/pkg/listers/networking/v1alpha3"

	// deployment "hcp-pkg/kube-resource/deployment"
	// "hcp-pkg/util/clusterManager"
	"hcp-scheduler/src/scheduler"
	// "strconv"
	"time"

	// appsv1 "k8s.io/api/apps/v1"
	// hpav2beta1 "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	// "k8s.io/apimachinery/pkg/api/errors"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/tools/cache"
	// cache "k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
)

const controllerAgentName = "hcp-deployment-controller"

const (
	// SuccessSynced is used as part of the Event 'reason' when a Foo is synced
	SuccessSynced = "Synced"
	// ErrResourceExists is used as part of the Event 'reason' when a Foo fails
	// to sync due to a Deployment of the same name already existing.
	ErrResourceExists = "ErrResourceExists"

	// MessageResourceExists is the message used for Events when a resource
	// fails to sync due to a Deployment already existing
	MessageResourceExists = "Resource %q already exists and is not managed by Foo"
	// MessageResourceSynced is the message used for an Event fired when a Foo
	// is synced successfully
	MessageResourceSynced = "Foo synced successfully"
)

type Controller struct {
	kubeclientset  kubernetes.Interface
	istioclientset istiov1alpha3clientset.Interface
	istioLister    lister.VirtualServiceLister
	istioSynced    cache.InformerSynced
	workqueue      workqueue.RateLimitingInterface
	recorder       record.EventRecorder
	scheduler      *scheduler.Scheduler
}

func NewController(
	kubeclientset kubernetes.Interface,
	istioclientset istiov1alpha3clientset.Interface,
	istioInformer Informer.VirtualServiceInformer) *Controller {
	utilruntime.Must(istiov1alpha3scheme.AddToScheme(scheme.Scheme))
	klog.V(4).Infof("Creating event broadcaster")
	eventBroadCaster := record.NewBroadcaster()
	eventBroadCaster.StartStructuredLogging(0)
	eventBroadCaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("hcp")})
	recorder := eventBroadCaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})
	sched := scheduler.NewScheduler()

	controller := &Controller{
		kubeclientset:  kubeclientset,
		istioclientset: istioclientset,
		istioLister:    istioInformer.Lister(),
		istioSynced:    istioInformer.Informer().HasSynced,
		workqueue:      workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "istio"),
		recorder:       recorder,
		scheduler:      sched,
	}

	klog.Infof("Setting up event handlers")

	istioInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: controller.enqueneistio,
		UpdateFunc: func(old, new interface{}) {
			controller.enqueneistio(new)
		},
	})

	return controller
}

func (c *Controller) enqueneistio(obj interface{}) {
	var key string
	var err error
	if key, err = cache.MetaNamespaceKeyFunc(obj); err != nil {
		utilruntime.HandleError(err)
		return
	}
	c.workqueue.Add(key)
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing Informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(workers int, stopCh <-chan struct{}) error {
	defer utilruntime.HandleCrash()
	defer c.workqueue.ShutDown()

	// Start the Informer factories to begin populating the Informer caches
	klog.Infof("Starting istio controller")
	// Wait for the caches to be synced before starting workers
	klog.Infof("Waiting for Informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.istioSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	klog.Infof("Starting workers")
	// Launch two workers to process Foo resources
	for i := 0; i < workers; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	klog.Infof("Started workers")
	<-stopCh
	klog.Infof("Shutting down workers")

	return nil
}

//

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.workqueue.Get()

	if shutdown {
		return false
	}
	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.workqueue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the Informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.workqueue.Forget(obj)
			utilruntime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			// Put the item back on the workqueue to handle any transient errors.
			c.workqueue.AddRateLimited(key)
			return fmt.Errorf("error syncing '%s': %s, requeuing", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.workqueue.Forget(obj)
		klog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		utilruntime.HandleError(err)
		return true
	}

	return true
}

func (c *Controller) syncHandler(key string) error {
	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		utilruntime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	_ = namespace
	_ = name
	return nil
}
