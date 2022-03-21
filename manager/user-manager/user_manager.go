package user_manager

import (
	"database/sql"
	"github.com/Masterminds/squirrel"
	"landing-page/models"
	"landing-page/pkg/constant"
	"landing-page/pkg/sql-tool"
)

type UserManager struct {
	DB *sql.DB
}

var userTable = sql_tool.Table{
	Name:      "user",
	AIColumns: []string{"id"},
}

func NewUserManager(db *sql.DB) UserManager {
	return UserManager{
		DB: db,
	}
}

func (m *UserManager) GetUser(request *models.UserRequest) (users []*models.User, err error) {
	sqlTool := sql_tool.NewSqlTool(m.DB, userTable, models.User{})
	qb := squirrel.
		Select(sqlTool.GetQueryColumnsList()...).
		From(userTable.Name)
	if request.Email != "" {
		qb = qb.Where(squirrel.Eq{"email": request.Email})
	}
	if request.PhoneNumber != "" {
		qb = qb.Where(squirrel.Eq{"phone_number": request.PhoneNumber})
	}
	if request.Id > 0 {
		qb = qb.Where(squirrel.Eq{"id": request.Id})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	err = sqlTool.Query(&users, query, args)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, sql.ErrNoRows
	}

	return users, nil
}

func (m *UserManager) CreateUser(request *models.User) (user *models.User, err error) {
	sqlTool := sql_tool.NewSqlTool(m.DB, userTable, models.User{})
	sqlTool.SetType(constant.CREATE)
	qb := squirrel.Insert(userTable.Name).
		Columns(sqlTool.GetQueryColumnsList()...).
		Values(sqlTool.GetFillValue(request))
	query, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}
	err = sqlTool.Query(&user, query, args)
	return user, err
}
