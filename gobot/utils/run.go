package utils

import (
    "encoding/json"
    "log"
    "net/http"

    "bot/config/settings"

    "github.com/PaulSonOfLars/gotgbot/v2"
    "github.com/PaulSonOfLars/gotgbot/v2/ext"
)

func BotInit() (*gotgbot.Bot, *ext.Dispatcher) {
    bot, err := gotgbot.NewBot(settings.Envs.BOT_TOKEN, nil)
    if err != nil {
        log.Fatal(err)
    }

    dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{
        Error: func(b *gotgbot.Bot, ctx *ext.Context, err error) ext.DispatcherAction {
            log.Println(err)
            return ext.DispatcherActionNoop
        },
    })

    return bot, dispatcher
}

func StartWebhook(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) {
    _, err := bot.SetWebhook(settings.Envs.WEBHOOK_URL, &gotgbot.SetWebhookOpts{})
    if err != nil {
        log.Fatal(err)
    }

    http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
        defer r.Body.Close()

        var update gotgbot.Update
        if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }

        if err := dispatcher.ProcessUpdate(bot, &update, nil); err != nil {
            log.Println(err)
        }

        w.WriteHeader(http.StatusOK)
    })

    log.Println("Webhook listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func StartPolling(bot *gotgbot.Bot, dispatcher *ext.Dispatcher) {
    log.Println("Starting polling...")

    updater := ext.NewUpdater(dispatcher, nil)

    err := updater.StartPolling(bot, nil)
    if err != nil {
        log.Fatal("Pollingni boshlashda xatolik:", err)
    }

    log.Println("Bot muvaffaqiyatli polling rejimida ishga tushdi.")
    updater.Idle()
}