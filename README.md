Very simple Hipchat bot library.

Docs: http://godoc.org/github.com/jonvaldes/hipbot

Example:

    package main

    import (
        "github.com/jonvaldes/hipbot"
    )

    func main() {

        bot, err := hipbot.NewBot("MYUSERID@chat.hipchat.com", "NICKNAME", "PASS")
        if err != nil {
            panic(err)
        }
        bot.JoinRoom("ROOMNAME@conf.hipchat.com")

        mux := hipbot.NewMux()

        mux.AddHandler(`^echo *`, func(b *hipbot.Bot, m hipbot.Msg) {
            b.Answer(m, "Echo "+m.Text)
        })

        bot.KeepAlive() // So Hipchat doesn't kick us out after 3 minutes

        if err := bot.Listen(mux); err != nil {
            panic(err)
        }
    }
