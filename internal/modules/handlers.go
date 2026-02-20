/*
  - This file is part of YukkiMusic.
    *

  - YukkiMusic — A Telegram bot that streams music into group voice chats with seamless playback and control.
  - Copyright (C) 2025 TheTeamVivek
    *
  - This program is free software: you can redistribute it and/or modify
  - it under the terms of the GNU General Public License as published by
  - the Free Software Foundation, either version 3 of the License, or
  - (at your option) any later version.
    *
  - This program is distributed in the hope that it will be useful,
  - but WITHOUT ANY WARRANTY; without even the implied warranty of
  - MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
  - GNU General Public License for more details.
    *
  - You should have received a copy of the GNU General Public License
  - along with this program. If not, see <https://www.gnu.org/licenses/>.
*/
package modules

import (
	"fmt"
	"log"

	"github.com/Laky-64/gologging"
	"github.com/amarnathcjd/gogram/telegram"

	"main/internal/config"
	"main/internal/core"
	"main/internal/database"
)

type MsgHandlerDef struct {
	Pattern string
	Handler telegram.MessageHandler
	Filters []telegram.Filter
}

type CbHandlerDef struct {
	Pattern string
	Handler telegram.CallbackHandler
	Filters []telegram.Filter
}

var handlers = []MsgHandlerDef{
	{Pattern: "json", Handler: jsonHandle},
	{
		Pattern: "eval",
		Handler: evalHandle,
		Filters: []telegram.Filter{ownerFilter},
	},
	{
		Pattern: "ev",
		Handler: evalCommandHandler,
		Filters: []telegram.Filter{ownerFilter},
	},
	{
		Pattern: "(bash|sh)",
		Handler: shellHandle,
		Filters: []telegram.Filter{ownerFilter},
	},
	{
		Pattern: "(restart|yenidenbaslat)",
		Handler: handleRestart,
		Filters: []telegram.Filter{ownerFilter, ignoreChannelFilter},
	},

	{
		Pattern: "(addsudo|addsudoer|sudoadd|sudoekle)",
		Handler: handleAddSudo,
		Filters: []telegram.Filter{ownerFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(delsudo|delsudoer|sudodel|remsudo|rmsudo|sudorem|dropsudo|unsudo|sudokaldır)",
		Handler: handleDelSudo,
		Filters: []telegram.Filter{ownerFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(sudoers|listsudo|sudolist|sudolistesi)",
		Handler: handleGetSudoers,
		Filters: []telegram.Filter{ignoreChannelFilter},
	},

	{
		Pattern: "(speedtest|spt|hiztesti)",
		Handler: sptHandle,
		Filters: []telegram.Filter{sudoOnlyFilter, ignoreChannelFilter},
	},

	{
		Pattern: "(broadcast|gcast|bcast|duyuru|yayin)",
		Handler: broadcastHandler,
		Filters: []telegram.Filter{ownerFilter, ignoreChannelFilter},
	},

	{
		Pattern: "(ac|active|activevc|activevoice|aktifler|aktifsesler)",
		Handler: activeHandler,
		Filters: []telegram.Filter{sudoOnlyFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(maintenance|maint|bakim)",
		Handler: handleMaintenance,
		Filters: []telegram.Filter{ownerFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(logger|logcu)",
		Handler: handleLogger,
		Filters: []telegram.Filter{sudoOnlyFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(autoleave|otocikis)",
		Handler: autoLeaveHandler,
		Filters: []telegram.Filter{sudoOnlyFilter, ignoreChannelFilter},
	},
	{
		Pattern: "(log|logs|kayitlar)",
		Handler: logsHandler,
		Filters: []telegram.Filter{sudoOnlyFilter, ignoreChannelFilter},
	},

	{
		Pattern: "(help|yardim)",
		Handler: helpHandler,
		Filters: []telegram.Filter{ignoreChannelFilter},
	},
	{
		Pattern: "ping",
		Handler: pingHandler,
		Filters: []telegram.Filter{ignoreChannelFilter},
	},
	{
		Pattern: "(start|basla)",
		Handler: startHandler,
		Filters: []telegram.Filter{ignoreChannelFilter},
	},
	{
		Pattern: "(stats|istatistikler)",
		Handler: statsHandler,
		Filters: []telegram.Filter{ignoreChannelFilter, sudoOnlyFilter},
	},
	{
		Pattern: "(bug|hata)",
		Handler: bugHandler,
		Filters: []telegram.Filter{ignoreChannelFilter},
	},
	{
		Pattern: "(lang|language|dil)",
		Handler: langHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},

	// SuperGroup & Admin Filters

	{
		Pattern: "(stream|yayinla)",
		Handler: streamHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(streamstop|yayindurdur)",
		Handler: streamStopHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(streamstatus|yayindurumu)",
		Handler: streamStatusHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{Pattern: "(rtmp|setrtmp|rtmpayarla)", Handler: setRTMPHandler},
	{
		Pattern: "(autoplay|otooynat)",
		Handler: autoplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	// play/cplay/vplay/fplay commands
	{
		Pattern: "(play|oynat)",
		Handler: playHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(fplay|playforce|foynat)",
		Handler: fplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cplay|koynat)",
		Handler: cplayHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(cfplay|fcplay|cplayforce|fkoynat)",
		Handler: cfplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(vplay|voynat)",
		Handler: vplayHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(fvplay|vfplay|vplayforce|fvoynat)",
		Handler: fvplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(vcplay|cvplay|vokoynat)",
		Handler: vcplayHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(fvcplay|fvcpay|vcplayforce|fvokoynat)",
		Handler: fvcplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},

	{
		Pattern: "(speed|setspeed|speedup|hiz|hizlandir)",
		Handler: speedHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(skip|atla|gec)",
		Handler: skipHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(pause|duraklat)",
		Handler: pauseHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(resume|devamet)",
		Handler: resumeHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(replay|tekrarlar)",
		Handler: replayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(mute|sustur)",
		Handler: muteHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(unmute|sesac)",
		Handler: unmuteHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(seek|sar)",
		Handler: seekHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(seekback|gerisar)",
		Handler: seekbackHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(jump|zipla)",
		Handler: jumpHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(position|pozisyon|sure)",
		Handler: positionHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(queue|sira|liste)",
		Handler: queueHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(clear|temizle)",
		Handler: clearHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(remove|sil|kaldir)",
		Handler: removeHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(move|tasi)",
		Handler: moveHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(shuffle|karistir)",
		Handler: shuffleHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(loop|setloop|dongu)",
		Handler: loopHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(end|stop|durdur|kapat)",
		Handler: stopHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(reload|yenile)",
		Handler: reloadHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},
	{
		Pattern: "(addauth|yetkiver)",
		Handler: addAuthHandler,
		Filters: []telegram.Filter{superGroupFilter, adminFilter},
	},
	{
		Pattern: "(delauth|yetkial)",
		Handler: delAuthHandler,
		Filters: []telegram.Filter{superGroupFilter, adminFilter},
	},
	{
		Pattern: "(authlist|yetkililer)",
		Handler: authListHandler,
		Filters: []telegram.Filter{superGroupFilter},
	},

	// CPlay commands
	{
		Pattern: "(cplay|cvplay|koynat|kvoynat)",
		Handler: cplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cfplay|fcplay|cforceplay|kfoynat)",
		Handler: cfplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cpause|kduraklat)",
		Handler: cpauseHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cresume|kdevamet)",
		Handler: cresumeHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cmute|ksustur)",
		Handler: cmuteHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cunmute|ksesac)",
		Handler: cunmuteHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cstop|cend|kdurdur)",
		Handler: cstopHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cqueue|ksira)",
		Handler: cqueueHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cskip|katla)",
		Handler: cskipHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cloop|csetloop|kdongu)",
		Handler: cloopHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cseek|ksar)",
		Handler: cseekHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cseekback|kgerisar)",
		Handler: cseekbackHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cjump|kzipla)",
		Handler: cjumpHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cremove|ksil)",
		Handler: cremoveHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cclear|ktemizle)",
		Handler: cclearHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cmove|ktasi)",
		Handler: cmoveHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(channelplay|kanaloynat)",
		Handler: channelPlayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cspeed|csetspeed|cspeedup|khiz)",
		Handler: cspeedHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(creplay|ktekrarlar)",
		Handler: creplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cposition|ksure)",
		Handler: cpositionHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cshuffle|kkaristir)",
		Handler: cshuffleHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(creload|kyenile)",
		Handler: creloadHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
	{
		Pattern: "(cautoplay|kotooynat)",
		Handler: cautoplayHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},

	{
		Pattern: "(nothumb|nothumbs|resimsiz)",
		Handler: nothumbHandler,
		Filters: []telegram.Filter{superGroupFilter, authFilter},
	},
}

