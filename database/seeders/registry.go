package seeders

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type ISeederRegistry interface {
	run()
}

func NewSeederRegistry(db *gorm.DB) ISeederRegistry {
	return &Registry{db: db}
}

func (s *Registry) run() {
	// Run all seeders here
	RunRoleSeeder(s.db)
	RunUserSeeder(s.db)
}
