package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
	mockdb "github.com/tpmdigital/simplebank/db/mock"
	db "github.com/tpmdigital/simplebank/db/sqlc"
	"github.com/tpmdigital/simplebank/util"
)

// Example of custom mock matcher
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreateUserParamsMatcher) Matches(x interface{}) bool {
	arg, ok := x.(db.CreateUserParams) // cast to db.CreateUserParams
	if !ok {
		return false
	}
	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}

	e.arg.HashedPassword = arg.HashedPassword

	// In case, some CreateUserParams are nil
	if (e.arg != db.CreateUserParams{} || arg != db.CreateUserParams{}) {
		return reflect.DeepEqual(e.arg, arg)
	}

	return false
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {
	user, password := randomUser(t)

	testCases := []struct {
		name          string
		body          gin.H
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recoder *httptest.ResponseRecorder)
	}{
		{
			name: "OK",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
		{
			name: "InvalidUsername",
			body: gin.H{
				"username":  user.Username + "#",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingUsername",
			body: gin.H{
				"username":  "",
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidPassword",
			body: gin.H{
				"username":  user.Username,
				"password":  "123",
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  "",
				"full_name": user.FullName,
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingFullName",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": "",
				"email":     user.Email,
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "MissingEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "",
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InvalidEmail",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     "tim@",
			},

			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}

				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)). // gomock.Eq(arg)
					Times(0).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
		{
			name: "InternalError",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(db.User{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "DuplicateUsername",
			body: gin.H{
				"username":  user.Username,
				"password":  password,
				"full_name": user.FullName,
				"email":     user.Email,
			},
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.CreateUserParams{
					Username: user.Username,
					// HashedPassword: hashedPassword,
					FullName: user.FullName,
					Email:    user.Email,
				}
				store.EXPECT().
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(db.User{}, &pq.Error{Code: "23505"})
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusForbidden, recorder.Code)
			},
		},
		// {
		// 	name: "DuplicateEmail",
		// 	body: gin.H{
		// 		"username":  user.Username,
		// 		"password":  password,
		// 		"full_name": user.FullName,
		// 		"email":     user.Email,
		// 	},
		// 	buildStubs: func(store *mockdb.MockStore) {
		// 		arg := db.CreateUserParams{
		// 			Username: user.Username,
		// 			// HashedPassword: hashedPassword,
		// 			FullName: user.FullName,
		// 			Email:    user.Email,
		// 		}
		// 		store.EXPECT().
		// 			CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
		// 			Times(1).
		// 			Return(db.User{}, &pq.Error{Code: "23505"})
		// 	},
		// 	checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
		// 		require.Equal(t, http.StatusForbidden, recorder.Code)
		// 	},
		// },
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			data, err := json.Marshal(tc.body)
			require.NoError(t, err)

			url := "/users"
			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}
}

func randomUser(t *testing.T) (db.User, string) {
	password := util.RandomString(6)
	// hashedPassword, err := util.HashedPassword(password)
	// require.NoError(t, err)

	return db.User{
		Username: util.RandomOwner(),
		// HashedPassword: hashedPassword,
		FullName: util.RandomString(12),
		Email:    util.RandomEmail(),
	}, password
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)

	var gotUser db.User
	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user, gotUser)
}