package example

import (
	"context"
	"errors"
	"testing"

	"github.com/godepo/groat"
	"github.com/godepo/groat/integration"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jaswdr/faker/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type RepoDeps struct {
	DB       *pgxpool.Pool `groat:"pgxpool"`
	Faker    faker.Faker
	MockDB   *MockDB
	MockRows *MockRows
}

type State struct {
	Faker         faker.Faker
	Payment       Payment
	ctx           context.Context
	ExpectError   error
	ResultError   error
	ResultPayment Payment
}

var suite *integration.Container[RepoDeps, State, *Repository]

func TestNew(t *testing.T) {
	tcs := suite.Case(t)
	tcs.Go()
}

func ArrangePayment(kind KindOfPayment, paymentState KindOfState) groat.Given[State] {
	return func(t *testing.T, state State) State {
		state.Payment = Payment{
			ID:       uuid.Must(uuid.NewV7()),
			UserID:   state.Faker.UUID().V4(),
			Kind:     kind,
			State:    paymentState,
			Currency: USD,
			Amount:   state.Faker.Int64Between(10, 1000),
		}
		return state
	}
}

func ArrangeError(t *testing.T, state State) State {
	state.ExpectError = errors.New(state.Faker.RandomStringWithLength(10))
	return state
}

func ActExpectErrorAtDBExec(t *testing.T, deps RepoDeps, state State) State {
	t.Helper()
	deps.MockDB.EXPECT().Exec(mock.Anything,
		createQuery,
		state.Payment.ID,
		state.Payment.UserID,
		state.Payment.Amount, state.Payment.Kind,
		state.Payment.Currency,
	).Return(pgconn.CommandTag{}, state.ExpectError)
	return state
}

func InjectMock(sut *Repository) groat.When[RepoDeps, State] {
	return func(t *testing.T, deps RepoDeps, state State) State {
		sut.db = deps.MockDB
		return state
	}
}

func AssertPaymentEquality(db *pgxpool.Pool) groat.Then[State] {
	return func(t *testing.T, state State) {
		t.Helper()
		rows, err := db.Query(state.ctx, "SELECT * FROM groat_pay.payments WHERE id=$1", state.Payment.ID)
		require.NoError(t, err)
		res, err := pgx.CollectOneRow[Payment](rows, pgx.RowToStructByName[Payment])
		require.NoError(t, err)
		exp := state.Payment
		if exp.State == "" {
			exp.State = KOSCreated
		}
		exp.CreatedAt = res.CreatedAt
		exp.ProcessedAt = res.ProcessedAt
		assert.Equal(t, exp, res)
	}
}

func AssertExpectError(t *testing.T, state State) {
	t.Helper()
	assert.ErrorIs(t, state.ResultError, state.ExpectError)
}

func TestRepository_Create(t *testing.T) {
	prepare := func(t *testing.T, kind KindOfPayment) *groat.Case[RepoDeps, State, *Repository] {
		t.Helper()
		tcs := suite.Case(t)
		tcs.
			Given(ArrangePayment(kind, "")).
			Then(AssertPaymentEquality(tcs.Deps.DB))
		return tcs
	}
	t.Run("should be able to create new outgoing payment", func(t *testing.T) {
		tcs := prepare(t, KOPOutgoing)
		require.NoError(t, tcs.SUT.Create(tcs.State.ctx, tcs.State.Payment))
	})

	t.Run("should be able to create new incoming payment", func(t *testing.T) {
		tcs := prepare(t, KOPIncoming)
		require.NoError(t, tcs.SUT.Create(tcs.State.ctx, tcs.State.Payment))
	})

	t.Run("should be able to create new internal payment", func(t *testing.T) {
		tcs := prepare(t, KOPInternal)
		require.NoError(t, tcs.SUT.Create(tcs.State.ctx, tcs.State.Payment))
	})
	t.Run("should be able error when create new payment", func(t *testing.T) {
		tcs := suite.Case(t)
		tcs.
			Given(
				ArrangePayment(KOPOutgoing, ""),
				ArrangeError,
			).
			When(
				InjectMock(tcs.SUT),
				ActExpectErrorAtDBExec,
			).
			Then(AssertExpectError)

		tcs.State.ResultError = tcs.SUT.Create(tcs.State.ctx, tcs.State.Payment)
	})
}

func PreparePayment(db *pgxpool.Pool, kind KindOfPayment, stateKind KindOfState) groat.Given[State] {
	return func(t *testing.T, state State) State {
		t.Helper()
		state = ArrangePayment(kind, stateKind)(t, state)
		if state.Payment.State == "" {
			state.Payment.State = KOSCreated
		}
		payment := state.Payment
		_, err := db.Exec(state.ctx, `
INSERT INTO 
    groat_pay.payments (id, user_id, amount, kind, currency, state ) 
VALUES 
    ($1, $2, $3, $4, $5, $6)`, payment.ID,
			payment.UserID, payment.Amount, payment.Kind, payment.Currency, payment.State)
		require.NoError(t, err)
		return state
	}
}

func ActExpectErrorAtDBQuery(t *testing.T, deps RepoDeps, state State) State {
	t.Helper()
	deps.MockDB.EXPECT().Query(mock.Anything,
		findByIDAndUserIDQuery,
		state.Payment.ID,
		state.Payment.UserID,
	).Return(nil, state.ExpectError)
	return state
}

func RequireNoError(t *testing.T, state State) {
	require.NoError(t, state.ResultError)
}

func ExpectResultEquality(t *testing.T, state State) {
	exp := state.Payment
	exp.CreatedAt = state.ResultPayment.CreatedAt
	exp.ProcessedAt = state.ResultPayment.ProcessedAt

	assert.Equal(t, exp, state.ResultPayment)
}

func TestRepository_FindByIDAndUserID(t *testing.T) {
	prepareCase := func(t *testing.T, kind KindOfPayment,
		stateKind KindOfState,
	) *groat.Case[RepoDeps, State, *Repository] {
		t.Helper()
		tcs := suite.Case(t)
		tcs.Given(PreparePayment(tcs.Deps.DB, KOPOutgoing, stateKind)).
			Then(RequireNoError, ExpectResultEquality)
		return tcs
	}

	t.Run("should be able to find outgoing payment", func(t *testing.T) {
		t.Run("in created state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, "")
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
		t.Run("in declined state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, KOSDeclined)
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
		t.Run("in succeeded state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, KOSSucceeded)
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
		t.Run("in processing state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, KOSProcessing)
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
		t.Run("in delayed state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, KOSDelayed)
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
		t.Run("in aborted state", func(t *testing.T) {
			tcs := prepareCase(t, KOPOutgoing, KOSAborted)
			tcs.State.ResultPayment, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
	})
	t.Run("should be failed", func(t *testing.T) {

		actExpectErrorAtFetchRows := func(t *testing.T, deps RepoDeps, state State) State {
			deps.MockRows.EXPECT().Next().Return(false)
			deps.MockRows.EXPECT().Close().Return()
			deps.MockRows.EXPECT().Err().Return(state.ExpectError)
			return state
		}

		actMockRows := func(t *testing.T, deps RepoDeps, state State) State {
			t.Helper()
			deps.MockDB.EXPECT().Query(mock.Anything, findByIDAndUserIDQuery,
				state.Payment.ID,
				state.Payment.UserID).Return(deps.MockRows, nil)
			return state
		}

		t.Run("when can't create query", func(t *testing.T) {
			tcs := suite.Case(t)
			tcs.Given(ArrangePayment(KOPOutgoing, ""), ArrangeError).
				When(InjectMock(tcs.SUT), ActExpectErrorAtDBQuery).
				Then(AssertExpectError)

			_, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})

		t.Run("when can't fetch result from rows", func(t *testing.T) {
			tcs := suite.Case(t)
			tcs.Given(ArrangePayment(KOPOutgoing, ""), ArrangeError).
				When(InjectMock(tcs.SUT), actMockRows, actExpectErrorAtFetchRows).
				Then(AssertExpectError)

			_, tcs.State.ResultError = tcs.SUT.
				FindByIDAndUserID(tcs.State.ctx, tcs.State.Payment.ID, tcs.State.Payment.UserID)
		})
	})
}
