package database

import (
	"fmt"
	"sync"
	"time"
)

// Database Adapter Interface

type Database interface {
	Name() string
	Connect(config *Config) error
	Close() error
	Ping() error
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) ([]Row, error)
	QueryRow(query string, args ...interface{}) (Row, error)
	Begin() (Transaction, error)
	Migrate(migrations []Migration) error
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
	Timeout  time.Duration
}

type Result struct {
	RowsAffected int64
	LastInsertID int64
}

type Row map[string]interface{}

type Transaction interface {
	Exec(query string, args ...interface{}) (Result, error)
	Query(query string, args ...interface{}) ([]Row, error)
	Commit() error
	Rollback() error
}

type Migration struct {
	Version int
	Name    string
	Up      string
	Down    string
}

// PostgreSQL Adapter

type PostgreSQL struct {
	config *Config
	connected bool
}

func NewPostgreSQL() *PostgreSQL {
	return &PostgreSQL{}
}

func (p *PostgreSQL) Name() string {
	return "postgresql"
}

func (p *PostgreSQL) Connect(config *Config) error {
	p.config = config
	p.connected = true
	return nil
}

func (p *PostgreSQL) Close() error {
	p.connected = false
	return nil
}

func (p *PostgreSQL) Ping() error {
	if !p.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (p *PostgreSQL) Exec(query string, args ...interface{}) (Result, error) {
	if !p.connected {
		return Result{}, fmt.Errorf("not connected")
	}
	return Result{RowsAffected: 1}, nil
}

func (p *PostgreSQL) Query(query string, args ...interface{}) ([]Row, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected")
	}
	return []Row{}, nil
}

func (p *PostgreSQL) QueryRow(query string, args ...interface{}) (Row, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected")
	}
	return Row{}, nil
}

func (p *PostgreSQL) Begin() (Transaction, error) {
	if !p.connected {
		return nil, fmt.Errorf("not connected")
	}
	return &PostgreSQLTx{db: p}, nil
}

func (p *PostgreSQL) Migrate(migrations []Migration) error {
	if !p.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

type PostgreSQLTx struct {
	db *PostgreSQL
}

func (t *PostgreSQLTx) Exec(query string, args ...interface{}) (Result, error) {
	return Result{RowsAffected: 1}, nil
}

func (t *PostgreSQLTx) Query(query string, args ...interface{}) ([]Row, error) {
	return []Row{}, nil
}

func (t *PostgreSQLTx) Commit() error {
	return nil
}

func (t *PostgreSQLTx) Rollback() error {
	return nil
}

// MySQL Adapter

type MySQL struct {
	config *Config
	connected bool
}

func NewMySQL() *MySQL {
	return &MySQL{}
}

func (m *MySQL) Name() string {
	return "mysql"
}

func (m *MySQL) Connect(config *Config) error {
	m.config = config
	m.connected = true
	return nil
}

func (m *MySQL) Close() error {
	m.connected = false
	return nil
}

func (m *MySQL) Ping() error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (m *MySQL) Exec(query string, args ...interface{}) (Result, error) {
	if !m.connected {
		return Result{}, fmt.Errorf("not connected")
	}
	return Result{RowsAffected: 1}, nil
}

func (m *MySQL) Query(query string, args ...interface{}) ([]Row, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	return []Row{}, nil
}

func (m *MySQL) QueryRow(query string, args ...interface{}) (Row, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	return Row{}, nil
}

func (m *MySQL) Begin() (Transaction, error) {
	if !m.connected {
		return nil, fmt.Errorf("not connected")
	}
	return &MySQLTx{db: m}, nil
}

func (m *MySQL) Migrate(migrations []Migration) error {
	if !m.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

type MySQLTx struct {
	db *MySQL
}

func (t *MySQLTx) Exec(query string, args ...interface{}) (Result, error) {
	return Result{RowsAffected: 1}, nil
}

func (t *MySQLTx) Query(query string, args ...interface{}) ([]Row, error) {
	return []Row{}, nil
}

func (t *MySQLTx) Commit() error {
	return nil
}

func (t *MySQLTx) Rollback() error {
	return nil
}

// SQLite Adapter

type SQLite struct {
	config *Config
	connected bool
}

func NewSQLite() *SQLite {
	return &SQLite{}
}

func (s *SQLite) Name() string {
	return "sqlite"
}

func (s *SQLite) Connect(config *Config) error {
	s.config = config
	s.connected = true
	return nil
}

func (s *SQLite) Close() error {
	s.connected = false
	return nil
}

func (s *SQLite) Ping() error {
	if !s.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

func (s *SQLite) Exec(query string, args ...interface{}) (Result, error) {
	if !s.connected {
		return Result{}, fmt.Errorf("not connected")
	}
	return Result{RowsAffected: 1}, nil
}

func (s *SQLite) Query(query string, args ...interface{}) ([]Row, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected")
	}
	return []Row{}, nil
}

func (s *SQLite) QueryRow(query string, args ...interface{}) (Row, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected")
	}
	return Row{}, nil
}

func (s *SQLite) Begin() (Transaction, error) {
	if !s.connected {
		return nil, fmt.Errorf("not connected")
	}
	return &SQLiteTx{db: s}, nil
}

func (s *SQLite) Migrate(migrations []Migration) error {
	if !s.connected {
		return fmt.Errorf("not connected")
	}
	return nil
}

type SQLiteTx struct {
	db *SQLite
}

func (t *SQLiteTx) Exec(query string, args ...interface{}) (Result, error) {
	return Result{RowsAffected: 1}, nil
}

func (t *SQLiteTx) Query(query string, args ...interface{}) ([]Row, error) {
	return []Row{}, nil
}

func (t *SQLiteTx) Commit() error {
	return nil
}

func (t *SQLiteTx) Rollback() error {
	return nil
}

// Database Manager

type Manager struct {
	databases map[string]Database
	mu        sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		databases: make(map[string]Database),
	}
}

func (m *Manager) Register(name string, db Database) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.databases[name] = db
}

func (m *Manager) Get(name string) (Database, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.databases[name]
	return db, ok
}

func (m *Manager) Remove(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.databases, name)
}

func (m *Manager) List() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.databases))
	for name := range m.databases {
		names = append(names, name)
	}
	return names
}

func (m *Manager) ConnectAll(configs map[string]*Config) error {
	for name, config := range configs {
		db, ok := m.Get(name)
		if !ok {
			continue
		}
		if err := db.Connect(config); err != nil {
			return fmt.Errorf("failed to connect to %s: %w", name, err)
		}
	}
	return nil
}

func (m *Manager) CloseAll() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for name, db := range m.databases {
		if err := db.Close(); err != nil {
			return fmt.Errorf("failed to close %s: %w", name, err)
		}
	}
	return nil
}

// Connection Pool

type Pool struct {
	maxOpen    int
	maxIdle    int
	maxLifetime time.Duration
	conns      chan Database
}

func NewPool(maxOpen, maxIdle int, maxLifetime time.Duration) *Pool {
	return &Pool{
		maxOpen:     maxOpen,
		maxIdle:     maxIdle,
		maxLifetime: maxLifetime,
		conns:       make(chan Database, maxOpen),
	}
}

func (p *Pool) Get() Database {
	select {
	case conn := <-p.conns:
		return conn
	default:
		return nil
	}
}

func (p *Pool) Put(conn Database) {
	select {
	case p.conns <- conn:
	default:
		conn.Close()
	}
}

func (p *Pool) Size() int {
	return len(p.conns)
}
