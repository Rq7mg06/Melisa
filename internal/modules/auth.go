func init() {
	helpTexts["/addauth"] = fmt.Sprintf(
		`<i>Aşkooo, gruptaki tatlış birine admin yapmadan müziği yönetme yetkisi vermek istersen bu komut tam senlik! 💅✨</i>

<u>Nasıl kullanılır bebeğim:</u>
<b>/addauth [mesajını yanıtla]</b> — Aşkonun mesajını yanıtlayarak yetki ver. 💖
<b>/addauth &lt;kullanıcı_id / @kullanıcı_adı&gt;</b> — Direkt ID veya kullanıcı adıyla ekle. 🎀

<b>⚙️ Minik Notlar:</b>
• Sadece <b>grup adminleri</b> yapabilir bunu tatlım, herkes değil! 💁‍♀️
• Yetki alan aşko müziği durdurabilir, geçebilir (<code>/pause</code>, <code>/skip</code> falan filan işte). 🎶
• 🤖 Botlara yetki veremiyoruz maalesef aşkım.
• 🔢 Her grupta en fazla <b>%d</b> kişiye bu özel yetkiyi verebilirsin. 🌸
• 👑 Sahibim, asistanım falan zaten doğuştan yetkili, onları listeye eklemene gerek yok bebeğim! ✨

Benzer şeyler için <code>/delauth</code> ve <code>/authlist</code> komutlarına da bakabilirsin kızzz.`,
		config.MaxAuthUsers,
	)

	helpTexts["/delauth"] = `<i>Ayyy yetki verdiğin biri canını mı sıktı? Hemen yetkisini alıyoruz aşko! 💅🚫</i>

<u>Nasıl kullanılır bebeğim:</u>
<b>/delauth [mesajını yanıtla]</b> — Mesajını yanıtlayıp şutla! 💁‍♀️
<b>/delauth &lt;kullanıcı_id / @kullanıcı_adı&gt; </b>— Direkt ID veya adıyla şutla! 💅

<b>⚙️ Minik Notlar:</b>
• Sadece <b>grup adminleri</b> bu komutu kullanabilir tatlım. 👑
• Kimde yetki var diye merak ediyorsan <code>/authlist</code> yazıp bakabilirsin aşkooo! 🌸`

	helpTexts["/authlist"] = `<u>Nasıl kullanılır bebeğim:</u>
<b>/authlist</b> - <i>Grupta kimlerin müziği yönetme yetkisi var, hepsini dökerim ortaya kızzz! 💅✨</i>

<b>⚙️ Minik Notlar:</b>
• Gruptaki herkes bakabilir buna tatlım. 👀
• Sadece sonradan yetki verilen tatlışları gösterir, zaten yetkili olan büyük patronları göstermez bebeğim. 💖`
}
