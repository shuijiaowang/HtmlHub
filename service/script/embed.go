package script

import _ "embed"

// syncHelperJS is embedded at build time to avoid runtime path issues.
//
//go:embed sync_helper.js
var syncHelperJS string

//go:embed publisher_force_sync.js
var publisherForceSyncJS string

func SyncHelperJS() string {
	return syncHelperJS
}

func PublisherForceSyncJS() string {
	return publisherForceSyncJS
}
