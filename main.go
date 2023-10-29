/*
Author: Vlad Ursolov
Project: LineClub.RU
Subproject: R&D
2023. All Rights Reserved
*/


package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/mattn/go-sqlite3"
)

const serverKeyPolicy = "EDLJDhHD9mzbloENzA9pjdteAhgE4VdClUWR2SPP7tfkkQFHb9APBW4STOTTlM7S"

type Block struct {
	Index     int
	Timestamp string
	Data      string
	Hash      string
	PrevHash  string
}

func createBlock(c *fiber.Ctx, db *sql.DB) error {
	serverKey := c.Query("server_key")
	if serverKey == "" || serverKey != serverKeyPolicy {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	data := c.Query("data")
	if data == "" {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	blocks := getBlocksFromDB(db)

	var prevHash string
	if len(blocks) > 0 {
		prevHash = blocks[len(blocks)-1].Hash
	}

	newBlock := Block{
		Index:     len(blocks),
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  prevHash,
	}

	newBlock.Hash = calculateHash(newBlock)

	insertBlockQuery := `INSERT INTO blocks (timestamp, data, hash, prev_hash) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(insertBlockQuery, newBlock.Timestamp, newBlock.Data, newBlock.Hash, newBlock.PrevHash)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(newBlock)
}

func calculateHash(b Block) string {
	record := fmt.Sprintf("%d%s%s%s", b.Index, b.Timestamp, b.Data, b.PrevHash)
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func getBlocksFromDB(db *sql.DB) []Block {
	query := "SELECT * FROM blocks"
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var blocks []Block
	for rows.Next() {
		var b Block
		err := rows.Scan(&b.Index, &b.Timestamp, &b.Data, &b.Hash, &b.PrevHash)
		if err != nil {
			log.Fatal(err)
		}
		blocks = append(blocks, b)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return blocks
}

func getBlockByHashFromDB(db *sql.DB, hash string) (Block, error) {
	query := "SELECT * FROM blocks WHERE hash = ?"
	row := db.QueryRow(query, hash)

	var b Block
	err := row.Scan(&b.Index, &b.Timestamp, &b.Data, &b.Hash, &b.PrevHash)
	if err != nil {
		if err == sql.ErrNoRows {
			return Block{}, fiber.NewError(fiber.StatusNotFound, "Block not found")
		}
		log.Fatal(err)
	}

	return b, nil
}

func main() {
	// Для анонимного подключения к БД, используйте эту строку
	//db, err := sql.Open("sqlite3", "blocks.db")

	// Для подключения к БД с логином и паролем используйте эту строку
	db, err := sql.Open("sqlite3", "blocks.db?_auth&_auth_user=dbu_gochain&_auth_pass=any123_password&here!")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableQuery := `CREATE TABLE IF NOT EXISTS blocks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        timestamp TEXT,
        data TEXT,
        hash TEXT,
        prev_hash TEXT
    )`
	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()

	app.Post("/blocks", func(c *fiber.Ctx) error {
		return createBlock(c, db)
	})

	app.Get("/blocks", func(c *fiber.Ctx) error {
		blocks := getBlocksFromDB(db)
		return c.JSON(blocks)
	})

	app.Get("/blocks/:hash", func(c *fiber.Ctx) error {
		hash := c.Params("hash")
		block, err := getBlockByHashFromDB(db, hash)
		if err != nil {
			return c.SendStatus(err.(*fiber.Error).Code)
		}
		return c.JSON(block)
	})

	app.Get("/check-integrity", func(c *fiber.Ctx) error {
		blocks := getBlocksFromDB(db)
		for i := 1; i < len(blocks); i++ {
			prevBlock := blocks[i-1]
			currBlock := blocks[i]
			if currBlock.PrevHash != prevBlock.Hash {
				return c.SendStatus(fiber.StatusBadRequest)
			}
		}
		return c.SendString("Целостность блокчейна не нарушена")
	})

	log.Fatal(app.Listen(":3000"))
}
