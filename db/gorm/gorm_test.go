package gorm

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type GormMysqlTest struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
}

func (g *GormMysqlTest) SetupSuite() {
	db, mock, err := sqlmock.New()
	g.Require().NoError(err)

	g.DB, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{
		ConnPool:        db,
		CreateBatchSize: 100000,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})

	// g.DB.DryRun = true
	g.mock = mock
}

func TestGORM(t *testing.T) {
	suite.Run(t, &GormMysqlTest{})
}

var (
	errQueryError = errors.New("query error")
)

type TbTest struct {
	Id   int
	Name string
}

func (g *GormMysqlTest) TestSelect() {
	tests := []struct {
		name    string
		raw     string
		before  func(raw string)
		wantErr error
		after   func(raw string) (any, error)
		wantRes any
	}{
		{
			name: "select *",
			before: func(raw string) {
				g.mock.ExpectQuery(raw).WillReturnError(errQueryError)
			},
			after: func(raw string) (any, error) {
				var res string
				err := g.DB.Raw(raw).Scan(&res).Error
				return res, err
			},
			wantErr: errQueryError,
		},
		{
			name: "select name by scan",
			raw:  "SELECT name FROM `test`",
			before: func(raw string) {
				rows := sqlmock.NewRows([]string{"name"})
				rows.AddRow("Zhangsan")
				g.mock.ExpectQuery(raw).WillReturnRows(rows)
			},
			after: func(raw string) (any, error) {
				var res string
				err := g.DB.Raw(raw).Scan(&res).Error
				return res, err
			},
			wantRes: "Zhangsan",
		},
		{
			name: "select name by take",
			raw:  "SELECT name FROM `test`",
			before: func(raw string) {
				rows := sqlmock.NewRows([]string{"name"})
				rows.AddRow("Zhangsan")
				g.mock.ExpectQuery(raw).WillReturnRows(rows)
			},
			after: func(raw string) (any, error) {
				var res string
				err := g.DB.Raw(raw).Take(&res).Error
				return res, err
			},
			wantRes: "Zhangsan",
		},
		{
			name: "select id,name",
			raw:  "SELECT `id`, `name` FROM `test`",
			before: func(raw string) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				rows.AddRow("999", "Zhangsan")
				g.mock.ExpectQuery(raw).WillReturnRows(rows)
			},
			after: func(raw string) (any, error) {
				res := TbTest{}
				err := g.DB.Raw(raw).Scan(&res).Error
				return res, err
			},
			wantRes: TbTest{
				Id:   999,
				Name: "Zhangsan",
			},
		},
	}

	for _, tt := range tests {
		g.Run(tt.name, func() {
			tt.before(tt.raw)
			res, err := tt.after(tt.raw)
			g.Equal(tt.wantErr, err)
			if err != nil {
				return
			}
			g.Equal(tt.wantRes, res)
		})
	}
}

type TbUser struct {
	Id   int    `gorm:"primaryKey;column:id;type:int;not null"`
	Name string `gorm:"column:name;type:varchar(64);not null"`
	Age  int    `gorm:"column:age;type:int;not null"`
}

// func (g *GormMysqlTest) Test00() {
// 	stmt := g.DB.Session(&gorm.Session{DryRun: true}).First(&TbUser{}, 1).Statement

// 	fmt.Println("sql:", stmt.SQL.String())
// 	fmt.Println("vars:", stmt.Vars)

// 	tx := g.DB.Begin()

// 	stmt = tx.Session(&gorm.Session{DryRun: true}).Model(&TbUser{
// 		Id: 2,
// 	}).Updates(&TbUser{
// 		Name: "hi",
// 		Age:  18,
// 	}).Statement

// 	tx.Commit()

// 	fmt.Println("sql:", stmt.SQL.String())
// 	fmt.Println("vars:", stmt.Vars)
// 	fmt.Println("err:", stmt.Error)
// }

// 测试针对主键有效的构建
func (g *GormMysqlTest) TestPrimarySQL() {
	tests := []struct {
		name    string
		raw     string
		before  func(raw string)
		after   func(raw string) string
		wantSQL string
	}{
		{
			name:   "update",
			before: func(raw string) {},
			after: func(raw string) string {
				return g.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
					tb := &TbUser{
						Id:   2,
						Name: "hi",
						Age:  18,
					}
					return tx.Updates(tb)
				})
			},
			wantSQL: "UPDATE `tb_user` SET `name`='hi',`age`=18 WHERE `id` = 2",
		},
		{
			name:   "take",
			before: func(raw string) {},
			after: func(raw string) string {
				return g.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
					tb := &TbUser{
						Id:   2,
						Name: "hi",
						Age:  18,
					}
					return tx.Take(tb)
				})
			},
			wantSQL: "SELECT * FROM `tb_user` WHERE `tb_user`.`id` = 2 LIMIT 1",
		},
		{
			name:   "scan err",
			before: func(raw string) {},
			after: func(raw string) string {
				return g.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
					tb := &TbUser{
						Id:   2,
						Name: "hi",
						Age:  18,
					}
					return tx.Scan(tb)
				})
			},
			// scan 不会去构建主键, 只是映射作用
			wantSQL: "SELECT * FROM `",
		},
		{
			name:   "scan correct",
			before: func(raw string) {},
			after: func(raw string) string {
				return g.DB.ToSQL(func(tx *gorm.DB) *gorm.DB {
					tb := &TbUser{
						Id:   2,
						Name: "hi",
						Age:  18,
					}
					return tx.Model(tb).Scan(tb)
				})
			},
			// scan 不会去构建主键, 只是映射作用, 借助model去构建
			wantSQL: "SELECT * FROM `tb_user` WHERE `tb_user`.`id` = 2",
		},
	}

	for _, tt := range tests {
		g.Run(tt.name, func() {
			tt.before(tt.raw)

			sql := tt.after(tt.raw)

			g.Equal(tt.wantSQL, sql)
		})
	}
}
