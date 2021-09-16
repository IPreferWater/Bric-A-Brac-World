package main

type State int

const (
        WaitPlayerAction State = iota
        AnimatePlayerAction
        WaitWorldAction
		AnimateWorldAction
)