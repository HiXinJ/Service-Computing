# 设计博客网站REST API
REST是一种为分布式超媒体系统设计的架构风格，由Roy Fielding在 2000 年的博士论文中提出。Fielding也是HTTP的主要设计者。REST内涵丰富，很难用几句话将REST讲清楚，这里只将一个重要的概念：**超文本驱动**。  

2008 年，Leonard Richardson 提议对 Web API 使用以下[成熟度模型](https://martinfowler.com/articles/richardsonMaturityModel.html)：  
- 级别 0：定义一个 URI，所有操作是对此 URI 发出的 POST 请求。
- 级别 1：为各个资源单独创建 URI。
- 级别 2：使用 HTTP 方法来定义对资源执行的操作。
- 级别 3：使用超媒体（HATEOAS，如下所述）。

根据 Fielding 的定义，级别 3 对应于某个真正的 RESTful API。
“超文本驱动”又名“将超媒体作为应用状态的引擎”（Hypermedia As The Engine Of Application State，来自 Fielding 博士论文中的一句话，缩写为 HATEOAS）将 Web 应用看作是一个由很多状态（应用状态）组成的有限状态机。资源之间通过超链接相互关联，超链接既代表资源之间的关系，也代表可执行的状态迁移。在超媒体之中不仅仅包含数据，还包含了状态迁移的语义。以超媒体作为引擎，驱动 Web 应用的状态迁移。通过超媒体暴露出服务器所提供的资源，服务器提供了哪些资源是在运行时通过解析超媒体发现的，而不是事先定义的。从面向服务的角度看，超媒体定义了服务器所提供服务的协议。客户端应该依赖的是超媒体的状态迁移语义，而不应该对于是否存在某个 URI 或 URI 的某种特殊构造方式作出假设。一切都有可能变化，只有超媒体的状态迁移语义能够长期保持稳定。   
HTTP就是以REST为指导原则设计的，从这也可以看出为什么HTTP被设计为transfer protocol（转移协议）而不是transport protocol（传输协议）。传输协议是用来传输无语义的比特流，比如传输层的协议，HTTP转移协议强调的是状态的转移。   

## 博客API概述
接下来设计一个简单的博客网站的REST API。   
博客网站有用户，博客，评论。需要用户认证才可以发布、删除、更新、评论博客，查看博客不需要用户认证。  

## 用户操作

    GET https://api.blog.com/users
response body:
```json
200
[  
  {
    "login": "mojombo",
    "id": 1,
    "url": "https://api.blog.com/users/mojombo",
    "html_url": "https://blog.com/mojombo",
    "type": "User",
    "site_admin": false
  },
//   ...
//   ...
//   ...
]

404	
    not found

405	
    Validation exception
```  

我们可以通过某个用户的url字段获取该用户详细的信息：  

    GET https://api.blog.com/users/{username}
例如{username}=mojombo  
response：
```json
200
{
  "login": "mojombo",
  "id": 1,
  "node_id": "MDQ6VXNlcjE=",
  "url": "https://api.blog.com/users/mojombo",
  "html_url": "https://blog.com/mojombo",
  "type": "User",
  "site_admin": false,
  "name": "Tom Preston-Werner",
  "company": null,
  "blog": "http://tom.preston-werner.com",
  "location": "San Francisco",
  "email": null,
  "created_at": "2007-10-20T05:24:19Z",
  "updated_at": "2019-11-20T16:56:13Z"
}

400	
    Invalid name supplied

404	
    not found
``` 

## 博客操作  
获取某个用户下的所有博客

    GET https://api.blog.com/users/{username}/blogs  
例如{username}=mojombo  
返回的是一组blog对象，通过blog对象中的url字段，我们可以获取该博客的详细信息。
```json
[
  {
    "id": 26899533,
    "node_id": "MDEwOlJlcG9zaXRvcnkyNjg5OTUzMw==",
    "name": "30daysoflaptops.blog.io",
    "full_name": "mojombo/30daysoflaptops.blog.io",
    "private": false,
    "owner": {
        "login": "mojombo",
        "id": 1,
        "url": "https://api.blog.com/users/mojombo",
        "html_url": "https://blog.com/mojombo",
        "type": "User",
        "site_admin": false
    },
    "html_url": "https://blog.com/mojombo/30daysoflaptops",
    "description": null,
    "url": "https://api.blog.com/blogs/mojombo/30daysoflaptops",
  },

  {
    "id": 17358646,
    "node_id": "MDEwOlJlcG9zaXRvcnkxNzM1ODY0Ng==",
    "name": "Optimize-for-Happiness",
    "full_name": "mojombo/Optimize-for-Happiness",
    "private": false,
    "owner": {
        "login": "mojombo",
        "id": 1,
        "url": "https://api.blog.com/users/mojombo",
        "html_url": "https://blog.com/mojombo",
        "type": "User",
        "site_admin": false
    },
    "html_url": "https://blog.com/mojombo/Optimize-for-Happiness",
    "description": "Destroy your Atom editor, Optimize-for-Happiness style!",
    "url": "https://api.blog.com/blogs/mojombo/Optimize-for-Happiness",
  },
]
```

获取用户的某一篇博客：   

    GET https://api.blog.com/blogs/{username}/{blog}
例如：

    GET https://api.blog.com/blogs/mojombo/Optimize-for-Happiness


在更复杂的系统中，我们往往提供 URI（例如 /users/someuser/blogs/someblog/reviews），使客户端能够通过多个关系级别进行导航。 但是，如果资源之间的关系在将来更改，此级别的复杂性可能很难维护并且不够灵活。 相反，请尽量让 URI 相对简单。 应用程序获取对某个资源的引用后，应该可以使用此引用去查找与该资源相关的项目。 可将前面的查询替换为 URI /users/someuser/blogs 以查找用户 1 的所有博客，然后替换为 /blogs/someuser/someblog 以查找某用户的博客。虽然在这个简单的博客API设计中难以看到这种灵活性，但是将来的扩展时就会体现出灵活性来。 

创建新的博客：  

    POST https://api.blog.com/users/{username}/blogs  

```json
{
  "owner":{
        "login": "mojombo",
        "id": 1,
        "url": "https://api.blog.com/users/mojombo",
        "html_url": "https://blog.com/mojombo",
        "type": "User",
        "site_admin": false
    },
  "category": {
    "id": 0,
    "name": "string"
  },
  "title": "string",
  "content":"this is a blog.",
  "reviewer_urls": [
    null,
  ],
  "tags": [
    {
      "id": 0,
      "name": "string"
    }
  ],
  "private":false
}
```

如果 POST 方法创建了新资源，则会返回 HTTP 状态代码 201（已创建）。 新资源的 URI 包含在响应的 Location 标头中。 响应正文包含资源的表示形式。
如果该方法执行了一些处理但未创建新资源，则可以返回 HTTP 状态代码 200，并在响应正文中包含操作结果。 或者，如果没有可返回的结果，该方法可以返回 HTTP 状态代码 204（无内容）但不返回任何响应正文。
如果客户端将无效数据放入请求，服务器应返回 HTTP 状态代码 400（错误的请求）。 

删除一篇博客

    DELETE https://api.blog.com/blogs/{username}/{title}  

如果删除操作成功，Web 服务器应以 HTTP 状态代码 204 做出响应，指示已成功处理该过程，但响应正文不包含其他信息。 如果资源不存在，Web 服务器可以返回 HTTP 404（未找到）。

注意：GET POST DELETE操作之前都需要用户认证。

## 评论操作  

列出一篇博客的评论：

    GET https://blogs/{username}/{title}/reviews

例如列出用户mojombo的博客Optimize-for-Happiness下的评论
```json
[
  {
    "url": "https://api.blog.com/blogs/mojombo/Optimize-for-Happiness/reviews/68297443",
    "html_url": "https://blog.com/mojombo/Optimize-for-Happiness/reviews/68297443",
    "id": 68297443,
    "node_id": "MDEyOklzc3VlQ29tbWVudDY4Mjk3NDQz",
    "user": {
        "login": "mojombo",
        "id": 1,
        "url": "https://api.blog.com/users/mojombo",
        "html_url": "https://blog.com/mojombo",
        "type": "User",
        "site_admin": false
    },
    "created_at": "2014-12-29T20:07:11Z",
    "updated_at": "2014-12-29T20:07:11Z",
    "author_association": "NONE",
    "body": "It's an excellent article\n"
  },
//   ...
//   ...
//   ...
]
```