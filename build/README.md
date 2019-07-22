# `/build`

包装和持续集成。

将您的云（AMI），容器（Docker），OS（deb，rpm，pkg）包配置和脚本放在/build/package目录中。

将CI（travis，circle，drone）配置和脚本放在/build/ci目录中。请注意，某些CI工具（例如，Travis CI）对其配置文件的位置非常挑剔。尝试将配置文件放在/build/ci将它们链接到CI工具所期望的位置的目录中（如果可能）。