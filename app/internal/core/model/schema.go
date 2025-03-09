package model

import "time"

type UserSchema struct {
	ID       int       `gorm:"column:id;primaryKey"`
	UUID     string    `gorm:"column:uuid"`
	Username string    `gorm:"column:username"`
	Password string    `gorm:"column:password"`
	CreateAt time.Time `gorm:"column:create_at"`
	UpdateAt time.Time `gorm:"column:update_at"`
	IsActive bool      `gorm:"column:is_active"`
}

type PokemonAbilitiesSchema struct {
	Abilities []struct {
		Ability struct {
			Name string `json:"name"`
		} `json:"ability"`
	} `json:"abilities"`
}

type PokemonSchema struct {
	Name      string `json:"name"      redis:"name"`
	Abilities []struct {
		Ability struct {
			Name string `json:"name" redis:"name"`
		} `json:"ability" redis:"ability"`
	} `json:"abilities" redis:"abilities"`
	Types []struct {
		Type struct {
			Name string `json:"name" redis:"name"`
		} `json:"type" redis:"type"`
	} `json:"types"     redis:"types"`
	Weight int `json:"weight"    redis:"weight"`
}

func (*UserSchema) TableName() string {
	return "tb_users"
}
