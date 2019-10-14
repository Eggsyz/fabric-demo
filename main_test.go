package main

import (
	"fabric-demo/service"
	"os"
	"testing"
)

var (
	configPath        = "config/config_e2e.yaml"
	channelId         = "mychannel"
	userName          = "User1"
	orgName           = "Org1"
	ordererOrgName    = "orderer.example.com"
	orgAdmin          = "Admin"
	chainCodeId       = "mycc"
	chainCodeVarsion  = "1.0"
	chaincodeGoPath   = os.Getenv("GOPATH")
	chaincodePath     = "fabric-demo/fabric/chaincode/chaincode_example02/go"
	channelConfigPath = "./fabric/channel-artifacts/channel.tx"
)

func TestInitClient(t *testing.T) {
	_, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
}

func TestCreateChannel(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	//创建通道
	err = fabClient.CreateChannel(channelId, channelConfigPath, orgAdmin, ordererOrgName)
	if err != nil {
		t.Errorf("创建通道失败: %v \n", err)
		return
	}
}

func TestJoinChannel(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	//加入通道
	err = fabClient.JoinChannel(channelId, ordererOrgName)
	if err != nil {
		t.Errorf("加入通道失败: %v \n", err)
		return
	}
}

func TestInstallChainCode(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	//安装链码
	err = fabClient.InstallChainCode(chainCodeId, chainCodeVarsion, chaincodePath, chaincodeGoPath)
	if err != nil {
		t.Errorf("安装链码失败: %v \n", err)
		return
	}
}

func TestInstantiateChainCode(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	//实例化链码
	policy := "OR ('Org1MSP.peer','Org2MSP.peer')"
	args := []string{"init", "a", "100", "b", "200"}
	peer := "peer0.org1.example.com"
	err = fabClient.InstantiateChainCode(channelId, chainCodeId, chainCodeVarsion, chaincodePath, policy, args, peer)
	if err != nil {
		t.Errorf("实例化链码失败: %v \n", err)
		return
	}
}
func TestQueryChainCode(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	ccClient := service.CcClient{}
	ccClient.Client, err = fabClient.GetChannelClient(channelId, userName, orgName)
	if err != nil {
		t.Errorf("实例化链码客户端失败: %v \n", err)
		return
	}
	err = ccClient.Query(chainCodeId, "query", []string{"a"})
	if err != nil {
		t.Errorf("查询链码失败: %v \n", err)
		return
	}
}

func TestInvokeChainCode(t *testing.T) {
	//初始化客户端
	fabClient, err := service.New(configPath, orgAdmin, orgName)
	if err != nil {
		t.Errorf("fabricClient实例化失败: %v \n", err)
		return
	}
	ccClient := service.CcClient{}
	ccClient.Client, err = fabClient.GetChannelClient(channelId, userName, orgName)
	if err != nil {
		t.Errorf("实例化链码客户端失败: %v \n", err)
		return
	}
	err = ccClient.Invoke(chainCodeId, "invoke", []string{"a", "b", "10"})
	if err != nil {
		t.Errorf("调用链码失败: %v \n", err)
		return
	}
}
