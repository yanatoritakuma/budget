package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"

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
	oauth2Config *oauth2.Config
	userRepo     user.UserRepository
	userUsecase  UserUsecase // JWT生成のために既存のUserUsecaseを利用
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
		userRepo:    userRepo,
		userUsecase: userUsecase,
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
	authURL := uc.oauth2Config.AuthCodeURL(state,
		oauth2.SetAuthURLParam("prompt", "consent"),
		oauth2.SetAuthURLParam("bot_prompt", "aggressive"))
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

	// IDトークンのパースと検証
	claims := &LineIDTokenClaims{}
	_, err = jwt.ParseWithClaims(rawIDToken, claims, func(token *jwt.Token) (interface{}, error) {
		// LINEのIDトークンはHS256で署名されている。
		// 実際にはLINEの公開鍵（JWKS URIから取得）で署名を検証すべきだが、
		// 今回はChannel SecretをHMACの鍵として利用する簡易的な検証を行う。
		// よりセキュアな実装のためには、LINEのOpenID Connect Discovery Endpointから
		// JWKS URIを取得し、公開鍵を使って検証するべき。
		return []byte(os.Getenv("LINE_CHANNEL_SECRET")), nil
	}, jwt.WithValidMethods([]string{"HS256"}),
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
