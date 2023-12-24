package userrepo

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"gitlab.com/AntYats/go_project/internal/model"
	"gitlab.com/AntYats/go_project/internal/repo"
	"log"
)

type userRepo struct {
	pgxPool *pgxpool.Pool
}

func (r *userRepo) conn(ctx context.Context) Conn {
	if tx, ok := ctx.Value(repo.CtxKeyTx).(pgx.Tx); ok {
		return tx
	}

	return r.pgxPool
}

func (r *userRepo) WithNewTx(ctx context.Context, f func(ctx context.Context) error) error {
	return r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		return f(context.WithValue(ctx, repo.CtxKeyTx, tx))
	})
}

func (r *userRepo) GetDrivers(ctx context.Context, radius float64, coords *model.UserData) ([]model.User, error) {
	rows, err := r.conn(ctx).Query(ctx, `SELECT user_id FROM users WHERE POWER(($1 - lat), 2) + POWER(($2 - lng), 2)) <= ($1 * $1)`, radius, coords.Lat, coords.Lng)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]model.User, 0)

	for rows.Next() {
		var user model.User

		if err := rows.Scan(&user.UserId); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepo) ChangeDriverInfo(ctx context.Context, id string, coords *model.UserData) error {
	_, err := r.conn(ctx).Exec(ctx, "UPDATE users SET lag = $1, lng = $2 WHERE user_id = $3", coords.Lat, coords.Lng, id)

	if err != nil {
		return err
	}

	return nil
}

func New(pgxPool *pgxpool.Pool) (repo.User, error) {
	r := &userRepo{
		pgxPool: pgxPool,
	}

	ctx := context.Background()

	err := r.pgxPool.BeginFunc(ctx, func(tx pgx.Tx) error {
		if _, err := r.GetDrivers(ctx, 0, &model.UserData{Lng: 0, Lat: 0}); err != nil {
			log.Fatal(err.Error())
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return r, nil
}
