package crlmgr

//nolint:gochecknoglobals // local global for looping all actions.
var eventTypes = []EventType{
	PreCheckAction,
	PreInstallAction,
	PostInstallAction,
	PostAction,
}

// EventType is the type of actions available.
type EventType string

const (
	// PreCheckAction is the action that is executed before a target is checked.
	PreCheckAction EventType = "precheck"

	// PreInstallAction is the action that is executed before a target is updated.
	PreInstallAction EventType = "preinstall"

	// PostInstallAction is the action executed after a target is updated.
	PostInstallAction EventType = "postinstall"

	// PostAction is the action executed after a target is checked and possibly updated.
	PostAction EventType = "post"
)

func (e EventType) String() string {
	return string(e)
}
