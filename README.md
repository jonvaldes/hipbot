Very simple Hipchat bot library.

Docs: http://godoc.org/github.com/jonvaldes/hipchat_bot

Example:

    package main

    import (
        "github.com/jonvaldes/hipchat_bot"
    )

    func main() {

        bot, err := hpbot.NewBot("MYUSERID@chat.hipchat.com", "NICKNAME", "PASS")
        if err != nil {
            panic(err)
        }
        bot.JoinRoom("ROOMNAME@conf.hipchat.com")

        mux := hpbot.NewMux()

        mux.AddHandler(`^echo *`, func(b *hpbot.Bot, m hpbot.Msg) {
            b.Answer(m, "Echo "+m.Text)
        })

        if err := bot.Listen(mux); err != nil {
            panic(err)
        }
    }
