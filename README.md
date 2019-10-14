# fabric-demo
## 测试流程
1. 进入fabric目录
2. 执行start.sh文件，完成fabric网络部署
```
./start.sh gm/sw 
默认sw算法
```
3. 进入fabric-demo 打开main_test.go源文件或者执行go test命令，根据流程执行测试方法
note：需要根据选择的加密算法GM/SW修改config/config_e2e.yaml文档
``` 
创建通道
go test -timeout 6000s -run TestCreateChannel -v -count=1
加入通道
go test -timeout 6000s -run TestJoinChannel -v -count=1
安装链码
go test -timeout 6000s -run TestInstallChainCode -v -count=1
实例化链码
go test -timeout 6000s -run TestInstantiateChainCode -v -count=1
查询链码
go test -timeout 6000s -run TestQueryChainCode -v -count=1
调用链码
go test -timeout 6000s -run TestInvokeChainCode -v -count=1
```
