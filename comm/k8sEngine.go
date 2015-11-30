package comm

import (
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/util"
	"k8s.io/kubernetes/pkg/api"
	"github.com/golang/glog"
)

var (
	kubeClient *client.Client
	NamespaceList = make(chan api.NamespaceList, 10)

	namespaceClient client.NamespaceInterface
	know_namespace = api.NamespaceList{}
)

func initKubeClient() bool {
	if kubeClient == nil{
		var err error
		kubeClient, err = client.NewInCluster()
		if err != nil{
			glog.Fatalf("Failed to create client: %v.\n", err)
			return false
		}
	}
	return true
}

func getNamespaceClient() client.NamespaceInterface {
	if namespaceClient == nil {
		if initKubeClient(){
			namespaceClient = kubeClient.Namespaces()
		}else{
			return nil
		}
	}

	return namespaceClient
}

func WatchResource()  {
	rateLimiter := util.NewTokenBucketRateLimiter(0.1, 1)
	for{
		rateLimiter.Accept()
		// namespace
//		go namespace()
	}
}

//func namespace(){
//	namespaces, err := namespaceClient.List(labels.Everything(), fields.Everything())
//	if err != nil || reflect.DeepEqual(namespaces.Items, know_namespace.Items) {
//		return
//	}
//
//	know_namespace = namespaces
//	NamespcaeList <- namespaces
//}
