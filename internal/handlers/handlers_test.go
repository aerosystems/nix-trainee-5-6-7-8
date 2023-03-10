package handlers

import (
	"database/sql"
	"github.com/aerosystems/nix-trainee-5-6-7-8/internal/models"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/aerosystems/nix-trainee-5-6-7-8/internal/storage"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DataStruct struct {
	Post    models.Post
	Comment models.Comment
}

type NewDataStruct struct {
	Post    models.Post
	Comment models.Comment
}

type testTable []struct {
	Name                     string
	Data                     DataStruct
	NewData                  NewDataStruct
	RequestPath              string
	RequestBody              string
	RequestHeaderContentType string
	RequestHeaderAccept      string
	ResponseStatusCode       int
	ResponseBody             string
}

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock

	basehandler BaseHandler
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	dialector := mysql.New(mysql.Config{
		Conn:                      db,
		DriverName:                "mysql",
		SkipInitializeWithVersion: true,
	})

	s.DB, err = gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	require.NoError(s.T(), err)

	s.basehandler.commentRepo = storage.NewCommentRepo(s.DB)
	s.basehandler.postRepo = storage.NewPostRepo(s.DB)
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}
