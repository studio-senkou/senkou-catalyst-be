package main

import (
	"fmt"
	"os"
	"senkou-catalyst-be/database/seeder"
	"senkou-catalyst-be/platform/config"
	"strings"

	"gorm.io/gorm"
)

type SeederClosure func(*gorm.DB) error

type SeederContext struct {
	seeders map[string]SeederClosure
}

func NewSeederContext() *SeederContext {
	return &SeederContext{
		seeders: make(map[string]SeederClosure),
	}
}

func (ctx *SeederContext) Register(name string, fn SeederClosure) {
	ctx.seeders[name] = fn
}

func (ctx *SeederContext) Run(name string, db *gorm.DB) error {
	fn, ok := ctx.seeders[name]
	if !ok {
		return fmt.Errorf("unknown seeder: %s", name)
	}

	return fn(db)
}

func main() {
	db := config.GetDB()

	ctx := NewSeederContext()
	register(ctx)

	seeder := "all"
	if len(os.Args) > 1 {
		seeder = strings.TrimSuffix(os.Args[1], ".go")
	}

	if err := ctx.Run(seeder, db); err != nil {
		panic(err)
	}
	fmt.Println("Database seeded successfully")
}

func register(c *SeederContext) {
	c.Register("user_seeder", seeder.SeedUsers)
	c.Register("subscription_seeder", seeder.SeedSubscriptions)

	seederPool := []SeederClosure{
		seeder.SeedSubscriptions,
		seeder.SeedUsers,
	}

	c.Register("all", func(db *gorm.DB) error {
		for _, seed := range seederPool {
			if err := seed(db); err != nil {
				panic(err)
			}
		}
		return nil
	})
}
