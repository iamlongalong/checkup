package checkup

import (
	"encoding/json"
	"fmt"

	"github.com/iamlongalong/checkup/storage/appinsights"
	"github.com/iamlongalong/checkup/storage/fs"
	"github.com/iamlongalong/checkup/storage/github"
	"github.com/iamlongalong/checkup/storage/mysql"
	"github.com/iamlongalong/checkup/storage/postgres"
	"github.com/iamlongalong/checkup/storage/s3"
	"github.com/iamlongalong/checkup/storage/sql"
)

func storageDecode(typeName string, config json.RawMessage) (Storage, error) {
	switch typeName {

	case mysql.Type:
		return mysql.New(config)
	case postgres.Type:
		return postgres.New(config)
	case s3.Type:
		return s3.New(config)
	case github.Type:
		return github.New(config)
	case fs.Type:
		return fs.New(config)
	case sql.Type:
		return sql.New(config)
	case appinsights.Type:
		return appinsights.New(config)
	default:
		return nil, fmt.Errorf(errUnknownStorageType, typeName)
	}
}
