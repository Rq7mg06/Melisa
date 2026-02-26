/*
  - This file is part of YukkiMusic.
  - Edited by KIYICI BOSS (@officialkiyici) - Aşko Kuşko Versiyonu 💅
*/
package modules

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/amarnathcjd/gogram/telegram"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"

	"main/internal/config"
	"main/internal/database"
	"main/internal/locales"
)

func init() {
	helpTexts["/stats"] = `<i>Ayyy benim içimde neler oluyor görmek istersen buraya bakıyoruz aşkooo! (Botun dedikodu haritası) 📊💅</i>

<u>Nasıl kullanılır bebeğim:</u>
<b>/stats</b> — Bütün sistem sırlarımı dökerim 🎀

<b>📊 Neler Var Neler Kızzz:</b>
• Sistem bilgileri (Beynim, hafızam falan) 🧠
• Go çalışma detayları (Hızım, tempom) 🏃‍♀️
• Sunucu gücüm (Enerjim ne alemde) ⚡
• Hizmet ettiğim gruplar (Hangi mekanlara akıyoruz) 💃
• Hizmet ettiğim kişiler (Kaç kişinin kalbini çaldım) 💖

<b>🔒 Ama Bir Şartım Var:</b>
• Sadece <b>Sudo (Patron) Aşkolar</b> görebilir, herkes değil tatlım! Özel hayatımı herkese açamam 💁‍♀️`
}

func statsHandler(m *telegram.NewMessage) error {
	var sb strings.Builder
	sb.Grow(512)
	chatID := m.ChannelID()

	sb.WriteString(getSystemStats(chatID))
	sb.WriteString(getGoMemStats(chatID))
	sb.WriteString(getServerStats(chatID))
	sb.WriteString(getServedStats(chatID))

	m.Reply(sb.String())
	return telegram.ErrEndGroup
}

// ---- Sub Functions ----

func getSystemStats(chatID int64) string {
	var sb strings.Builder

	sb.WriteString(F(chatID, "stats_system_header") + "\n")
	sb.WriteString(F(chatID, "stats_system_os_arch", locales.Arg{
		"os":   runtime.GOOS,
		"arch": runtime.GOARCH,
	}) + "\n")
	sb.WriteString(F(chatID, "stats_system_cpu_goroutines", locales.Arg{
		"cpus":       runtime.NumCPU(),
		"goroutines": runtime.NumGoroutine(),
	}) + "\n\n")

	return sb.String()
}

func getGoMemStats(chatID int64) string {
	var sb strings.Builder
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	uptime := time.Since(config.StartTime).Minutes()
	gcPerMin := float64(memStats.NumGC) / uptime

	gcEmoji := "🟢"
	switch {
	case gcPerMin > 20:
		gcEmoji = "🔴"
	case gcPerMin > 10:
		gcEmoji = "🟠"
	}

	sb.WriteString(F(chatID, "stats_go_mem_header") + "\n")

	sb.WriteString(F(chatID, "stats_go_alloc", locales.Arg{
		"alloc": memStats.Alloc / 1024 / 1024,
	}) + "\n")
	sb.WriteString(F(chatID, "stats_go_sys", locales.Arg{
		"sys": memStats.Sys / 1024 / 1024,
	}) + "\n")
	sb.WriteString(F(chatID, "stats_go_gc", locales.Arg{
		"gc_count": memStats.NumGC,
		"emoji":    gcEmoji,
		"gc_rate":  fmt.Sprintf("%.1f", gcPerMin),
	}) + "\n\n")

	return sb.String()
}

func getServerStats(chatID int64) string {
	var sb strings.Builder

	sysMem, _ := mem.VirtualMemory()
	cpuPercent, _ := cpu.Percent(0, false)
	diskStat, _ := disk.Usage("/")

	cpuEmoji := "🟢"
	if len(cpuPercent) > 0 {
		switch {
		case cpuPercent[0] > 70:
			cpuEmoji = "🔴"
		case cpuPercent[0] > 40:
			cpuEmoji = "🟡"
		}
	}

	ramUsagePercent := (float64(sysMem.Used) / float64(sysMem.Total)) * 100
	ramEmoji := "🟢"
	switch {
	case ramUsagePercent > 80:
		ramEmoji = "🔴"
	case ramUsagePercent > 50:
		ramEmoji = "🟡"
	}

	sb.WriteString(F(chatID, "stats_server_header") + "\n")

	sb.WriteString(F(chatID, "stats_server_cpu", locales.Arg{
		"emoji": cpuEmoji,
		"cpu":   fmt.Sprintf("%.2f", cpuPercent[0]),
	}) + "\n")

	sb.WriteString(F(chatID, "stats_server_ram", locales.Arg{
		"emoji":     ramEmoji,
		"used_gib":  fmt.Sprintf("%.2f", float64(sysMem.Used)/1073741824),
		"total_gib": fmt.Sprintf("%.2f", float64(sysMem.Total)/1073741824),
	}) + "\n")

	sb.WriteString(F(chatID, "stats_server_storage", locales.Arg{
		"used_gib":  fmt.Sprintf("%.2f", float64(diskStat.Used)/1073741824),
		"total_gib": fmt.Sprintf("%.2f", float64(diskStat.Total)/1073741824),
	}) + "\n\n")

	return sb.String()
}

func getServedStats(chatID int64) string {
	var sb strings.Builder

	servedChats, err1 := database.GetServed()
	servedUsers, err2 := database.GetServed(true)

	sb.WriteString(F(chatID, "stats_served_header") + "\n")

	if err1 != nil {
		sb.WriteString(F(chatID, "stats_served_chats_err", locales.Arg{
			"error": err1.Error(),
		}) + "\n")
	} else {
		sb.WriteString(F(chatID, "stats_served_chats", locales.Arg{
			"count": len(servedChats),
		}) + "\n")
	}

	if err2 != nil {
		sb.WriteString(F(chatID, "stats_served_users_err", locales.Arg{
			"error": err2.Error(),
		}) + "\n")
	} else {
		sb.WriteString(F(chatID, "stats_served_users", locales.Arg{
			"count": len(servedUsers),
		}) + "\n")
	}

	return sb.String()
}
