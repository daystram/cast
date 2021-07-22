module github.com/daystram/cast/cast-be

go 1.15

require (
	github.com/astaxie/beego v1.12.3
	github.com/aws/aws-sdk-go v1.34.28
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/disintegration/imaging v1.6.2
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.4.3 // indirect
	github.com/golang/snappy v0.0.2 // indirect
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/gorilla/websocket v1.4.2
	github.com/klauspost/compress v1.11.7 // indirect
	github.com/nareix/joy4 v0.0.0-20200507095837-05a4ffbb5369
	github.com/shiena/ansicolor v0.0.0-20200904210342-c7312218db18 // indirect
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	go.mongodb.org/mongo-driver v1.5.1
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777 // indirect
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a // indirect
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace github.com/gorilla/websocket => github.com/daystram/websocket v1.4.3
