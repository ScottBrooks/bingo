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

// Index using a template build the index page
func (hw *Hotwire) Index(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

// Greeting process greetings using a query parameter
func (hw *Hotwire) Greeting(c echo.Context) error {
	name := c.QueryParam("person")

	return c.Render(http.StatusOK, "greeting.html", map[string]interface{}{
		"person": name,
	})
}

func (hw *Hotwire) Stats(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	ticker := time.NewTicker(3 * time.Second)

	for range ticker.C {
		buf := new(bytes.Buffer)

		err := c.Echo().Renderer.Render(buf, "stats.turbo-stream.html", hw.stats, c)
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

type cardUpdate struct {
	Msg  string
	Card *Card
}

func (hw *Hotwire) Listen() (chan cardUpdate, func()) {
	update := make(chan cardUpdate)
	hw.clients = append(hw.clients, update)

	return update, func() {
		hw.unlisten(update)
	}
}

func (hw *Hotwire) unlisten(update chan cardUpdate) {
	remove := -1
	for i, c := range hw.clients {
		if c == update {
			remove = i
		}
	}
	if remove != -1 {
		hw.clients[remove] = hw.clients[len(hw.clients)-1]
		hw.clients[len(hw.clients)-1] = nil
		hw.clients = hw.clients[:len(hw.clients)-1]
	}

}

func (hw *Hotwire) NotifyCardChanged(cardID int) {
	for _, c := range hw.clients {
		c <- cardUpdate{Msg: "hihi", Card: GetCard(cardID)}
	}
}

// Card process bingo cards
func (hw *Hotwire) Card(c echo.Context) error {
	c.SetCookie(&http.Cookie{Name: "fpid", Value: c.FormValue("fpid")})
	return c.Render(http.StatusOK, "card.html", nil)
}

// CardSocket process bingo cards
func (hw *Hotwire) CardSocket(c echo.Context) error {
	log := log.Ctx(c.Request().Context())
	fpidCookie, err := c.Cookie("fpid")
	if err != nil {
		return err
	}
	fpid, err := strconv.Atoi(fpidCookie.Value)
	if err != nil {
		return err
	}
	spots := getRandCards(fpid, 5, 5)

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	hw.stats.Cards++
	defer func() {
		hw.stats.Cards--
	}()

	buf := new(bytes.Buffer)

	err = c.Echo().Renderer.Render(buf, "cards.turbo-stream.html", map[string]interface{}{
		"spots":   spots,
		"cols":    5,
		"lastCol": 4,
		"rows":    5,
	}, c)
	if err != nil {
		log.Error().Err(err).Msg("failed to build message")
		return nil
	}

	err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
	if err != nil {
		log.Error().Err(err).Msg("send failed")
		return nil
	}

	updateChan, cancel := hw.Listen()
	defer cancel()

	for update := range updateChan {
		log.Printf("Card updated: %+v", update)
		err = c.Echo().Renderer.Render(buf, "event.turbo-stream.html", update.Card, c)
		if err != nil {
			log.Error().Err(err).Msg("failed to build message")
			continue
		}

		err = ws.WriteMessage(websocket.TextMessage, buf.Bytes())
		if err != nil {
			log.Error().Err(err).Msg("failed to write to socket")
			break
		}
	}

	return nil
}
