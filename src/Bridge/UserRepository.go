package bridge

import (
	"database/sql"
	"fmt"
	"reflect"
)

type UserRepository struct {
	db     *sql.DB
	hasher Hasher
}

func NewUserRepository(db *sql.DB, hasher Hasher) *UserRepository {
	return &UserRepository{
		db:     db,
		hasher: hasher,
	}
}

func (r *UserRepository) getUserEntityByUserCredentials(username, password, grantType string, client *Client) (*User, error) {
	provider := client.Provider
	if provider == "" {
		provider = config.Get("auth.guards.api.provider")
	}

	model := config.Get("auth.providers." + provider + ".model")
	if model == nil {
		return nil, fmt.Errorf("Unable to determine authentication model from configuration.")
	}

	user := &model{}

	// Check for the 'findAndValidateForPassport' method
	findAndValidateMethod, exists := reflect.TypeOf(model).MethodByName("findAndValidateForPassport")
	if exists {
		// Use the 'findAndValidateForPassport' method to find and validate the user
		args := []reflect.Value{reflect.ValueOf(username), reflect.ValueOf(password)}
		methodValue := reflect.ValueOf(findAndValidateMethod)
		returnValues := methodValue.Call(args)

		if returnValues[0].IsNil() {
			return nil, nil
		}

		user = returnValues[0].Interface().(*model)
	} else {
		// Check for the 'findForPassport' method
		findForPassportMethod, exists := reflect.TypeOf(model).MethodByName("findForPassport")
		if exists {
			// Use the 'findForPassport' method to find the user
			args := []reflect.Value{reflect.ValueOf(username)}
			methodValue := reflect.ValueOf(findForPassportMethod)
			returnValues := methodValue.Call(args)

			if returnValues[0].IsNil() {
				return nil, nil
			}

			user = returnValues[0].Interface().(*model)
		} else {
			// Use the 'Where' method to find the user based on their email address
			query := `SELECT * FROM ` + model.Table + ` WHERE email = ?`
			rows, err := r.db.Query(query, username)
			if err != nil {
				return nil, err
			}
			defer rows.Close()

			for rows.Next() {
				err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.AuthPassword)
				if err != nil {
					return nil, err
				}
			}

			if err := rows.Err(); err != nil {
				return nil, err
			}
		}
	}

	// Validate the user's credentials
	if method_exists(user, "validateForPassportPasswordGrant") {
		if err := user.ValidateForPassportPasswordGrant(password); err != nil {
			return nil, err
		}
	} else if !r.hasher.Check(password, user.AuthPassword) {
		return nil, fmt.Errorf("Invalid password")
	}

	return user, nil
}
