package seeders

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type ISeederRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) ISeederRegistry {
	return &Registry{db: db}
}

func (s *Registry) Run() {
	// Run all seeders here
	RunRoleSeeder(s.db)
	RunUserSeeder(s.db)
}
