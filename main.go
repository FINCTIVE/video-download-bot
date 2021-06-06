package main

import (
	"context"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"mvdan.cc/xurls/v2"
	"os"
	"os/exec"
)

func main() {
	Launch(func(bot *tb.Bot) {
		var stopCurrentDownload context.CancelFunc
		bot.Handle(tb.OnText, func(m *tb.Message) {
			pass := CheckUser(m.Sender)
			if !pass {
				return
			}

			var urls []string

			// find urls from formatted urls
			for i := range m.Entities {
				urlFromFormat := m.Entities[i].URL
				if len(urlFromFormat) != 0 {
					urls = append(urls, urlFromFormat)
				}
			}

			// find urls from text
			if len(urls) == 0 {
				rx := xurls.Strict()
				urls = rx.FindAllString(m.Text, -1)
			}

			if len(urls) == 0 {
				return
			}

			log.Println("verbose: download list: ", urls)

			// external dependence, get this scripts from:
			// https://github.com/FINCTIVE/download-videos
			dlPath := "/root/download-videos.sh"
			// cancel function
			var ctx context.Context
			ctx, stopCurrentDownload = context.WithCancel(context.Background())
			cmd := exec.CommandContext(ctx, dlPath, urls...)
			cmd.Env = os.Environ()
			cmd.Dir = "/root/downloads"
			done := RunCommand(m.Sender, cmd, &tb.SendOptions{ReplyTo: m})
			go func() {
				<-done
				Send(m.Sender, "done!", &tb.SendOptions{ReplyTo: m})
			}()
		})

		bot.Handle("/stop", func(m *tb.Message) {
			if stopCurrentDownload != nil {
				stopCurrentDownload()
			}
		})

		bot.Handle("/ping", func(m *tb.Message) {
			var ctx context.Context
			ctx, stopCurrentDownload = context.WithCancel(context.Background())
			cmd := exec.CommandContext(ctx, "ping", "baidu.com")
			_ = RunCommand(m.Sender, cmd, &tb.SendOptions{ReplyTo: m})
		})

		_ = bot.SetCommands([]tb.Command{
			{"/stop", "stop the running task"},
			{"/ping", "ping baidu"},
		})
	})
}
