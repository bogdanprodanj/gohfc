/*
Copyright: Cognition Foundry. All Rights Reserved.
License: Apache License Version 2.0
*/
package gohfc

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/orderer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

// Orderer provides an atomic broadcast ordering service for consumption by the peers
type Orderer interface {
	Broadcast(envelope *common.Envelope) (*orderer.BroadcastResponse, error)
	Deliver(ctx context.Context, envelope *common.Envelope) (*common.Block, error)
	getGenesisBlock(ctx context.Context, identity Identity, crypto CryptoSuite, channelId string) (*common.Block, error)
}

// OrdererService expose API's to communicate with orderers.
type OrdererService struct {
	Name   string
	Uri    string
	Opts   []grpc.DialOption
	caPath string
	con    *grpc.ClientConn
	client orderer.AtomicBroadcastClient
}

// Broadcast Broadcast envelope to orderer for execution.
func (o *OrdererService) Broadcast(envelope *common.Envelope) (*orderer.BroadcastResponse, error) {
	if o.con == nil {
		c, err := grpc.Dial(o.Uri, o.Opts...)
		if err != nil {
			return nil, fmt.Errorf("cannot connect to orderer: %s err is: %v", o.Name, err)
		}
		o.con = c
		o.client = orderer.NewAtomicBroadcastClient(o.con)
	}
	bcc, err := o.client.Broadcast(context.Background())
	if err != nil {
		return nil, err
	}
	defer bcc.CloseSend()
	bcc.Send(envelope)
	response, err := bcc.Recv()
	if err != nil {
		return nil, err
	}
	if response.Status != common.Status_SUCCESS {
		return nil, fmt.Errorf("unexpected status: %v", response.Status)
	}

	return response, err
}

// Deliver delivers envelope to orderer. Please note that new connection will be created on every call of Deliver.
func (o *OrdererService) Deliver(ctx context.Context, envelope *common.Envelope) (*common.Block, error) {
	connection, err := grpc.Dial(o.Uri, o.Opts...)
	if err != nil {
		return nil, fmt.Errorf("cannot connect to orderer: %s err is: %v", o.Name, err)
	}
	defer connection.Close()

	dk, err := orderer.NewAtomicBroadcastClient(connection).Deliver(ctx)
	if err != nil {
		return nil, err
	}
	if err := dk.Send(envelope); err != nil {
		return nil, err
	}
	var block *common.Block
	response, err := dk.Recv()
	if err != nil {
		return nil, err
	}
	switch t := response.Type.(type) {
	// Seek operation success, no more responses
	case *orderer.DeliverResponse_Status:
		if t.Status != common.Status_SUCCESS {
			return nil, fmt.Errorf("orderer response with status: %d", t.Status)
		}
		return block, nil
	// Response is a requested block
	case *orderer.DeliverResponse_Block:
		block = response.GetBlock()
		return block, nil
	// Unknown response
	default:
		return nil, fmt.Errorf("unknown response type from orderer: %s", t)
	}
}

func (o *OrdererService) getGenesisBlock(ctx context.Context, identity Identity, crypto CryptoSuite, channelId string) (*common.Block, error) {
	seekInfo := &orderer.SeekInfo{
		Start:    &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: 0}}},
		Stop:     &orderer.SeekPosition{Type: &orderer.SeekPosition_Specified{Specified: &orderer.SeekSpecified{Number: 0}}},
		Behavior: orderer.SeekInfo_BLOCK_UNTIL_READY,
	}
	seekInfoBytes, err := proto.Marshal(seekInfo)
	if err != nil {
		return nil, err
	}

	creator, err := marshalProtoIdentity(identity)
	if err != nil {
		return nil, err
	}
	txId, err := newTransactionId(creator)
	if err != nil {
		return nil, err
	}

	headerBytes, err := channelHeader(common.HeaderType_DELIVER_SEEK_INFO, txId, channelId, 0, nil)
	signatureHeaderBytes, err := signatureHeader(creator, txId)
	if err != nil {
		return nil, err
	}
	header := header(signatureHeaderBytes, headerBytes)
	payloadBytes, err := payload(header, seekInfoBytes)
	if err != nil {
		return nil, err
	}
	payloadSignedBytes, err := crypto.Sign(payloadBytes, identity.PrivateKey)
	if err != nil {
		return nil, err
	}
	env := &common.Envelope{Payload: payloadBytes, Signature: payloadSignedBytes}
	return o.Deliver(ctx, env)
}

// NewOrdererFromConfig create new OrdererService from config
func NewOrdererFromConfig(conf OrdererConfig) (*OrdererService, error) {
	o := OrdererService{Uri: conf.Host, caPath: conf.TlsPath}
	if !conf.UseTLS {
		o.Opts = []grpc.DialOption{grpc.WithInsecure()}
	} else if o.caPath != "" {
		creds, err := credentials.NewClientTLSFromFile(o.caPath, "")
		if err != nil {
			return nil, fmt.Errorf("cannot read orderer %s credentials err is: %v", o.Name, err)
		}
		o.Opts = append(o.Opts, grpc.WithTransportCredentials(creds))
	}
	o.Opts = append(o.Opts,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize)))
	return &o, nil
}
