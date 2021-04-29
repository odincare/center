package odin

import (
	"github.com/nacos-group/nacos-sdk-go/vo"
	"time"
)

type Instance struct {
	Valid       bool              `json:"valid"`
	Marked      bool              `json:"marked"`
	InstanceId  string            `json:"instanceId"`
	Port        uint64            `json:"port"`
	Ip          string            `json:"ip"`
	Weight      float64           `json:"weight"`
	Metadata    map[string]string `json:"metadata"`
	ClusterName string            `json:"clusterName"`
	ServiceName string            `json:"serviceName"`
	Enable      bool              `json:"enabled"`
	Healthy     bool              `json:"healthy"`
	Ephemeral   bool              `json:"ephemeral"`
}

type Service struct {
	Dom             string            `json:"dom"`
	CacheMillis     uint64            `json:"cacheMillis"`
	UseSpecifiedURL bool              `json:"useSpecifiedUrl"`
	Hosts           []Instance        `json:"hosts"`
	Checksum        string            `json:"checksum"`
	LastRefTime     uint64            `json:"lastRefTime"`
	Env             string            `json:"env"`
	Clusters        string            `json:"clusters"`
	Metadata        map[string]string `json:"metadata"`
	Name            string            `json:"name"`
}

type ServiceDetail struct {
	Service  ServiceInfo `json:"service"`
	Clusters []Cluster   `json:"clusters"`
}

type ServiceInfo struct {
	App              string            `json:"app"`
	Group            string            `json:"group"`
	HealthCheckMode  string            `json:"healthCheckMode"`
	Metadata         map[string]string `json:"metadata"`
	Name             string            `json:"name"`
	ProtectThreshold float64           `json:"protectThreshold"`
	Selector         ServiceSelector   `json:"selector"`
}

type ServiceSelector struct {
	Selector string
}

type Cluster struct {
	ServiceName      string               `json:"serviceName"`
	Name             string               `json:"name"`
	HealthyChecker   ClusterHealthChecker `json:"healthyChecker"`
	DefaultPort      uint64               `json:"defaultPort"`
	DefaultCheckPort uint64               `json:"defaultCheckPort"`
	UseIPPort4Check  bool                 `json:"useIpPort4Check"`
	Metadata         map[string]string    `json:"metadata"`
}

type ClusterHealthChecker struct {
	Type string `json:"type"`
}

type SubscribeService struct {
	ClusterName string            `json:"clusterName"`
	Enable      bool              `json:"enable"`
	InstanceId  string            `json:"instanceId"`
	Ip          string            `json:"ip"`
	Metadata    map[string]string `json:"metadata"`
	Port        uint64            `json:"port"`
	ServiceName string            `json:"serviceName"`
	Valid       bool              `json:"valid"`
	Weight      float64           `json:"weight"`
}

type BeatInfo struct {
	Ip          string            `json:"ip"`
	Port        uint64            `json:"port"`
	Weight      float64           `json:"weight"`
	ServiceName string            `json:"serviceName"`
	Cluster     string            `json:"cluster"`
	Metadata    map[string]string `json:"metadata"`
	Scheduled   bool              `json:"scheduled"`
	Period      time.Duration     `json:"-"`
	State       int32             `json:"-"`
}

type ExpressionSelector struct {
	Type       string `json:"type"`
	Expression string `json:"expression"`
}

type ServiceList struct {
	Count int64    `json:"count"`
	Doms  []string `json:"doms"`
}

type (
	RegisterInstanceParam        vo.RegisterInstanceParam
	DeregisterInstanceParam      vo.DeregisterInstanceParam
	GetServiceParam              vo.GetServiceParam
	SelectAllInstancesParam      vo.SelectAllInstancesParam
	SelectInstancesParam         vo.SelectInstancesParam
	SelectOneHealthInstanceParam vo.SelectOneHealthInstanceParam
	SubscribeParam               vo.SubscribeParam
	GetAllServiceInfoParam       vo.GetAllServiceInfoParam
)

//----注册发现---开始--

func (o *OdinCenter) Register() {

}
func (o *OdinCenter) DeRegister() {

}

func (o *OdinCenter) RegisterServiceInstance(param RegisterInstanceParam) (bool, error) {
	sParam := vo.RegisterInstanceParam{}
	assignParam(&sParam, &param)
	return o.rClient.RegisterInstance(sParam)
}

func (o *OdinCenter) DeRegisterServiceInstance(param DeregisterInstanceParam) (bool, error) {
	sParam := vo.DeregisterInstanceParam{}
	assignParam(&sParam, &param)
	return o.rClient.DeregisterInstance(sParam)
}

func (o *OdinCenter) GetService(param GetServiceParam) *Service {
	sParam := vo.GetServiceParam{}
	rSrv := Service{}
	assignParam(&sParam, &param)
	service, err := o.rClient.GetService(sParam)
	if err != nil {
		return nil
	}
	assignResult(&rSrv, &service)
	return &rSrv
}

func (o *OdinCenter) SelectAllInstances(param SelectAllInstancesParam) []Instance {
	var result []Instance
	sParam := vo.SelectAllInstancesParam{}
	assignParam(&sParam, &param)
	instances, _ := o.rClient.SelectAllInstances(sParam)
	assignResult(&result, &instances)
	return result
}

func (o *OdinCenter) SelectInstances(param SelectInstancesParam) []*Instance {
	var result []*Instance
	sParam := vo.SelectInstancesParam{}
	assignParam(&sParam, &param)
	instances, _ := o.rClient.SelectInstances(sParam)
	assignResult(&result, &instances)
	return result
}

func (o *OdinCenter) SelectOneHealthyInstance(param SelectOneHealthInstanceParam) *Instance {
	var result Instance
	sParam := vo.SelectOneHealthInstanceParam{}
	assignParam(&sParam, &param)
	instances, _ := o.rClient.SelectOneHealthyInstance(sParam)
	assignResult(&result, &instances)
	return &result
}

func (o *OdinCenter) Subscribe(param SubscribeParam) error {
	sParam := vo.SubscribeParam{}
	assignParam(&sParam, &param)
	return o.rClient.Subscribe(&sParam)
}

func (o *OdinCenter) UnSubscribe(param SubscribeParam) error {
	sParam := vo.SubscribeParam{}
	assignParam(&sParam, &param)
	return o.rClient.Unsubscribe(&sParam)
}

func (o *OdinCenter) GetAllService(param GetAllServiceInfoParam) ServiceList {
	sParam := vo.GetAllServiceInfoParam{}
	assignParam(&sParam, &param)
	if sParam.PageSize == 0 {
		sParam.PageSize = 20
	}
	service, _ := o.rClient.GetAllServicesInfo(sParam)
	result := ServiceList{}
	assignResult(&result, &service)
	return result

}

//----注册发现---结束--
