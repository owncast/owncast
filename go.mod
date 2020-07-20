module github.com/gabek/owncast

go 1.14

require (
	github.com/aws/aws-sdk-go v1.32.1
	github.com/kr/pretty v0.2.0 // indirect
	github.com/mattn/go-sqlite3 v1.14.0
	github.com/mssola/user_agent v0.5.2
	github.com/nareix/joy5 v0.0.0-20200710135721-d57196c8d506
	github.com/radovskyb/watcher v1.0.7
	github.com/sirupsen/logrus v1.6.0
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
	golang.org/x/sys v0.0.0-20200420163511-1957bb5e6d1f // indirect
	golang.org/x/text v0.3.2 // indirect
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/nareix/joy4 v0.0.0 => github.com/Seize/joy4 v0.0.0
