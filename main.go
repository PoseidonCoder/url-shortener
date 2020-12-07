package main

import(
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/prologic/bitcask"
	"log"
)

var db *bitcask.Bitcask

type shorten struct {
	URL string
}

func resolveURL (ctx *fiber.Ctx) error {
	URL, error := db.Get([]byte(ctx.Params("URL")))
	if error != nil {
		return error
	}

	return ctx.Redirect(string(URL))
}

func shortenHandler(ctx *fiber.Ctx) error {
	set := new(shorten)

	err := ctx.BodyParser(&set)
	if err != nil {
		return err
	}

	var id string
	for {
		id = uuid.New().String()[:5]
		if _, err := db.Get([]byte(id)); err != nil {
			//If this ID doesn't already exist, break from this loop
			break
		}
	}

	err = db.Put([]byte(id), []byte(set.URL))
	if err != nil {
		return err
	}

	return ctx.SendString(id)
}

func main () {
	db, _ = bitcask.Open("db")
	app := fiber.New()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendFile("public/index.html")
	})

	app.Post("/api/v1", shortenHandler)
	app.Get("/@:URL", resolveURL)

	log.Fatal(app.Listen(":8080"))
}