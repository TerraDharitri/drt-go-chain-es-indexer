package transport

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/TerraDharitri/drt-go-chain-es-indexer/core"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/core/request"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/metrics"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/mock"
	"github.com/stretchr/testify/require"
)

func TestNewMetricsTransport(t *testing.T) {
	t.Parallel()

	transportHandler, err := NewMetricsTransport(nil)
	require.Nil(t, transportHandler)
	require.Equal(t, core.ErrNilMetricsHandler, err)

	metricsHandler := metrics.NewStatusMetrics()
	transportHandler, err = NewMetricsTransport(metricsHandler)
	require.Nil(t, err)
	require.NotNil(t, transportHandler)
}

func TestMetricsTransport_NilRequest(t *testing.T) {
	metricsHandler := metrics.NewStatusMetrics()
	transportHandler, _ := NewMetricsTransport(metricsHandler)

	_, err := transportHandler.RoundTrip(nil)
	require.Equal(t, errNilRequest, err)
}

func TestMetricsTransport_RoundTripNilResponseShouldWork(t *testing.T) {
	t.Parallel()

	metricsHandler := metrics.NewStatusMetrics()
	transportHandler, _ := NewMetricsTransport(metricsHandler)

	testErr := errors.New("test")
	transportHandler.transport = &mock.TransportMock{
		Response: nil,
		Err:      testErr,
	}

	testTopic := "test"
	contextWithValue := context.WithValue(context.Background(), request.ContextKey, testTopic)
	req, _ := http.NewRequestWithContext(contextWithValue, http.MethodGet, "dummy", bytes.NewBuffer([]byte("test")))

	_, _ = transportHandler.RoundTrip(req)

	metricsMap := metricsHandler.GetMetrics()
	require.Equal(t, uint64(1), metricsMap[testTopic].OperationsCount)
	require.Equal(t, uint64(1), metricsMap[testTopic].TotalErrorsCount)
	require.Equal(t, uint64(4), metricsMap[testTopic].TotalData)
}

func TestMetricsTransport_RoundTrip(t *testing.T) {
	t.Parallel()

	metricsHandler := metrics.NewStatusMetrics()
	transportHandler, _ := NewMetricsTransport(metricsHandler)

	transportHandler.transport = &mock.TransportMock{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
		Err: nil,
	}

	testTopic := "test"
	contextWithValue := context.WithValue(context.Background(), request.ContextKey, testTopic)
	req, _ := http.NewRequestWithContext(contextWithValue, http.MethodGet, "dummy", bytes.NewBuffer([]byte("test")))

	_, _ = transportHandler.RoundTrip(req)

	metricsMap := metricsHandler.GetMetrics()
	require.Equal(t, uint64(1), metricsMap[testTopic].OperationsCount)
	require.Equal(t, uint64(0), metricsMap[testTopic].TotalErrorsCount)
	require.Equal(t, uint64(4), metricsMap[testTopic].TotalData)
}

func TestMetricsTransport_RoundTripNoValueInContextShouldNotAddMetrics(t *testing.T) {
	t.Parallel()

	metricsHandler := metrics.NewStatusMetrics()
	transportHandler, _ := NewMetricsTransport(metricsHandler)

	transportHandler.transport = &mock.TransportMock{
		Response: &http.Response{
			StatusCode: http.StatusOK,
		},
		Err: nil,
	}

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet, "dummy", bytes.NewBuffer([]byte("test")))
	_, _ = transportHandler.RoundTrip(req)

	metricsMap := metricsHandler.GetMetrics()
	require.Len(t, metricsMap, 0)
}
