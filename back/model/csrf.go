package model

import (
	"sync"
	"time"
)

// CSRFToken はトークンと有効期限を管理する構造体
type CSRFToken struct {
	Token     string
	ExpiresAt time.Time
}

// TokenStore はCSRFトークンを保存するためのストア
type TokenStore struct {
	tokens map[string]CSRFToken
	mutex  sync.RWMutex
}

// NewTokenStore は新しいTokenStoreを作成します
func NewTokenStore() *TokenStore {
	return &TokenStore{
		tokens: make(map[string]CSRFToken),
	}
}

// SaveToken はトークンを保存します
func (s *TokenStore) SaveToken(sessionID string, token CSRFToken) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tokens[sessionID] = token
}

// ValidateToken はトークンを検証します
func (s *TokenStore) ValidateToken(sessionID, token string) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	storedToken, exists := s.tokens[sessionID]
	if !exists {
		return false
	}

	// 有効期限チェック
	if time.Now().After(storedToken.ExpiresAt) {
		delete(s.tokens, sessionID)
		return false
	}

	return storedToken.Token == token
}

// GetToken は保存されているトークンを取得します
// 変更
func (s *TokenStore) GetToken(sessionID string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	storedToken, exists := s.tokens[sessionID]
	if !exists {
		return "", false
	}

	// 有効期限チェック
	if time.Now().After(storedToken.ExpiresAt) {
		delete(s.tokens, sessionID)
		return "", false
	}

	return storedToken.Token, true
}

// DeleteToken は指定されたセッションIDのトークンを削除します
func (s *TokenStore) DeleteToken(sessionID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.tokens, sessionID)
} // CleanupExpiredTokens は期限切れのトークンを削除します
func (s *TokenStore) CleanupExpiredTokens() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now()
	for sessionID, token := range s.tokens {
		if now.After(token.ExpiresAt) {
			delete(s.tokens, sessionID)
		}
	}
}
