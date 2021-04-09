package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int32
	Coin      int32
}

// InsertUser データベースをレコードを登録する
func InsertUser(record *User) error {
	stmt, err := db.Conn.Prepare("INSERT INTO `user`(`id`, `auth_token`, `name`, `high_score`, `coin`) VALUES (?,?,?,?,?);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return err
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func SelectUserByAuthToken(authToken string) (*User, error) {
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `auth_token`=?", authToken)
	return convertToUser(row)
}

// SelectUserByUserID userIDを条件にレコードを取得する
func SelectUserByUserID(userID string) (*User, error) {
	row := db.Conn.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`=?", userID)
	return convertToUser(row)
}

// SelectUserByUserIDTx トランザクションで排他制御してuserIDを条件にレコードを取得する
func SelectUserByUserIDTx(userID string, tx *sql.Tx) (*User, error) {
	row := tx.QueryRow("SELECT `id`, `auth_token`, `name`, `high_score`, `coin` FROM `user` WHERE `id`=? FOR UPDATE", userID)
	return convertToUser(row)
}

// UpdateUserByUserID userIDを条件にレコードを更新する
func UpdateUserByUserID(record *User) error {
	stmt, err := db.Conn.Prepare("UPDATE `user` SET `auth_token`=?, `name`=?, `high_score`=?, `coin`=? WHERE `id`=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.AuthToken, record.Name, record.HighScore, record.Coin, record.ID)
	return err
}

// UpdateUserByUserIDTx トランザクションでuserID条件にレコードを更新する
func UpdateUserByUserIDTx(record *User, tx *sql.Tx) error {
	stmt, err := tx.Prepare("UPDATE `user` SET `auth_token`=?, `name`=?, `high_score`=?, `coin`=? WHERE `id`=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.AuthToken, record.Name, record.HighScore, record.Coin, record.ID)
	return err
}

// UpdateUserCoinHighScoreByUserID userIDを条件にレコードのコイン、ハイスコアを更新する
func UpdateUserCoinHighScoreByUserID(record *User) error {
	stmt, err := db.Conn.Prepare("UPDATE `user` SET `coin`=?,`high_score`=? WHERE `id`=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.Coin, record.HighScore, record.ID)
	return err
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

// SelectUsersByRanking 指定された位置からconstant.RankingNumber個のデータをuserテーブルから取り出す
func SelectUsersByRanking(startRank int) ([]*User, error) {

	rows, err := db.Conn.Query("SELECT id, auth_token, name, high_score, coin FROM user ORDER by high_score DESC LIMIT ?,?", startRank-1, constant.RankingNumber)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	userSlice := make([]*User, 0, constant.RankingNumber)
	for rows.Next() {
		us := &User{}
		if err := rows.Scan(&us.ID, &us.AuthToken, &us.Name, &us.HighScore, &us.Coin); err != nil {
			log.Println(err)
			return nil, err
		}
		userSlice = append(userSlice, us)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return userSlice, nil

}
