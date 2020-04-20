package memory

import (
	"github.com/victornm/es-backend/pkg/auth/internal"
	"github.com/victornm/es-backend/pkg/store"
	"github.com/victornm/es-backend/pkg/store/memory"
)

type AuthUserRepository struct {
	gw *memory.UserGateway
}

func (r *AuthUserRepository) FindUserByID(id int) (*internal.User, error) {
	row, err := r.gw.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return toUserModel(row), nil
}

func (r *AuthUserRepository) FindUserByEmail(email string) (*internal.User, error) {
	row, err := r.gw.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return toUserModel(row), nil
}

func (r *AuthUserRepository) FindUserByUsername(username string) (*internal.User, error) {
	row, err := r.gw.FindUserByUsername(username)
	if err != nil {
		return nil, err
	}

	return toUserModel(row), nil
}

func (r *AuthUserRepository) CreateUser(u *internal.User) (int, error) {
	id, err := r.gw.CreateUser(toUserRow(u))
	if err != nil {
		return 0, err
	}

	return id, nil
}

func NewRepository(gw *memory.UserGateway) *AuthUserRepository {
	return &AuthUserRepository{
		gw: gw,
	}
}

func toUserModel(row *store.UserRow) *internal.User {
	return &internal.User{
		ID:             row.ID,
		Email:          row.Email,
		Username:       row.Username,
		HashedPassword: row.HashedPassword,
		FullName:       row.FullName,
		CreatedAt:      row.CreatedAt,
		IsActive:       row.IsActive,
		IsSuperAdmin:   row.IsSuperAdmin,
		ActivationKey:  row.ActivationKey,
		Provider:       row.OAuth2Provider,
	}
}

func toUserRow(model *internal.User) *store.UserRow {
	return &store.UserRow{
		ID:             model.ID,
		Email:          model.Email,
		Username:       model.Username,
		HashedPassword: model.HashedPassword,
		FullName:       model.FullName,
		CreatedAt:      model.CreatedAt,
		IsActive:       model.IsActive,
		IsSuperAdmin:   model.IsSuperAdmin,
		ActivationKey:  model.ActivationKey,
		OAuth2Provider: model.Provider,
	}
}
