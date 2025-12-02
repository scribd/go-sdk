package instrumentation

import (
	"fmt"

	redistrace "github.com/DataDog/dd-trace-go/contrib/redis/go-redis.v9/v2"
	"github.com/redis/go-redis/v9"
)

const (
	redisServiceNameSuffix = "cache-redis"
)

func InstrumentRedis(client redis.UniversalClient, applicationName string) {
	serviceName := fmt.Sprintf("%s-%s", applicationName, redisServiceNameSuffix)

	redistrace.WrapClient(client, redistrace.WithService(serviceName))
}
