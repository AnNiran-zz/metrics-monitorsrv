package counter

import (
	"testing"
	"reflect"
	"math"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/arbitrary"
)

func createTestMetric() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&TestMetric{}), map[string]gopter.Gen {
		"Data:": gen.Struct(reflect.TypeOf(&TestMetricData{}), map[string]gopter.Gen {
			"Name":     gen.AlphaString(),
			"ValueSet": gen.MapOf(gen.Int64Range(0, math.MaxInt64), gen.IntRange(0, 10000)),
		}),
	})
}

func createTestMetrics() gopter.Gen {
	return gen.Struct(reflect.TypeOf(&TestMetrics{}), map[string]gopter.Gen {
		"LastRefreshTime": gen.Int64Range(0, math.MaxInt64),
		"Data"           : gen.Struct(reflect.TypeOf(&TestMetricsData{}), map[string]gopter.Gen {
			"FirstRecord" : gen.Int64Range(0, math.MaxInt64),
			"LastRecord"  : gen.Int64Range(0, math.MaxInt64),
			"List"        : gen.MapOf(gen.AlphaString(), gen.Struct(reflect.TypeOf(&TestMetric{}), map[string]gopter.Gen {
				"Data:": gen.Struct(reflect.TypeOf(&TestMetricData{}), map[string]gopter.Gen {
					"Name":     gen.AlphaString(),
					"ValueSet": gen.MapOf(gen.Int64Range(0, math.MaxInt64), gen.IntRange(0, 10000)),
				}),
			})),
		}),
	})
}

func testMetrics(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.Rng.Seed(1234)

	properties := gopter.NewProperties(parameters)
	arbitraries := arbitrary.DefaultArbitraries()
	arbitraries.RegisterGen(createTestMetric())
	arbitraries.RegisterGen(createTestMetrics())

	properties.Property("metric all empty", arbitraries.ForAll(
		func (tm *TestMetric) bool {
			for _, value := range tm.Data.ValueSet {
				if value != 0 {
					return false
				}
			}
			return true
		},
	))

	properties.Property("metrics all empty", arbitraries.ForAll(
		func (tm *TestMetrics) bool {
			for _, value := range tm.Data.List {
				if value != nil {
					return false
				}
			}
			return true
		},
	))

	properties.Run(gopter.ConsoleReporter(true))
	properties.TestingRun(t)
}
