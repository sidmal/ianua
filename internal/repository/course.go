package repository

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/sidmal/ianua/pkg"
	"go.uber.org/zap"
	"time"
)

type courseRepository repository

func newCourseRepository(db *sqlx.DB, cacheLifetime int, logger *zap.Logger) CourseRepositoryInterface {
	repository := &courseRepository{
		db:            db,
		logger:        logger,
		cacheLifetime: cacheLifetime,
		cache:         make(Cached),
	}
	return repository
}

func (m *courseRepository) GetCourseRate(ctx context.Context, from, to string) (float32, error) {
	cacheKey := from + to
	cache, ok := m.cache[cacheKey]
	current := time.Now()

	if ok && cache.expire.After(current) {
		return cache.value.(float32), nil
	}

	rate := float32(0)
	query := "SELECT `value` FROM courses WHERE `from` = $1 AND `to` = $2 AND `date` <= $3 AND `deleted_at` IS NULL " +
		"ORDER BY date"
	args := []interface{}{from, to, current}
	err := m.db.GetContext(ctx, &rate, query, args...)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, pkg.ErrorCourseNotFound
		}

		m.logger.Error(
			pkg.ErrorDatabaseQueryFailed,
			zap.Error(err),
			zap.String(pkg.ErrorDatabaseFieldFilter, query),
			zap.Any(pkg.ErrorDatabaseFieldArguments, args),
		)
		return 0, pkg.ErrorUnknown
	}

	if m.cacheLifetime > 0 {
		m.mx.Lock()
		m.cache[cacheKey] = &CachedValue{
			value:  rate,
			expire: current.Add(time.Duration(m.cacheLifetime) * time.Second),
		}
		m.mx.Unlock()
	}

	return rate, nil
}

func (m *courseRepository) GetAllCached() Cached {
	return m.cache
}

func (m *courseRepository) RemoveCachedByKey(key interface{}) {
	m.mx.Lock()
	delete(m.cache, key)
	m.mx.Unlock()
}

func (m *courseRepository) RemoveAllCached() {
	m.mx.Lock()
	m.cache = make(Cached)
	m.mx.Unlock()
}
