package client

import (
	"sync"
	"time"

	"github.com/Ovler-Young/MiraiGo/message"
)

// globalRKeyManager is a shared RKeyManager across all clients
// since rkey is not client-specific
var globalRKeyManager = NewRKeyManager()

func init() {
	// Register the callback to capture rkey from parsed messages
	message.OnRKeyUpdate = func(rkey string, isGroup bool) {
		if isGroup {
			globalRKeyManager.SetGroupRKey(rkey)
		} else {
			globalRKeyManager.SetPrivateRKey(rkey)
		}
	}

	// Register the callback to provide rkey for URL construction
	message.GetRKey = func(isGroup bool) string {
		if isGroup {
			return globalRKeyManager.GetGroupRKey()
		}
		return globalRKeyManager.GetPrivateRKey()
	}
}

// GetGlobalRKeyManager returns the shared RKeyManager instance
func GetGlobalRKeyManager() *RKeyManager {
	return globalRKeyManager
}

// RKeyInfo holds cached rkey data for image downloads
type RKeyInfo struct {
	PrivateRKey string    // For appid=1406 (friend images)
	GroupRKey   string    // For appid=1407 (group images)
	ExpiredTime int64     // Unix timestamp when rkey expires
	UpdatedTime time.Time // When this info was last updated
}

// RKeyManager handles rkey caching for image URL construction
type RKeyManager struct {
	mu   sync.RWMutex
	info *RKeyInfo
}

// NewRKeyManager creates a new RKeyManager
func NewRKeyManager() *RKeyManager {
	return &RKeyManager{
		info: &RKeyInfo{},
	}
}

// GetGroupRKey returns the cached group rkey, or empty string if not available/expired
func (m *RKeyManager) GetGroupRKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.info == nil || m.info.GroupRKey == "" {
		return ""
	}

	// Check if expired (with 60 second buffer)
	if m.info.ExpiredTime > 0 && time.Now().Unix() > m.info.ExpiredTime-60 {
		return ""
	}

	return m.info.GroupRKey
}

// GetPrivateRKey returns the cached private rkey, or empty string if not available/expired
func (m *RKeyManager) GetPrivateRKey() string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.info == nil || m.info.PrivateRKey == "" {
		return ""
	}

	// Check if expired (with 60 second buffer)
	if m.info.ExpiredTime > 0 && time.Now().Unix() > m.info.ExpiredTime-60 {
		return ""
	}

	return m.info.PrivateRKey
}

// SetGroupRKey updates the cached group rkey (extracted from message)
func (m *RKeyManager) SetGroupRKey(rkey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.info == nil {
		m.info = &RKeyInfo{}
	}

	m.info.GroupRKey = rkey
	m.info.UpdatedTime = time.Now()
	// Set expiry to 1 hour from now (rkeys are typically valid for ~1 hour)
	m.info.ExpiredTime = time.Now().Add(time.Hour).Unix()
}

// SetPrivateRKey updates the cached private rkey (extracted from message)
func (m *RKeyManager) SetPrivateRKey(rkey string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.info == nil {
		m.info = &RKeyInfo{}
	}

	m.info.PrivateRKey = rkey
	m.info.UpdatedTime = time.Now()
	// Set expiry to 1 hour from now
	m.info.ExpiredTime = time.Now().Add(time.Hour).Unix()
}

// GetInfo returns a copy of the current rkey info
func (m *RKeyManager) GetInfo() *RKeyInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.info == nil {
		return nil
	}

	// Return a copy
	return &RKeyInfo{
		PrivateRKey: m.info.PrivateRKey,
		GroupRKey:   m.info.GroupRKey,
		ExpiredTime: m.info.ExpiredTime,
		UpdatedTime: m.info.UpdatedTime,
	}
}
