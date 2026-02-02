module github.com/owncast/owncast

go 1.24.4

toolchain go1.24.11

require (
	github.com/CAFxX/httpcompression v0.0.9
	github.com/SherClockHolmes/webpush-go v1.4.0
	github.com/TwiN/go-away v1.8.1
	github.com/andybalholm/cascadia v1.3.3
	github.com/aws/aws-sdk-go v1.55.8
	github.com/go-chi/chi/v5 v5.2.3
	github.com/go-fed/activity v1.0.1-0.20220119073622-b14b50eecad0
	github.com/go-fed/httpsig v1.1.0
	github.com/gorilla/websocket v1.5.3
	github.com/grafov/m3u8 v0.12.1
	github.com/hashicorp/go-retryablehttp v0.7.8
	github.com/jellydator/ttlcache/v3 v3.4.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-sqlite3 v1.14.33
	github.com/microcosm-cc/bluemonday v1.0.27
	github.com/mssola/user_agent v0.6.0
	github.com/nakabonne/tstorage v0.3.6
	github.com/nareix/joy5 v0.0.0-20210317075623-2c912ca30590
	github.com/oapi-codegen/runtime v1.1.2
	github.com/oschwald/geoip2-golang v1.13.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.23.2
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/schollz/sqlite3dump v1.3.1
	github.com/shirou/gopsutil/v4 v4.25.8
	github.com/sirupsen/logrus v1.9.4
	github.com/stretchr/testify v1.11.1
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	github.com/yuin/goldmark v1.7.13
	github.com/yuin/goldmark-emoji v1.0.6
	golang.org/x/crypto v0.47.0
	golang.org/x/mod v0.31.0
	golang.org/x/net v0.48.0
	golang.org/x/time v0.14.0
	gopkg.in/evanphx/json-patch.v5 v5.9.11
	mvdan.cc/xurls/v2 v2.6.0
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/ebitengine/purego v0.8.4 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jonboulle/clockwork v0.5.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/lestrrat-go/strftime v1.1.1 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oschwald/maxminddb-golang v1.13.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/power-devops/perfstat v0.0.0-20240221224432-82ca36839d55 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.66.1 // indirect
	github.com/prometheus/procfs v0.16.1 // indirect
	github.com/tklauser/go-sysconf v0.3.15 // indirect
	github.com/tklauser/numcpus v0.10.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.33.0 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/go-fed/activity => github.com/owncast/activity v1.0.1-0.20260122170223-675f6eb53e71
