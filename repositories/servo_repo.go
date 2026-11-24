package repositories

import (
	"bab7/entities"
	"database/sql"
)

func InitProj(db *sql.DB) (err error) {
	query := "INSERT INTO servo (id, status_servo) VALUES (1, 0) ON CONFLICT (id) DO NOTHING"
	_, err = db.Exec(query)
	return err
}

func GetStatus(db *sql.DB) (result []entities.Servo, err error) {
	query := "SELECT id, status_servo FROM servo"
	rows, err := db.Query(query)

	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Servo

		err = rows.Scan(&data.Id, &data.StatusServo)
		if err != nil {
			return
		}

		result = append(result, data)
	}

	err = rows.Err()
	return
}

func UpdateStatus(db *sql.DB, status entities.Servo) (err error) {
	query := "UPDATE servo SET status_servo = $1 WHERE id = 1"
	_, err = db.Exec(query, status.StatusServo)
	return
}
