package pgsql_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Soapstone-Services/go-template-2024/pkg/api/user/platform/pgsql"
	"github.com/Soapstone-Services/go-template-2024/pkg/utl/mock"
)

func TestCreate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		req      stems.User
		wantData stems.User
	}{
		{
			name:    "Fail on insert duplicate ID",
			wantErr: true,
			req: stems.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: stems.Base{
					ID: 1,
				},
			},
		},
		{
			name: "Success",
			req: stems.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: stems.Base{
					ID: 2,
				},
			},
			wantData: stems.User{
				Email:      "newtomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "newtomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "pass",
				Base: stems.Base{
					ID: 2,
				},
			},
		},
		{
			name:    "User already exists",
			wantErr: true,
			req: stems.User{
				Email:    "newtomjones@mail.com",
				Username: "newtomjones",
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &stems.Role{}, &stems.User{})

	err := mock.InsertMultiple(db,
		&stems.Role{
			ID:          1,
			AccessLevel: 1,
			Name:        "SUPER_ADMIN",
		},
		&stems.User{
			Email:    "nottomjones@mail.com",
			Username: "nottomjones",
			Base: stems.Base{
				ID: 1,
			},
		})
	if err != nil {
		t.Error(err)
	}

	udb := pgsql.User{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := udb.Create(db, tt.req)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				if resp.ID == 0 {
					t.Error("expected data, but got empty struct.")
					return
				}
				tt.wantData.CreatedAt = resp.CreatedAt
				tt.wantData.UpdatedAt = resp.UpdatedAt
				assert.Equal(t, tt.wantData, resp)
			}
		})
	}
}

func TestView(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		id       int
		wantData stems.User
	}{
		{
			name:    "User does not exist",
			wantErr: true,
			id:      1000,
		},
		{
			name: "Success",
			id:   2,
			wantData: stems.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: stems.Base{
					ID: 2,
				},
				Role: &stems.Role{
					ID:          1,
					AccessLevel: 1,
					Name:        "SUPER_ADMIN",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &stems.Role{}, &stems.User{})

	if err := mock.InsertMultiple(db, &stems.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.User{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			user, err := udb.View(db, tt.id)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				if user.ID == 0 {
					t.Errorf("response was empty due to: %v", err)
				} else {
					tt.wantData.CreatedAt = user.CreatedAt
					tt.wantData.UpdatedAt = user.UpdatedAt
					assert.Equal(t, tt.wantData, user)
				}
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      stems.User
		wantData stems.User
	}{
		{
			name: "Success",
			usr: stems.User{
				Base: stems.Base{
					ID: 2,
				},
				FirstName: "Z",
				LastName:  "Freak",
				Address:   "Address",
				Phone:     "123456",
				Mobile:    "345678",
				Username:  "newUsername",
			},
			wantData: stems.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Z",
				LastName:   "Freak",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Address:    "Address",
				Phone:      "123456",
				Mobile:     "345678",
				Base: stems.Base{
					ID: 2,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &stems.Role{}, &stems.User{})

	if err := mock.InsertMultiple(db, &stems.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[0].usr); err != nil {
		t.Error(err)
	}

	udb := pgsql.User{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := udb.Update(db, tt.wantData)
			if tt.wantErr != (err != nil) {
				fmt.Println(tt.wantErr, err)
			}
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData.ID != 0 {
				user := stems.User{
					Base: stems.Base{
						ID: tt.usr.ID,
					},
				}
				if err := db.Select(&user); err != nil {
					t.Error(err)
				}
				tt.wantData.UpdatedAt = user.UpdatedAt
				tt.wantData.CreatedAt = user.CreatedAt
				tt.wantData.LastLogin = user.LastLogin
				tt.wantData.DeletedAt = user.DeletedAt
				assert.Equal(t, tt.wantData, user)
			}
		})
	}
}

func TestList(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		qp       *stems.ListQuery
		pg       stems.Pagination
		wantData []stems.User
	}{
		{
			name:    "Invalid pagination values",
			wantErr: true,
			pg: stems.Pagination{
				Limit: -100,
			},
		},
		{
			name: "Success",
			pg: stems.Pagination{
				Limit:  100,
				Offset: 0,
			},
			qp: &stems.ListQuery{
				ID:    1,
				Query: "company_id = ?",
			},
			wantData: []stems.User{
				{
					Email:      "tomjones@mail.com",
					FirstName:  "Tom",
					LastName:   "Jones",
					Username:   "tomjones",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "newPass",
					Base: stems.Base{
						ID: 2,
					},
					Role: &stems.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
				},
				{
					Email:      "johndoe@mail.com",
					FirstName:  "John",
					LastName:   "Doe",
					Username:   "johndoe",
					RoleID:     1,
					CompanyID:  1,
					LocationID: 1,
					Password:   "hunter2",
					Base: stems.Base{
						ID: 1,
					},
					Role: &stems.Role{
						ID:          1,
						AccessLevel: 1,
						Name:        "SUPER_ADMIN",
					},
					Token: "loginrefresh",
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &stems.Role{}, &stems.User{})

	if err := mock.InsertMultiple(db, &stems.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[1].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.User{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			users, err := udb.List(db, tt.qp, tt.pg)
			assert.Equal(t, tt.wantErr, err != nil)
			if tt.wantData != nil {
				for i, v := range users {
					tt.wantData[i].CreatedAt = v.CreatedAt
					tt.wantData[i].UpdatedAt = v.UpdatedAt
				}
				assert.Equal(t, tt.wantData, users)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	cases := []struct {
		name     string
		wantErr  bool
		usr      stems.User
		wantData stems.User
	}{
		{
			name: "Success",
			usr: stems.User{
				Base: stems.Base{
					ID:        2,
					DeletedAt: mock.TestTime(2018),
				},
			},
			wantData: stems.User{
				Email:      "tomjones@mail.com",
				FirstName:  "Tom",
				LastName:   "Jones",
				Username:   "tomjones",
				RoleID:     1,
				CompanyID:  1,
				LocationID: 1,
				Password:   "newPass",
				Base: stems.Base{
					ID: 2,
				},
			},
		},
	}

	dbCon := mock.NewPGContainer(t)
	defer dbCon.Shutdown()

	db := mock.NewDB(t, dbCon, &stems.Role{}, &stems.User{})

	if err := mock.InsertMultiple(db, &stems.Role{
		ID:          1,
		AccessLevel: 1,
		Name:        "SUPER_ADMIN"}, &cases[0].wantData); err != nil {
		t.Error(err)
	}

	udb := pgsql.User{}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {

			err := udb.Delete(db, tt.usr)
			assert.Equal(t, tt.wantErr, err != nil)

			// Check if the deleted_at was set
		})
	}
}
