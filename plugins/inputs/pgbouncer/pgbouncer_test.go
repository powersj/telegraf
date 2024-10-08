package pgbouncer

import (
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/influxdata/telegraf/config"
	"github.com/influxdata/telegraf/plugins/common/postgresql"
	"github.com/influxdata/telegraf/testutil"
)

func TestPgBouncerGeneratesMetricsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	postgresServicePort := "5432"
	pgBouncerServicePort := "6432"

	backend := testutil.Container{
		Image:        "postgres:alpine",
		ExposedPorts: []string{postgresServicePort},
		Env: map[string]string{
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	require.NoError(t, backend.Start(), "failed to start container")
	defer backend.Terminate()

	container := testutil.Container{
		Image:        "z9pascal/pgbouncer-container:1.23.1-latest",
		ExposedPorts: []string{pgBouncerServicePort},
		Env: map[string]string{
			"PG_ENV_POSTGRESQL_USER": "pgbouncer",
			"PG_ENV_POSTGRESQL_PASS": "pgbouncer",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(pgBouncerServicePort)),
			wait.ForLog("LOG process up"),
		),
	}
	require.NoError(t, container.Start(), "failed to start container")
	defer container.Terminate()

	addr := fmt.Sprintf(
		"host=%s user=pgbouncer password=pgbouncer dbname=pgbouncer port=%s sslmode=disable",
		container.Address,
		container.Ports[pgBouncerServicePort],
	)

	p := &PgBouncer{
		Config: postgresql.Config{
			Address:     config.NewSecret([]byte(addr)),
			IsPgBouncer: true,
		},
	}
	require.NoError(t, p.Init())

	var acc testutil.Accumulator
	require.NoError(t, p.Start(&acc))
	defer p.Stop()
	require.NoError(t, p.Gather(&acc))

	intMetricsPgBouncer := []string{
		"total_received",
		"total_sent",
		"total_query_time",
		"avg_query_count",
		"avg_query_time",
		"avg_wait_time",
	}

	intMetricsPgBouncerPools := []string{
		"cl_active",
		"cl_waiting",
		"sv_active",
		"sv_idle",
		"sv_used",
		"sv_tested",
		"sv_login",
		"maxwait",
	}

	metricsCounted := 0

	for _, metric := range intMetricsPgBouncer {
		require.True(t, acc.HasInt64Field("pgbouncer", metric))
		metricsCounted++
	}

	for _, metric := range intMetricsPgBouncerPools {
		require.True(t, acc.HasInt64Field("pgbouncer_pools", metric))
		metricsCounted++
	}

	require.Positive(t, metricsCounted)
	require.Equal(t, len(intMetricsPgBouncer)+len(intMetricsPgBouncerPools), metricsCounted)
}

func TestPgBouncerGeneratesMetricsIntegrationShowCommands(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	postgresServicePort := "5432"
	pgBouncerServicePort := "6432"

	backend := testutil.Container{
		Image:        "postgres:alpine",
		ExposedPorts: []string{postgresServicePort},
		Env: map[string]string{
			"POSTGRES_HOST_AUTH_METHOD": "trust",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
	}
	err := backend.Start()
	require.NoError(t, err, "failed to start container")
	defer backend.Terminate()

	container := testutil.Container{
		Image:        "z9pascal/pgbouncer-container:1.23.1-latest",
		ExposedPorts: []string{pgBouncerServicePort},
		Env: map[string]string{
			"PG_ENV_POSTGRESQL_USER": "pgbouncer",
			"PG_ENV_POSTGRESQL_PASS": "pgbouncer",
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort(nat.Port(pgBouncerServicePort)),
			wait.ForLog("LOG process up"),
		),
	}
	err = container.Start()
	require.NoError(t, err, "failed to start container")
	defer container.Terminate()

	addr := fmt.Sprintf(
		"host=%s user=pgbouncer password=pgbouncer dbname=pgbouncer port=%s sslmode=disable",
		container.Address,
		container.Ports[pgBouncerServicePort],
	)

	p := &PgBouncer{
		Config: postgresql.Config{
			Address:     config.NewSecret([]byte(addr)),
			IsPgBouncer: true,
		},
		ShowCommands: []string{"pools", "lists", "databases"},
	}
	require.NoError(t, p.Init())

	var acc testutil.Accumulator
	require.NoError(t, p.Start(&acc))
	defer p.Stop()
	require.NoError(t, p.Gather(&acc))

	intMetricsPgBouncerPools := []string{
		"cl_active",
		"cl_waiting",
		"sv_active",
		"sv_idle",
		"sv_used",
		"sv_tested",
		"sv_login",
		"maxwait",
	}

	intMetricsPgBouncerLists := []string{
		"databases",
		"users",
		"pools",
		"free_clients",
		"used_clients",
		"login_clients",
		"free_servers",
		"used_servers",
		"dns_names",
		"dns_zones",
		"dns_queries",
	}

	intMetricsPgBouncerDatabases := []string{
		"pool_size",
		"min_pool_size",
		"reserve_pool",
		"max_connections",
		"current_connections",
		"paused",
		"disabled",
	}

	metricsCounted := 0

	for _, metric := range intMetricsPgBouncerPools {
		require.True(t, acc.HasInt64Field("pgbouncer_pools", metric))
		metricsCounted++
	}

	for _, metric := range intMetricsPgBouncerLists {
		require.True(t, acc.HasInt64Field("pgbouncer_lists", metric))
		metricsCounted++
	}

	for _, metric := range intMetricsPgBouncerDatabases {
		require.True(t, acc.HasInt64Field("pgbouncer_databases", metric))
		metricsCounted++
	}

	require.Positive(t, metricsCounted)
	require.Equal(t, len(intMetricsPgBouncerPools)+len(intMetricsPgBouncerLists)+len(intMetricsPgBouncerDatabases), metricsCounted)
}
