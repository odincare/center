package center

import (
	"fmt"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	newClientTest()
}

func newClientTest() {
	client := NewOdinCenter()
	getConfig(client)
	listenConfig(client)
	time.Sleep(time.Second * 2)
	testPublishConfig(client)
	//deleteConfig(client)

	time.Sleep(time.Second * 30)
	cancelListen(client)
	time.Sleep(time.Second * 200)
}

func testPublishConfig(client *OdinCenter) {
	fmt.Println("---##发布配置###---")
	client.PublishConfig(ConfigParam{
		DataId: "test-data-2",
		//Group:   "test-group",
		Content: "发布测试数据",
	})
}

func getConfig(client *OdinCenter) {
	fmt.Println("---##读取配置###---")
	res := client.GetConfig(ConfigParam{DataId: "odin"})
	fmt.Println("---#####---", res)
}
func listenConfig(client *OdinCenter) {
	fmt.Println("---##监听配置###---")

	client.ListenConfig(ConfigParam{
		DataId: "test-data-2",
		Group:  "test-group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("config changed dataId:" + dataId + ", content:" + data)
		},
	})
}
func cancelListen(client *OdinCenter) {
	fmt.Println("---##取消监听配置###---")

	client.CancelListenConfig(ConfigParam{
		DataId: "tasks",
	})
}
func deleteConfig(client *OdinCenter) {
	fmt.Println("---##删除配置###---")

	client.DeleteConfig(ConfigParam{
		DataId: "test-data-2",
		Group:  "test-group",
	})
}
