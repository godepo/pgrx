package example

import (
	"context"
	"testing"

	"github.com/godepo/groat/integration"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type RepoDeps struct {
	FirstDB  *pgxpool.Pool `groat:"first_pool"`
	SecondDB *pgxpool.Pool `groat:"second_pool"`
	Faker    faker.Faker
}

type State struct {
	Faker faker.Faker
	ctx   context.Context
}

var suite *integration.Container[RepoDeps, State, *Repository]

func TestNew(t *testing.T) {
	tcs := suite.Case(t)
	require.NotNil(t, tcs.SUT.firstDB)
	require.NotNil(t, tcs.SUT.secondDB)

	rows, err := tcs.SUT.firstDB.Query(t.Context(), "SELECT id FROM first")
	require.NoError(t, err)

	result, err := pgx.CollectOneRow[int](rows, pgx.RowTo[int])
	require.NoError(t, err)
	assert.Equal(t, 1, result)

	rows, err = tcs.SUT.secondDB.Query(t.Context(), "SELECT id FROM second")
	require.NoError(t, err)

	result, err = pgx.CollectOneRow[int](rows, pgx.RowTo[int])
	require.NoError(t, err)
	assert.Equal(t, 2, result)

	_, err = tcs.SUT.firstDB.Query(t.Context(), "SELECT id FROM second")
	require.ErrorContains(t, err, "relation \"second\" does not exist")

	_, err = tcs.SUT.secondDB.Query(t.Context(), "SELECT id FROM first")
	require.ErrorContains(t, err, "relation \"first\" does not exist")
}
