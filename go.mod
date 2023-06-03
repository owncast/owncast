module github.com/owncast/owncast

go 1.20

require (
	github.com/aws/aws-sdk-go v1.44.273
	github.com/go-fed/activity v1.0.1-0.20210803212804-d866ba75dd0f
	github.com/go-fed/httpsig v1.1.0
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/gorilla/websocket v1.5.0
	github.com/grafov/m3u8 v0.12.0
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/microcosm-cc/bluemonday v1.0.24
	github.com/nareix/joy5 v0.0.0-20210317075623-2c912ca30590
	github.com/oschwald/geoip2-golang v1.8.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/schollz/sqlite3dump v1.3.1
	github.com/sirupsen/logrus v1.9.2
	github.com/teris-io/shortid v0.0.0-20220617161101-71ec9f2aa569
	github.com/yuin/goldmark v1.5.4
	golang.org/x/mod v0.10.0
	golang.org/x/time v0.3.0
)

require (
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/tklauser/go-sysconf v0.3.11 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/net v0.10.0
	golang.org/x/sys v0.8.0 // indirect
)

require github.com/prometheus/client_golang v1.15.1

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/klauspost/compress v1.16.4 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.9.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
)

require (
	github.com/nakabonne/tstorage v0.3.6
	github.com/shirou/gopsutil/v3 v3.23.3
)

require github.com/SherClockHolmes/webpush-go v1.2.0

require (
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/oschwald/maxminddb-golang v1.10.0 // indirect
	github.com/shoenig/go-m1cpu v0.1.4 // indirect
)

require (
	github.com/CAFxX/httpcompression v0.0.8
	github.com/andybalholm/cascadia v1.3.2
	gopkg.in/evanphx/json-patch.v5 v5.6.0
	mvdan.cc/xurls v1.1.0
)

replace github.com/go-fed/activity => github.com/owncast/activity v1.0.1-0.20211229051252-7821289d4026
