# Standard Go Project Layout

This is a basic layout for Go application projects. It's not an official standard defined by the core Go dev team; however, it is a set of common historical and emerging project layout patterns in the Go ecosystem. Some of these patterns are more popular than others. It also has a number of small enhancements along with several supporting directories common to any large enough real world application.

If you are trying to learn Go or if you are building a PoC or a toy project for yourself this project layout is an overkill. Start with something really simple (a single `main.go` file is more than enough). As your project grows keep in mind that it'll be important to make sure your code is well structured otherwise you'll end up with a messy code with lots of hidden dependencies and global state. When you have more people working on the project you'll need even more structure. That's when it's important to introduce a common way to manage packages/libraries. When you have an open source project or when you know other projects import the code from your project repository that's when it's important to have private (aka `internal`) packages and code. Clone the repository, keep what you need and delete everything else! Just because it's there it doesn't mean you have to use it all. None of these patterns are used in every single project. Even the `vendor` pattern is not universal.

This project layout is intentionally generic and it doesn't try to impose a specific Go package structure.

This is a community effort. Open an issue if you see a new pattern or if you think one of the existing patterns needs to be updated.

If you need help with naming, formatting and style start by running [`gofmt`](https://golang.org/cmd/gofmt/) and [`golint`](https://github.com/golang/lint). Also make sure to read these Go code style guidelines and recommendations:
* https://talks.golang.org/2014/names.slide
* https://golang.org/doc/effective_go.html#names
* https://blog.golang.org/package-names
* https://github.com/golang/go/wiki/CodeReviewComments
* [Style guideline for Go packages](https://rakyll.org/style-packages) (rakyll/JBD)

See [`Go Project Layout`](https://medium.com/golang-learn/go-project-layout-e5213cdcfaa2) for additional background information.

More about naming and organizing packages as well as other code structure recommendations:
* [GopherCon EU 2018: Peter Bourgon - Best Practices for Industrial Programming](https://www.youtube.com/watch?v=PTE4VJIdHPg)
* [GopherCon Russia 2018: Ashley McNamara + Brian Ketelsen - Go best practices.](https://www.youtube.com/watch?v=MzTcsI6tn-0)
* [GopherCon 2017: Edward Muller - Go Anti-Patterns](https://www.youtube.com/watch?v=ltqV6pDKZD8)
* [GopherCon 2018: Kat Zien - How Do You Structure Your Go Apps](https://www.youtube.com/watch?v=oL6JBUk6tj0)

## logcollect

### 小项目:将多个机台的日志数据包统一存放在服务器数据库中

### 机台数据 基本都是.zip压缩包

### 服务器 sql server 2012


## Go Directories

### /cmd
该项目的主要应用。

每个应用程序的目录名称应与您想要的可执行文件的名称相匹配（例如/cmd/myapp）。

不要在应用程序目录中放入大量代码。如果您认为代码可以导入并在其他项目中使用，那么它应该存在于/pkg目录中。 如果代码不可重用或者您不希望其他人重用它，请将该代码放在/internal目录中。你会惊讶于别人会做什么，所以要明确你的意图！

通常有一个小main函数可以从/internal和/pkg目录中导入和调用代码，而不是其他任何东西。

请参阅/cmd目录以获取示例。

### /internal
私有应用程序和库代码。这是您不希望其他人在其应用程序或库中导入的代码。

将您的实际应用程序代码放在/internal/app目录（例如/internal/app/myapp）和/internal/pkg目录中这些应用程序共享的代码（例如/internal/pkg/myprivlib）。

### /pkg
可以由外部应用程序使用的库代码（例如/pkg/mypubliclib）。其他项目将导入这些库，期望它们可以工作，所以在你把东西放在这里之前要三思而后行:-)

当你的根目录包含许多非Go组件和目录时，它也可以在一个地方将Go代码分组，从而更容易运行各种Go工具（如Best Practices for Industrial ProgrammingGopherCon EU 2018中所述）。

/pkg如果您想查看哪个热门的Go repos使用此项目布局模式，请查看该目录。这是一种常见的布局模式，但它并未被普遍接受，Go社区中的一些人不推荐它。

### /vendor
应用程序依赖项（手动管理或由您喜欢的依赖管理工具管理dep）。

如果要构建库，请不要提交应用程序依赖项。

## 服务应用程序目录
### /api
OpenAPI/Swagger规范，JSON模式文件，协议定义文件。

请参阅/api目录以获取示例。

## Web应用程序目录
### /web
特定于Web应用程序的组件：静态Web资产，服务器端模板和SPA。

## 常见应用程序目录
### /configs
配置文件模板或默认配置。

将您的confd或consul-template模板文件放在这里。

### /init
系统初始化（systemd，upstart，sysv）和进程管理器/主管（runit，supervisord）配置。

### /scripts
脚本执行各种构建，安装，分析等操作。

这些脚本使根级Makefile保持简洁（例如https://github.com/hashicorp/terraform/blob/master/Makefile）。

请参阅/scripts目录以获取示例。

### /build
包装和持续集成。

将您的云（AMI），容器（Docker），OS（deb，rpm，pkg）包配置和脚本放在/build/package目录中。

将CI（travis，circle，drone）配置和脚本放在/build/ci目录中。请注意，某些CI工具（例如，Travis CI）对其配置文件的位置非常挑剔。尝试将配置文件放在/build/ci将它们链接到CI工具所期望的位置的目录中（如果可能）。

### /deployments
IaaS，PaaS，系统和容器编排部署配置和模板（docker-compose，kubernetes / helm，mesos，terraform，bosh）。

### /test
其他外部测试应用和测试数据。您可以随意构建/test目录。对于更大的项目，有一个数据子目录是有意义的。例如，您可以拥有/test/data或者/test/testdata如果需要Go来忽略该目录中的内容。请注意，Go也会忽略以“。”开头的目录或文件。或“_”，因此您在命名测试数据目录方面具有更大的灵活性。

请参阅/test目录以获取示例。

## 其他目录
### /docs
设计和用户文档（除了你的godoc生成的文档）。

请参阅/docs目录以获取示例。

### /tools
该项目的支持工具。请注意，这些工具可以从/pkg和/internal目录中导入代码。

请参阅/tools目录以获取示例。

### /examples
应用程序和/或公共库的示例。

请参阅/examples目录以获取示例。

### /third_party
外部帮助工具，分叉代码和其他第三方实用程序（例如，Swagger UI）。

### /githooks
Git钩子。

### /assets
与您的存储库一起使用的其他资产（图像，徽标等）。

### /website
如果您不使用Github页面，这是放置项目的网站数据的地方。


## Badges

* [Go Report Card](https://goreportcard.com/) - It will scan your code with `gofmt`, `go vet`, `gocyclo`, `golint`, `ineffassign`, `license` and `misspell`. Replace `github.com/golang-standards/project-layout` with your project reference.

* [GoDoc](http://godoc.org) - It will provide online version of your GoDoc generated documentation. Change the link to point to your project.

* Release - It will show the latest release number for your project. Change the github link to point to your project.

[![Go Report Card](https://goreportcard.com/badge/github.com/golang-standards/project-layout?style=flat-square)](https://goreportcard.com/report/github.com/golang-standards/project-layout)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/golang-standards/project-layout)
[![Release](https://img.shields.io/github/release/golang-standards/project-layout.svg?style=flat-square)](https://github.com/golang-standards/project-layout/releases/latest)

## Notes

A more opinionated project template with sample/reusable configs, scripts and code is a WIP.

