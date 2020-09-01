package main

import (
	"log"
	"strings"
	"testing"
)

func Test_getStreamText(t *testing.T) {
	task := "https://valipl-vip.cp31.ott.cibntv.net/67756D6080932713CFC02204E/05000900005D9111939E30803BAF2B466540E2-97BE-4090-B8A0-60B67B9C4892-00279.ts?ccode=0502&duration=2780&expire=18000&psid=b7e1ed2a7cc5c4786f1eedd3598ecd2244151&ups_client_netip=70613d76&ups_ts=1598922016&ups_userid=1523172795&utid=QzG%2FFyFiSjECAXBh9XUWDoBH&vid=XMjgzODY4NjEy&sm=1&operate_type=1&dre=u38&si=78&eo=0&dst=1&iv=1&s=14eaa592286e11e097c0&type=mp5hdv3&bc=2&rid=2000000096C19E6F80D425B10D0F648B36192BF902000000&vkey=Ba46931a2b22ceec5af423c87c450588b"
	name := task[strings.LastIndex(task, "/")+1:]
	if strings.Contains(name, "?") {
		name = name[:strings.Index(name, "?")]
	}
	log.Println(name)
}
