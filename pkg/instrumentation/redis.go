package instrumentation

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	redistrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/redis/go-redis.v9"
)

const (
	redisServiceNameSuffix = "cache-redis"
)

func InstrumentRedis(client redis.UniversalClient, applicationName string) {
	serviceName := fmt.Sprintf("%s-%s", applicationName, redisServiceNameSuffix)

	redistrace.WrapClient(client, redistrace.WithServiceName(serviceName))
}
