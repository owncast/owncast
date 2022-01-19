package utils

import "testing"

func TestUserAgent(t *testing.T) {
	testAgents := []string{
		"Pleroma 1.0.0-1168-ge18c7866-pleroma-dot-site; https://pleroma.site info@pleroma.site",
		"Mastodon 1.2.3 Bot",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_5) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.1.1 Safari/605.1.15 (Applebot/0.1; +http://www.apple.com/go/applebot)",
		"WhatsApp",
	}

	for _, agent := range testAgents {
		if !IsUserAgentABot(agent) {
			t.Error("Incorrect parsing of useragent", agent)
		}
	}
}

func TestSanitizeString(t *testing.T) {
	targetString := "this is annoying"
	testStrings := []string{
		"𝓉𝒽𝒾𝓈 𝒾𝓈 𝒶𝓃𝓃𝑜𝓎𝒾𝓃𝑔",
		"𝒕𝒉𝒊𝒔 𝒊𝒔 𝒂𝒏𝒏𝒐𝒚𝒊𝒏𝒈",
		"𝖙𝖍𝖎𝖘 𝖎𝖘 𝖆𝖓𝖓𝖔𝖞𝖎𝖓𝖌",
		"t̸̰̰̪̤̲͕̯̳̰͆̐h̶̙͉̝̲͈̘̜̯̖̺͌͘i̷̢̦͓̪̱͝͠ș̴̢́̓ ̴̡͕̺͎̹̽͊i̵̡̳̟̙͔͗́̔̎̾͜s̸̞͍̭̞̙̥͑̑͊͜͜ ̴̮̝̔̐́͑̀̐̒́á̶̪̣̝̝͝ṋ̸̨̱̖̖̥̝̈́͗̓͑̓̏͘̚͝͠n̶̠̓̅͂́̽͛͘ő̶͇̮̹͇̭͕͋͆̋̓̔̓́̈́͘͜ỳ̷̛͉̺̪̯͚͛̋͝ì̴̞̹̑̂̐͂͝n̵̳̞͇̘͔̣͌̈́̀͝g̸̢̢̡̢̛̜̬̤͋͆̈̎̓̀̌̚",
		"t҉h҉i҉s҉ ҉i҉s҉ ҉a҉n҉n҉o҉y҉i҉n҉g҉",
	}

	for _, s := range testStrings {
		r := SanitizeString(s)
		if r != targetString {
			t.Error("Incorrect sanitization of string", s, "got", r)
		}
	}

	zwspStr := "str1​str2"
	zwspStrExpected := "str1str2"
	r := SanitizeString(zwspStr)
	if r != zwspStrExpected {
		t.Error("Incorrect sanitization of string", zwspStr, "got", r)
	}
}
