package center

import (
	"encoding/json"
	"errors"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/odincare/center/options"
	"gopkg.in/yaml.v2"
	"os"
	"reflect"
)

type OdinCenter struct {
	cClient config_client.IConfigClient
	rClient naming_client.INamingClient
}

type ConfigParam struct {
	DataId   string `param:"dataId"`  //required
	Group    string `param:"group"`   //required
	Content  string `param:"content"` //required
	DatumId  string `param:"datumId"`
	OnChange func(namespace, group, dataId, data string)
}

func NewOdinCenter(opts ...options.Option) *OdinCenter {
	if len(opts) == 0 {
		var ConfigAddr, NamespaceID, LogLevel string
		ConfigAddr = os.Getenv("CONF_ADDR")
		NamespaceID = os.Getenv("NAMESPACE_ID")
		LogLevel = os.Getenv("LOG_LEVEL")
		if ConfigAddr == "" || NamespaceID == "" {
			panic("初始化配置错误")
			return nil
		}
		opts = append(opts, options.Addrs(ConfigAddr))
		opts = append(opts, options.NameSpaceId(NamespaceID))
		opts = append(opts, options.LogLevel(LogLevel))
	}
	conf := options.OptionInit(opts...)
	return &OdinCenter{
		cClient: initNacosConfigCenter(conf),
		rClient: initNacosRegistryCenter(conf),
	}
}

// PublishConfig /**----配置相关---开始--**/
func (o *OdinCenter) PublishConfig(param ConfigParam) error {
	sParam := vo.ConfigParam{}
	assignParam(&sParam, &param)

	_, err := o.cClient.PublishConfig(sParam)
	if err != nil {
		logger.Error("发布配置失败,错误信息:", err)
		return err
	}
	return nil
}

//获取配置
//func (o *OdinCenter) GetConfig(param ConfigParam) string {
//	sParam := vo.ConfigParam{}
//	assignParam(&sParam, &param)
//
//	result, err := o.cClient.GetConfig(sParam)
//	if err != nil {
//		logger.Error("获取配置失败,错误信息:", err)
//		return ""
//	}
//	return result
//}
func (o *OdinCenter) GetConfig(param ConfigParam, obj interface{}) (err error) {
	sParam := vo.ConfigParam{}
	assignParam(&sParam, &param)

	result, err := o.cClient.GetConfig(sParam)
	if err != nil {
		logger.Error()
		return errors.New("获取配置失败,错误信息:" + err.Error())
	}

	err = yaml.Unmarshal([]byte(result), obj)
	if err != nil {
		return errors.New("解析配置失败,错误信息:" + err.Error())
	}

	// 获取入参的类型
	t := reflect.TypeOf(obj)

	// 入参类型校验
	if t.Kind() != reflect.Ptr || t.Elem().Kind() != reflect.Struct {
		return errors.New("参数应该为结构体指针")
	}
	v := reflect.ValueOf(obj).Elem()

	err = checkRequired(v)

	return err
}
func checkRequired(v reflect.Value) error {

	for i := 0; i < v.NumField(); i++ {

		fieldInfo := v.Type().Field(i)
		label := fieldInfo.Tag.Get("required")
		name := fieldInfo.Name
		sub := v.FieldByName(name)

		if sub.Kind() != reflect.Struct && sub.String() == "" && label == "true" {
			return errors.New("配置参数有误:" + name + " is undefined or empty")
		}

		if sub.Kind() == reflect.Struct {
			err := checkRequired(sub)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

//监听配置
func (o *OdinCenter) ListenConfig(param ConfigParam) error {
	sParam := vo.ConfigParam{}
	assignParam(&sParam, &param)

	err := o.cClient.ListenConfig(sParam)
	if err != nil {
		logger.Error("监听配置失败,错误信息:", err)
		return err
	}
	return nil
}

//取消配置监听
func (o *OdinCenter) CancelListenConfig(param ConfigParam) error {
	sParam := vo.ConfigParam{}
	assignParam(&sParam, &param)

	err := o.cClient.CancelListenConfig(sParam)
	if err != nil {
		logger.Error("获取配置失败,错误信息:", err)
		return err
	}
	return nil
}

//删除配置
func (o *OdinCenter) DeleteConfig(param ConfigParam) error {
	sParam := vo.ConfigParam{}
	assignParam(&sParam, &param)

	_, err := o.cClient.DeleteConfig(sParam)
	if err != nil {
		logger.Error("获取配置失败,错误信息:", err)
		return err
	}
	return nil
}

//----配置相关---结束--

//初始化配置

func initNacosConfigCenter(opt vo.NacosClientParam) config_client.IConfigClient {

	cClient, err := clients.NewConfigClient(opt)
	if err != nil {
		logger.Error("启动nacos配置客户端失败,错误信息:", err)
		panic(err)
	}
	return cClient
}
func initNacosRegistryCenter(opt vo.NacosClientParam) naming_client.INamingClient {
	rClient, err := clients.NewNamingClient(opt)
	if err != nil {
		logger.Error("启动nacos服务发现客户端失败,错误信息:", err)
		panic(err)
	}
	return rClient
}

//参数转换辅助方法通过json
func assignResult(toValue interface{}, fromValue interface{}) {
	//fVal := reflect.ValueOf(fromValue).Elem()
	//if ok := fVal.FieldByName("Group").IsValid(); ok {
	//	gv := fVal.FieldByName("Group").Interface().(string)
	//	if gv == "" {
	//		fVal.FieldByName("Group").Set(reflect.ValueOf("DEFAULT_GROUP"))
	//	}
	//}
	d, _ := json.Marshal(fromValue)
	_ = json.Unmarshal(d, toValue)

}

//参数转换辅助方法通过反射
func assignParam(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem()
	vVal := reflect.ValueOf(value).Elem()
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			if vVal.Field(i).Kind() == reflect.Slice {
				//for si := 0; si < vVal.Field(i).Len(); si++ {
				//}
			} else {
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
			}
		}
		if name == "Group" {
			gv := vVal.Field(i).Interface().(string)
			if gv == "" {
				bVal.FieldByName(name).Set(reflect.ValueOf("DEFAULT_GROUP"))
			}
		}
	}
}
