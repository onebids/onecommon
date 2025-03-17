module github.com/onebids/onecommon

go 1.22.4

replace github.com/apache/thrift => github.com/apache/thrift v0.13.0

require (
	aidanwoods.dev/go-paseto v1.5.2
	github.com/apache/thrift v0.20.0
	github.com/bytedance/gopkg v0.1.1
	github.com/bytedance/sonic v1.12.2
	github.com/cloudwego/hertz v0.9.3
	github.com/cloudwego/kitex v0.11.3
	github.com/hashicorp/consul/api v1.26.1
	github.com/hertz-contrib/gzip v0.0.3
	github.com/hertz-contrib/paseto v0.0.0-20230508023022-71af6635a26c
	github.com/jinzhu/copier v0.4.0
	github.com/kitex-contrib/monitor-prometheus v0.2.0
	github.com/kitex-contrib/obs-opentelemetry v0.2.9
	github.com/kitex-contrib/obs-opentelemetry/logging/logrus v0.0.0-20241120035129-55da83caab1b
	github.com/kitex-contrib/obs-opentelemetry/logging/zap v0.0.0-20241120035129-55da83caab1b
	github.com/kitex-contrib/registry-consul v0.1.0
	github.com/prometheus/client_golang v1.20.4
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.30.0
	go.opentelemetry.io/otel/sdk v1.32.0
	go.opentelemetry.io/otel/trace v1.32.0
	go.uber.org/zap v1.27.0
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.5.7
	gorm.io/gorm v1.25.12
	gorm.io/plugin/opentelemetry v0.1.11
	github.com/redis/go-redis/v9 v9.6.1
)
