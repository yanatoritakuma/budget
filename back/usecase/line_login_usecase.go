package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/yanatoritakuma/budget/back/domain/user"
	"golang.org/x/oauth2"
)

// LineLoginUsecase はLINEログインに関するユースケースのインターフェースです。
type LineLoginUsecase interface {
	GetLineAuthURL(ctx context.Context, state string) (string, error)
	LineLoginCallback(ctx context.Context, code, state string) (string, error) // JWTを返す
}

// LineLoginUsecaseImpl はLineLoginUsecaseの実装です。
type LineLoginUsecaseImpl struct {
	oauth2Config    *oauth2.Config
	userRepo        user.UserRepository
	userUsecase     UserUsecase // JWT生成のために既存のUserUsecaseを利用
	jwksCache       map[string]interface{}
	jwksCacheMutex  sync.RWMutex
	jwksCacheExpiry time.Time
}

// NewLineLoginUsecaseImpl はLineLoginUsecaseImplの新しいインスタンスを生成します。
func NewLineLoginUsecaseImpl(userRepo user.UserRepository, userUsecase UserUsecase) LineLoginUsecase {
	return &LineLoginUsecaseImpl{
		oauth2Config: &oauth2.Config{
			ClientID:     os.Getenv("LINE_CHANNEL_ID"),
			ClientSecret: os.Getenv("LINE_CHANNEL_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
				TokenURL: "https://api.line.me/oauth2/v2.1/token",
			},
			RedirectURL: os.Getenv("LINE_REDIRECT_URI"),
			Scopes:      []string{"openid", "profile", "email"},
		},
		userRepo:        userRepo,
		userUsecase:     userUsecase,
		jwksCache:       make(map[string]interface{}),
		jwksCacheExpiry: time.Now(),
	}
}

// GenerateState はCSRF対策のためのランダムなstate文字列を生成します。
func GenerateState() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("failed to generate random state: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetLineAuthURL はLINEの認証URLを生成します。
func (uc *LineLoginUsecaseImpl) GetLineAuthURL(ctx context.Context, state string) (string, error) {
	if uc.oauth2Config.ClientID == "" || uc.oauth2Config.ClientSecret == "" || uc.oauth2Config.RedirectURL == "" {
		return "", fmt.Errorf("LINE環境変数が設定されていません")
	}
	// stateはCSRF対策のためにセッション等で管理される必要があります。
	// ここでは引数として受け取ったstateをそのまま使用します。
	authURL := uc.oauth2Config.AuthCodeURL(state)
	return authURL, nil
}

// JWTクレームの構造体
type LineIDTokenClaims struct {
	jwt.RegisteredClaims
	Name    string `json:"name,omitempty"`
	Picture string `json:"picture,omitempty"`
}

// LineLoginCallback はLINEからのコールバックを処理し、JWTを返します。
func (uc *LineLoginUsecaseImpl) LineLoginCallback(ctx context.Context, code, state string) (string, error) {
	// ここでstateの検証を行う必要がありますが、セッション管理の実装に依存するため、
	// 一旦は引数として受け取ったstateを使用するのみとします。

	// 認可コードとアクセストークンを交換
	token, err := uc.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange auth code for token: %w", err)
	}

	// IDトークンの検証とユーザー情報の取得
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return "", fmt.Errorf("id_token not found")
	}

	// IDトークンのパースと署名検証（JWKS使用）
	claims := &LineIDTokenClaims{}
	unverifiedClaims := &LineIDTokenClaims{}

	// トークンをパース（署名検証なし）して kid を取得
	_, _, err = new(jwt.Parser).ParseUnverified(rawIDToken, unverifiedClaims)
	if err != nil {
		return "", fmt.Errorf("failed to parse ID token: %w", err)
	}

	// kid を使用してキーを取得し、署名を検証
	_, err = jwt.ParseWithClaims(rawIDToken, claims, func(token *jwt.Token) (interface{}, error) {
		// トークンのアルゴリズムを確認
		switch token.Method.(type) {
		case *jwt.SigningMethodRSA:
			// RSA (RS256) の場合は公開鍵を取得
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return nil, fmt.Errorf("kid not found in token header")
			}
			publicKey, err := uc.getLinePublicKey(kid)
			if err != nil {
				return nil, fmt.Errorf("failed to get LINE public key: %w", err)
			}
			return publicKey, nil
		case *jwt.SigningMethodHMAC:
			// HMAC (HS256) の場合はチャネルシークレットを使用
			return []byte(os.Getenv("LINE_CHANNEL_SECRET")), nil
		default:
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
	},
		jwt.WithAudience(uc.oauth2Config.ClientID),
		jwt.WithIssuer("https://access.line.me"),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to parse or validate ID token: %w", err)
	}

	lineUserID := claims.Subject // subクレームがLINEユーザーID
	userName := claims.Name
	userImage := claims.Picture

	// ユーザーの登録または取得
	targetUser, err := uc.userUsecase.CreateUserForLine(lineUserID, userName, userImage)
	if err != nil {
		return "", fmt.Errorf("failed to create or get user for LINE login: %w", err)
	}

	// アプリケーションのJWTを生成
	jwtToken, err := uc.userUsecase.GenerateToken(targetUser)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT for LINE user: %w", err)
	}

	return jwtToken, nil
}

// getLinePublicKey はLINEのJWKSから指定された kid に対応する公開鍵を取得します。
func (uc *LineLoginUsecaseImpl) getLinePublicKey(kid string) (interface{}, error) {
	// キャッシュをチェック
	uc.jwksCacheMutex.RLock()
	if time.Now().Before(uc.jwksCacheExpiry) && len(uc.jwksCache) > 0 {
		if key, exists := uc.jwksCache[kid]; exists {
			uc.jwksCacheMutex.RUnlock()
			return key, nil
		}
	}
	uc.jwksCacheMutex.RUnlock()

	// JWKSを取得
	jwksURL := "https://api.line.me/oauth2/v2.1/certs"
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch JWKS from LINE: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch JWKS: status=%d, body=%s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read JWKS response: %w", err)
	}

	// JWKS JSON をパース
	var jwks struct {
		Keys []map[string]interface{} `json:"keys"`
	}
	if err := json.Unmarshal(body, &jwks); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JWKS: %w", err)
	}

	// キャッシュを更新
	uc.jwksCacheMutex.Lock()
	uc.jwksCache = make(map[string]interface{})
	for _, key := range jwks.Keys {
		if keyID, ok := key["kid"].(string); ok {
			uc.jwksCache[keyID] = key
		}
	}
	uc.jwksCacheExpiry = time.Now().Add(24 * time.Hour) // 24時間キャッシュ
	uc.jwksCacheMutex.Unlock()

	// 要求された kid に対応する鍵を取得
	if key, exists := uc.jwksCache[kid]; exists {
		// key は map[string]interface{} で、jwt.ParseWithClaims では RSA 公開鍵が必要
		// jwt ライブラリが自動で処理できるよう、raw JSON を返す
		return key, nil
	}

	return nil, fmt.Errorf("public key with kid %s not found in JWKS", kid)
}
