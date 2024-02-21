package postgres

import (
	"context"
	"fmt"
	"strings"

	"kitab/config"
	"kitab/storage"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database"          //database is needed for migration
	_ "github.com/golang-migrate/migrate/v4/database/postgres" //postgres is used for database
	_ "github.com/golang-migrate/migrate/v4/source/file"       //file is needed for migration url
	_ "github.com/lib/pq"
)

type BookStore struct {
	pool *pgxpool.Pool
	cfg  config.Config
}

func New(ctx context.Context, cfg config.Config) (storage.IStorage, error) {
	url := fmt.Sprintf(
		`postgres://%s:%s@%s:%s/%s?sslmode=disable`,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDB,
	)

	poolConfig, err := pgxpool.ParseConfig(url)
	if err != nil {
		fmt.Println("error while parsing config", err)
		return nil, err
	}

	poolConfig.MaxConns = 100

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		fmt.Println("error while connecting to db", err)
		return nil, err
	}

	//migration
	m, err := migrate.New("file://migrations/postgres/", url)
	if err != nil {
		fmt.Println("error while migrating", err)
		return nil, err
	}

	if err = m.Up(); err != nil {
		fmt.Println("migration up", err)
		if !strings.Contains(err.Error(), "no change") {
			fmt.Println("entered")
			version, dirty, err := m.Version()
			fmt.Println("version and dirty", "version", version, "dirty", dirty)
			if err != nil {
				fmt.Println("err in checking version and dirty", err)
				return nil, err
			}

			if dirty {
				version--
				if err = m.Force(int(version)); err != nil {
					fmt.Println("ERR in making force", err)
					return nil, err
				}
			}
			fmt.Println("WARNING in migrating", err)
			return nil, err
		}
	}

	return BookStore{
		pool: pool,
		cfg:  cfg,
	}, nil
}

func (b BookStore) Close() {
	b.pool.Close()
}

func (b BookStore) Book() storage.IBookStorage {
	return NewBookRepo(b.pool)
}
