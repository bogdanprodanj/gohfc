/*
Copyright: Cognition Foundry. All Rights Reserved.
License: Apache License Version 2.0
*/
package gohfc

import (
	"context"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/protos/peer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
)

type Peer interface {
	Endorse(ctx context.Context, resp chan *PeerResponse, prop *peer.SignedProposal)
	GetURI() string
	GetOpts() []grpc.DialOption
}

// Peer expose API's to communicate with peer
type PeerClient struct {
	Name   string
	URI    string
	MspId  string
	Opts   []grpc.DialOption
	caPath string
	conn   *grpc.ClientConn
	client peer.EndorserClient
}

// PeerResponse is response from peer transaction request
type PeerResponse struct {
	Response *peer.ProposalResponse
	Err      error
	Name     string
}

// Endorse sends single transaction to single peer.
func (p *PeerClient) Endorse(ctx context.Context, resp chan *PeerResponse, prop *peer.SignedProposal) {
	if p.conn == nil {
		conn, err := grpc.DialContext(ctx, p.URI, p.Opts...)
		if err != nil {
			resp <- &PeerResponse{Response: nil, Err: err, Name: p.Name}
			return
		}
		p.conn = conn
		p.client = peer.NewEndorserClient(p.conn)
	}

	proposalResp, err := p.client.ProcessProposal(ctx, prop)
	if err != nil {
		resp <- &PeerResponse{Response: nil, Name: p.Name, Err: err}
		return
	}
	resp <- &PeerResponse{Response: proposalResp, Name: p.Name, Err: nil}
}

func (p *PeerClient) GetURI() string {
	return p.URI
}

func (p *PeerClient) GetOpts() []grpc.DialOption {
	return p.Opts
}

// NewPeerFromConfig creates new peer from provided config
func NewPeerFromConfig(conf PeerConfig) (*PeerClient, error) {
	p := PeerClient{URI: conf.Host, caPath: conf.TlsPath}
	if !conf.UseTLS {
		p.Opts = []grpc.DialOption{grpc.WithInsecure()}
	} else if p.caPath != "" {
		creds, err := credentials.NewClientTLSFromFile(p.caPath, "")
		if err != nil {
			return nil, fmt.Errorf("cannot read peer %s credentials err is: %v", p.Name, err)
		}
		p.Opts = append(p.Opts, grpc.WithTransportCredentials(creds))
	}

	p.Opts = append(p.Opts,
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                time.Minute,
			Timeout:             20 * time.Second,
			PermitWithoutStream: true,
		}),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxRecvMsgSize),
			grpc.MaxCallSendMsgSize(maxSendMsgSize)))
	return &p, nil
}
