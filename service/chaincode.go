package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

type CcClient struct {
	*channel.Client
}

// 查询链码
func (c *CcClient) Query(chainCodeId, fcn string, args []string) error {
	req := channel.Request{
		ChaincodeID: chainCodeId,
		Fcn:         fcn,
		Args:        ConvertArgs(args),
	}
	response, err := c.Client.Query(req)
	if err != nil {
		fmt.Printf("查询链码失败: %v \n", err)
		return err
	}
	fmt.Printf("查询结果: %s\n", string(response.Payload))
	return nil
}

// 调用链码
func (c *CcClient) Invoke(chainCodeId, fcn string, args []string) error {
	req := channel.Request{
		ChaincodeID: chainCodeId,
		Fcn:         fcn,
		Args:        ConvertArgs(args),
	}
	response, err := c.Client.Execute(req)
	if err != nil {
		fmt.Printf("调用链码失败: %v \n", err)
		return err
	}
	fmt.Printf("调用链码成功，TransactionID: %s\n", response.TransactionID)
	return nil
}
