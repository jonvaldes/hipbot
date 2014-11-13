package hpbot

import (
	"regexp"
	"strings"

	"github.com/mattn/go-xmpp"
)

type Msg xmpp.Chat

type Bot struct {
	c        *xmpp.Client
	fullName string
}

func NewBot(userJabberId, roomNickname, password string) (*Bot, error) {
	b := &Bot{}
	b.fullName = roomNickname

	options := xmpp.Options{
		Host:          "chat.hipchat.com:5222",
		User:          userJabberId,
		Password:      password,
		NoTLS:         true,
		StartTLS:      true,
		Debug:         false,
		Session:       false,
		Status:        "",
		StatusMessage: "",
		Resource:      "bot",
	}
	var err error
	b.c, err = options.NewClient()
	return b, err
}

func (b *Bot) JoinRoom(roomJabberId string) {
	b.c.JoinMUC(roomJabberId + "/" + b.fullName)
}

func (b *Bot) SendDirect(userJabberId, message string) {
	b.c.Send(xmpp.Chat{Type: "chat", Remote: userJabberId, Text: message})
}

func (b *Bot) SendGroup(roomJabberId, message string) {
	b.c.Send(xmpp.Chat{Type: "groupchat", Remote: roomJabberId, Text: message})
}

func (b *Bot) Answer(m Msg, text string) {
	remote := m.Remote
	if m.Type == "groupchat" {
		remote = strings.Split(remote, "/")[0]
	}
	b.c.Send(xmpp.Chat{Type: m.Type, Remote: remote, Text: text})
}

func (b *Bot) Listen(l Listener) error {
	for {
		m, err := b.c.Recv()
		if err != nil {
			return err
		}

		switch v := m.(type) {
		case xmpp.Chat:
			l.HandleMsg(b, Msg{Remote: v.Remote, Type: v.Type, Text: v.Text, Other: v.Other})
		case xmpp.Presence:
			// do nothing
		}
	}
}

type HandleFunc func(b *Bot, m Msg)

type Listener interface {
	HandleMsg(b *Bot, m Msg)
}

type Mux struct {
	matches map[*regexp.Regexp]HandleFunc
}

func NewMux() *Mux {
	result := &Mux{}
	result.matches = make(map[*regexp.Regexp]HandleFunc)
	return result
}

func (m *Mux) HandleMsg(c *Bot, msg Msg) {
	for regex, callback := range m.matches {
		if regex.MatchString(msg.Text) {
			callback(c, msg)
		}
	}
}

func (m *Mux) AddHandler(regex string, f HandleFunc) {
	r := regexp.MustCompile(regex)
	m.matches[r] = f
}
