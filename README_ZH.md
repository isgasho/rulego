# RuleGo

--------
[English](README.md)| 中文

<img src="doc/imgs/logo.png" width="100">  

`RuleGo`是一个基于`Go`语言的轻量级、高性能、嵌入式的规则引擎。也一个灵活配置和高度定制化的事件处理框架。可以对输入消息进行过滤、转换、丰富和执行各种动作。

本项目很大程度受[thingsboard](https://github.com/thingsboard/thingsboard) 启发。参考其规则链思想，但在架构上做了重大调整，来满足以下场景:
* 在资源占用和性能上做了很大程度优化，使其更适用边缘计算场景。
* 不停机、不需要重新编译，动态编排业务，满足高度定制化和高度变化的业务需求。
* 非侵入式嵌入到现有项目。
* 提供更加灵活的接口和回调钩子。
* 更加开放的组件生态。可以利用社区提供的组件或者把业务封装成组件，通过搭积木方式快速动态实现业务需求。

## 特性

--------

* 开发语言：Go 1.8
* 轻量级：无外部中间件依赖，在低成本设备中也能高效对数据进行处理和联动，适用于物联网边缘计算。
* 高性能：得益于`Go`的高性能特性，另外`RuleGo`采用协程池和对象池等技术。对10W条数据进行`JS脚本过滤->JS脚本数据处理->HTTP推送` 处理,平均用时9秒。
* 嵌入式：支持把`RuleGo`嵌入到现有项目，非入侵式利用其特性。
* 组件化：所有业务逻辑都是组件，并能灵活配置和重用它们。
* 规则链：可以灵活地组合和重用不同的组件，实现高度定制化和可扩展性的业务流程。
* 流程编排：支持对规则链进行动态编排，你可以把业务地封装成`RuleGo`组件，然后通过搭积木方式实现你高度变化的业务需求。
* 扩展简单：提供丰富灵活的扩展接口和钩子，如：自定义组件、组件注册管理、规则链DSL解析器、协程池、规则节点消息流入/流出回调、规则链处理结束回调。
* 动态加载：支持通过`Go plugin` 动态加载组件和扩展组件。
* 内置常用组件：`消息类型Switch`,`JavaScript Switch`,`JavaScript过滤器`,`JavaScript转换器`,`HTTP推送`，`MQTT推送`，`发送邮件`，`日志记录`
  等组件。可以自行扩展其他组件。
* 上下文隔离机制：可靠的上下文隔离机制，无需担心高并发情况下的数据串流。

## 安装

--------
使用`go get`命令安装`RuleGo`：

```bash
go get github.com/parki/rulego
```

## 使用

--------
使用Json格式定义规则链DSL：      

以下例子定义3个规则节点，规则链逻辑如下图：（更多例子参考[testcases/](testcases)）
![img_1.png](doc/imgs/rulechain/img_1.png)

```json
{
  "ruleChain": {
    "name": "测试规则链",
    "root": true
  },
  "metadata": {
    "nodes": [
      {
        "id": "s1",
        "type": "jsFilter",
        "name": "过滤",
        "debugMode": true,
        "configuration": {
          "jsScript": "return msg!='bb';"
        }
      },
      {
        "id": "s2",
        "type": "jsTransform",
        "name": "转换",
        "debugMode": true,
        "configuration": {
          "jsScript": "metadata['test']='test02';\n metadata['index']=50;\n msgType='TEST_MSG_TYPE2';\n var msg2=JSON.parse(msg);\n msg2['aa']=66;\n return {'msg':msg2,'metadata':metadata,'msgType':msgType};"
        }
      },
      {
        "id": "s3",
        "type": "restApiCall",
        "name": "推送数据",
        "debugMode": true,
        "configuration": {
          "restEndpointUrlPattern": "http://192.168.216.21:9099/api/socket/msg",
          "requestMethod": "POST",
          "maxParallelRequestsCount": 200
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "s2",
        "type": "True"
      },
      {
        "fromId": "s2",
        "toId": "s3",
        "type": "Success"
      }
    ],
    "ruleChainConnections": null
  }
}
```

字段说明：

- `ruleChain`: 规则链定义的根对象，包含以下字段：
  - `name`: 规则链的名称，可以是任意字符串。
  - `root`: 一个布尔值，表示这个规则链是根规则链还是子规则链。每个规则引擎实例只允许有一个根规则链。
- `metadata`: 一个对象，包含了规则链中节点和连接的信息，有以下字段：
  - `nodes`: 一个对象数组，每个对象代表规则链中的一个规则节点。每个节点对象有以下字段：
    - `id`: 节点的唯一标识符，可以是任意字符串。
    - `type`: 节点的类型，决定了节点的逻辑和行为。它应该与规则引擎中注册的节点类型之一匹配。
    - `name`: 节点的名称，可以是任意字符串。
    - `debugMode`: 一个布尔值，表示这个节点是否处于调试模式。如果为真，当节点处理消息时，会触发调试回调函数。
    - `configuration`: 一个对象，包含了节点的配置参数，具体内容取决于节点类型。例如，一个JS过滤器节点可能有一个`jsScript`字段，定义了过滤逻辑，而一个REST API调用节点可能有一个`restEndpointUrlPattern`字段，定义了要调用的URL。
  - `connections`: 一个对象数组，每个对象代表规则链中两个节点之间的连接。每个连接对象有以下字段：
    - `fromId`: 连接的源节点的id，应该与nodes数组中的某个节点id匹配。
    - `toId`: 连接的目标节点的id，应该与nodes数组中的某个节点id匹配。
    - `type`: 连接的类型，决定了什么时候以及如何把消息从一个节点发送到另一个节点。它应该与源节点类型支持的连接类型之一匹配。例如，一个JS过滤器节点可能支持两种连接类型："True"和"False"，表示消息是否通过或者失败过滤条件。
  - `ruleChainConnections`: 一个对象数组，每个对象代表规则链中一个节点和一个子规则链之间的连接。每个规则链连接对象有以下字段：
    - `fromId`: 连接的源节点的id，应该与nodes数组中的某个节点id匹配。
    - `toId`: 连接的目标子规则链的id，应该与规则引擎中注册的子规则链之一匹配。
    - `type`: 连接的类型，决定了什么时候以及如何把消息从一个节点发送到另一个节点。它应该与源节点类型支持的连接类型之一匹配。
    

导入`RuleGo`包并创建一个规则引擎实例：

```go
import "github.com/yourname/rulego"

//创建一个规则引擎实例，每个规则引擎实例有且只有一个根规则链
ruleEngine, err := rulego.New("rule01", []byte(ruleFile))
```

把消息、消息类型、消息元数据交给规则引擎实例处理：

```go
//定义消息元数据
metaData := types.NewMetadata()
metaData.PutValue("productType", "test01")
//定义消息和消息类型
msg := types.NewMsg(0, "TELEMETRY_MSG", types.JSON, metaData, "{\"temperature\":35}")

//把消息交给规则引擎处理
ruleEngine.OnMsg(msg)

//需要得到结束回调方式
ruleEngine.OnMsgWithEndFunc(msg, func (msg types.RuleMsg, err error) {
//规则链异步回调结果 
//注意：规则链如果有多个分支结束点，会调用多次
})
```

添加子规则链：

```go
//规则引擎实例化时和子规则链一起创建
ruleEngine, err := rulego.New("rule01", []byte(ruleFile), rulego.WithAddSubChain("rule_chain_test", subRuleFile))
//或者通过创建或者更新的方式
ruleEngine.ReloadChild(types.EmptyRuleNodeId, types.RuleNodeId{Id: "rule_chain_test", Type: types.CHAIN}, subRuleFile)

```

更新规则链

```go
//更新根规则链
err := ruleEngine.ReloadSelf([]byte(ruleFile))
//更新子规则链
ruleEngine.ReloadChild(types.EmptyRuleNodeId, types.RuleNodeId{Id: "rule_chain_test", Type: types.CHAIN}, subRuleFile)
//更新规则链某个节点,详情看以下方法
ruleEngine.ReloadChild(chainId types.RuleNodeId, ruleNodeId types.RuleNodeId, dls []byte)

```

规则引擎实例管理：

```go
//通过ID获取已经创建的规则引擎实例
ruleEngine, ok := rulego.Get("rule01")
//删除已经创建的规则引擎实例
rulego.Del("rule01")
```

### 配置

详见`types.Config`

```go
//创建一个默认的配置
config := rulego.NewConfig()
//调试节点回调，节点配置必须配置debugMode:true 才会触发调用
//节点入和出信息都会调用该回调函数
config.OnDebug = func (flowType string, nodeId string, msg types.RuleMsg, relationType string, err error) {
}
//全局的规则链结束回调
//如果只是想要某一次消息调用，使用ruleEngine.OnMsgWithEndFunc方式
//注意：规则链如果有多个分支结束点，会调用多次
config.OnEnd = func (msg types.RuleMsg, err error) {
}
//配置使用
ruleEngine, err := rulego.New("rule01", []byte(ruleFile), rulego.WithConfig(config))
```

## 关于规则链

--------

### 规则节点

规则节点是规则引擎的基本组件，它一次处理单个传入消息并生成一个或多个传出消息。规则节点是规则引擎的主要逻辑单元。规则节点可以过滤，丰富，转换传入消息，执行操作或与外部系统通信。 你可以把业务很方便地封装成`RuleGo`
节点组件，然后灵活配置和复用它们，像搭积木一样实现你的业务需求。 定义`rulego`自定义节点组件方式：

* 方式一：实现`types.Node` 接口，参考[components](components)例子 例如：

```go
//定义Node组件
//UpperNode A plugin that converts the message data to uppercase
type UpperNode struct{}

func (n *UpperNode) Type() string {
return "test/upper"
}
func (n *UpperNode) New() types.Node {
return &UpperNode{}
}
func (n *UpperNode) Init(ruleConfig types.Config, configuration types.Configuration) error {
// Do some initialization work
return nil
}
//处理消息
func (n *UpperNode) OnMsg(ctx types.RuleContext, msg types.RuleMsg) error {
msg.Data = strings.ToUpper(msg.Data)
// Send the modified message to the next node
ctx.TellSuccess(msg)
return nil
}

func (n *UpperNode) Destroy() {
// Do some cleanup work
}
//注册到rulego默认注册器
rulego.Registry.Register(&MyNode{})
```

* 方式二：使用`go plugin` 实现接口 `types.PluginRegistry`接口。并导出变量名称:`Plugins`,参考[testcases/plugin](testcases/plugin/)
  例如：

```go
// plugin entry point
var Plugins MyPlugins

type MyPlugins struct{}

func (p *MyPlugins) Init() error {
return nil
}
func (p *MyPlugins) Components() []types.Node {
//一个插件可以提供多个组件
return []types.Node{&UpperNode{}, &TimeNode{}, &FilterNode{}}
}
//go build -buildmode=plugin -o plugin.so plugin.go # 编译插件，生成plugin.so文件，需要在mac或者linux环境下编译
//注册到rulego默认注册器
rulego.Registry.RegisterPlugin("test", "./plugin.so")
```

然后在规则链DSL文件使用您的组件

```json
{
  "ruleChain": {
    "name": "测试规则链",
    "root": true,
    "debugMode": false
  },
  "metadata": {
    "nodes": [
      {
        "id": "s1",
        "type": "test/upper",
        "name": "名称",
        "debugMode": true,
        "configuration": {
          "field1": "组件定义的配置参数",
          "....": "..."
        }
      }
    ],
    "connections": [
      {
        "fromId": "s1",
        "toId": "连接下一个组件ID",
        "type": "与组件的连接关系"
      }
    ],
    "ruleChainConnections": null
  }
}
```

`rulego`内置部分常用组件分为以下几种：

* 过滤组件：对消息进行过滤。
* 转换组件：对消息进行转换和增强。
* 动作组件：执行某些动作，或者和外部系统联动。

### 规则链

规则链是`规则节点`及其`关系`的逻辑组。接收来自节点的出站消息将其通过指定`关系`发送至下一个节点。以下是一些常用的规则链例子：

* 顺序执行：
  ![img_1.png](doc/imgs/rulechain/img_1.png)
* 异步+顺序执行
  ![img_3.png](doc/imgs/rulechain/img_2.png)
* 使用子规则链方式
  ![img_4.png](doc/imgs/rulechain/img_3.png)
* 一些复杂例子：
  ![img_5.png](doc/imgs/rulechain/img_4.png)

## 性能

--------
`rulego` 几乎不会额外增加系统开销，资源占用极低，因为使用了对象协程池和对象池，甚至比直接调用业务的方式性能还高，特别适合在边缘服务器运行。

--------
机器：树莓派2(900MHz Cortex-A7*4,1GB LPDDR2)  
数据大小：260B   
规则链：JS脚本过滤->JS复杂转换->HTTP推送   
测试结果：100并发和500并发，内存占用变化不大都在19M左右

## 贡献

--------

欢迎任何形式的贡献，包括提交问题、建议、文档、测试或代码。请遵循以下步骤：

* 克隆项目仓库到本地
* 创建一个新的分支并进行修改
* 提交一个合并请求到主分支
* 等待审核和反馈

## 许可

--------
`RuleGo`使用Apache 2.0许可证，详情请参见[LICENSE](LICENSE)文件。