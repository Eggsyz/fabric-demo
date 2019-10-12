package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	promsp "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
)

type FabricClient struct {
	rcClient  *resmgmt.Client //fabric 资源管理客户端
	mspClient *msp.Client     //msp客户端
	sdk       *fabsdk.FabricSDK
}

//实例化fabric客户端
func New(configFile, orgAdmin, orgName string) (*FabricClient, error) {
	sdk, err := fabsdk.New(config.FromFile(configFile))
	if err != nil {
		return nil, fmt.Errorf("实例化sdk失败, err: %v", err)
	}
	clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	if clientContext == nil {
		return nil, fmt.Errorf("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}
	rcClient, err := resmgmt.New(clientContext)
	if err != nil {
		return nil, fmt.Errorf("根据指定的资源管理客户端Context创建通道管理客户端失败: %v", err)
	}
	mspClient, err := msp.New(sdk.Context(), msp.WithOrg(orgName))
	if err != nil {
		return nil, fmt.Errorf("根据指定的组织名称创建OrgMSP客户端实例失败: %v", err)
	}
	return &FabricClient{sdk: sdk, rcClient: rcClient, mspClient: mspClient}, nil
}

//创建通道
func (client *FabricClient) CreateChannel(channelID, channelConfigPath, orgAdmin, ordererOrgName string) error {
	adminIdentity, err := client.mspClient.GetSigningIdentity(orgAdmin)
	if err != nil {
		return fmt.Errorf("获取指定id的签名标识失败: %v", err)
	}
	channelReq := resmgmt.SaveChannelRequest{
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigPath,
		SigningIdentities: []promsp.SigningIdentity{adminIdentity},
	}
	response, err := client.rcClient.SaveChannel(channelReq, resmgmt.WithOrdererEndpoint(ordererOrgName))
	if err != nil {
		return fmt.Errorf("创建应用通道失败: %v", err)
	}
	fmt.Printf("创建通道成功. TransactionID: %s \n", response.TransactionID)
	return nil
}

// 加入通道
func (client *FabricClient) JoinChannel(channelID, ordererOrgName string) error {
	err := client.rcClient.JoinChannel(channelID, resmgmt.WithOrdererEndpoint(ordererOrgName), resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("Peers加入通道失败: %v", err)
	}
	fmt.Println("peers 已成功加入通道.")
	return nil
}

// 安装链码
func (client *FabricClient) InstallChainCode(chainCodeId, chainCodeVersion, chainCodePath, chainCodeGoPath string) error {
	ccPkg, err := gopackager.NewCCPackage(chainCodePath, chainCodeGoPath)
	if err != nil {
		return fmt.Errorf("创建链码包失败: %v", err)
	}
	req := resmgmt.InstallCCRequest{
		Name:    chainCodeId,
		Path:    chainCodePath,
		Version: chainCodeVersion,
		Package: ccPkg,
	}
	response, err := client.rcClient.InstallCC(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("安装链码失败: %v", err)
	}
	fmt.Printf("安装链码成功. %s \n", response)
	return nil
}

// 实例化链码
func (client *FabricClient) InstantiateChainCode(channelId, chainCodeId, chainCodeVersion, chainCodePath, policy string, args []string, peer string) error {
	ccPolicy, err := cauthdsl.FromString(policy)
	if err != nil {
		return fmt.Errorf("获取背书策略失败，err: %v \n", err)
	}

	req := resmgmt.InstantiateCCRequest{
		Name:    chainCodeId,
		Path:    chainCodePath,
		Version: chainCodeVersion,
		Policy:  ccPolicy,
		Args:    ConvertArgs(args),
	}
	response, err := client.rcClient.InstantiateCC(channelId, req, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargetEndpoints(peer))
	if err != nil {
		return fmt.Errorf("实例化链码失败: %v", err)
	}
	fmt.Printf("实例化链码成功: %s \n", response.TransactionID)
	return nil
}

func ConvertArgs(args []string) [][]byte {
	var ccArgs [][]byte
	for i := 0; i < len(args); i++ {
		ccArgs = append(ccArgs, []byte(args[i]))
	}
	return ccArgs
}

func (client *FabricClient) GetChannelClient(channelID, userName, orgName string) (*channel.Client, error) {
	clientChannelContext := client.sdk.ChannelContext(channelID, fabsdk.WithUser(userName), fabsdk.WithOrg(orgName))
	channelClient, err := channel.New(clientChannelContext)
	if err != nil {
		return nil, fmt.Errorf("创建应用通道客户端失败: %v", err)
	}
	fmt.Println("通道客户端创建成功，可以利用此客户端调用链码进行查询或执行事务.")
	return channelClient, nil
}
