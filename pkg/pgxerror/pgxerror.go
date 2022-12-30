package pgxerror

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/rteja-library3/rapperror"
)

type Functions struct {
	ConstraintFilter func(pgErr *pgconn.PgError) *rapperror.AppError
}

type DetailFunctions func(pgErr *pgconn.PgError) map[string]interface{}

func FromPgxError(err error, msg, src string) *rapperror.AppError {
	defDetailFn := func(pgErr *pgconn.PgError) map[string]interface{} {
		return map[string]interface{}{
			"detail": pgErr.Detail,
		}
	}

	return FromPgxErrorWithConstraintFilter(err, msg, src, Functions{}, defDetailFn)
}

func FromPgxErrorWithConstraintFilter(err error, msg, src string, fnFilter Functions, fnDetail DetailFunctions) *rapperror.AppError {
	var pgErr *pgconn.PgError

	if err == pgx.ErrNoRows {
		return rapperror.ErrNotFound(
			"",
			"Data not found",
			src,
			nil,
		)
	}

	if errors.As(err, &pgErr) {
		// if pgerrcode.IsSyntaxErrororAccessRuleViolation(pgErr.Code) {
		// 	return rapperror.ErrBadRequest(
		// 		"",
		// 		"Invalid role",
		// 		src,
		// 		fnDetail(pgErr),
		// 	)
		// }

		switch pgErr.Code {
		case pgerrcode.ForeignKeyViolation:
			return rapperror.ErrNotFound(
				"",
				"Data not found",
				src,
				fnDetail(pgErr),
			)
		case pgerrcode.UniqueViolation:
			return rapperror.ErrConflict(
				"",
				"Duplicate data",
				src,
				fnDetail(pgErr),
			)
		}

		return rapperror.ErrInternalServerError(
			"",
			msg,
			src,
			fnDetail(pgErr),
		)
	}

	return rapperror.ErrInternalServerError(
		"",
		msg,
		src,
		nil,
	)
}
