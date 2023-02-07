//go:build !inmemory
// +build !inmemory

package repository

import (
	"database/sql"
	"pomo2/pomodoro"
	"sync"
	"time"
)

const createTableInterval = `create table if not exists "interval" (
"id" integer,
"start_time" datetime not null,
"planned_duration" integer default 0,
"actual_duration" integer default 0,
"category" text not null,
"state" integer default 1,
primary key("id")
);
`

type dbRepo struct {
	db *sql.DB
	sync.RWMutex
}

func NewSQLite3Repo(dbfile string) (*dbRepo, error) {
	db, err := sql.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(1)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec(createTableInterval); err != nil {
		return nil, err
	}

	return &dbRepo{
		db: db,
	}, nil
}

func (r *dbRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	// prepare insert statement
	insStmt, err := r.db.Prepare("insert into interval values(null, ?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer insStmt.Close()

	res, err := insStmt.Exec(i.StartTime, i.PlannedDuration,
		i.ActualDuration, i.Category, i.State)
	if err != nil {
		return 0, err
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return id, nil
}

// Update entry in the repository
func (r *dbRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()

	// prepare update statement
	updStmt, err := r.db.Prepare(
		"update interval set start_time=?, actual_duration=?, state=? where id=?")
	if err != nil {
		return err
	}
	defer updStmt.Close()

	// exec update statement
	res, err := updStmt.Exec(i.StartTime, i.ActualDuration, i.State, i.ID)
	if err != nil {
		return err
	}

	// update results
	_, err = res.RowsAffected()
	return err
}

// ByID search item in the repository by ID
func (r *dbRepo) ByID(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	row := r.db.QueryRow("select * from interval where id=?", id)
	i := pomodoro.Interval{}
	err := row.Scan(&i.ID, &i.StartTime, &i.PlannedDuration,
		&i.ActualDuration, &i.Category, &i.State)

	return i, err
}

// Last search the last item in the repository
func (r *dbRepo) Last() (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	last := pomodoro.Interval{}
	err := r.db.QueryRow("select * from INTERVAL order by id desc limit 1").Scan(
		&last.ID, &last.StartTime, &last.PlannedDuration,
		&last.ActualDuration, &last.Category, &last.State,
	)

	if err == sql.ErrNoRows {
		return last, pomodoro.ErrNoIntervals
	}
	if err != nil {
		return last, err
	}
	return last, nil
}

// Breaks search last n items of type break in the repository
func (r *dbRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()

	stmt := `select * from interval where category like '%Break'
order by id desc limit ?`
	rows, err := r.db.Query(stmt, n)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var data []pomodoro.Interval
	for rows.Next() {
		i := pomodoro.Interval{}
		err = rows.Scan(&i.ID, &i.StartTime, &i.PlannedDuration,
			&i.ActualDuration, &i.Category, &i.State)
		if err != nil {
			return nil, err
		}
		data = append(data, i)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return data, nil
}

// CategorySummary return a daily summary
func (r *dbRepo) CategorySummary(day time.Time, filter string) (time.Duration, error) {
	r.RLock()
	defer r.RUnlock()

	stmt := `select sum(actual_duration) from interval
where category like ?
and strftime('%Y-%m-%d', start_time, 'localtime')=
strftime('%Y-%m-%d', ?, 'localtime')`

	var ds sql.NullInt64
	err := r.db.QueryRow(stmt, filter, day).Scan(&ds)

	var d time.Duration
	if ds.Valid {
		d = time.Duration(ds.Int64)
	}

	return d, err
}
