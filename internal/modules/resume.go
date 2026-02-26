/*
  - This file is part of YukkiMusic.
  - Edited by KIYICI BOSS (@officialkiyici) - Aşko Kuşko Versiyonu 💅
*/
package modules

import (
	"fmt"
	"html"

	"github.com/amarnathcjd/gogram/telegram"

	"main/internal/locales"
	"main/internal/utils"
)

func init() {
	helpTexts["/resume"] = `<i>Ayyy şarkı yarım mı kaldı? Durdurulan müziği kaldığı yerden devam ettirir aşkooo! 🎶💅</i>

<u>Nasıl kullanılır bebeğim:</u>
<b>/resume</b> — Bekleyen şarkıyı canlandırır, koptuğumuz yerden devam! 💖

<b>⚙️ Neler Yapabiliyorum Kızzz:</b>
• Şarkı tam nerede kaldıysa oradan başlar, hiçbir detayı kaçırmazsın tatlım 🌸
• Kendi kendine başlama süresi falan varsa iptal eder, ipler senin elinde aşkım 💁‍♀️

<b>⚠️ Minik Notlar:</b>
• Zaten bangır bangır çalan şarkıyı devam ettiremem kız, delirtme beni! Sadece durdurulmuşsa işe yarar. 💅
• Sen durdurduğunda saniyesi saniyesine aklımda tutarım, unutmam 🎀
• Şarkı hızını değiştirdiysen o ayarların aynen kalır, modumuz asla bozulmaz ✨`
}

func resumeHandler(m *telegram.NewMessage) error {
	return handleResume(m, false)
}

func cresumeHandler(m *telegram.NewMessage) error {
	return handleResume(m, true)
}

func handleResume(m *telegram.NewMessage, cplay bool) error {
	chatID := m.ChannelID()

	r, err := getEffectiveRoom(m, cplay)
	if err != nil {
		m.Reply(err.Error())
		return telegram.ErrEndGroup
	}

	if !r.IsActiveChat() {
		m.Reply(F(chatID, "room_no_active"))
		return telegram.ErrEndGroup
	}

	if !r.IsPaused() {
		m.Reply(F(chatID, "resume_already_playing"))
		return telegram.ErrEndGroup
	}

	t := r.Track()
	if _, err := r.Resume(); err != nil {
		m.Reply(F(chatID, "resume_failed", locales.Arg{
			"error": err,
		}))
	} else {
		title := html.EscapeString(utils.ShortTitle(t.Title, 25))
		pos := formatDuration(r.Position())
		total := formatDuration(t.Duration)
		mention := utils.MentionHTML(m.Sender)

		speedLine := ""
		if sp := r.GetSpeed(); sp != 1.0 {
			speedLine = F(chatID, "speed_line", locales.Arg{
				"speed": fmt.Sprintf("%.2f", r.GetSpeed()),
			})
		}

		m.Reply(F(chatID, "resume_success", locales.Arg{
			"title":      title,
			"position":   pos,
			"duration":   total,
			"user":       mention,
			"speed_line": speedLine,
		}))
	}

	return telegram.ErrEndGroup
}
}
