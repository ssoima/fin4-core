package assetservice

import (
	"fin4-core/server/datatype"
	"fin4-core/server/decimaldt"
	"github.com/lytics/logrus"
)

// DepositBalance cache user balance in the SQL db
func (db *Service) DepositBalance(
	userID datatype.ID,
	assetID datatype.ID,
	amount decimaldt.Decimal,
) error {
	stmt, err := db.Prepare(
		`INSERT INTO asset_user_balance
      SET
        userId = ?, assetId = ?, balance = ?, reserved = 0
      ON DUPLICATE KEY UPDATE
        balance = balance + ?`,
	)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"userID":  userID,
			"assetID": assetID,
			"amount":  amount.String(),
		}).Error("user_balance:DepositBalance:e1")
		return err
	}
	_, err = stmt.Exec(userID, assetID, amount, amount)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err.Error(),
			"userID":  userID,
			"assetID": assetID,
			"amount":  amount.String(),
		}).Error("user_balance:DepositBalance:e2")
		return err
	}
	return nil
}
