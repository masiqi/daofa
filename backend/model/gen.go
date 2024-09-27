package model

import (
	"gorm.io/gen"
	"gorm.io/gorm"
)

func GenerateModel(db *gorm.DB) {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("subjects"),
		g.GenerateModel("knowledge_points"),
		g.GenerateModel("exercises"),
		g.GenerateModel("questions"),
	)

	g.Execute()
}