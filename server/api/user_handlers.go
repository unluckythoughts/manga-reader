package api

// NOTE: User authentication and management is now handled by the auth package from go-microservice.
// To use authentication in your application:
// 1. Import: "github.com/unluckythoughts/go-microservice/tools/auth"
// 2. Create an auth.Service instance with auth.New(options)
// 3. Use the auth handlers for login, register, update user, etc.
//
// The auth package provides:
// - LoginHandler: POST /auth/login
// - GetRoleUserRegister: POST /auth/register
// - LogoutHandler: POST /auth/logout
// - GetUser: GET /auth/me
// - UpdateUserHandler: PUT /auth/me
// - ChangePasswordHandler: POST /auth/change-password
// - SendVerificationHandler: POST /auth/send-verification
// - VerifyHandler: GET /auth/verify/:token
// - UpdatePasswordHandler: POST /auth/update-password
// - Google OAuth handlers for social login
//
// For admin operations (listing all users, etc.), you can keep using the db and service methods:

import (
	"strconv"

	"github.com/unluckythoughts/go-microservice/v2/tools/web"
)

// ListUsers retrieves a paginated list of all users (admin operation)
func (api *api) ListUsers(r web.Request) (any, error) {
	pageSize := r.GetURLParam("limit")
	pageNumber := r.GetURLParam("page")

	limit, err := strconv.Atoi(pageSize)
	if err != nil || limit <= 0 {
		limit = 20 // default limit
	}
	page, err := strconv.Atoi(pageNumber)
	if err != nil || page <= 0 {
		page = 1 // default page
	}

	users, total, err := api.s.GetUsers(page, limit)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"items": users,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	}, nil
}

// GetUser retrieves a single user by ID (admin operation)
func (api *api) GetUser(r web.Request) (any, error) {
	idText := r.GetRouteParam("id")
	id, err := strconv.Atoi(idText)
	if err != nil {
		return nil, err
	}

	return api.s.GetUserByID(id)
}

// The following handlers are commented out as they should be replaced with auth package handlers:
// - CreateUser -> Use auth.GetRoleUserRegister
// - UpdateUser -> Use auth.UpdateUserHandler
// - DeleteUser -> Can be implemented as admin operation if needed
