package repository

import (
	"context"
	"github.com/Eugene-Usachev/logger"
	"github.com/jackc/pgx/v5"
)

type PostgresLogger struct {
	logger *logger.FastLogger
}

func NewPostgresLogger(logger *logger.FastLogger) *PostgresLogger {
	return &PostgresLogger{
		logger: logger,
	}
}

func (l *PostgresLogger) TraceQueryStart(ctx context.Context, _ *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	if len(data.Args) > 0 {
		l.logger.FormatInfo("%s, with args: %v\n", data.SQL, data.Args)
	} else {
		l.logger.Info(data.SQL)
	}
	return ctx
}

func (l *PostgresLogger) TraceQueryEnd(_ context.Context, _ *pgx.Conn, data pgx.TraceQueryEndData) {
	if data.Err != nil {
		l.logger.Error("tag: " + data.CommandTag.String() + " err: " + data.Err.Error())
	}
}
