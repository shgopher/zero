## `cmd` 目录介绍

`cmd` 目录用来存放应用的 main 文件。该目录下存放 3 类应用文件：

- 代码生成工具：以 `gen-` 开头；
- 应用程序：以 `zero-` 开头；
- 代码静态检查工具：以 `lint-`开头。

## 为什么要在一个 `cmd` 目录下单独存放一个 main 文件？

在 Go 项目中，通常会在项目根目录下创建一个名为 `cmd` 的目录，用于存放项目的可执行文件。在 `cmd` 目录下，又会为每个可执行文件创建一个单独的目录，并在该目录下创建一个名为 `main.go` 的文件，用于定义该可执行文件的入口函数 `main()`。

这样做的原因如下：
- 方便管理和维护项目代码：将不同的可执行文件分别放在不同的目录中，可以使代码结构更加清晰，便于查找和修改；
- 依赖关系更加明确：在每个目录中单独存放一个 `main.go` 文件，可以使得各个可执行文件之间的依赖关系更加明确，也方便了构建和部署；
- 方便构建：采用这种方式还可以方便地将不同的可执行文件分别编译成独立的二进制文件，以便将它们部署到不同的服务器上，或者将它们打包成容器镜像进行分发。

> 注意：这种方式并不是强制性的，对于一些小型应用程序或者是单文件应用程序，也可以直接将 `main()` 函数定义在项目根目录下的 `main.go` 文件中。

## 为什么建议 `main.go` 文件内容简单，把具体实现放在其他目录下？

在 Go 语言中，建议将 `cmd` 目录下的 `main.go` 文件保持逻辑简单的原因是为了方便代码的维护和扩展。

main.go 文件通常包含程序的入口函数，负责解析命令行参数、初始化配置等基础操作，具有很高的可复用性。

将具体的实现存放在其他目录下，可以使代码结构更加清晰，方便代码的组织和管理。此外，将具体的实现存放在其他目录下也有利于代码重用，可以避免代码的重复编写，提高代码的可维护性和可扩展性。

另外，这样做还可以使用 `go install cmd/zero-usercenter/xxx.go` 这种方式来安装软件。

一个反例：如果 `cmd/zero-usercenter/{xxx.go, yyy.go}`，`yyy.go` 中有不可导出的变量/函数，则 `go install` 会失败。


