package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, args TransferTxParams) (TransferTxresult, error)
}

type SQLstore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLstore{
		db:      db,
		Queries: New(db),
	}
}

func (s *SQLstore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := New(tx)

	err = fn(query)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxresult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one acount to another
// It creates a transfer record, adds account entries and updates account balance within a single database ransaction
func (s *SQLstore) TransferTx(ctx context.Context, args TransferTxParams) (TransferTxresult, error) {
	var result TransferTxresult

	err := s.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}

		//get account -> update its balance
		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = transferMoney(args.FromAccountID, args.ToAccountID, args.Amount, q, ctx)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(args.ToAccountID, args.FromAccountID, -args.Amount, q, ctx)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func transferMoney(fromAccountID int64, toAccountID int64, amount int64, q *Queries, ctx context.Context) (updatedFromAccount Account, updatedToAccount Account, err error) {
	updatedFromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     fromAccountID,
		Amount: -amount,
	})
	if err != nil {
		return
	}

	updatedToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     toAccountID,
		Amount: amount,
	})
	if err != nil {
		return
	}

	return
}
