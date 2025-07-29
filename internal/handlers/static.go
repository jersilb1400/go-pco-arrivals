package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type StaticHandler struct{}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (h *StaticHandler) ServeIndex(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Index page"})
}

func (h *StaticHandler) ServeAdmin(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Admin page"})
}

func (h *StaticHandler) ServeBillboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Billboard page"})
}

func (h *StaticHandler) ServeLocationBillboard(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Location billboard page"})
}

func (h *StaticHandler) ServeLogin(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Login page"})
}

func (h *StaticHandler) ServeOffline(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Offline page"})
}

func (h *StaticHandler) ServeManifest(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Manifest"})
}

func (h *StaticHandler) ServeServiceWorker(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "Service worker"})
}

func (h *StaticHandler) Handle404(c *fiber.Ctx) error {
	return c.Status(404).JSON(fiber.Map{"error": "Not found"})
}
