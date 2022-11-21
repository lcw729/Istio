package main

import (
	"context"
	"flag"
	"fmt"

	"hcp-pkg/util/clusterManager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog/v2"

	// informers "Hybrid_Cloud/pkg/client/resource/v1alpha1/informers/externalversions"
	icv "istio.io/client-go/pkg/clientset/versioned"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()

	cm, err := clusterManager.NewClusterManager()
	if err != nil {
		klog.Errorln(err)
	}

	cm.HCPResource_Client.HcpV1alpha1().HCPDeployments("hcp").Get(context.TODO(), "productpage", metav1.GetOptions{})
	ic := icv.NewForConfigOrDie(cm.Host_config)
	target, err := ic.NetworkingV1alpha3().VirtualServices("default").Get(context.TODO(), "vs-hello", metav1.GetOptions{})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(target.Name)

	// stopCh := signals.SetupSignalHandler()
	// kubeInformerFactory := kubeinformers.NewSharedInformerFactory(cm.Host_kubeClient, time.Second*30)
	// IstioInformerFactory := informers.NewVirtualServiceInformer(ic, "hcp", time.Second*30, nil)

	// controller := controller.NewController(cm.Host_kubeClient, ic, IstioInformerFactory)
	// kubeInformerFactory.Start(stopCh)
	// IstioInformerFactory.Start(stopCh)
	// if err := controller.Run(2, stopCh); err != nil {
	// 	klog.Fatalf("Error running controller: %s", err.Error())
	// }
}
