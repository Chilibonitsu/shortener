package postgres

import (
	"database/sql"
	"fmt"
	"short-url-api/storage"

	"errors"

	"github.com/jmoiron/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type StorageSQLX struct {
	db *sqlx.DB
}
type ConnectorDBPostgre struct {
	db *gorm.DB
}

func NewConnectorPostgreSQL() (*ConnectorDBPostgre, error) {
	const op = "storage.postgres.NewConnectorPostgreSQL"

	dsn := "host=localhost port=5432 user=postgres password=aboba dbname=urlsdb sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{PrepareStmt: true})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	db.AutoMigrate(&storage.Url{})
	return &ConnectorDBPostgre{db: db}, nil
}
func (c *ConnectorDBPostgre) Show() {
	var url storage.Url
	query := c.db.First(&url)
	fmt.Println(url, query)
}
func (c *ConnectorDBPostgre) SaveUrl(urlToSave string, alias string) (storage.Url, error) {
	const op = "storage.postgres.SaveUrl"
	url := storage.Url{
		Url:   urlToSave,
		Alias: alias,
	}
	//не обработан кейс urlExists
	query := c.db.Create(&url)

	if query.Error != nil {
		return url, fmt.Errorf("%s: %w", op, query.Error)
	}

	return url, nil
}
func (c *ConnectorDBPostgre) GetUrl(alias string) (string, error) {
	const op = "storage.postgres.Geturl"
	url := storage.Url{Alias: alias}

	query := c.db.First(&url, "alias = ?", url.Alias)
	if errors.Is(query.Error, sql.ErrNoRows) {
		return "", storage.ErrURLNotFound
	}

	if query.Error != nil {
		return "", fmt.Errorf("%s: %w", op, query.Error)
	}
	return url.Url, nil
}
func (c *ConnectorDBPostgre) GetAll() ([]storage.Url, error) {
	const op = "storage.postgres.GetAll"
	urls := []storage.Url{}

	query := c.db.Find(&urls)

	if errors.Is(query.Error, sql.ErrNoRows) {
		return nil, storage.ErrURLNotFound
	}

	if query.Error != nil {
		return nil, fmt.Errorf("%s: %w", op, query.Error)
	}
	return urls, nil
}
func (c *ConnectorDBPostgre) GetById(id int64) (storage.Url, error) {
	const op = "storage.postgres.GetById"
	url := storage.Url{}

	query := c.db.First(&url, id)

	if errors.Is(query.Error, gorm.ErrRecordNotFound) {
		return url, storage.ErrURLNotFound
	}

	if query.Error != nil {
		return url, fmt.Errorf("%s: %w", op, query.Error)
	}
	return url, nil

}
func (c *ConnectorDBPostgre) DeleteUrl(id int64) error {
	const op = "storage.postgres.DeleteUrl"
	url := storage.Url{ID: id}

	query := c.db.Delete(&url, "ID = ?", url.ID)

	if query.RowsAffected == 0 {
		return storage.ErrNoRows
	}

	if query.Error != nil {
		return fmt.Errorf("%s: %w", op, query.Error)
	}

	return nil
}

func NewStorageSQLX(db *sqlx.DB) *StorageSQLX {
	return &StorageSQLX{db: db}
}
func (s *StorageSQLX) ConnSQLX() (*sqlx.DB, error) {
	const op = "storage.postgres.Conn"

	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return db, nil
}
