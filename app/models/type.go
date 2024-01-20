package models

import (
	"database/sql"
	"strings"

	"github.com/kevincobain2000/action-coveritup/db"
)

type Type struct {
	ID     int64  `gorm:"column:id;primary_key;AUTO_INCREMENT" json:"id"`
	Name   string `gorm:"column:name;NOT NULL;size:255" json:"name"`
	Metric string `gorm:"column:metric;NOT NULL;size:32" json:"metric"`
}

func (Type) TableName() string {
	return "types"
}

const (
	TYPE_AVERAGE_PR_DAYS        = "average-pr-days"
	TYPE_NUMBER_OF_CONTRIBUTORS = "number-of-contributors"
)

func (t *Type) Get(name string) (Type, error) {
	var ret Type

	query := `SELECT * FROM types WHERE name = @name LIMIT 1`
	err := db.Db().Raw(
		query,
		sql.Named("name", name)).
		Scan(&ret).Error

	ret.Metric = strings.TrimSpace(ret.Metric)
	ret.Name = strings.TrimSpace(ret.Name)

	return ret, err
}

func (t *Type) Create(name string, metric string) (Type, error) {
	var ret Type
	query := `INSERT INTO types (name, metric) VALUES (@name, @metric)`
	err := db.Db().Raw(
		query,
		sql.Named("name", name),
		sql.Named("metric", metric)).
		Scan(&ret).Error

	return ret, err
}

func (t *Type) GetTypesFor(orgName string, repoName string) ([]Type, error) {
	var ret []Type

	query := `SELECT t.* FROM types t
				LEFT JOIN
					coverages c ON t.id = c.type_id
			LEFT JOIN
				repos r ON c.repo_id = r.id
			LEFT JOIN
				orgs o ON c.org_id = o.id
			WHERE
				o.name = @orgName
			AND
				r.name = @repoName
			GROUP BY t.id
			LIMIT @limit`

	err := db.Db().Raw(query,
		sql.Named("orgName", orgName),
		sql.Named("repoName", repoName),
		sql.Named("limit", SAFE_LIMIT_TYPES)).
		Scan(&ret).Error
	if err != nil {
		return ret, err
	}
	for i := range ret {
		ret[i].Name = strings.TrimSpace(ret[i].Name)
		ret[i].Metric = strings.TrimSpace(ret[i].Metric)
	}

	return ret, err
}
