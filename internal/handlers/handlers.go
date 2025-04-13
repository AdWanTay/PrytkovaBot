package handlers

import (
	"PrytkovaBot/internal/services"
	"PrytkovaBot/internal/utils"
	"fmt"
	tele "gopkg.in/telebot.v4"
	"log"
	"net/url"
	"strconv"
	"strings"
)

var AdminId int64

const (
	DietButton         = "buy_item"
	RegisterButton     = "register"
	ProgramsButton     = "programs"
	AcceptRegistration = "accept_registration"
	FirstDiet          = "first_diet"
	SecondDiet         = "second_diet"
	Pay                = "pay"
	FirstDietAmount    = 990
	SecondDietAmount   = 1590
	CheckPayment       = "check_payment"
)

func RegisterHandlers(b *tele.Bot, adminId int64) {
	AdminId = adminId
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
	btnPrograms := selector.Data("Программы", ProgramsButton)
	selector.Inline(
		selector.Row(btnBuy),
		selector.Row(btnRegister),
		selector.Row(btnPrograms),
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
	case ProgramsButton:
		return programsHandler(c)
	case AcceptRegistration:
		if len(parts) < 2 {
			return c.Send("Ошибка: некорректный формат данных кнопки")
		}
		return handleRegistrationResponse(c, parts[1])
	case DietButton:
		return dietHandler(c)
	case FirstDiet, SecondDiet:
		return handleDietButtons(c, action)
	case Pay:
		switch parts[1] {
		case FirstDiet:
			return handlePay(c, FirstDietAmount)
		case SecondDiet:
			return handlePay(c, SecondDietAmount)
		}
	case CheckPayment:
		return handleCheckPayment(c, parts[1], parts[2])

	}
	return c.Respond()
}

func handleCheckPayment(c tele.Context, paymentId string, amount string) error {
	state, err := services.GetState(paymentId)
	if err != nil {
		return err
	}
	if state == "CONFIRMED" {
		err = c.Send("Платеж подтвержден!")
		var document *tele.Document
		a, _ := strconv.ParseFloat(amount, 64)
		switch a {
		case FirstDietAmount:
			document = &tele.Document{File: tele.FromDisk("diet1.pdf"), FileName: "Диета на 2 недели.pdf"}
		case SecondDietAmount:
			document = &tele.Document{File: tele.FromDisk("diet2.pdf"), FileName: "Диета на 4 недели.pdf"}
		}
		_, err = c.Bot().Send(c.Sender(), document, tele.Protected)
	}
	return nil
}

func handlePay(c tele.Context, amount float64) error {
	paymentUrl, paymentId, err := services.InitTransaction(amount)
	if err != nil {
		return err
	}
	selector := &tele.ReplyMarkup{}
	btnCheckPayment := selector.Data("Проверить платеж", fmt.Sprintf("%s@%s@%f", CheckPayment, paymentId, amount))
	selector.Inline(selector.Row(btnCheckPayment))

	return c.Edit("Оплатить: "+paymentUrl, selector)
}

func programsHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.URL("Хочу на программу!", utils.GetWhatsAppString()),
		),
	)
	return c.Send("Я буду работать с тобой шаг за шагом, отслеживая результат и корректируя программу\\. "+
		"Обновленная программа включает новые материалы, чек\\-листы и рекомендации\\. Вместе мы проработаем твои цели "+
		"и сделаем так, чтобы процесс был максимально комфортным\\.\n"+
		"🔹*Что ты получишь:*\n\n• Персонализированная поддержка\n• Чек\\-листы для контроля\n• Мотивационные материалы\n",
		tele.ModeMarkdownV2, selector)
}

func registerHandler(c tele.Context) error {
	selector := &tele.ReplyMarkup{}
	btnAccept := selector.Data("✅", fmt.Sprintf("%s@%d", AcceptRegistration, c.Sender().ID))
	selector.Inline(
		selector.Row(btnAccept),
	)

	err := c.Send("Похоже тебе нужно разобраться в питании, получить советы по планированию рациона или просто разобраться в том, что мешает двигаться вперед, моя консультация — это то, что тебе нужно\\. Мы разберемся, как сделать твое питание здоровым и удобным\\."+
		"\n*Заявка отправлена\\. Ожидайте подтверждения\\.*", tele.ModeMarkdownV2)

	_, err = c.Bot().Send(tele.ChatID(AdminId), fmt.Sprintf("Новый клиент: @%s (ID: %d) хочет записаться.\nПодтвердить?", c.Sender().Username, c.Sender().ID), selector)
	return err
}

func handleRegistrationResponse(c tele.Context, userIDStr string) error {
	greeting := "Здравствуйте, я хочу к вам на консультацию."
	userToId, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Send("Ошибка: некорректный ID пользователя")
	}

	selector := &tele.ReplyMarkup{}
	whatsAppUrl, err := url.Parse("https://wa.me/79659413788")
	if err != nil {
		log.Fatal(err)
	}
	params := url.Values{}
	params.Add("text", greeting)
	whatsAppUrl.RawQuery = params.Encode()
	selector.Inline(
		selector.Row(
			selector.URL("WhatsApp", whatsAppUrl.String()),
		),
		selector.Row(
			selector.URL("Telegram", fmt.Sprintf("https://t.me/ReginaUspeshnaya?text=%s", greeting)),
		),
	)
	responseText := "✅ Запись подтверждена!"
	_, err = c.Bot().Send(tele.ChatID(userToId), fmt.Sprintf("%s\nПожалуйста, свяжитесь со мной👇🏻", responseText), selector)
	err = c.Edit(strings.ReplaceAll(c.Message().Text, "Подтвердить?", ""))
	err = c.Reply(responseText)
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

	return c.Send("Есть два варианта приобретения РАЦИОНА — они сбалансированные, доступные и легкие\\! "+
		"Это отличный старт для тех, кто хочет наладить питание без жестких ограничений\\. Ты будешь удивлена, "+
		"как просто и вкусно можно питаться\\!\n\n🔹*Преимущества:*\n• Простой и понятный формат\n• Доступные продукты\n• Подробные рекомендации и рецепты\n", tele.ModeMarkdownV2, selector)
}

func handleDietButtons(c tele.Context, action string) error {
	selector := &tele.ReplyMarkup{}
	unique := ""
	description := ""
	if action == FirstDiet {
		unique = Pay + "@" + FirstDiet
		description = "Описание диеты на 2 недели"
	} else if action == SecondDiet {
		unique = Pay + "@" + SecondDiet
		description = "Описание диеты на 4 недели"
	}
	btnBuy := selector.Data("Оплатить", unique)
	selector.Inline(
		selector.Row(btnBuy),
	)
	return c.Send(description, selector)
}
