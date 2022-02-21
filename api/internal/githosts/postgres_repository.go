package githosts

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
)

const postgresGitHostsTableName string = "git_hosts"
const postgresGitHostsCredentialsTableName string = "git_hosts_credentials"

type postgresGitHost struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
	Endpoint string `db:"endpoint"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type postgresGitHostCredentials struct {
	ID string `db:"id"`
	Token string `db:"token"`
	Active bool `db:"active"`
	GitHostID string `db:"git_host_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PostgresRepository struct {
	logger zerolog.Logger
	cm connManager
}

func (r *PostgresRepository) toPostgresGitHost(gh *GitHost) *postgresGitHost {
	return &postgresGitHost{
		ID: gh.GetID().String(),
		Name: gh.GetName(),
		Type: string(gh.GetType()),
		Endpoint: gh.GetEndpoint(),
		CreatedAt: gh.GetCreationTime(),
		UpdatedAt: gh.GetLastUpdatedTime(),
	}
}

func (r *PostgresRepository) toPostgresGitHostCredentials(creds *Credentials, gh *GitHost) *postgresGitHostCredentials {
	return &postgresGitHostCredentials{
		ID: creds.GetID().String(),
		Token: creds.GetToken(),
		Active: creds.IsActive(),
		GitHostID: gh.GetID().String(),
		CreatedAt: creds.GetCreationTime(),
		UpdatedAt: creds.GetLastUpdatedTime(),
	}
}

func (r *PostgresRepository) doCreateGitHostCredentals(tx *sqlx.Tx, creds *Credentials, gh *GitHost) error {
	dbGhC := r.toPostgresGitHostCredentials(creds, gh)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, type, endpoint, created_at, updated_at) 
VALUES 
(:id, :name, :type, :endpoint, :created_at, :updated_at);`,
postgresGitHostsCredentialsTableName)

	_, err := tx.NamedExec(insert, dbGhC)

	return err
}

func (r *PostgresRepository) Create(gh *GitHost, cred *Credentials) error {
	conn, err := r.cm.GetConnection()

	if err != nil {
		return err
	}

	tx, err := conn.Beginx()

	if err != nil {
		return err
	}

	dbGh := r.toPostgresGitHost(gh)

	insert := fmt.Sprintf(`
INSERT INTO %s 
(id, name, type, endpoint, created_at, updated_at) 
VALUES 
(:id, :name, :type, :endpoint, :created_at, :updated_at);`,
postgresGitHostsTableName)


	if _, err := tx.NamedExec(insert, dbGh); err != nil {
		return err
	}

	if err := r.doCreateGitHostCredentals(tx, cred, gh); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil  {
			return err
		}

		return err
	}

	return nil
}