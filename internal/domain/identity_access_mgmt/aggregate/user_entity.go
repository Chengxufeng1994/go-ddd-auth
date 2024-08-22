package aggregate

import (
	"fmt"
	"time"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/valueobject"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/assert"
)

// Aggregate root
type User struct {
	ID        int
	Username  string
	Password  string
	RoleID    *valueobject.RoleID
	CreatedAt time.Time
	UpdatedAt time.Time

	snapshot *User
}

// create new user
func NewUser(username, hashedPassword string, roleID *valueobject.RoleID) (*User, error) {
	user := &User{}
	if err := user.SetUsername(username); err != nil {
		return nil, err
	}
	if err := user.SetPassword(hashedPassword); err != nil {
		return nil, err
	}
	user.RoleID = roleID

	return user, nil
}

// rebuild user
func Hydrate(id int, username, password string, roleID *valueobject.RoleID) *User {
	return &User{
		ID:       id,
		Username: username,
		Password: password,
		RoleID:   roleID,
	}
}

func (u *User) Create() {
	u.CreatedAt = time.Now().UTC()
}

func (u *User) SetUsername(aUsername string) error {
	if err := assert.AssertArgumentNotEmpty(aUsername); err != nil {
		return err
	}

	if err := assert.AssertArgumentLength(aUsername, 1, 255); err != nil {
		return err
	}

	u.Username = aUsername
	return nil
}

func (u *User) SetPassword(aPassword string) error {
	u.Password = aPassword
	return nil
}

func (u *User) ChangeUsername(aUsername string) error {
	if err := assert.AssertArgumentEmpty(aUsername); err == nil {
		return err
	}

	if err := assert.AssertArgumentNotEqual(u.Username, aUsername); err != nil {
		return err
	}

	u.Username = aUsername

	return nil
}

func (u *User) ChangePassword(hashedPassword string) error {
	if err := assert.AssertArgumentEmpty(hashedPassword); err == nil {
		return err
	}

	if err := assert.AssertArgumentNotEqual(u.Password, hashedPassword); err != nil {
		return err
	}

	u.Password = hashedPassword
	return nil
}

func (u *User) Update(opts UpdateUserOpt) {
	if err := u.ChangeUsername(*opts.Username); err != nil {
	}

	if err := u.ChangePassword(*opts.Password); err != nil {
	}

	u.UpdatedAt = time.Now().UTC()
}

func (u *User) DeepCopy() *User {
	return &User{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		RoleID:    u.RoleID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u *User) Attach() {
	if u.snapshot == nil || u.snapshot.ID == u.ID {
		u.snapshot = u.DeepCopy()
	}
}

func (u *User) Detach() {
	if u.snapshot != nil && u.snapshot.ID == u.ID {
		u.snapshot = nil
	}
}

func (u *User) DetectChanges() *UserDiff {
	if u.snapshot == nil {
		return nil
	}
	userChanged := false
	if u.Password != u.snapshot.Password {
		userChanged = true
	}
	if u.Username != u.snapshot.Username {
		userChanged = true
	}
	if u.RoleID != u.snapshot.RoleID {
		userChanged = true
	}

	return &UserDiff{
		UserChanged: userChanged,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("User{ID: %d, Username: %s, RoleID: %v, CreatedAt: %s, UpdatedAt: %s}", u.ID, u.Username, u.RoleID, u.CreatedAt, u.UpdatedAt)
}

type UpdateUserOpt struct {
	Username *string
	Password *string
	RoleID   *int
}

type SearchUserOpts struct {
	Page    int
	PerPage int
	OrderBy string
	SortBy  string
	Term    string
}

type UserDiff struct {
	UserChanged bool
}
