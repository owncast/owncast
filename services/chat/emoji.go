package chat

import (
	"bytes"
	"strings"
	"sync"
	"text/template"
	"time"

	emojiDef "github.com/yuin/goldmark-emoji/definition"
)

// implements the emojiDef.Emojis interface but uses case-insensitive search.
// the .children field isn't currently used, but could be used in a future
// implementation of say, emoji packs where a child represents a pack.
type emojis struct {
	list     []emojiDef.Emoji
	names    map[string]*emojiDef.Emoji
	children []emojiDef.Emojis

	emojiMu           sync.Mutex
	emojiDefs         emojiDef.Emojis
	emojiHTML         map[string]string
	emojiModTime      time.Time
	emojiHTMLFormat   string
	emojiHTMLTemplate *template.Template
}

// return a new Emojis set.
func newEmojis(emotes ...emojiDef.Emoji) *emojis {
	loadEmoji()

	e := &emojis{
		list:     emotes,
		names:    map[string]*emojiDef.Emoji{},
		children: []emojiDef.Emojis{},

		emojiMu:           sync.Mutex{},
		emojiHTML:         make(map[string]string),
		emojiModTime:      time.Now(),
		emojiHTMLFormat:   `<img src="{{ .URL }}" class="emoji" alt=":{{ .Name }}:" title=":{{ .Name }}:">`,
		emojiHTMLTemplate: template.Must(template.New("emojiHTML").Parse(emojiHTMLFormat)),
	}

	for i := range e.list {
		emoji := &e.list[i]
		for _, s := range emoji.ShortNames {
			e.names[s] = emoji
		}
	}

	return e
}

func (self *emojis) Get(shortName string) (*emojiDef.Emoji, bool) {
	v, ok := self.names[strings.ToLower(shortName)]
	if ok {
		return v, ok
	}

	for _, child := range self.children {
		v, ok := child.Get(shortName)
		if ok {
			return v, ok
		}
	}

	return nil, false
}

func (self *emojis) Add(emotes emojiDef.Emojis) {
	self.children = append(self.children, emotes)
}

func (self *emojis) Clone() emojiDef.Emojis {
	clone := &emojis{
		list:     self.list,
		names:    self.names,
		children: make([]emojiDef.Emojis, len(self.children)),
	}

	copy(clone.children, self.children)

	return clone
}

var (
	emojiMu sync.Mutex
	// emojiDefs         = newEmojis()
	// emojiHTML         = make(map[string]string)
	emojiModTime      time.Time
	emojiHTMLFormat   = `<img src="{{ .URL }}" class="emoji" alt=":{{ .Name }}:" title=":{{ .Name }}:">`
	emojiHTMLTemplate = template.Must(template.New("emojiHTML").Parse(emojiHTMLFormat))
)

type emojiMeta struct {
	emojiDefs emojiDef.Emojis
	emojiHTML map[string]string
}

func loadEmoji() emojiDef.Emojis {
	modTime, err := data.UpdateEmojiList(false)
	if err != nil {
		return
	}
	emojiArr := make([]emojiDef.Emoji, 0)

	if modTime.After(emojiModTime) {
		emojiMu.Lock()
		defer emojiMu.Unlock()

		emojiHTML := make(map[string]string)

		emojiList := data.GetEmojiList()

		for i := 0; i < len(emojiList); i++ {
			var buf bytes.Buffer
			err := emojiHTMLTemplate.Execute(&buf, emojiList[i])
			if err != nil {
				return
			}
			emojiHTML[strings.ToLower(emojiList[i].Name)] = buf.String()

			emoji := emojiDef.NewEmoji(emojiList[i].Name, nil, strings.ToLower(emojiList[i].Name))
			emojiArr = append(emojiArr, emoji)
		}

	}
	emojiDefs := newEmojis(emojiArr...)
	return emojiDefs
}
