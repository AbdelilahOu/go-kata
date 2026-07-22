package main

import (
	"context"
	"errors"
	"fmt"
	"os"
)

type DB interface {
	BeginTx(ctx context.Context) (Tx, error)
	Close() error
}

type Tx interface {
	QueryRows(ctx context.Context) (Rows, error)
	Commit() error
	Rollback() error
}

type Rows interface {
	Next() bool
	Scan() (string, error)
	Close() error
}

type mockRows struct {
	count    int
	current  int
	closeErr error
}

func (r *mockRows) Next() bool {
	if r.current >= r.count {
		return false
	}
	r.current++
	return true
}

func (r *mockRows) Scan() (string, error) {
	return fmt.Sprintf("row-%d", r.current), nil
}

func (r *mockRows) Close() error {
	return r.closeErr
}

type mockTx struct {
	rows        *mockRows
	commitErr   error
	rollbackErr error
}

func (t *mockTx) QueryRows(ctx context.Context) (Rows, error) {
	return t.rows, nil
}

func (t *mockTx) Commit() error {
	return t.commitErr
}

func (t *mockTx) Rollback() error {
	return t.rollbackErr
}

type mockDB struct {
	tx       *mockTx
	beginErr error
	closeErr error
}

func (d *mockDB) BeginTx(ctx context.Context) (Tx, error) {
	if d.beginErr != nil {
		return nil, d.beginErr
	}
	return d.tx, nil
}

func (d *mockDB) Close() error {
	return d.closeErr
}

var openDB = func(_ context.Context, dbURL string) (DB, error) {
	return nil, fmt.Errorf("not implimented: connect to %s", dbURL)
}

func BackupDatabase(ctx context.Context, dbURL, filename string) (err error) {
	f, ferr := os.Create(filename)
	if ferr != nil {
		return fmt.Errorf("open file: %w", ferr)
	}
	defer func() {
		if cerr := f.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("close file: %w", cerr))
		}
	}()

	db, derr := openDB(ctx, dbURL)
	if derr != nil {
		return fmt.Errorf("connect db: %w", derr)
	}
	defer func() {
		if cerr := db.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("close db: %w", cerr))
		}
	}()

	tx, terr := db.BeginTx(ctx)
	if terr != nil {
		return fmt.Errorf("begin tx: %w", terr)
	}

	committed := false
	defer func() {
		if !committed {
			if rerr := tx.Rollback(); rerr != nil {
				err = errors.Join(err, fmt.Errorf("rollback: %w", rerr))
			}
		}
	}()

	rows, qerr := tx.QueryRows(ctx)
	if qerr != nil {
		return fmt.Errorf("query rows: %w", qerr)
	}
	defer func() {
		if cerr := rows.Close(); cerr != nil {
			err = errors.Join(err, fmt.Errorf("close rows: %w", cerr))
		}
	}()

	for rows.Next() {
		line, serr := rows.Scan()
		if serr != nil {
			return fmt.Errorf("scan row: %w", serr)
		}
		if _, werr := fmt.Fprintln(f, line); werr != nil {
			return fmt.Errorf("write row: %w", werr)
		}
	}

	if cerr := tx.Commit(); cerr != nil {
		return fmt.Errorf("commit: %w", cerr)
	}

	committed = true

	return nil
}

func main() {
      fmt.Println("=== Scenario 1: happy path ===")
      openDB = func(ctx context.Context, dbURL string) (DB, error) {
              return &mockDB{
                      tx: &mockTx{
                              rows: &mockRows{count: 3},
                      },
              }, nil
      }
      err := BackupDatabase(context.Background(), "fake://db", "backup1.txt")
      fmt.Println("err:", err)

      fmt.Println("=== Scenario 2: Begin() fails ===")
      openDB = func(ctx context.Context, dbURL string) (DB, error) {
              return &mockDB{
                      beginErr: errors.New("begin failed: connection reset"),
              }, nil
      }
      err = BackupDatabase(context.Background(), "fake://db", "backup2.txt")
      fmt.Println("err:", err)

      fmt.Println("=== Scenario 3: Commit() fails, Rollback() also fails ===")
      openDB = func(ctx context.Context, dbURL string) (DB, error) {
              return &mockDB{
                      tx: &mockTx{
                              rows:        &mockRows{count: 2},
                              commitErr:   errors.New("commit failed: deadlock detected"),
                              rollbackErr: errors.New("rollback failed: connection already closed"),
                      },
              }, nil
      }
      err = BackupDatabase(context.Background(), "fake://db", "backup3.txt")
      fmt.Println("err:", err)
      fmt.Println("as list:", errors.Unwrap(err))

      fmt.Println("=== Scenario 4: 1000 runs, checking for leaks ===")
      for i := 0; i < 1000; i++ {
	            openDB = func(_ context.Context, _ string) (DB, error) {
	                    return &mockDB{tx: &mockTx{rows: &mockRows{count: 1}}}, nil
	            }
	             if err := BackupDatabase(context.Background(), "fake://db", fmt.Sprintf("backup_loop_%d.txt", i)); err != nil {
	                      fmt.Println("unexpected error on run", i, ":", err)
	                      break
	              }
	    }
	    fmt.Println("completed 1000 runs without error")
}
