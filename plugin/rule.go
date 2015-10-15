package plugin

import (
	"log"
	"strconv"
	"strings"

	"github.com/jqs7/Jqs7Bot/conf"
	"github.com/jqs7/bb"
)

type Rule struct{ bb.Base }

func (r *Rule) Run() {
	chatIDStr := strconv.Itoa(r.ChatID)
	rule := conf.Redis.Get("tgGroupRule:" + chatIDStr).Val()
	if rule != "" {
		r.NewMessage(r.ChatID, rule).Send()
	} else {
		r.NewMessage(r.ChatID,
			conf.List2StringInConf("rules")).Send()
	}
}

type SetRule struct{ Default }

func (s *SetRule) Run() {
	if len(s.Args) < 2 || !s.FromGroup {
		return
	}
	rule := strings.Join(s.Args[1:], " ")
	if s.isAuthed() {
		chatIDStr := strconv.Itoa(s.ChatID)
		log.Printf("setting rule %s to %s\n", rule, chatIDStr)
		conf.Redis.Set("tgGroupRule:"+chatIDStr, rule, -1)
		s.NewMessage(s.ChatID,
			"新的群组规则Get！✔️\n以下是新的规则：\n\n"+rule).
			Send()
	} else {
		s.sendQuestion()
	}
}

type AutoRule struct{ Default }

func (s *AutoRule) Run() {
	if s.FromGroup {
		chatIDStr := strconv.Itoa(s.ChatID)
		if conf.Redis.Exists("tgGroupAutoRule:" + chatIDStr).Val() {
			conf.Redis.Del("tgGroupAutoRule:" + chatIDStr)
			s.NewMessage(s.ChatID, "AutoRule Disabled!").Send()
		} else {
			conf.Redis.Set("tgGroupAutoRule:"+chatIDStr,
				strconv.FormatBool(true), -1)
			s.NewMessage(s.ChatID, "AutoRule Enabled!").Send()
		}
	}
}