package Postgres

import (
	fonstruct "Fonbet/json"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

//

func Result(fonbet *fonstruct.FonbetResult, db *pgxpool.Pool, logger *logrus.Logger) (err error) {
	var sum, count = 0, 0

	for i := 0; i < len(fonbet.Events); i++ {
		exist, err := db.Query(context.Background(), `SELECT coalesce((sum(CASE WHEN $1 IN ("stringname") and $2 in ("starttime")THEN 1 ELSE 0 END)),0) FROM results ;`, fonbet.Events[i].Name, fonbet.Events[i].StartTime)
		if err != nil {
			log.Println(err)
		}
		for exist.Next() {
			err := exist.Scan(&sum)
			if err != nil {
				fmt.Println(err)
			}
		}

		if strings.ContainsAny(fonbet.Events[i].Name, "-––") && fonbet.Events[i].Status == 3 {

			strarray := strings.Split(fonbet.Events[i].Score, " ")
			resultarray := strings.Split(strarray[0], ":")

			if len(resultarray) >= 2 && resultarray[0] != " " && resultarray[1] != " " {
				if sum == 0 {
					_, err = db.Exec(context.Background(), "INSERT INTO results (stringname, starttime, score,team1,team2) VALUES ($1, $2, $3, $4, $5)", fonbet.Events[i].Name, fonbet.Events[i].StartTime, fonbet.Events[i].Score, resultarray[0], resultarray[1])
					j := &count
					*j++
					if err != nil {
						fmt.Println(err)
					}
				}

			}
		}

	}
	logger.Infof("Total Result rows: %v", count)
	return
}
