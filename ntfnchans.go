// Copyright (c) 2017, Jonathan Chappelow
// See LICENSE for details.

package main

import (
	"github.com/decred/dcrd/chaincfg/chainhash"
	"github.com/decred/dcrutil"
)

const (
	// blockConnChanBuffer is the size of the block connected channel buffer.
	blockConnChanBuffer = 8

	// newTxChanBuffer is the size of the new transaction channel buffer, for
	// ANY transactions are added into mempool.
	newTxChanBuffer = 2000

	// relevantMempoolTxChanBuffer is the size of the new transaction channel
	// buffer, for relevant transactions that are added into mempool.
	relevantMempoolTxChanBuffer = 512
)

// BlockWatchedTx contains, for a certain block, the transactions for certain
// watched addresses
type BlockWatchedTx struct {
	BlockHeight   int64
	TxsForAddress map[string][]*dcrutil.Tx
}

// Channels are package-level variables for simplicity
var ntfnChans struct {
	connectChan                       chan *chainhash.Hash
	connectChanStkInf                 chan int32
	spendTxBlockChan, recvTxBlockChan chan *BlockWatchedTx
	relevantTxMempoolChan             chan *dcrutil.Tx
	newTxChan                         chan *chainhash.Hash
}

func makeNtfnChans(cfg *config) {
	// If we're monitoring for blocks OR collecting block data, these channels
	// are necessary to handle new block notifications. Otherwise, leave them
	// as nil so that both a send (below) blocks and a receive (in spy.go,
	// blockConnectedHandler) block. default case makes non-blocking below.
	// quit channel case manages blockConnectedHandlers.
	if !cfg.NoCollectBlockData && !cfg.NoMonitor {
		ntfnChans.connectChan = make(chan *chainhash.Hash, blockConnChanBuffer)
		//ntfnChans.stakeDiffChan = make(chan int64, blockConnChanBuffer)
	}

	// Like connectChan for block data, connectChanStkInf is used when a new
	// block is connected, but to signal the stake info monitor.
	if !cfg.NoCollectStakeInfo && !cfg.NoMonitor {
		ntfnChans.connectChanStkInf = make(chan int32, blockConnChanBuffer)
	}

	// watchaddress
	if len(cfg.WatchAddresses) > 0 && !cfg.NoMonitor {
		// recv/spendTxBlockChan come with connected blocks
		ntfnChans.recvTxBlockChan = make(chan *BlockWatchedTx, blockConnChanBuffer)
		ntfnChans.spendTxBlockChan = make(chan *BlockWatchedTx, blockConnChanBuffer)
		ntfnChans.relevantTxMempoolChan = make(chan *dcrutil.Tx, relevantMempoolTxChanBuffer)
	}

	if cfg.MonitorMempool {
		ntfnChans.newTxChan = make(chan *chainhash.Hash, newTxChanBuffer)
	}
}

func closeNtfnChans() {
	// if ntfnChans.stakeDiffChan != nil {
	// 	close(ntfnChans.stakeDiffChan)
	// }
	if ntfnChans.connectChan != nil {
		close(ntfnChans.connectChan)
	}
	if ntfnChans.connectChanStkInf != nil {
		close(ntfnChans.connectChanStkInf)
	}

	if ntfnChans.newTxChan != nil {
		close(ntfnChans.newTxChan)
	}
	if ntfnChans.relevantTxMempoolChan != nil {
		close(ntfnChans.relevantTxMempoolChan)
	}

	if ntfnChans.spendTxBlockChan != nil {
		close(ntfnChans.spendTxBlockChan)
	}
	if ntfnChans.recvTxBlockChan != nil {
		close(ntfnChans.recvTxBlockChan)
	}
}