var cbHandlers = []CbHandlerDef{
	{Pattern: "start", Handler: startCB},
	{Pattern: "help_cb", Handler: helpCB},
	{Pattern: "^lang:[a-z]", Handler: langCallbackHandler},
	{Pattern: `^help:(.+)`, Handler: helpCallbackHandler},

	{Pattern: "^close$", Handler: closeHandler},
	{Pattern: "^cancel$", Handler: cancelHandler},
	{Pattern: "^bcast_cancel$", Handler: broadcastCancelCB},

	{Pattern: `^room:(\w+)$`, Handler: roomHandle},
	{Pattern: "progress", Handler: emptyCBHandler},
}

func Init(bot *telegram.Client, assistants *core.AssistantManager) {
	bot.UpdatesGetState()
	assistants.ForEach(func(a *core.Assistant) {
		a.Client.UpdatesGetState()
	})

	for _, h := range handlers {
		bot.AddCommandHandler(h.Pattern, SafeMessageHandler(h.Handler), h.Filters...).
			SetGroup(100)
	}

	for _, h := range cbHandlers {
		bot.AddCallbackHandler(h.Pattern, SafeCallbackHandler(h.Handler), h.Filters...).
			SetGroup(90)
	}

	bot.On("edit:/eval", evalHandle).SetGroup(80)
	bot.On("edit:/ev", evalCommandHandler).SetGroup(80)

	bot.On("participant", handleParticipantUpdate).SetGroup(70)

	bot.AddActionHandler(handleActions).SetGroup(60)

	assistants.ForEach(func(a *core.Assistant) {
		a.Ntg.OnStreamEnd(streamEndHandler)
	})

	go MonitorRooms()

	if is, _ := database.GetAutoLeave(); is {
		go startAutoLeave()
	}

	if config.SetCmds && config.OwnerID != 0 {
		go setBotCommands(bot)
	}

	cplayCommands := []string{
		"/cfplay", "/vcplay", "/fvcplay",
		"/cpause", "/cresume", "/cskip", "/cstop",
		"/cmute", "/cunmute", "/cseek", "/cseekback",
		"/cjump", "/cremove", "/cclear", "/cmove",
		"/cspeed", "/creplay", "/cposition", "/cshuffle",
		"/cloop", "/cqueue", "/creload", "/cautoplay",
	}

	for _, cmd := range cplayCommands {
		baseCmd := "/" + cmd[2:] // Remove 'c' prefix
		if baseHelp, exists := helpTexts[baseCmd]; exists {
			helpTexts[cmd] = fmt.Sprintf(`<i>Channel play variant of %s</i>

<b>⚙️ Requires:</b>
First configure channel using: <code>/channelplay --set [channel_id]</code>

%s

<b>💡 Note:</b>
This command affects the linked channel's voice chat, not the current group.`, baseCmd, baseHelp)
		}
	}
}

