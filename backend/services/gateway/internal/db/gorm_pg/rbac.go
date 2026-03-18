package gorm_pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/chibx/vuecom/backend/shared/models/db/users"
	"github.com/chibx/vuecom/backend/shared/rbac"
	"github.com/chibx/vuecom/backend/shared/types"
	"gorm.io/gorm"
)

type rbacRepository struct {
	db *gorm.DB
}

func (rb *rbacRepository) GetUsersWithPermission(ctx context.Context, permission string, pageData ...types.Pagination) ([]users.BackendUser, error) {
	var users []users.BackendUser

	var limit = -1
	var offset = -1

	if len(pageData) > 0 {
		limit = pageData[0].PageSize
		offset = limit * pageData[0].Page
	}

	err := rb.db.WithContext(ctx).
		Where("additional_permissions @> ARRAY[?]", permission).
		Limit(limit).Offset(offset).
		Find(&users).Error

	return users, err
}

func (rb *rbacRepository) GetChildren(ctx context.Context, parentID int) ([]users.BackendUser, error) {
	var users []users.BackendUser

	err := rb.db.WithContext(ctx).
		Where("created_by = ?", parentID).
		Order("username ASC").
		Find(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get children: %w", err)
	}

	return users, nil
}

// Permission operations
func (rb *rbacRepository) GetEffectivePermissions(ctx context.Context, userID int) (rbac.PermissionSet, error)
func (rb *rbacRepository) GrantIndividualPermission(ctx context.Context, userID int, perm rbac.Permission, grantedBy int) error {

	return nil
}
func (rb *rbacRepository) RevokeIndividualPermission(ctx context.Context, userID int, perm rbac.Permission) error {
	return nil
}

func (rb *rbacRepository) GetChildrenRecursive(ctx context.Context, parentID int) ([]users.BackendUser, error) {
	var users []users.BackendUser

	// Recursive CTE to get all descendants
	err := rb.db.WithContext(ctx).Raw(`
		WITH RECURSIVE descendants AS (
			-- Anchor: immediate children
			SELECT
				id, username, email, role_id, created_by, 1 as depth
			FROM backend_users
			WHERE created_by = ?

			UNION ALL

			-- Recursive: their children
			SELECT
				u.id, u.username, u.email, u.role_id,
				u.created_by, d.depth + 1
			FROM backend_users u
			INNER JOIN descendants d ON u.created_by = d.id
			WHERE d.depth < 50  -- Safety limit to prevent infinite loops
		)
		SELECT * FROM descendants ORDER BY depth, username
	`, parentID).Scan(&users).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recursive children: %w", err)
	}

	return users, nil
}

type UserNode struct {
	User     *users.BackendUser
	Children []*UserNode
	Depth    int
}

func (rb *rbacRepository) GetChildrenCount(ctx context.Context, parentID int) (int64, error) {
	var count int64
	err := rb.db.WithContext(ctx).
		Model(&users.BackendUser{}).
		Where("created_by = ?", parentID).
		Count(&count).Error

	if err != nil {
		return 0, fmt.Errorf("failed to count children: %w", err)
	}

	return count, nil
}

func (rb *rbacRepository) GetChildrenPaginated(
	ctx context.Context,
	parentID int,
	pagination types.Pagination,
) ([]users.BackendUser, error) {

	offset := (pagination.Page - 1) * pagination.PageSize

	var gormUsers []users.BackendUser
	err := rb.db.WithContext(ctx).
		Where("created_by = ?", parentID).
		Order("created_at DESC").
		Offset(offset).
		Limit(pagination.PageSize).
		Find(&gormUsers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get paginated children: %w", err)
	}

	return gormUsers, nil
}

func (rb *rbacRepository) GetChildrenWithRoleFilter(
	ctx context.Context,
	parentID int,
	roleName string,
) ([]users.BackendUser, error) {

	var backendUsers []users.BackendUser

	err := rb.db.WithContext(ctx).
		Joins("JOIN backend_roles ON backend_users.role_id = backend_roles.id").
		Where("backend_users.created_by = ? AND backend_roles.name = ?", parentID, roleName).
		Find(&backendUsers).Error

	if err != nil {
		return nil, fmt.Errorf("failed to filter children by role: %w", err)
	}

	return backendUsers, nil
}

// Role operations
func (rb *rbacRepository) GetRole(ctx context.Context, id int) (*users.BackendRole, error) {
	role := &users.BackendRole{ID: uint(id)}
	err := rb.db.WithContext(ctx).First(role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, server.ErrDBRecordNotFound
		}
		return nil, err
	}

	return role, nil
}

func (rb *rbacRepository) GetRoleByName(ctx context.Context, name string) (*users.BackendRole, error) {
	role := &users.BackendRole{Name: name}
	err := rb.db.WithContext(ctx).First(role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, server.ErrDBRecordNotFound
		}
		return nil, err
	}

	return role, nil
}

func (rb *rbacRepository) CreateRole(ctx context.Context, role *users.BackendRole) error {
	role.ID = 0
	return rb.db.WithContext(ctx).Create(role).Error
}

func (rb *rbacRepository) GetRolePermissions(ctx context.Context, roleID int) ([]rbac.Permission, error) {
	perms := make([]rbac.Permission, 0)
	role := &users.BackendRole{ID: uint(roleID)}
	err := rb.db.WithContext(ctx).Select("allowed_permissions").First(role).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, server.ErrDBRecordNotFound
		}
		return nil, err
	}

	perms = append(perms, role.AllowedPerms...)

	return perms, nil
}

func (rb *rbacRepository) IsDescendant(ctx context.Context, ancestorID, descendantID int) (bool, error) {
	var exists bool
	err := rb.db.WithContext(ctx).Raw(`
		WITH RECURSIVE ancestors AS (
			SELECT id, created_by FROM backend_users WHERE id = ?
			UNION ALL
			SELECT u.id, u.created_by FROM backend_users u
			JOIN ancestors a ON u.id = a.created_by
		)
		SELECT EXISTS(SELECT 1 FROM ancestors WHERE id = ? AND id != ?)
	`, descendantID, ancestorID, descendantID).Scan(&exists).Error

	return exists, err
}
