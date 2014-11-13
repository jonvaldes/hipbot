package hpbot

import (
	"regexp"
	"strings"

	"github.com/mattn/go-xmpp"
)

// Msg represents a message send or received
type Msg struct {
	Remote string
	Type   string
	Text   string
}

// Bot is the main object used to interact with Hipchat
type Bot struct {
	c        *xmpp.Client
	fullName string
}

// NewBot creates a new Bot instance. Please note that Jabber ID and
// nickname must match the ones in the Hipchat acoount settings
func NewBot(userJabberID, nickname, password string) (*Bot, error) {
	b := &Bot{}
	b.fullName = nickname

	options := xmpp.Options{
		Host:          "chat.hipchat.com:5222",
		User:          userJabberID,
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

// JoinRoom tries to make the bot join a specific room. Room ID must
// match the one reported in the hipchat website
func (b *Bot) JoinRoom(roomJabberID string) {
	b.c.JoinMUC(roomJabberID + "/" + b.fullName)
}

// SendDirect sends a message directly to one user
func (b *Bot) SendDirect(userJabberID, message string) {
	b.c.Send(xmpp.Chat{Type: "chat", Remote: userJabberID, Text: message})
}

// SendRoom sends a message to a room
func (b *Bot) SendRoom(roomJabberID, message string) {
	b.c.Send(xmpp.Chat{Type: "groupchat", Remote: roomJabberID, Text: message})
}

// Answer sends a message back to a room if the original message was a chatgroup,
// or directly to a user otherwise
func (b *Bot) Answer(m Msg, text string) {
	remote := m.Remote
	if m.Type == "groupchat" {
		remote = strings.Split(remote, "/")[0]
	}
	b.c.Send(xmpp.Chat{Type: m.Type, Remote: remote, Text: text})
}

// Listener is the interface used to handle incoming messages
type Listener interface {
	HandleMsg(b *Bot, m Msg)
}

// Listen waits forever for incoming messages and relays them to
// the passed-in listener
func (b *Bot) Listen(l Listener) error {
	for {
		m, err := b.c.Recv()
		if err != nil {
			return err
		}

		switch v := m.(type) {
		case xmpp.Chat:
			l.HandleMsg(b, Msg{Remote: v.Remote, Type: v.Type, Text: v.Text})
		case xmpp.Presence:
			// do nothing
		}
	}
}

// HandleFunc is the type of functions
type HandleFunc func(b *Bot, m Msg)

// Mux works pretty similarly to gorilla-mux and similar constructs.
// You can register pairs of regex-function, so each time an incoming
// message matches a regex, the relevant function will be called
type Mux struct {
	matches map[*regexp.Regexp]HandleFunc
}

// NewMux creates a new Mux instance
func NewMux() *Mux {
	result := &Mux{}
	result.matches = make(map[*regexp.Regexp]HandleFunc)
	return result
}

// HandleMsg receives the incoming messages
func (m *Mux) HandleMsg(c *Bot, msg Msg) {
	for regex, callback := range m.matches {
		if regex.MatchString(msg.Text) {
			callback(c, msg)
		}
	}
}

// AddHandler registers a function to be called when an
// incoming message matches the regex
func (m *Mux) AddHandler(regex string, f HandleFunc) {
	r := regexp.MustCompile(regex)
	m.matches[r] = f
}
