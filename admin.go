package bingo

import (
	"bytes"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// Admin displays all the bingo cards
func (hw *Hotwire) Admin(c echo.Context) error {
	state := struct {
		Events []Card
	}{
		getEvents(),
	}

	return c.Render(http.StatusOK, "admin.html", state)
}

// AdminCards show current bingo card stats
func (hw *Hotwire) AdminCards(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		buf := new(bytes.Buffer)

		err := c.Echo().Renderer.Render(buf, "admin.turbo-stream.html", hw.stats, c)
		if err != nil {
			log.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to build message")
			break
		}

		err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			log.Ctx(c.Request().Context()).Error().Err(err).Msg("send failed")
			break
		}
	}

	return nil
}

// AdminEvents show current bingo card stats
func (hw *Hotwire) AdminEvents(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	state := struct {
		Events []Card
	}{
		getEvents(),
	}
	buf := new(bytes.Buffer)
	err = c.Echo().Renderer.Render(buf, "admin-events.turbo-stream.html", state, c)
	if err != nil {
		log.Ctx(c.Request().Context()).Error().Err(err).Msg("failed to build message")
		return err
	}

	err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
	if err != nil {
		log.Ctx(c.Request().Context()).Error().Err(err).Msg("send failed")
		return err
	}

	return nil
}

// AdminEvent updates that an event happened internally, and returns the timestamp
func (hw *Hotwire) AdminEvent(c echo.Context) error {
	log := log.Ctx(c.Request().Context())
	c.Response().Header().Set(echo.HeaderContentType, turboStreamMedia)
	values, err := c.FormParams()
	if err != nil {
		return err
	}

	log.Info().Interface("values", values).Msg("omg")
	cardId, err := strconv.Atoi(c.FormValue("id"))
	if err != nil {
		return err
	}
	card := GetCard(cardId)

	if c.FormValue("woops") == "true" {
		card.SetHappened(false)
	} else {
		card.SetHappened(true)
	}

	hw.NotifyCardChanged(cardId)

	c.Render(http.StatusOK, "admin-event.turbo-stream.html", card)
	return nil
}
