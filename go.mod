module github.com/owncast/owncast

go 1.17

require (
	github.com/amalfra/etag v0.0.0-20190921100247-cafc8de96bc5
	github.com/aws/aws-sdk-go v1.42.0
	github.com/go-fed/activity v1.0.1-0.20210803212804-d866ba75dd0f
	github.com/go-fed/httpsig v1.1.0
	github.com/gorilla/websocket v1.4.2
	github.com/grafov/m3u8 v0.11.1
	github.com/lestrrat-go/file-rotatelogs v2.4.0+incompatible
	github.com/mattn/go-sqlite3 v1.14.10
	github.com/microcosm-cc/bluemonday v1.0.17
	github.com/mssola/user_agent v0.5.3
	github.com/nareix/joy5 v0.0.0-20200712071056-a55089207c88
	github.com/oschwald/geoip2-golang v1.5.0
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/schollz/sqlite3dump v1.3.1
	github.com/shirou/gopsutil v3.21.11+incompatible
	github.com/sirupsen/logrus v1.8.1
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	github.com/yuin/goldmark v1.4.4
	golang.org/x/mod v0.5.1
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	mvdan.cc/xurls v1.1.0
)

require (
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/lestrrat-go/strftime v1.0.4 // indirect
	github.com/mvdan/xurls v1.1.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/tklauser/go-sysconf v0.3.5 // indirect
	github.com/yusufpapurcu/wmi v1.2.2 // indirect
	golang.org/x/crypto v0.0.0-20211215153901-e495a2d5b3d3 // indirect
)

require github.com/SherClockHolmes/webpush-go v1.1.3

require github.com/twilio/twilio-go v0.19.0

require (
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-test/deep v1.0.4 // indirect
	github.com/golang-jwt/jwt v3.2.2+incompatible // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/gorilla/css v1.0.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/oschwald/maxminddb-golang v1.8.0 // indirect
	github.com/tklauser/numcpus v0.2.2 // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sys v0.0.0-20211019181941-9d821ace8654 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/go-fed/activity => github.com/owncast/activity v1.0.1-0.20211229051252-7821289d4026
