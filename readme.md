## ESD-ROUTER-PREVIEW ESD消息接驳器

基础组件


### 项目结构
```
custom
 - 
core // 最底层 采用单例模式 能被所有组件公用的核心库（如日志、性能监控等）
 - 
adapter // 符合适配器的 公用组件（orm、cache等）
 - componentName // 抽象型组件类（如Cache、Orm等等）
    - thirdPartyPluginName // 实际组件(第三方库,如Redis、Gorm等)
       - plguin.go // 组件实现方式(init函数里向componentType.Register("w",newObject) 注册)
    - componentName.go // 定义IComponent接口 提供了Register New(IComponent) 方法,供业务代码或上层代码使用
 - adapter.go // 定义适配器的基本树形,(可定义父类接口，如初始化、卸载等)
```

```
    实际在业务代码使用基础组件时.可 componentName.New() 创建对应功能对象.并用该实例进行操作.
    在使用时 需要 import _ "adapter/componentName/thirdPartyPluginName" 执行实例init 进行初始化.
```

### 消息定义
#### 从设备获取的消息逻辑上分为两部分，一是处理前，二是处理后
* `MESSAGE_DEAL_BEFORE_xxx` 处理前，属于原始数据.(如客户要求插入自己逻辑的数据)
* `MESSAGE_DEAL_AFTER_xxx` 处理后，属于处理后数据.(如携带数据库验证的数据)

#### 这样的话,adapter事件接口适配器 可以执行单一职责的功能，如(发送/接收).
#### 为了提高性能以及不必要的创建(如tcp池等),在adapter插件的具体实现的内部采用单例模式.



程序会分发两个版本.
* 一个是没有`Custom`文件夹的版本,不会触发 `MESSAGE_DEAL_AFTER_xxx` 事件,让客户负责采集到他们的系统.
* 一个是有`Custom`文件夹的版本，可能会触发 `MESSAGE_DEAL_AFTER_xxx` 事件,客户可以采取，进行下一步的定制.
    可用于 客户需要系统，并希望将数据转发一份到客户公司数据库的情况.

