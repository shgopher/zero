// Copyright 2022 Lingfei Kong <colin404@foxmail.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/superproj/zero.
//

//nolint:errchkjson,errcheck
package ws

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"

	"golang.org/x/net/websocket"

	"github.com/superproj/zero/internal/toyblc/blc"
	"github.com/superproj/zero/pkg/log"
)

type ByIndex []*blc.Block

func (b ByIndex) Len() int           { return len(b) }
func (b ByIndex) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
func (b ByIndex) Less(i, j int) bool { return b[i].Index < b[j].Index }

func ConnectToPeers(bs *blc.BlockSet, ss *Sockets, peers []string) {
	for _, peer := range peers {
		if peer == "" {
			continue
		}

		ws, err := websocket.Dial(peer, "", peer)
		if err != nil {
			log.Errorw(err, "Dial to peer", "peer", peer)
			continue
		}

		go WSHandler(bs, ss, ws)

		log.Debugw("Query latest block")
		ws.Write(bs.LatestMessage())
	}
}

func WSHandler(bs *blc.BlockSet, ss *Sockets, ws *websocket.Conn) {
	var (
		resp = &blc.ResponseBlockchain{}
		peer = ws.LocalAddr().String()
	)

	ss.Add(ws)

	for {
		var msg []byte
		if err := websocket.Message.Receive(ws, &msg); err != nil {
			if errors.Is(err, io.EOF) {
				log.Warnw("p2p peer shutdown, remove it form peers pool", "peer", peer)
				break
			}

			log.Errorw(err, "Can't receive p2p msg from ", "peer", peer)
			break
		}

		log.Debugw("Received message", "peer", peer, "message", msg)
		if err := json.Unmarshal(msg, resp); err != nil {
			log.Warnw("Invalid p2p message", "err", err)
		}

		switch resp.Type {
		case blc.QueryLatestAction:
			resp.Type = blc.ResponseAction

			message := bs.LatestMessage()
			log.Debugw("Response latest message", "message", message)
			ws.Write(message)

		case blc.QueryAllAction:
			resp.Type = blc.ResponseAction
			resp.Data, _ = bs.MarshalJSON()
			data, _ := json.Marshal(resp)
			log.Debugw("Response chain message", "message", data)
			ws.Write(data)

		case blc.ResponseAction:
			ResponseBlockchain(bs, ss, resp.Data)
		}
	}
}

func ResponseBlockchain(bs *blc.BlockSet, ss *Sockets, msg []byte) {
	receivedBlocks := []*blc.Block{}

	if err := json.Unmarshal(msg, &receivedBlocks); err != nil {
		log.Warnw("Invalid blockchain", "err", err)
	}

	sort.Sort(ByIndex(receivedBlocks))

	latestBlockReceived := receivedBlocks[len(receivedBlocks)-1]
	latestBlockHeld := bs.Latest()
	if latestBlockReceived.Index <= latestBlockHeld.Index {
		log.Infow("received blockchain is not longer than current blockchain. Do nothing.")
		return
	}

	log.Warnf("blockchain possibly behind. We got: %d Peer got: %d", latestBlockHeld.Index, latestBlockReceived.Index)
	if latestBlockHeld.Hash == latestBlockReceived.PreviousHash {
		log.Infof("We can append the received block to our chain.")
		bs.Add(latestBlockReceived)
	} else if len(receivedBlocks) == 1 {
		log.Infow("We have to query the chain from our peer.")
		ss.Broadcast(queryAllMsg())
	} else {
		log.Infow("Received blockchain is longer than current blockchain.")
		replaceBlocks(receivedBlocks, bs, ss)
	}
}

func queryAllMsg() []byte {
	return []byte(fmt.Sprintf("{\"type\": %d}", blc.QueryAllAction))
}

func replaceBlocks(src []*blc.Block, dst *blc.BlockSet, ss *Sockets) {
	if !blc.IsValidChain(src) || len(src) <= dst.Len() {
		log.Errorf("Received blockchain invalid.")
		return
	}

	log.Debugw("Received blockchain is valid. Replacing current blockchain with received blockchain.")
	dst.SetBlocks(src)
	ss.Broadcast(dst.LatestMessage())
}
