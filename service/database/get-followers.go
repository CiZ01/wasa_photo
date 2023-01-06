package database

var query_GETFOLLOWERS = `SELECT userID, username, userPropicURL FROM User WHERE userID IN (SELECT followerID FROM Follow WHERE followedID=? LIMIT ?, ?)`

func (db *appdbimpl) GetFollowers(followedID uint32, offset uint32, limit int32) ([]User, error) {
	// Get the followers from the database
	rows, err := db.c.Query(query_GETFOLLOWERS, followedID, offset, limit)
	if err != nil {
		return nil, err
	}
	defer func() { err = rows.Close() }()

	var followers []User

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		var follower User

		// Get follower data
		err = rows.Scan(&follower.UserID, &follower.Username, &follower.UserPropicURL)
		if err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, err
}
