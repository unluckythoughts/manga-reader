# Authentication Integration Guide

The book-reader application now uses the `auth` package from `go-microservice` for user authentication and management.

## What Changed

1. **User Model**: The local `User` model has been removed. We now use `auth.User` from `github.com/unluckythoughts/go-microservice/v2/tools/auth`.

2. **Database Schema**: The user table now includes authentication fields:
   - `email` (unique, required)
   - `mobile` (unique, optional)
   - `password` (hashed)
   - `email_verified` / `mobile_verified`
   - `role` (integer, default: 1)
   - `verify_token` (for email/mobile verification)
   - `google_id` / `google_avatar` (for OAuth)

3. **User Endpoints**: 
   - Removed: `POST /api/v1/users` (create)
   - Removed: `PUT /api/v1/users/:id` (update)
   - Removed: `DELETE /api/v1/users/:id` (delete)
   - Kept: `GET /api/v1/users` (list all - admin)
   - Kept: `GET /api/v1/users/:id` (get by ID - admin)

## How to Integrate Auth Service

### Step 1: Initialize the Auth Service

In your main application file, create an auth service instance:

```go
import (
    "github.com/unluckythoughts/go-microservice/v2/tools/auth"
    "gorm.io/gorm"
)

func setupAuth(db *gorm.DB, router web.Router) *auth.Service {
    authService := auth.New(auth.Options{
        DB:                db,
        Logger:           logger,
        JwtKey:           os.Getenv("AUTH_JWT_KEY"), // Required: secret key for JWT
        TokenValidInHours: 24, // Token validity duration
        IgnoreRoutes: []string{
            "/api/v1/auth/login",
            "/api/v1/auth/register",
        },
        UserRoles: map[auth.UserRole]string{
            1:  "user",
            99: "admin",
        },
        DefaultMobileCountryCode: "+1",
        GoogleOauth: &auth.GoogleOauth{
            ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
            ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
            RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        },
    })
    
    // Run migrations for auth tables
    authService.Migrate()
    
    return authService
}
```

### Step 2: Register Auth Routes

Add authentication routes to your router:

```go
func registerAuthRoutes(router web.Router, authService *auth.Service) {
    // Authentication endpoints
    router.POST("/api/v1/auth/register", authService.GetRoleUserRegister(1)) // role 1 = user
    router.POST("/api/v1/auth/login", authService.LoginHandler)
    router.POST("/api/v1/auth/logout", authService.LogoutHandler)
    
    // User management (requires authentication)
    router.GET("/api/v1/auth/me", authService.GetUser)
    router.PUT("/api/v1/auth/me", authService.UpdateUserHandler)
    router.POST("/api/v1/auth/change-password", authService.ChangePasswordHandler)
    
    // Email/Mobile verification
    router.POST("/api/v1/auth/send-verification", authService.SendVerificationHandler)
    router.GET("/api/v1/auth/verify/:token", authService.VerifyHandler)
    router.POST("/api/v1/auth/update-password", authService.UpdatePasswordHandler)
    
    // Google OAuth (optional)
    router.GET("/api/v1/auth/google/login", authService.GoogleLoginHandler)
    router.GET("/api/v1/auth/google/callback", authService.GoogleCallbackHandler)
}
```

### Step 3: Add Authentication Middleware

Use the auth middleware to protect routes:

```go
// In your router setup
func setupRouter(router web.Router, authService *auth.Service) {
    // Public routes
    router.GET("/api/v1/books", api.ListBooks)
    router.GET("/api/v1/books/:id", api.GetBook)
    
    // Protected routes - require authentication
    router.Use(authService.AuthenticationMiddleware())
    router.POST("/api/v1/favorites", api.CreateFavorite)
    router.PUT("/api/v1/favorites/:id", api.UpdateFavorite)
    
    // Admin-only routes - require admin role
    router.Use(authService.AuthorizationMiddleware(99)) // 99 = admin role
    router.GET("/api/v1/users", api.ListUsers)
}
```

### Step 4: Access Authenticated User

In your handlers, you can get the authenticated user:

```go
func (api *api) CreateFavorite(r web.Request) (any, error) {
    // Get the authenticated user
    user, err := auth.GetAuthenticatedUser(r)
    if err != nil {
        return nil, err
    }
    
    // Use user.ID for creating favorite
    favorite := &models.Favorite{
        UserID: int(user.ID),
        BookID: bookID,
    }
    // ... rest of the logic
}
```

## Available Auth Handlers

The auth package provides these handlers out of the box:

- **LoginHandler**: Authenticate with email/mobile and password
- **GetRoleUserRegister(role)**: Register a new user with specified role
- **LogoutHandler**: Logout (clear session/token)
- **GetUser**: Get current authenticated user
- **UpdateUserHandler**: Update current user's profile
- **ChangePasswordHandler**: Change password (requires old password)
- **SendVerificationHandler**: Send verification email/SMS
- **VerifyHandler**: Verify email/mobile with token
- **UpdatePasswordHandler**: Reset password with verification token
- **GoogleLoginHandler**: Initiate Google OAuth flow
- **GoogleCallbackHandler**: Handle Google OAuth callback

## Environment Variables

Set these environment variables:

```bash
AUTH_JWT_KEY=your-secret-jwt-key-here
AUTH_TOKEN_VALID=24  # hours
AUTH_DEFAULT_MOBILE_COUNTRY_CODE=+1
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

## Request Examples

### Register a new user

```bash
POST /api/v1/auth/register
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!",
  "mobile": "+1234567890"
}
```

### Login

```bash
POST /api/v1/auth/login
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}

# Response
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### Get current user

```bash
GET /api/v1/auth/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```

### Update profile

```bash
PUT /api/v1/auth/me
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
{
  "name": "John Smith",
  "email": "john.smith@example.com"
}
```

## Password Requirements

Passwords must:
- Be at least 10 characters long
- Contain at least one uppercase letter
- Contain at least one lowercase letter
- Contain at least one digit
- Contain at least one special character

## User Roles

Default roles:
- `1`: Regular user
- `99`: Admin

Admins can access all user resources. You can define custom roles as needed.
