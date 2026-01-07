package middleware

import (
	"github.com/carllix/matchaciee-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func RoleMiddleware(allowedRoles ...models.UserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get role from context
		roleValue := c.Locals("role")
		if roleValue == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Access denied: no role found",
			})
		}

		userRole, ok := roleValue.(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"success": false,
				"error":   "Access denied: invalid role type",
			})
		}

		// Check if user role is in allowed roles
		for _, allowedRole := range allowedRoles {
			if userRole == string(allowedRole) {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"success": false,
			"error":   "Access denied: insufficient permissions",
		})
	}
}
