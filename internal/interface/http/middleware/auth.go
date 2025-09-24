package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type AuthMiddleware struct {
	authService domain.AuthService
	logger      *zap.Logger
}

func NewAuthMiddleware(authService domain.AuthService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
		logger:      logger,
	}
}

// RequireAuth validates JWT token and sets user claims in context
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			m.logger.Warn("Missing Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			m.logger.Warn("Invalid Authorization header format")
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := m.authService.ValidateAccessToken(token)
		if err != nil {
			m.logger.Warn("Invalid access token", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Set("org_id", claims.OrgID)
		c.Set("claims", claims)
		
		m.logger.Info("User authenticated and claims set",
			zap.String("user_id", claims.UserID.String()),
			zap.String("user_role", string(claims.Role)),
			zap.Any("org_id", claims.OrgID))

		c.Next()
	}
}

// RequireRole checks if user has the required role
func (m *AuthMiddleware) RequireRole(requiredRoles ...domain.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			m.logger.Error("User role not found in context")
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		var role domain.UserRole
		switch v := userRole.(type) {
		case domain.UserRole:
			role = v
		case string:
			role = domain.UserRole(v)
		default:
			m.logger.Error("Invalid user role type in context", 
				zap.Any("user_role", userRole),
				zap.String("user_role_type", fmt.Sprintf("%T", userRole)))
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		m.logger.Info("Checking role permissions",
			zap.String("user_role", string(role)),
			zap.Any("required_roles", requiredRoles))
		
		for _, requiredRole := range requiredRoles {
			if role == requiredRole {
				hasRole = true
				m.logger.Info("Role match found", 
					zap.String("user_role", string(role)),
					zap.String("matched_role", string(requiredRole)))
				break
			}
		}

		if !hasRole {
			m.logger.Warn("Insufficient permissions",
				zap.String("user_role", string(role)),
				zap.Any("required_roles", requiredRoles))
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireOrgScope ensures user can only access resources within their organization
func (m *AuthMiddleware) RequireOrgScope() gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		role := userRole.(domain.UserRole)

		// Admin can access everything
		if role == domain.UserRoleAdmin {
			c.Next()
			return
		}

		// Get user's org_id from claims
		orgID, exists := c.Get("org_id")
		if !exists || orgID == nil {
			// Users without org_id can only access their own resources
			c.Set("scope_org_id", nil)
			c.Next()
			return
		}

		// Set org scope for filtering
		c.Set("scope_org_id", orgID.(*uuid.UUID))
		c.Next()
	}
}

// Helper functions to extract user info from context
func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	id, ok := userID.(uuid.UUID)
	return id, ok
}

func GetUserRole(c *gin.Context) (domain.UserRole, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := userRole.(domain.UserRole)
	return role, ok
}

func GetOrgID(c *gin.Context) (*uuid.UUID, bool) {
	orgID, exists := c.Get("org_id")
	if !exists {
		return nil, false
	}
	id, ok := orgID.(*uuid.UUID)
	return id, ok
}

func GetScopeOrgID(c *gin.Context) (*uuid.UUID, bool) {
	orgID, exists := c.Get("scope_org_id")
	if !exists {
		return nil, false
	}
	if orgID == nil {
		return nil, true
	}
	id, ok := orgID.(*uuid.UUID)
	return id, ok
}

