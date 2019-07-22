# `/pkg`

可以由外部应用程序使用的库代码（例如/pkg/mypubliclib）。其他项目将导入这些库，期望它们可以工作，所以在你把东西放在这里之前要三思而后行:-)

当你的根目录包含许多非Go组件和目录时，它也可以在一个地方将Go代码分组，从而更容易运行各种Go工具（如Best Practices for Industrial ProgrammingGopherCon EU 2018中所述）。

/pkg如果您想查看哪个热门的Go repos使用此项目布局模式，请查看该目录。这是一种常见的布局模式，但它并未被普遍接受，Go社区中的一些人不推荐它。