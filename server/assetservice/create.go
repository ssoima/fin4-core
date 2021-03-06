package assetservice

import (
	"errors"
	"strings"
	"time"

	"fin4-core/server/apperrors"
	"fin4-core/server/datatype"
)

// Create insert asset
func (db *Service) Create(
	sc datatype.ServiceContainer,
	userID datatype.ID,
	name string,
	symbol string,
	description string,
	ethereumAddress string,
	ethereumTransactionAddress string,
) (*datatype.Asset, error) {
	name = strings.TrimSpace(name)
	symbol = strings.TrimSpace(symbol)
	description = strings.TrimSpace(description)
	e := db.ValidateAsset(name, symbol, description)
	if e != nil {
		return nil, e
	}
	{ // check if asset name and symbol dont exist already
		a, err := db.FindByName(name)
		if err != nil {
			apperrors.Critical("assetservice:create:1", err)
			return nil, datatype.ErrServerError
		}
		if a != nil {
			return nil, errors.New("Asset with (" + name + ") name already exists. Please try a different name")
		}
		a, err = db.FindBySymbol(symbol)
		if err != nil {
			apperrors.Critical("assetservice:create:2", err)
			return nil, datatype.ErrServerError
		}
		if a != nil {
			return nil, errors.New("Asset with (" + symbol + ") symbol already exists. Please try a different symbol")
		}
	}
	res, err := db.Exec(
		`INSERT INTO asset SET
	        name = ?,
					symbol = ?,
					description = ?,
					supply = 0,
					creatorId = ?,
					minersCounter = 0,
					favoritesCounter = 0,
					ethereumAddress = ?,
					ethereumTransactionAddress = ?,
					createdAt= ?
	      `,
		name,
		symbol,
		description,
		userID,
		ethereumAddress,
		ethereumTransactionAddress,
		time.Now(),
	)
	if err != nil {
		apperrors.Critical("assetservice:create:4", err)
		return nil, datatype.ErrServerError
	}
	assetID, err := res.LastInsertId()
	if err != nil {
		apperrors.Critical("assetservice:create:5", err)
		return nil, datatype.ErrServerError
	}
	asset, err := db.FindByID(datatype.ID(assetID))
	if err != nil {
		apperrors.Critical("assetservice:create:10", err)
		return nil, datatype.ErrServerError
	}
	return asset, nil
}
