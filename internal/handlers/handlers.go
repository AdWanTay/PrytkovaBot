package handlers

import (
	"PrytkovaBot/internal/services"
	st "PrytkovaBot/internal/storage"
	"PrytkovaBot/internal/utils"
	"fmt"
	tele "gopkg.in/telebot.v4"
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
	Book               = "book"
	MyBooks            = "Мои записи"
	Back               = "back"
)

func RegisterHandlers(b *tele.Bot, adminId int64) {

	AdminId = adminId
	b.Handle("/start", startHandler)
	b.Handle(tele.OnCallback, inlineButtonHandler)
	b.Handle(tele.OnText, textHandler)
}

func textHandler(c tele.Context) error {
	if c.Text() == MyBooks {
		return handleMyBooks(c)
	}

	return nil
}

func startHandler(c tele.Context) error {
	if c.Sender().ID == AdminId {
		keyboard := &tele.ReplyMarkup{ResizeKeyboard: true}
		btn := keyboard.Text(MyBooks)
		keyboard.Reply(keyboard.Row(btn))
		return c.Send("Вы администратор", keyboard)
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
	_ = c.Respond()
	data := strings.TrimSpace(c.Callback().Data)
	parts := strings.Split(data, "@")

	if len(parts) < 1 {
		return c.Respond()
	}

	action := parts[0]

	switch action {
	case Back:

		switch parts[1] {
		case "start":
			return startHandler(c)
		case "diets":
			return dietHandler(c)
		}
	case RegisterButton:
		return registerHandler(c)
	case ProgramsButton:
		return programsHandler(c)
	case AcceptRegistration:
		if len(parts) < 2 {
			return c.Send("Ошибка: некорректный формат данных кнопки")
		}
		return handleRegistrationResponse(c, parts[1], parts[2], parts[3])
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
	case Book:
		return handleBook(c, parts[1])
	case MyBooks:
		return handleMyBooks(c)
	}
	return c.Respond()
}

func handleMyBooks(c tele.Context) error {
	message, err := services.FormatBookedSlots(st.Db)
	if err != nil {
		return err
	}
	return c.Send(message)
}

func handleBook(c tele.Context, timeId string) error {
	slotId, err := strconv.Atoi(timeId)
	username := c.Sender().Username
	userId := c.Sender().ID
	if err != nil {
		return err
	}

	slotTime, err := st.GetTimeBySlotId(st.Db, int64(slotId))
	if err != nil {
		return err
	}

	//responseString := utils.PrepareForMarkdown(fmt.Sprintf("Вы записались на консультацию %s в %s. \nПожауйста, заполните [анкету](https://forms.gle/NnZ7XoQkKxS2FmDk9)", slotTime.Format("02.01"), slotTime.Format("15:04")))
	responseString := fmt.Sprintf("Вы записались на консультацию %s в %s. \nПожауйста, заполните [анкету](https://forms.gle/NnZ7XoQkKxS2FmDk9).", slotTime.Format("02.01"), slotTime.Format("15:04"))
	fmt.Println(responseString)
	err = c.Edit(responseString, tele.ModeMarkdown)
	if err != nil {
		return err
	}

	selector := &tele.ReplyMarkup{}
	btnAccept := selector.Data("✅", fmt.Sprintf("%s@%d@%s@%d", AcceptRegistration, c.Sender().ID, c.Sender().Username, slotId))
	selector.Inline(
		selector.Row(btnAccept),
	)

	_, err = c.Bot().Send(
		tele.ChatID(AdminId),
		fmt.Sprintf("Новый клиент: @%s (ID: %d) хочет записаться на консультацию на %s в %s.\nПодтвердить?",
			username, userId, slotTime.Format("02.01 "), slotTime.Format("15:04"),
		),
		selector)
	return err
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
	btnBack := selector.Data("⬅️ Назад", Back+"@start")

	selector.Inline(
		selector.Row(
			selector.URL("Хочу на программу!", utils.GetEncodedString("https://wa.me/79659413788", "Здравствуйте, хочу к вам на программу.")),
		),
		selector.Row(btnBack),
	)

	return c.Send("Я буду работать с тобой шаг за шагом, отслеживая результат и корректируя программу\\. "+
		"Обновленная программа включает новые материалы, чек\\-листы и рекомендации\\. Вместе мы проработаем твои цели "+
		"и сделаем так, чтобы процесс был максимально комфортным\\.\n"+
		"🔹*Что ты получишь:*\n\n• Персонализированная поддержка\n• Чек\\-листы для контроля\n• Мотивационные материалы\n",
		tele.ModeMarkdownV2, selector)
}

func registerHandler(c tele.Context) error {
	btns, err := services.FormatAvailableSlots(st.Db)
	if err != nil {
		return err
	}
	timeSelector := &tele.ReplyMarkup{}
	btnBack := tele.InlineButton{
		Text: "⬅️ Назад",
		Data: Back + "@start",
	}

	btns = append(btns, []tele.InlineButton{btnBack})

	timeSelector.InlineKeyboard = btns

	err = c.Send("Похоже тебе нужно разобраться в питании, получить советы по планированию рациона или просто разобраться в том, что мешает двигаться вперед, моя консультация — это то, что тебе нужно\\. Мы разберемся, как сделать твое питание здоровым и удобным\\."+
		"\n*Выбери удобное время:*", tele.ModeMarkdownV2, timeSelector)

	return err
}

func handleRegistrationResponse(c tele.Context, userIDStr, username, slotId string) error {
	userToId, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Send("Ошибка: некорректный ID пользователя")
	}
	greeting := "Здравствуйте, я хочу к вам на консультацию."
	selector := &tele.ReplyMarkup{}
	selector.Inline(
		selector.Row(
			selector.URL("WhatsApp", utils.GetEncodedString("https://wa.me/79659413788", greeting)),
		),
		selector.Row(
			selector.URL("Telegram", utils.GetEncodedString("https://t.me/ReginaUspeshnaya", greeting)),
		),
	)

	slotIdInt, err := strconv.Atoi(slotId)
	if err != nil {
		return err
	}
	err = st.BookSlot(st.Db, slotIdInt, userToId, username)
	if err != nil {
		return err
	}

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
	btnBack := selector.Data("⬅️ Назад", Back+"@start")
	selector.Inline(
		selector.Row(btnFirst),
		selector.Row(btnSecond),
		selector.Row(btnBack),
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
		description = "🥗 Рацион на 2 недели\nДля мягкого старта и быстрого результата\n\n💚 Хочешь начать питаться правильно, но не знаешь с чего начать?\nЭтот рацион — как надёжная опора: всё продумано за тебя. Без голода, без БАДов, без жёстких ограничений.\n\n🔸 2 недели разнообразного, вкусного и простого питания\n🔸 5 приёмов пищи в день: завтрак, обед, ужин, перекусы\n🔸 Подходит для снижения веса, нормализации аппетита, стабилизации энергии\n🔸 Меню без «странных» продуктов — только то, что можно купить в любом магазине\n🔸 Сбалансировано по белкам, углеводам, жирам\n🔸 Можно готовить на всю семью — они тоже будут в восторге!\n\n🌿 Особенно подойдёт тем, кто:\n— хочет перезапустить питание без стресса\n— пробовал \"ПП\", но всё время срывался\n— хочет убрать вздутие, тягу к сладкому и начать снижать вес\n\n📥 Готов загрузиться заботой? Этот рацион — твой первый шаг к себе.\n"
	} else if action == SecondDiet {
		unique = Pay + "@" + SecondDiet
		description = "🥑 Рацион на 4 недели\nДля тех, кто готов к устойчивому результату\n\n🌸 4 недели — это не просто меню. Это настоящая перезагрузка тела и питания.\n\nТы не просто ешь.\nТы начинаешь заботиться о себе — осознанно, вкусно, стабильно.\n\n🔹 4 недели разнообразного, лёгкого и вкусного питания\n🔹 Чёткий план без скуки и однообразия\n🔹 Поддержка гормонального фона, ЖКТ и уровня энергии\n🔹 Еда, которая насыщает, а не «разгоняет» аппетит\n🔹 Рацион составлен с учётом физиологии, без дефицитов\n🔹 Без БАДов, без экзотики, без страха перед едой\n\n🌿 Особенно подойдёт тем, кто:\n— хочет не просто похудеть, а перейти на питание, которое держится\n— устал от диет и откатов\n— готов заложить прочную основу для здоровья\n\n🧡 Этот рацион — как личная забота, которая рядом каждый день.\nТы не один/а. Ты в потоке. Ты начинаешь жить легче.\n"
	}
	btnBuy := selector.Data("Оплатить", unique)
	btnBack := selector.Data("⬅️ Назад", Back+"@diets")

	selector.Inline(
		selector.Row(btnBuy),
		selector.Row(btnBack),
	)
	return c.Send(description, selector)
}
