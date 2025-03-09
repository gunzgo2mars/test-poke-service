package info

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/gunzgo2mars/test-poke-service/app/internal/constants"
	"github.com/gunzgo2mars/test-poke-service/app/internal/core/service/pokemon"
	"github.com/gunzgo2mars/test-poke-service/app/pkg/response"
)

type IInfoHandler interface {
	GetPokemonByName(c echo.Context) error
	GetPokemonAbilities(c echo.Context) error
	Random(c echo.Context) error
}

type infoHandler struct {
	pokemonSvc pokemon.IPokemonService
}

func New(pokemonSvc pokemon.IPokemonService) IInfoHandler {
	return &infoHandler{
		pokemonSvc: pokemonSvc,
	}
}

func (h *infoHandler) GetPokemonByName(c echo.Context) error {
	ctx := c.Request().Context()
	res := response.NewJSONMessage()

	resp, err := h.pokemonSvc.RequestPokemonInfo(ctx, c.Param("name"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			res.AddCode(constants.ERROR_INTERNAL_CODE).AddMessage(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.BuildMessageWithData(
			resp,
			res.AddCode(constants.SUCCESS_CODE).AddMessage(constants.SUCCESS_MSG),
		),
	)
}

func (h *infoHandler) GetPokemonAbilities(c echo.Context) error {
	ctx := c.Request().Context()
	res := response.NewJSONMessage()

	resp, err := h.pokemonSvc.RequestPokemonAbilities(ctx, c.Param("name"))
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			res.AddCode(constants.ERROR_INTERNAL_CODE).AddMessage(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.BuildMessageWithData(
			resp,
			res.AddCode(constants.SUCCESS_CODE).AddMessage(constants.SUCCESS_MSG),
		),
	)
}

func (h *infoHandler) Random(c echo.Context) error {
	ctx := c.Request().Context()
	res := response.NewJSONMessage()

	resp, err := h.pokemonSvc.RandomPokemon(ctx)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			res.AddCode(constants.ERROR_INTERNAL_CODE).AddMessage(err.Error()),
		)
	}

	return c.JSON(
		http.StatusOK,
		response.BuildMessageWithData(
			resp,
			res.AddCode(constants.SUCCESS_CODE).AddMessage(constants.SUCCESS_MSG),
		),
	)
}
