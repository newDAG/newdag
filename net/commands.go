package net

import "github.com/newdag/consensus"

type SyncRequest struct {
	FromID int
	Known  map[int]int
}

type SyncResponse struct {
	FromID    int
	SyncLimit bool
	Events    []consensus.WireEvent
	Known     map[int]int
}

//++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++

type EagerSyncRequest struct {
	FromID int
	Events []consensus.WireEvent
}

type EagerSyncResponse struct {
	FromID  int
	Success bool
}
