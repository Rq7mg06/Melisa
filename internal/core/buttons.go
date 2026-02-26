/*
  - This file is part of YukkiMusic.
  - Edited by KIYICI BOSS (@officialkiyici) - Aşko Kuşko Versiyonu 💅
*/
package core

import (
	"fmt"

	"github.com/Laky-64/gologging"
	tg "github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	"main/internal/locales"
	"main/internal/utils"
)

var GetChatLanguage func(chatID int64) (string, error) // overwritten from main.go

func AddMeMarkup(chatID int64) tg.ReplyMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL(F(chatID, "ADD_ME_BTN"),
				"https://t.me/"+BUser.Username+"?startgroup&admin=invite_users",
			),
		).
		Build()
}

func GetCancelKeyboard(chatID int64) *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.Data(F(chatID, "DOWNLOAD_CANCEL_BTN"), "cancel"),
		).
		Build()
}

func GetBroadcastCancelKeyboard(chatID int64) *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.Data(F(chatID, "BROADCAST_CANCEL_BTN"), "bcast_cancel"),
		).
		Build()
}

func SuppMarkup(chatID int64) tg.ReplyMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL(F(chatID, "SUPPORT_BTN"), config.SupportChat),
		).
		Build()
}

func GetPlayMarkup(chatID int64, r *RoomState, queued bool) tg.ReplyMarkup {
	btn := tg.NewKeyboard()
	prefix := "room:"
	if r.CplayID() != 0 {
		prefix = "croom:"
	}
	track := r.Track()
	duration := 0
	if track != nil {
		duration = track.Duration
	}

	// İlerleme çubuğunu daha tatlış yaptık aşkooo ✨
	progress := utils.GetProgressBar(r.Position(), duration)
	progress = "🌸 " + formatDuration(r.Position()) + " " + progress + " " + formatDuration(duration) + " 🌸"

	if !queued {
		btn.AddRow(
			tg.Button.Data(progress, "progress"),
		)
	}
	
	// Sembolleri daha modern ve aşko kuşko yaptık 💅
	btn.AddRow(
		tg.Button.Data("▶️", prefix+"resume"),
		tg.Button.Data("⏸️", prefix+"pause"),
		tg.Button.Data("⏭️", prefix+"skip"),
		tg.Button.Data("⏹️", prefix+"stop"),
	)

	btn.AddRow(
		tg.Button.Data("⏪ 15s", "room:seekback_15"),
		tg.Button.Data("🔁 Baştan", "room:replay"),
		tg.Button.Data("15s ⏩", "room:seek_15"),
	)

	btn.AddRow(
		tg.Button.Data(F(chatID, "CLOSE_BTN"), "close"),
	)

	return btn.Build()
}

func GetGroupHelpKeyboard(chatID int64) *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL(F(chatID, "GC_HELP_BTN"), "https://t.me/"+BUser.Username+"?start=pm_help"),
		).
		Build()
}

func GetStartMarkup(chatID int64) tg.ReplyMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.URL(
				F(chatID, "ADD_ME_BTN"),
				"https://t.me/"+BUser.Username+"?startgroup&admin=invite_users",
			),
		).
		AddRow(
			tg.Button.Data(
				F(chatID, "START_HELP_BTN"),
				"help_cb",
			),
		).
		AddRow(
			tg.Button.URL(
				F(chatID, "UPDATES_BTN"),
				config.SupportChannel,
			),
			tg.Button.URL(
				F(chatID, "SUPPORT_BTN"),
				config.SupportChat,
			),
		).
		Build()
}

func GetHelpKeyboard(chatID int64) *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.Data(
				F(chatID, "HELP_ADMINS_BTN"),
				"help:admins",
			),
			tg.Button.Data(
				F(chatID, "HELP_PUBLIC_BTN"),
				"help:public",
			),
		).
		AddRow(
			tg.Button.Data(
				F(chatID, "HELP_OWNER_BTN"),
				"help:owner",
			),
			tg.Button.Data(
				F(chatID, "HELP_SUDOERS_BTN"),
				"help:sudoers",
			),
		).
		AddRow(
			tg.Button.Data(
				F(chatID, "BACK_BTN"),
				"start",
			),
		).
		Build()
}

func GetBackKeyboard(chatID int64) *tg.ReplyInlineMarkup {
	return tg.NewKeyboard().
		AddRow(
			tg.Button.Data(
				F(chatID, "BACK_BTN"),
				"help:main",
			),
		).
		Build()
}

func formatDuration(sec int) string {
	h := sec / 3600
	m := (sec % 3600) / 60
	s := sec % 60

	if h > 0 {
		return fmt.Sprintf("%d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

func F(chatID int64, key string, values ...locales.Arg) string {
	lang := config.DefaultLang
	if GetChatLanguage != nil {
		l, err := GetChatLanguage(chatID)
		if err != nil {
			gologging.Error("Ayy dilini bulamadım tatlım: " + err.Error())
		} else {
			lang = l
		}
	}

	var val locales.Arg
	if len(values) > 0 {
		val = values[0]
	}
	return locales.Get(lang, key, val)
}
