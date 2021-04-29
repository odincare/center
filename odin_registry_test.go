package odin

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/sanity-io/litter"
	"testing"
	"time"
)

func TestServiceRegistry(t *testing.T) {
	client := NewOdinCenter(
	//options.Addrs("d.nacos.sixlai.com:80/nacos"),
	//options.NameSpaceId("3597fb11-b3a3-4649-a1e9-13623739e469"),
	)
	client.RegisterServiceInstance(RegisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "test_nacos_service",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
	})
	time.Sleep(600 * time.Second)
}

func TestServiceRegistry1(t *testing.T) {
	client := NewOdinCenter()
	client.RegisterServiceInstance(RegisterInstanceParam{
		Ip:          "10.0.0.99",
		Port:        8848,
		ServiceName: "test_nacos_service_XXXoptions",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "beijing"},
	})
	time.Sleep(600 * time.Second)
}

func TestServiceDeRegistry(t *testing.T) {
	client := NewOdinCenter()
	client.DeRegisterServiceInstance(DeregisterInstanceParam{
		Ip:          "10.0.0.10",
		Port:        8848,
		ServiceName: "demo.go",
		Ephemeral:   true,
	})
	time.Sleep(100 * time.Second)
}
func TestGetService(t *testing.T) {
	client := NewOdinCenter()
	srv := client.GetService(GetServiceParam{
		ServiceName: "test_nacos_service",
	})
	litter.Dump(srv)
}
func TestSelectAllInstances(t *testing.T) {
	client := NewOdinCenter()
	res := client.SelectAllInstances(SelectAllInstancesParam{
		GroupName:   "DEFAULT_GROUP",
		ServiceName: "test_nacos_service",
	})
	litter.Dump(res)
}
func TestSelectInstances(t *testing.T) {
	client := NewOdinCenter()
	res := client.SelectInstances(SelectInstancesParam{
		GroupName:   "DEFAULT_GROUP",
		ServiceName: "test_nacos_service",
	})
	litter.Dump(res)
}
func TestStruct(t *testing.T) {

	as := vo.ConfigParam{}
	bs := ConfigParam{
		DataId: "test-data-2",
		//Group:   "test-group",
		Content: "发布测试数据",
	}
	fmt.Println(as)
	assignParam(&as, &bs)
	fmt.Println(as)
}