func setBotCommands(bot *telegram.Client) {
	// Set commands for normal users in private chats
	if _, err := bot.BotsSetBotCommands(&telegram.BotCommandScopeUsers{}, "", AllCommands.PrivateUserCommands); err != nil {
		gologging.Error("Failed to set PrivateUserCommands " + err.Error())
	}

	// Set commands for normal users in group chats
	if _, err := bot.BotsSetBotCommands(&telegram.BotCommandScopeChats{}, "", AllCommands.GroupUserCommands); err != nil {
		gologging.Error("Failed to set GroupUserCommands " + err.Error())
	}

	// Set commands for chat admins
	if _, err := bot.BotsSetBotCommands(
		&telegram.BotCommandScopeChatAdmins{},
		"",
		append(AllCommands.GroupUserCommands, AllCommands.GroupAdminCommands...),
	); err != nil {
		gologging.Error("Failed to set GroupAdminCommands " + err.Error())
	}

	// Set commands for sudo users in their private chat
	sudoers, err := database.GetSudoers()
	if err != nil {
		log.Printf("Failed to get sudoers for setting commands: %v", err)
	} else {
		sudoCommands := append(AllCommands.PrivateUserCommands, AllCommands.PrivateSudoCommands...)
		for _, sudoer := range sudoers {
			if _, err := bot.BotsSetBotCommands(&telegram.BotCommandScopePeer{
				Peer: &telegram.InputPeerUser{UserID: sudoer, AccessHash: 0},
			},
				"",
				sudoCommands,
			); err != nil {
				gologging.Error("Failed to set PrivateSudoCommands " + err.Error())
			}
		}
	}

	ownerCommands := append(
		AllCommands.PrivateUserCommands,
		AllCommands.PrivateSudoCommands...)
	ownerCommands = append(ownerCommands, AllCommands.PrivateOwnerCommands...)
	if _, err := bot.BotsSetBotCommands(&telegram.BotCommandScopePeer{
		Peer: &telegram.InputPeerUser{UserID: config.OwnerID, AccessHash: 0},
	}, "", ownerCommands); err != nil {
		gologging.Error("Failed to set PrivateOwnerCommands " + err.Error())
	}
}
