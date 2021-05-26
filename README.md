# istiiomcp

istio mcp 示例

# meshconfig

istio 所需的配置文件

```
cat meshconfig
configSources:
  - address: xds://192.168.11.94:9901
```

# 运行

源码编译本地运行

```
ENABLE_CA_SERVER=false PILOT_ENABLE_ANALYSIS=false PILOT_ENABLED_SERVICE_APIS=false ./pilot-discovery discovery --meshConfig ./meshconfig  --registries ""
```


# 检查是否生效

```
curl 127.0.0.1:8080/debug/configz
```
