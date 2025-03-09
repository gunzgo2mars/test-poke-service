package ext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/redis/go-redis/v9"

	"github.com/gunzgo2mars/test-poke-service/app/internal/core/model"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/configurer"
)

type IExtRepository interface {
	// Http request
	GetPokemonInfo(
		ctx context.Context,
		name string,
	) (*model.PokemonSchema, error)
	GetPokemonAbilities(
		ctx context.Context,
		name string,
	) (*model.PokemonAbilitiesSchema, error)

	// cache
	CacheSetPokemonData(
		ctx context.Context,
		stampData *model.PokemonSchema,
	) error
	CacheGetPokemonData(
		ctx context.Context,
		name string,
	) (*model.PokemonSchema, error)
	CacheSetPokemonAbilities(
		ctx context.Context,
		name string,
		stampData *model.PokemonAbilitiesSchema,
	) error
	CacheGetPokemonAbilities(
		ctx context.Context,
		name string,
	) (*model.PokemonAbilitiesSchema, error)
}

type extRepository struct {
	conf    *configurer.AppConfig
	request *resty.Client
	cache   *redis.Client
}

func New(resty *resty.Client, cache *redis.Client, conf *configurer.AppConfig) IExtRepository {
	r := resty.
		SetBaseURL(conf.Http.PokeAPI.BaseUrl).
		SetTimeout(conf.Http.PokeAPI.Timeout).
		SetRetryCount(conf.Http.PokeAPI.RetryCount).
		SetRetryWaitTime(conf.Http.PokeAPI.RetryWaitTime).
		SetRetryMaxWaitTime(conf.Http.PokeAPI.RetryMaxWaitTime)

	return &extRepository{
		request: r,
		conf:    conf,
		cache:   cache,
	}
}

func (r *extRepository) GetPokemonInfo(
	ctx context.Context,
	name string,
) (*model.PokemonSchema, error) {
	res, err := r.request.R().SetHeaders(
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	).Get(fmt.Sprintf("/%s", name))
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, errors.New("pokemon not found.")
	}

	var resp *model.PokemonSchema
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *extRepository) GetPokemonAbilities(
	ctx context.Context,
	name string,
) (*model.PokemonAbilitiesSchema, error) {
	res, err := r.request.R().SetHeaders(
		map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
	).Get(fmt.Sprintf("/%s", name))
	if err != nil {
		return nil, err
	}

	if res.StatusCode() == 404 {
		return nil, errors.New("pokemon ability not found.")
	}

	var resp *model.PokemonAbilitiesSchema
	if err := json.Unmarshal(res.Body(), &resp); err != nil {
		return nil, err
	}

	return resp, nil
}
