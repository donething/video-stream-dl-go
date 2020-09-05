package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"testing"
)

func Test_getStreamText(t *testing.T) {
	url := "https://valipl.cp31.ott.cibntv.net/69755C40DA84171B1D477628A/05000A00005EBE00E18BB780000000C3D87831-D26E-4D4D-A0F1-4EEBB10686FF.m3u8?ccode=0502&duration=2805&expire=18000&psid=e8d908afabe25cb143ab0c060bdd3f8f434af&ups_client_netip=7061d0d7&ups_ts=1599270746&ups_userid=1146161140&utid=KzK%2FF8TpxlACAXBh9XWeR4S4&vid=XMjgzODcwMzE2&vkey=Bed64b7033fec6d890968fd0b12589a78&sm=1&operate_type=1&dre=u37&si=73&eo=1&dst=1&iv=0&s=14eaa592286e11e097c0&type=mp5hd2v3&bc=2&rid=20000000B0E1204CC64A245917DD48DE0BB5BBB702000000"
	resp, err := http.Get(url)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err.Error())
	}
	log.Println(string(bs))
}
