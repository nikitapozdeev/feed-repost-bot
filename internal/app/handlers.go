package app

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/nikitapozdeev/feed-repost-bot/internal/storage"
	tele "gopkg.in/telebot.v3"
)

// handlers set bot handlers
func (a *app) handlers() {
	a.bot.Handle("/start", a.handlerStart)
	a.bot.Handle("/ping", a.handlerPing)
	a.bot.Handle("/add", a.handlerAdd)
	a.bot.Handle("/remove", a.handlerRemove)
	a.bot.Handle(tele.OnText, a.handlerUnknown)
}

// handlerStart handles bot /start command
func (a *app) handlerStart(c tele.Context) error {
	return c.Send("TODO: send usage")
}

// handlerPing handles bot /ping command
func (a *app) handlerPing(c tele.Context) error {
	return c.Send("pong")
}

// handlerAdd handles bot /add command
func (a *app) handlerAdd(c tele.Context) error {
	subscription := storage.Subscription{
		ClientID: c.Sender().ID,
		FeedLink: c.Args()[0],
	}

	if err := a.storage.Add(subscription); err != nil {
		fmt.Println("[ERROR]: adding subscription failed: %w", err)
		c.Send("fail")
	}

	url1 := "https://sun9-50.userapi.com/s/v1/ig2/Fj5REsc9AMSWaejtjVSzJYyJI5-DFkkS2u-5lOAwCCDDtbpaklRn-r6Z4Jwqy-gQ5BhHM7r7X0XrgWS84EYG2YWv.jpg?size=604x603&quality=95&type=album"
	url2 := "https://sun9-45.userapi.com/s/v1/ig2/XMyMA0-YmKAaAZDaMb_v4q6pf5TdQSM5N-j74RZ8BEDcZX5FZe4uzldThYxPPE5LI5zV-6eKu0hRR1_dugbRI9hk.jpg?size=1080x1079&quality=95&type=album"

	msgTemplate, err := template.New("msg").Parse(`
		{{ .Text }}
		{{range .Photos}}<a href="{{ . }}">&#8205;</a>&#8205;{{end}}
	`)
	if err != nil {
		return err
	}

	data := struct {
		Text   string
		Photos []string
	}{
		Text: "Джон Херси: Хиросима\nЦитата: «Поначалу, пробираясь между рядами разрушенных зданий, они никак не могли понять, где находятся; оживленный город с населением в 245 тысяч человек за одно утро превратился в обгоревшие руины с неясными очертаниями — и эта перемена была столь же разительной, сколь и внезапной.». \nСтраниц: 224 \n \nПосле книги про Японию захотелось больше узнать о людях, которые воочию наблюдали и переживали ужасы взрыва атомной бомбы в 1945 году. Журналист Джон Херси проследил, что было с шестью выжившими до, в момент и после взрыва. Таких людей называют хибакуся, среди них наши герои: две женщины, два врача и два священника. Все они находились в разных местах во время взрыва, и смогли выжить благодаря случайному стечению обстоятельств. Тяжело читать про улицы усеянные трупами, развитие лучевой болезни, страдания людей, поломанные судьбы. Никто не понес наказание за этот американский \"эксперимент\" над мирными жителями, а ядерное оружие начало множиться по всей планете. Оставшиеся в живых хибакуся были опрошены в 1984 году, и 54,3% из них заявили, что по их мнению ядерное оружие будет использовано снова. Надеюсь, они ошиблись в своих предположениях. \n \nОценка: 5/5 \n \n#книги #нонфикшн@nookisbook",
		Photos: []string{
			url1,
			url2,
		},
	}

	var htmlMsg bytes.Buffer
	if err := msgTemplate.Execute(&htmlMsg, data); err != nil {
		return err
	}

	msg := htmlMsg.String()
	fmt.Println(msg)
	return c.Send(msg, &tele.SendOptions{
		ParseMode:             tele.ModeHTML,
		DisableWebPagePreview: false,
	})

}

// handlerRemove handles bot /remove command
func (a *app) handlerRemove(c tele.Context) error {
	return c.Send("remove")
}

// handlerUnknown handles unknown/unrecognized command
func (a *app) handlerUnknown(c tele.Context) error {
	return c.Send("unknown command, see usage")
}