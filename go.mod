module github.com/scribd/go-sdk

go 1.23

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/DATA-DOG/go-txdb v0.1.9
	github.com/DataDog/datadog-go v4.8.2+incompatible
	github.com/aws/aws-sdk-go v1.44.327
	github.com/aws/aws-sdk-go-v2 v1.32.3
	github.com/aws/aws-sdk-go-v2/config v1.28.1
	github.com/aws/aws-sdk-go-v2/credentials v1.17.42
	github.com/aws/aws-sdk-go-v2/service/s3 v1.66.2
	github.com/aws/aws-sdk-go-v2/service/sagemakerruntime v1.32.3
	github.com/aws/aws-sdk-go-v2/service/sfn v1.33.3
	github.com/aws/aws-sdk-go-v2/service/sqs v1.36.3
	github.com/aws/aws-sdk-go-v2/service/sts v1.32.3
	github.com/aws/smithy-go v1.22.0
	github.com/getsentry/sentry-go v0.12.0
	github.com/go-kit/kit v0.13.0
	github.com/go-kit/log v0.2.0
	github.com/go-sql-driver/mysql v1.7.1
	github.com/google/uuid v1.6.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/magefile/mage v1.15.0
	github.com/redis/go-redis/v9 v9.5.1
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/viper v1.10.1
	github.com/statsig-io/go-sdk v1.26.0
	github.com/stretchr/testify v1.9.0
	github.com/twmb/franz-go v1.17.0
	github.com/twmb/franz-go/pkg/kmsg v1.8.0
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.33.0
	gopkg.in/DataDog/dd-trace-go.v1 v1.69.1
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gorm.io/driver/mysql v1.4.6
	gorm.io/driver/sqlite v1.4.4
	gorm.io/gorm v1.25.3
)

require (
	github.com/DataDog/appsec-internal-go v1.8.0 // indirect
	github.com/DataDog/datadog-agent/pkg/obfuscate v0.48.0 // indirect
	github.com/DataDog/datadog-agent/pkg/remoteconfig/state v0.57.0 // indirect
	github.com/DataDog/datadog-go/v5 v5.3.0 // indirect
	github.com/DataDog/go-libddwaf/v3 v3.4.0 // indirect
	github.com/DataDog/go-tuf v1.1.0-0.5.2 // indirect
	github.com/DataDog/gostackparse v0.7.0 // indirect
	github.com/DataDog/sketches-go v1.4.5 // indirect
	github.com/Microsoft/go-winio v0.6.1 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.6.6 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.3.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.6.22 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.1 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.3.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.36.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.187.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/eventbridge v1.35.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.12.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.4.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/endpoint-discovery v1.10.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.12.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.18.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/kinesis v1.32.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sns v1.33.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.24.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.28.3 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/eapache/queue/v2 v2.0.0-20230407133247-75960ed334e4 // indirect
	github.com/ebitengine/purego v0.6.0-alpha.5 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/go-logfmt/logfmt v0.6.0 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/pprof v0.0.0-20230817174616-7a8ec2ada47b // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.7 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/hcl v1.0.1-vault-5 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.17.8 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-sqlite3 v1.14.22 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/outcaste-io/ristretto v0.2.3 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240612014219-fbbf4953d986 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/richardartoul/molecule v1.0.1-0.20240531184615-7ca0df43c0b3 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/secure-systems-lab/go-securesystemslib v0.7.0 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.9.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/statsig-io/ip3country-go v0.2.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/tinylib/msgp v1.2.1 // indirect
	github.com/ua-parser/uap-go v0.0.0-20211112212520-00c877edfe0f // indirect
	go.uber.org/atomic v1.11.0 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sync v0.10.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240108191215-35c7eff3a6b1 // indirect
	gopkg.in/ini.v1 v1.66.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
