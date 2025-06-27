package bot

import (
	"fmt"
	"log/slog"
	"time"

	"imageBot/internal/config"
	"imageBot/internal/model"
	"imageBot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	bot     *tgbotapi.BotAPI
	service *service.Service
	cfg     config.Bot
	prompt  string
}

func New(cfg config.Bot, prompt string, service *service.Service) *Bot {
	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		slog.Info(err.Error())
		return nil
	}
	return &Bot{service: service, cfg: cfg, bot: bot, prompt: prompt}
}
func (b *Bot) Start() {
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	b.bot = bot
	if err != nil {
		slog.Info(err.Error())
		return
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)
	go func() {
		for update := range updates {
			if update.PollAnswer != nil {
				user_id := update.PollAnswer.User.ID
				poll_id := update.PollAnswer.PollID
				options_id := update.PollAnswer.OptionIDs
				b.service.Poll.Vote(user_id, poll_id, options_id[0])
			}
		}
	}()
	timeout := time.Duration(b.cfg.Timeout) * time.Hour
	ticker := time.NewTicker(timeout).C
	//prompt := "Роскошный мужской пиджак, сшитый из удивительной рыбьей шерсти, с шелковой подкладкой. Макросъемка, демонстрирующая уникальную текстуру ткани. Стиль: элитная модная реклама, высокая детализация, мягкое светотеневое оформление. Можешь добавить какие-то случайные элементы"
	//i := 1

	for t := range ticker {
		//if i%24 == 0 && i < 12 {
		photo, message_id, err := b.PostPhoto(b.prompt, 1024, 1024)
		if err != nil {
			slog.Info(err.Error())
			continue
		}
		err = b.service.SaveImage(photo)
		if err != nil {
			slog.Info(fmt.Sprintf("Error saving image: %s", err.Error()))
			continue
		}
		err = b.service.Image.SaveImageMessage(message_id, photo.ID)
		if err != nil {
			slog.Info(fmt.Sprintf("Error saving image: %s", err.Error()))
			continue
		}

		//}

		//if i%(24*7) == 0 && i > 1 {
		//	images, err := b.service.Image.GetImage(7)
		//	if err != nil {
		//		slog.Info(err.Error())
		//		continue
		//	}
		//
		//	poll := tgbotapi.NewPoll(
		//		b.cfg.ChatId,
		//		"Изображение недели",
		//	)
		//	poll.Options = b.GenAns(i, 7, images)
		//	m, err := b.bot.Send(poll)
		//	if err != nil {
		//		slog.Info(err.Error())
		//		continue
		//	}
		//
		//	err = b.service.Poll.SavePoll(m.Poll.ID, m.MessageID, "week_top")
		//	if err != nil {
		//		slog.Info(err.Error())
		//		continue
		//	}
		//}
		//if i-12%(24*7) == 0 && i > 12 {
		//	poll, err := b.service.Poll.GetPoll("week_top", 7)
		//	if err != nil {
		//		slog.Info(err.Error())
		//		continue
		//	}
		//	answer_id, answer_count, err := b.service.Poll.GetPollResults(poll)
		//
		//	pollConfig := tgbotapi.StopPollConfig{}
		//	pollConfig.ChatID = b.cfg.ChatId
		//	pollConfig.MessageID = poll.MessageId
		//	resp, err := b.bot.Request(pollConfig)
		//	if err != nil {
		//		slog.Info(err.Error())
		//
		//	}
		//	slog.Info(fmt.Sprintf("%s-%s", answer_id, answer_count, resp))
		//	//poll := tgbotapi.NewPoll(
		//	//	b.cfg.ChatId,
		//	//	"Изображение недели",
		//	//)
		//	//poll.Options = GenAns(i, 7, ima)
		//	//poll.IsAnonymous = false
		//}
		//i++
		slog.Info(fmt.Sprintf("%v", t))
	}

}
func (b *Bot) PostPhoto(prompt string, w, h int) (*model.Image, int, error) {
	image, err := b.service.Image.GenerateImage(b.prompt, w, h)
	if err != nil {
		return nil, 0, err
	}
	T := time.Now()
	file := tgbotapi.FileBytes{
		Name:  fmt.Sprintf("%d-%s-%d-%s", T.Day(), T.Month().String(), T.Year(), T.Month(), prompt),
		Bytes: image.Content,
	}
	msg := tgbotapi.NewPhotoToChannel(b.cfg.Channel, file)
	m, err := b.bot.Send(msg)
	if err != nil {
		return nil, 0, err
	}

	b.cfg.ChatId = m.Chat.ID
	if len(m.Photo) > 0 {
		image.ID = m.Photo[len(m.Photo)-1].FileID
		return &image, m.MessageID, nil
	}
	return nil, 0, fmt.Errorf("Photo do not exist")

}
func (b *Bot) GenAns(i int, delta int, images []model.Image) []string {
	k := 0
	if len(images) < delta {
		delta = len(images)
	}
	ans := make([]string, delta)

	for photo_i := i - delta; photo_i < i; photo_i++ {
		//ans[k] = fmt.Sprintf("https://t.me/c/%s/%d", b.cfg.Channel, images[k].MessageId)
		ans[k] = fmt.Sprintf("Изображение №%d", k)
		k++
	}
	return ans
}
