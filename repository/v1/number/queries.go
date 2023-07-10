package number

import(
	//"fmt"
	"context"
	"github.com/negeek/short-access/utils"
	"github.com/negeek/short-access/db"
	"github.com/jackc/pgx/v4"
)

func (n *Number) CreateOrUpdate() error {
	utils.Time(n,false)
	query:="INSERT INTO numbers (id, number, date_updated) VALUES ($1, $2, $3) ON CONFLICT (id) DO UPDATE SET number = numbers.number + $2 RETURNING number"
	err := db.PostgreSQLDB.QueryRow(context.Background(), query, n.Id, n.Step, n.DateCreated).Scan(&n.Number)
	if err != nil {
		return err
	}
	return nil
}

func (n *Number) FindById() (error, bool){
	query:="SELECT number FROM numbers WHERE id = $1"
	err := db.PostgreSQLDB.QueryRow(context.Background(),query,n.Id).Scan(&n.Number)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, false
		}
		return err,false
	}
	return nil, true

}