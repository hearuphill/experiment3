package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/glebarez/go-sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	pcmsClient "pcms3/client"
)

var db *sql.DB

func init() {
	// 初始化数据库，建表
	var err error
	db, err = sql.Open("sqlite", "./pcms.db")
	if err != nil {
		log.Fatal(err)
	}

	var sqlite_version string
	db.QueryRow("select sqlite_version()").Scan(&sqlite_version)
	fmt.Println(sqlite_version)

	db.Exec(`CREATE TABLE IF NOT EXISTS pcms (name TEXT PRIMARY KEY, address TEXT, phone TEXT);`)
	db.Exec(`INSERT INTO pcms values("name_test", "address_test", "123456789")`)

}

type Contact struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func main() {
	db.QueryRow("select sqlite_version()")

	app := fiber.New()

	app.Use("/", filesystem.New(filesystem.Config{
		Root: http.FS(pcmsClient.StaticFS),
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World 👋!")
	})

	api := app.Group("/api")
	api.Get("/get", func(c *fiber.Ctx) error {
		rows, _ := db.Query(`SELECT * FROM pcms where name NOT LIKE ""`)
		defer rows.Close()
		contacts := make([]Contact, 0)
		for rows.Next() {
			contact := Contact{}
			rows.Scan(&contact.Name, &contact.Address, &contact.Phone)
			contacts = append(contacts, contact)
		}
		c.Status(200)
		return c.JSON(contacts)
	})

	api.Post("/post", func(c *fiber.Ctx) error {
		contact := new(Contact)
		err := c.BodyParser(contact)
		if err != nil {
			c.Status(500)
			c.WriteString(err.Error())
			return err
		}

		rst, err := db.Exec(`INSERT INTO pcms values(?, ?, ?)`, contact.Name, contact.Address, contact.Phone)
		if err != nil {
			c.Status(500)
			c.WriteString(err.Error())
			return err
		}

		a, err := rst.RowsAffected()
		c.WriteString(strconv.Itoa(int(a)))
		return err
	})

	api.Delete("/delete", func(c *fiber.Ctx) error {
		contact := new(Contact)
		err := c.BodyParser(contact)
		if err != nil {
			c.Status(500)
			c.WriteString(err.Error())
			return err
		}

		rst, err := db.Exec(`DELETE FROM pcms where name=? AND address=? AND phone=?`, contact.Name, contact.Address, contact.Phone)
		if err != nil {
			c.Status(500)
			c.WriteString(err.Error())
			return err
		}

		a, err := rst.RowsAffected()
		c.WriteString(strconv.Itoa(int(a)))
		return err
	})

	app.Listen(":3000")
}
