package database

const SaveUserQuery = `
	INSERT INTO users (name, code) VALUES ($1, $2)
`

const ListUsersQuery = `SELECT u.name, u.code FROM users u`

const FindUserQuery = `SELECT EXISTS(SELECT u.name, u.code FROM users u WHERE u.code = $1)`

const FindUserByCodeQuery = `SELECT u.name, u.code FROM users u WHERE u.code = $1`
