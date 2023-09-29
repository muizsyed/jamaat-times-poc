package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	engine.Reload(true)

	appRevision := os.Getenv("VCS_REVISION")

	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":    "Hello, World!",
			"Revision": appRevision,
		}, "layouts/main")
	})

	app.Listen(":8081")
}

type Masjid struct {
	Name     string
	City     string
	Postcode string
}

type Prayer struct {
	Start  string
	Jamaat string
}

type Month struct {
	Name string
	Days []Day
}

type Day struct {
	Fajr    Prayer
	Zuhr    Prayer
	Asr     Prayer
	Maghrib Prayer
	Esha    Prayer
}
