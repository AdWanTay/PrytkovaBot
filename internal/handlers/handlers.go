package handlers

import (
	"fmt"
	tele "gopkg.in/telebot.v4"
	"strconv"
	"strings"
)

var AdminId int64

const (
	DietButton         = "buy_item"
	RegisterButton     = "register"
	AcceptRegistration = "accept_registration"
	RejectRegistration = "reject_registration"
	FirstDiet          = "first_diet"
	SecondDiet         = "second_diet"
	Buy                = "buy"
)

func RegisterHandlers(b *tele.Bot, adminId int64) {
	AdminId = adminId
	//b.Handle(tele.OnText, registerHandler)
	b.Handle("/start", startHandler)
	b.Handle(tele.OnCallback, inlineButtonHandler)
}

func startHandler(c tele.Context) error {
	if c.Sender().ID == AdminId {
		return c.Send("Вы администратор")
	}
	selector := &tele.ReplyMarkup{}
	btnBuy := selector.Data("Рационы", DietButton)
	btnRegister := selector.Data("Записаться на консультацию", RegisterButton)
	selector.Inline(
		selector.Row(btnBuy),
		selector.Row(btnRegister),
	)
	return c.Send("На связи помощник Регины Прытковой👩‍💼.\nРадa видеть тебя в на этой страничке!👋 Тут ты найдешь всё, что нужно для правильного питания, похудения и здорового образа жизни для всей семьи😋.\n\n🎯 Смотри,  у меня есть несколько продуктов, которые помогут тебе достичь своих целей. Выбери что тебя сейчас интересует?", selector)
}

func inlineButtonHandler(c tele.Context) error {
	data := strings.TrimSpace(c.Callback().Data)
	parts := strings.Split(data, "@")

	if len(parts) < 1 {
		return c.Respond()
	}

	action := parts[0]

	switch action {
	case RegisterButton:
		return registerHandler(c)
	case AcceptRegistration, RejectRegistration:
		if len(parts) < 2 {
			return c.Send("Ошибка: некорректный формат данных кнопки")
		}
		return handleRegistrationResponse(c, action, parts[1])
	case DietButton:
		return dietHandler(c)
	case FirstDiet, SecondDiet:
		return handleDietButtons(c, action)
	}
	return c.Respond()
}

func registerHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnAccept := selector.Data("✅", fmt.Sprintf("%s@%d", AcceptRegistration, c.Sender().ID))
	btnReject := selector.Data("❌", fmt.Sprintf("%s@%d", RejectRegistration, c.Sender().ID))
	selector.Inline(
		selector.Row(btnReject, btnAccept),
	)
	err := c.Send("Заявка отправлена. Ожидайте подтверждения.")

	if err != nil {
		return err
	}
	_, err = c.Bot().Send(tele.ChatID(AdminId), fmt.Sprintf("Новый клиент: @%s (ID: %d) хочет записаться.\nПодтвердить?", c.Sender().Username, c.Sender().ID), selector)
	return err
}

func handleRegistrationResponse(c tele.Context, action, userIDStr string) error {
	userToId, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Send("Ошибка: некорректный ID пользователя")
	}
	var responseText string
	if action == AcceptRegistration {
		responseText = "✅ Запись подтверждена!"
		_, err = c.Bot().Send(tele.ChatID(userToId), fmt.Sprintf("%s\n%s", responseText, "Свяжитесь со мной: @kizzzzzaaaaa"))
	} else {
		responseText = "❌ Запись отклонена."
		_, err = c.Bot().Send(tele.ChatID(userToId), fmt.Sprintf("%s", responseText))
	}
	err = c.Reply(responseText)

	if err != nil {
		return c.Send("Ошибка при редактировании сообщения у администратора.")
	}
	//if action == AcceptRegistration {
	//	_, err = c.Bot().Send(tele.ChatID(userToId), fmt.Sprintf("%s\n%s", responseText, "Свяжитесь со мной: @kizzzzzaaaaa"))
	//} else {
	//	_, err = c.Bot().Send(tele.ChatID(userToId), fmt.Sprintf("%s", responseText))
	//}
	return err
}

func dietHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnFirst := selector.Data("на 2 недели за 990 рублей 🔥", FirstDiet)
	btnSecond := selector.Data("на 4 недели за 1590 рублей 🔥", SecondDiet)
	selector.Inline(
		selector.Row(btnFirst),
		selector.Row(btnSecond),
	)
	return c.Send("Есть два варианта приобретения РАЦИОНА — они сбалансированные, доступные и легкие! "+
		"Это отличный старт для тех, кто хочет наладить питание без жестких ограничений. Ты будешь удивлена, "+
		"как просто и вкусно можно питаться!\n", selector)
}

func handleDietButtons(c tele.Context, action string) error {
	selector := &tele.ReplyMarkup{}
	unique := ""
	description := ""
	if action == FirstDiet {
		unique = Buy + "@" + FirstDiet
		description = "Описание диеты на 2 недели"
	} else if action == SecondDiet {
		unique = Buy + "@" + FirstDiet
		description = "Описание диеты на 4 недели"
	}
	btnBuy := selector.Data("Оплатить", unique)
	selector.Inline(
		selector.Row(btnBuy),
	)
	return c.Send(description, selector)
}

func echoHandler(c tele.Context) error {
	if AdminId == c.Sender().ID {
		return c.Send(c.Message().Text)
	}
	return nil
}
