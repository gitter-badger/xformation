package migrations

import (
	"testing"

	"github.com/go-xorm/xorm"
	. "github.com/xformation/xformation/pkg/services/sqlstore/migrator"
	"github.com/xformation/xformation/pkg/services/sqlstore/sqlutil"

	. "github.com/smartystreets/goconvey/convey"
	//"github.com/xformation/xformation/pkg/log"
)

var indexTypes = []string{"Unknown", "INDEX", "UNIQUE INDEX"}

func TestMigrations(t *testing.T) {
	testDBs := []sqlutil.TestDB{
		sqlutil.TestDB_Sqlite3,
	}

	for _, testDB := range testDBs {
		sql := `select count(*) as count from migration_log`
		r := struct {
			Count int64
		}{}

		Convey("Initial "+testDB.DriverName+" migration", t, func() {
			x, err := xorm.NewEngine(testDB.DriverName, testDB.ConnStr)
			So(err, ShouldBeNil)

			sqlutil.CleanDB(x)

			has, err := x.SQL(sql).Get(&r)
			So(err, ShouldNotBeNil)

			mg := NewMigrator(x)
			AddMigrations(mg)

			err = mg.Start()
			So(err, ShouldBeNil)

			has, err = x.SQL(sql).Get(&r)
			So(err, ShouldBeNil)
			So(has, ShouldBeTrue)
			expectedMigrations := mg.MigrationsCount() - 2 //we currently skip to migrations. We should rewrite skipped migrations to write in the log as well. until then we have to keep this
			So(r.Count, ShouldEqual, expectedMigrations)

			mg = NewMigrator(x)
			AddMigrations(mg)

			err = mg.Start()
			So(err, ShouldBeNil)

			has, err = x.SQL(sql).Get(&r)
			So(err, ShouldBeNil)
			So(has, ShouldBeTrue)
			So(r.Count, ShouldEqual, expectedMigrations)
		})
	}
}