package comm

import (
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/labels"
	"k8s.io/kubernetes/pkg/fields"
	"github.com/go-martini/martini"
	"github.com/golang/glog"
)

func Namespaces(r martini.Router){
	r.Get("/", all)
	r.Get("/:name", namespace)
}

func all() *api.NamespaceList {
	client := getNamespaceClient()
	glog.Info("namespaceclient ", client)
	if client != nil{
		list, err := client.List(labels.Everything(), fields.Everything())
		glog.Info("namespaces list", list)
		if err == nil{
			return list
		}
		glog.Errorln("====>>", err)
	}
	return &api.NamespaceList{}
}

func namespace(params martini.Params) string {
	return "Hello " + params["name"]
}
