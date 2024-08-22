package dao

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/Chengxufeng1994/go-ddd-auth/internal/domain/identity_access_mgmt/aggregate"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/persistence/po"
	"github.com/Chengxufeng1994/go-ddd-auth/internal/infrastructure/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// The layer cake is as follows: (From bottom to top)
// SQL layer
// |
// Retry layer
// |
// Search layer
// |
// Timer layer
// |
// Cache layer

var UserSearchAll = []string{"Username"}

type IUserDao interface {
	Save(ctx context.Context, tx *gorm.DB, po *po.User) error
	DeleteByID(ctx context.Context, tx *gorm.DB, id int) error

	GetByID(ctx context.Context, tx *gorm.DB, id int) (*po.User, error)
	GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*po.User, error)
	Search(ctx context.Context, tx *gorm.DB, opts *aggregate.SearchUserOpts) ([]*po.User, *util.Pagination, error)
}

type UserDao struct {
	db *gorm.DB
}

var _ IUserDao = (*UserDao)(nil)

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (dao *UserDao) Save(ctx context.Context, tx *gorm.DB, userPO *po.User) error {
	return tx.WithContext(ctx).
		Model(&po.User{}).
		Where("id = ?", userPO.ID).
		Save(&userPO).Error
}

func (dao *UserDao) GetByID(ctx context.Context, tx *gorm.DB, id int) (*po.User, error) {
	var user po.User
	err := tx.WithContext(ctx).
		Model(&po.User{}).
		Preload("Role").
		Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (dao *UserDao) DeleteByID(ctx context.Context, tx *gorm.DB, id int) error {
	return tx.WithContext(ctx).
		Model(&po.User{}).
		Where("id = ?", id).
		Delete(&po.User{}).Error
}

func (dao *UserDao) GetByUsername(ctx context.Context, tx *gorm.DB, username string) (*po.User, error) {
	var user po.User
	err := tx.Debug().WithContext(ctx).
		Model(&po.User{}).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (dao *UserDao) Search(ctx context.Context, tx *gorm.DB, opts *aggregate.SearchUserOpts) ([]*po.User, *util.Pagination, error) {
	var totalCount int64
	tx.WithContext(ctx).Model(&po.User{}).Count(&totalCount)
	f := math.Ceil(float64(totalCount) / float64(opts.PerPage))
	totalPages := int64(f)

	var users []*po.User
	term := util.SanitizeSearchTerm(opts.Term, "*")
	query := tx.Debug().WithContext(ctx).
		Model(&po.User{}).
		Offset((opts.Page - 1) * opts.PerPage).
		Limit(opts.PerPage).
		Order(clause.OrderByColumn{
			Column: clause.Column{Name: opts.OrderBy},
			Desc:   opts.SortBy == "desc"})
	query = dao.generateSearchQuery(query, strings.Fields(term), UserSearchAll)
	err := query.Find(&users).Error
	if err != nil {
		return nil, nil, err
	}

	return users, &util.Pagination{
		TotalCount: totalCount,
		TotalPages: totalPages,
		Page:       opts.Page,
		Size:       opts.PerPage,
	}, nil
}

func (dao *UserDao) generateSearchQuery(query *gorm.DB, terms []string, fields []string) *gorm.DB {
	for _, term := range terms {
		searchFields := []string{}
		termArgs := []any{}
		for _, field := range fields {
			searchFields = append(searchFields, fmt.Sprintf("lower(%s) LIKE lower(?) escape '*' ", field))
			termArgs = append(termArgs, fmt.Sprintf("%%%s%%", strings.TrimLeft(term, "@")))
		}
		// searchFields = append(searchFields, "Id = ?")
		// termArgs = append(termArgs, strings.TrimLeft(term, "@"))
		query = query.Where(fmt.Sprintf("(%s)", strings.Join(searchFields, " OR ")), termArgs...)
	}

	return query
}
