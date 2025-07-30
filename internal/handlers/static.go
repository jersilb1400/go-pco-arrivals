package handlers

import (
	"github.com/gofiber/fiber/v2"
)

type StaticHandler struct{}

func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (h *StaticHandler) ServeIndex(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeAdmin(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeBillboard(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeLocationBillboard(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeLogin(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeOffline(c *fiber.Ctx) error {
	return c.SendFile("./web/static/index.html")
}

func (h *StaticHandler) ServeManifest(c *fiber.Ctx) error {
	return c.SendFile("./web/static/manifest.json")
}

func (h *StaticHandler) ServeServiceWorker(c *fiber.Ctx) error {
	return c.SendFile("./web/static/sw.js")
}

func (h *StaticHandler) Handle404(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./web/static/index.html")
}
