# Gee

A [gin](https://github.com/gin-gonic/gin)-like web framework

## 技术点总结

最主要的技术点：**前缀树，中间件，分组控制**

### 顶层Engine

- 最顶层是一个Engine（是一种http.Handler，因为实现了ServeHTTP方法）。Engine在main里最后使用的Run(":9999")实际上是调用的是http.ListenAndServe(":9999", engine)

    > ServeHTTP()的逻辑如下：
    >
    > 对Engine里记录的所有分组进行判断，若传入的req.URL.Path里有该分组的前缀，则将该分组的中间件都放入当前ResponseWriter和Request组成的上下文Context中，并调用c.handle()开始执行中间件

### 路由分组RouterGroup

- Engine也是一个分组（即为RouterGroup），Engine记录了所有的RouterGroup。

- 每个RouterGroup结构里都有一个指向唯一Engine的指针。RouterGroup的所有操作其实都是调用其中的Engine来完成的。

- 比如addRoute()操作。其实group.addRoute("GET", "/date", HandlerFunc) 调用的是group.engine.router.addRoute("GET", "/date", HandlerFunc)。

- (Group).addRoute是在GET、POST注册路由时使用，对传入的路由片段part与Group的prefix拼接，形成完整路由，再交由group.engine.router调用真正的addRoute进行注册

### 上下文Context

- route是简历连接的。context是处理这次会话内容的：书写响应、执行中间件。

- 初始化Context时的输入是http.ResponseWriter和*http.Request。Context的作用是进行封装，提取出重要信息，防止写大量重复冗余的代码，比如：
  - http.ResponseWriter：响应response的消息头要设置状态码StatusCode和消息类型ContentType
  - *http.Request：提取出URL的path和Method（比如是GET？还是POST？）
  - 还写了快速构造String/Data/JSON/HTML响应的方法。这些方法的输入参数都是`(状态码，内容)`。JSON使用json.NewEncoder解析内容
- handlers：该访问都需要执行哪些中间件
- index记录当前处理到第几个中间件了

```go
func Next() {
    c.index++ // 重要! 因为并不是每个中间件都会在其结尾调用Next()
	for (若c.index没超过所有中间件; index++) {  // 即所有中间件还没执行完
		执行这第c.index个中间件
	}
}
```

### 前缀树的定义

#### tri.insert(完整path, 分割后的parts数组, 第几层index)

- 逐层递归
- 若递归了len(parts)层（即len(parts)==index），则函数return
- 对每行进行匹配，若没匹配到，则插入一个节点
- **注意！！只有path的最末端节点才设置`node.pattern` !**

#### tri.search(分割后的parts数组, 第几层index)

- 也是递归
- 若匹配到了len(parts)层或匹配到了*，但该节点没有设置`node.pattern`，则返回匹配失败
- 否则返回匹配的这个末尾节点

### 前缀树的使用

- 路由的 **注册** 和 **匹配** 是最重要的

- 维护了一个map(名为root)。比如map[GET]=一棵树, map[POST]=另一个树。
- **注册addRoute函数**：调用parsePattern解析传入的完整路由分割成小的parts数组--->在GET或POST对应的树中调用`tri.insert`插入节点
- **匹配getRoute函数**：调用parsePattern解析传入的完整路由分割成小的parts数组--->在GET或POST对应的树中调用`tri.search`寻找节点，找到后若有`:`或`*`，返回`:name`或`*filepath`的真实映射关系
