## go-easy-scada 工业控制数据采集器

- [开发日志](develop-log.md)

```
    这是一个 用go语言编写, 赋能于 工业数据采集、转换以及分发、接驳其他系统 的集成框架。
    这个项目还在不断迭代中。
```
### 项目初衷
```
    和不同工厂的各种需求周旋后，发现安装到工厂的传统软件，
    安装、调试存在困难。（工厂一般不开放外网）
    过度集成化，迭代困难。（程序需要定制或者修复BUG时，打包文件巨大）
    过于灵活，安装困难。（需要专业人员安装部署）
    可替换/定制的基础组件、低学习(人力)成本维护(实习生也能解决业务逻辑) 成为这个框架的追求的目标
```

### 目前架构功能
- 配置注入
  - 可以通过外部配置文件,启动时会加载 定义的插件 和 消息的发布者、订阅者
    ```
    {
        "plugin": {
            "inPlugin": {
                "module": "mqinput",
                "param": "addr=nats://127.0.0.1:4222;topic=ESD_CLIENT"
            },
            "outPlugin": {
                "module": "mqinput",
                "param": "addr=nats://127.0.0.1:4222;topic=ESD_SERVER"
            }
        },
        "message": {
            "MESSAGE_SEND_TO_MONITOR": {
                "pub": ["inPlugin"],
                "sub": ["outPlugin"],
            }
        }
    }
    ```
    引擎初始化时，会根据配置文件中的`plugin`初始化`module`指定的`adapter`插件.(在项目目录adapter/transmit)
    
    当工厂有定制需求时,可以根据各种的情况,将通讯插件插在不同的事件上.
    
    如客户需要设备原生数据,或者是需要处理后的数据.

    自己使用的话 可以将数据采集器和管理系统后端进行解耦 
    
    如后端对数据库进行增删改后，可定义一个消息，以便数据采集器重新更新缓存.
  - 插件架构如下
    ```
                业务逻辑代码   (业务逻辑代码,直接与引擎抽象层和引擎实例相关,将基础组件和业务逻辑解耦)
                    ||
                 引擎抽象层    (引擎[程序]生命期时，对功能要求的接口定义)
                    ||
                 组件抽象层    (组件抽象层，以便第三方组件或者自己写的基础组件 突然灵光一现)
        //          ||           \
     基础组件      基础组件      基础组件
    ```


### 项目结构
```
custom                           // 业务代码
core                             // 最底层 采用单例模式 能被所有组件公用的核心库（如日志、性能监控等）
adapter                          // 符合适配器的 公用组件（orm、cache等）
 - componentName                 // 抽象型组件类（如Cache、Orm等等）
    - thirdPartyPluginName       // 实际组件(第三方库,如Redis、Gorm等)
       - plguin.go               // 组件实现方式(init函数里向componentType.Register("w",newObject) 注册)
    - componentName.go           // 定义IComponent接口 提供了Register New(IComponent) 方法,供业务代码或上层代码使用
 - adapter.go                    // 定义适配器的基本树形,(可定义父类接口，如初始化、卸载等)
engine                           // 单例引擎 (目前设计上可能和自身的业务耦合,以后尽量优化成通用，即插即用)
main 
utilio                           // 常用的工具函数
```

```
    实际在业务代码使用基础组件时.可 componentName.New() 创建默认的实例对象.并用该实例进行操作.
    在使用时 需要 import _ "adapter/componentName/thirdPartyPluginName" 执行实例init 进行初始化.
```

### 消息定义
#### 从设备获取的消息逻辑上分为两部分，一是处理前，二是处理后
* `MESSAGE_DEAL_BEFORE_xxx` 处理前，属于原始数据.(如客户要求插入自己逻辑的数据)
* `MESSAGE_DEAL_AFTER_xxx` 处理后，属于处理后数据.(如携带数据库的数据)

#### 这样的话,adapter事件接口适配器 可以执行单一职责的功能，如(发布者/订阅者).
#### 为了提高性能以及不必要的创建(如tcp池等),在adapter插件功能内部实现 在包装一层,采用单例模式.



程序会分发两个版本.
* 一个是没有`Custom`文件夹的版本,不会触发 `MESSAGE_DEAL_AFTER_xxx` 事件,让客户负责采集到他们的系统.
* 一个是有`Custom`文件夹的版本，可能会触发 `MESSAGE_DEAL_AFTER_xxx` 事件,Custom代码中可以进行拦截并传递，进行下一步的定制.
    可用于 客户需要公司制作的系统，并希望将数据转发一份到客户公司数据库的情况.

