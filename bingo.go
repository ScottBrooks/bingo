package bingo

import (
	"bytes"
	"html/template"
	"log"
	"math/rand"
	"sort"
	"time"
)

type Card struct {
	Id        int
	Card      string
	CSSClass  string
	Timestamp time.Time
}

func (c *Card) SetHappened(happened bool) {
	if happened {
		c.Timestamp = time.Now()
		c.CSSClass = "uk-card-secondary"
	} else {
		c.Timestamp = time.Time{}
		c.CSSClass = "uk-card-default"
	}
}

var eventTemplates = []string{
	"Win!", "Top 5", "Top 10", "Josh landing", "{{.}} - complains about a cheater",
	"Find a cheater", "Pull a Team Josh", "Josh Math", "{{.}} - does a Pancake Pete",
	"{{.}} - pulls a frage", "{{.}} dies", "{{.}} - most kills at end of a match", "{{.}} - most deaths at end of match", "Strat gets left behind",
	"Frage runs off on his own", "{{.}} - complains about their FPS", "{{.}} - gets stratted", "Look left... I mean right", "Look right... I mean left",
}

var players = []string{
	"Pete", "Gil", "Shaun", "Andrew", "Scott",
}

var allEvents = []Card{}

var eventMap map[int]*Card

func init() {
	uniqEvents := map[string]bool{}
	for _, tpl := range eventTemplates {
		t := template.Must(template.New("").Parse(tpl))
		for _, p := range players {
			var evald bytes.Buffer
			err := t.Execute(&evald, p)
			if err != nil {
				log.Fatal(err)
			}
			uniqEvents[evald.String()] = true
		}
	}

	i := 0
	for evt := range uniqEvents {
		i++
		allEvents = append(allEvents, Card{Id: i, Card: evt, CSSClass: "uk-card-default"})
	}
	sort.Slice(allEvents, func(i int, j int) bool {
		return allEvents[i].Card < allEvents[j].Card
	})

	eventMap = make(map[int]*Card)
	for i, c := range allEvents {
		eventMap[c.Id] = &allEvents[i]
	}

}

func GetCard(i int) *Card {
	return eventMap[i]
}

func getEvents() []Card {
	out := make([]Card, len(allEvents))
	copy(out, allEvents)
	return out
}

func getRandCards(seed int, rows int, cols int) []Card {
	out := make([]Card, len(allEvents))
	copy(out, allEvents)
	rand.Seed(int64(seed))
	rand.Shuffle(len(out), func(i, j int) { out[i], out[j] = out[j], out[i] })

	return out[0:25]
}
