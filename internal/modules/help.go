/*
  - This file is part of YukkiMusic.
  - Edited by KIYICI BOSS (@officialkiyici) - Aşko Kuşko Versiyonu 💅
*/
package modules

import (
	"fmt"
	"strings"

	tg "github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	"main/internal/core"
)

func init() {
	helpTexts["/help"] = fmt.Sprintf(`ℹ️ <b>Yardım Menüsü Aşkooo</b> 💅✨
<i>Benden ne istediğini bilemiyorsan veya komutların dedikodusunu merak ediyorsan burası tam senlik tatlım! 💖</i>

<u>Nasıl kullanılır bebeğim:</u>
<code>/yardim</code> — Tüm sırlarımı ve ana menümü dökerim önüne kızzz. 🌸  
<code>/yardim &lt;komut&gt;</code> — O komut ne işe yarıyor hemen fısıldarım kulağına. 🤫🎀

<b>💡 Tatlış Bir İpucu:</b> İstediğin komutun sonuna <code>-h</code> veya <code>--help</code> eklersen direkt sana özel açıklarım bebeğim, misal: <code>/oynat -h</code> 💅

<b>⚠️ Minik Bir Uyarı Aşko:</b> Bazı komutlar herkese açık değil tatlım; sadece <b>Gruplara</b>, <b>Adminlere</b>, <b>Sudo Aşkolara</b> veya <b>Büyük Patrona (Kurucu)</b> özel. Özel alanıma girme yani! 🚫💁‍♀️  
Eğer yetkin olmayan bir yerde bunları denersen hiç oralı olmam, trip atarım bilesin. 💅🙄  
Yine de inadım inat, o komutu öğreneceğim diyorsan şunu yaz aşkım:
<code>/yardim &lt;komut&gt;</code>

Daha fazla dedikodu ve yardım için <a href="%s">Pembiş Destek Grubumuza</a> gelmeyi unutma kızzz! ☕💖`, config.SupportChat)

	// Türkçe komutu yardım menüsüne eşliyoruz
	helpTexts["/yardim"] = helpTexts["/help"]
}

func helpHandler(m *tg.NewMessage) error {
	args := strings.Fields(m.Text())
	if len(args) > 1 {
		cmd := args[1]
		if cmd != "pm_help" {
			if !strings.HasPrefix(cmd, "/") {
				cmd = "/" + cmd
			}
			return showHelpFor(m, cmd)
		}
	}

	if m.ChatType() != tg.EntityUser {
		m.Reply(
			F(m.ChannelID(), "help_private_only"),
			&tg.SendOptions{
				ReplyMarkup: core.GetGroupHelpKeyboard(m.ChannelID()),
			},
		)
		return tg.ErrEndGroup
	}

	m.Reply(
		F(m.ChannelID(), "help_main"),
		&tg.SendOptions{ReplyMarkup: core.GetHelpKeyboard(m.ChannelID())},
	)
	return tg.ErrEndGroup
}

func helpCB(c *tg.CallbackQuery) error {
	c.Edit(
		F(c.ChannelID(), "help_main"),
		&tg.SendOptions{ReplyMarkup: core.GetHelpKeyboard(c.ChannelID())},
	)
	c.Answer("")
	return tg.ErrEndGroup
}

func helpCallbackHandler(c *tg.CallbackQuery) error {
	data := c.DataString()
	c.Answer("")
	if data == "" {
		return tg.ErrEndGroup
	}
	chatID := c.ChannelID()
	parts := strings.SplitN(data, ":", 2)
	if len(parts) < 2 {
		return tg.ErrEndGroup
	}

	var text string
	btn := core.GetBackKeyboard(chatID)

	switch parts[1] {
	case "admins":
		text = F(chatID, "help_admin")
	case "sudoers":
		text = F(chatID, "help_sudo")
	case "owner":
		text = F(chatID, "help_owner")
	case "public":
		text = F(chatID, "help_public")
	case "main":
		text = F(chatID, "help_main")
		btn = core.GetHelpKeyboard(chatID)
	}

	c.Edit(text, &tg.SendOptions{ReplyMarkup: btn})
	return tg.ErrEndGroup
}
