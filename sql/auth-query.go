package sql

const CreateUser = `INSERT INTO users (name, email, password, gender_id, role_id, coin, is_active, create_at, update_at, image)
					VALUES ($1, $2, $3, $4, 1, 2, false, NOW(), NOW(), '-')`

const VerifyCredential = `SELECT id, name, email, password, gender_id, role_id, coin, is_active, create_at, update_at, image FROM users WHERE email = $1`

const GetLastId = `SELECT COUNT(*) FROM users`

const GetByEmail = `SELECT COUNT(*) FROM users WHERE email = $1`

const GetByName = `SELECT COUNT(*) FROM name WHERE name = $1`

const RegisterVal = `SELECT * FROM createValidate($1, $2)`

const UpdateActive = `UPDATE users SET is_active = true WHERE email = $1`
