package entity

import (
	"crypto/rand"
	"errors"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
	"github.com/paulmach/orb"
)

var (
	// ErrInvalidName はユーザー名が無効な場合のエラー
	ErrInvalidName = errors.New("name must not be empty")
	// ErrInvalidEmail はメールアドレスが無効な場合のエラー
	ErrInvalidEmail = errors.New("email must not be empty")
)

// UserID はユーザーの一意な識別子
type UserID string

// NewUserID は新しいUserIDを生成します
func NewUserID() UserID {
	return UserID(ulid.MustNew(ulid.Now(), rand.Reader).String())
}

// String はUserIDの文字列表現を返します
func (id UserID) String() string {
	return string(id)
}

// Geometry はドメイン層でのジオメトリ型（PostGISのgeometryに対応）
type Geometry struct {
	orb.Geometry
}

// User はユーザーエンティティを表します
type User struct {
	id                 UserID
	name               string
	highlightedPhotoID *int64
	locale             *string
	description        *string
	locality           *string
	administrativeArea *string
	countryCode        *string
	postalCode         *string
	geom               *Geometry
	firstName          *string
	lastName           *string
	email              *string
	hasSetLocation     bool
}

// 新しいユーザーを作成
func NewUser(
	name string,
	email *string,
	firstName *string,
	lastName *string,
) (*User, error) {
	// バリデーション
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}

	if email != nil && strings.TrimSpace(*email) == "" {
		return nil, ErrInvalidEmail
	}

	return &User{
		id:             NewUserID(),
		name:           name,
		email:          email,
		firstName:      firstName,
		lastName:       lastName,
		hasSetLocation: false,
	}, nil
}

// ReconstructUser は既存のユーザーを再構築します（リポジトリからの取得時に使用）
func ReconstructUser(
	id UserID,
	name string,
	highlightedPhotoID *int64,
	locale *string,
	description *string,
	locality *string,
	administrativeArea *string,
	countryCode *string,
	postalCode *string,
	geom *Geometry,
	firstName *string,
	lastName *string,
	email *string,
	hasSetLocation bool,
) (*User, error) {
	// 再構築時の最小限のバリデーション
	if strings.TrimSpace(name) == "" {
		return nil, ErrInvalidName
	}

	return &User{
		id:                 id,
		name:               name,
		highlightedPhotoID: highlightedPhotoID,
		locale:             locale,
		description:        description,
		locality:           locality,
		administrativeArea: administrativeArea,
		countryCode:        countryCode,
		postalCode:         postalCode,
		geom:               geom,
		firstName:          firstName,
		lastName:           lastName,
		email:              email,
		hasSetLocation:     hasSetLocation,
	}, nil
}

// ゲッターメソッド

func (u *User) ID() UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) HighlightedPhotoID() *int64 {
	return u.highlightedPhotoID
}

func (u *User) Locale() *string {
	return u.locale
}

func (u *User) Description() *string {
	return u.description
}

func (u *User) Locality() *string {
	return u.locality
}

func (u *User) AdministrativeArea() *string {
	return u.administrativeArea
}

func (u *User) CountryCode() *string {
	return u.countryCode
}

func (u *User) PostalCode() *string {
	return u.postalCode
}

func (u *User) Geom() *Geometry {
	return u.geom
}

func (u *User) FirstName() *string {
	return u.firstName
}

func (u *User) LastName() *string {
	return u.lastName
}

func (u *User) Email() *string {
	return u.email
}

func (u *User) HasSetLocation() bool {
	return u.hasSetLocation
}

// ビジネスロジックメソッド

// ユーザーのプロフィール情報を更新します
func (u *User) UpdateProfile(
	name *string,
	description *string,
	firstName *string,
	lastName *string,
) error {
	if name != nil {
		if strings.TrimSpace(*name) == "" {
			return ErrInvalidName
		}
		u.name = *name
	}

	if description != nil {
		u.description = description
	}

	if firstName != nil {
		u.firstName = firstName
	}

	if lastName != nil {
		u.lastName = lastName
	}

	return nil
}

// ユーザーの位置情報を設定
func (u *User) SetLocation(
	locality *string,
	administrativeArea *string,
	countryCode *string,
	postalCode *string,
	geom *Geometry,
) {
	u.locality = locality
	u.administrativeArea = administrativeArea
	u.countryCode = countryCode
	u.postalCode = postalCode
	u.geom = geom
	u.hasSetLocation = true
}

// SetHighlightedPhoto はハイライト写真を設定します
func (u *User) SetHighlightedPhoto(photoID int64) {
	u.highlightedPhotoID = &photoID
}

// ClearHighlightedPhoto はハイライト写真をクリアします
func (u *User) ClearHighlightedPhoto() {
	u.highlightedPhotoID = nil
}

// SetLocale はロケールを設定します
func (u *User) SetLocale(locale string) {
	u.locale = &locale
}

// String はユーザーの文字列表現を返します
func (u *User) String() string {
	return fmt.Sprintf("User{id=%s, name=%s, email=%v}", u.id, u.name, u.email)
}
