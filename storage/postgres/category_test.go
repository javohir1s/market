package postgres

import (
	"context"
	"fmt"
	"market_system/config"
	"market_system/models"
	"reflect"
	"testing"

	faker "github.com/bxcodec/faker/v4"
)

func Test_categoryRepo_Create(t *testing.T) {

	cfg := config.Load()

	stroge, err := NewConnectionPostgres(&cfg)
	if err != nil {
		panic("no connect to postgres:" + err.Error())
	}

	type args struct {
		ctx context.Context
		req *models.CreateCategory
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Category
		wantErr bool
	}{
		// TODO: Add test cases.
	}

	for i := 1; i <= 10000; i++ {
		var title = faker.FirstName()

		tests = append(tests, struct {
			name    string
			args    args
			want    *models.Category
			wantErr bool
		}{
			name: fmt.Sprintf("case %d", i),
			args: args{
				context.Background(),
				&models.CreateCategory{
					Title: title,
				},
			},
			want: &models.Category{
				Title: title,
			},
			wantErr: false,
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stroge.Category().Create(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("categoryRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Title != tt.want.Title {
				t.Errorf("categoryRepo.Create() = %v, want %v", got.Title, tt.want.Title)
			}
		})
	}
}

func Test_categoryRepo_GetByID(t *testing.T) {
	cfg := config.Load()

	_, err := NewConnectionPostgres(&cfg)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	type args struct {
		ctx context.Context
		req *models.CategoryPrimaryKey
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Category
		wantErr bool
	}{
		{
			name: "Valid ID",
			args: args{
				ctx: context.Background(),
				req: &models.CategoryPrimaryKey{Id: "ca57750e-bce0-41c3-a6fc-db245a719024"},
			},
			want: &models.Category{
				Id:        "ca57750e-bce0-41c3-a6fc-db245a719024",
				Title:     "Jarvis",
				CreatedAt: "2023-12-13T22:13:29.989701Z",
				UpdatedAt: "2023-12-13T22:13:29.989701Z",
			},
			wantErr: false,
		},
		{
			name: "Invalid ID",
			args: args{
				ctx: context.Background(),
				req: &models.CategoryPrimaryKey{Id: "invalid-id"},
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stroge, err := NewConnectionPostgres(&cfg)
			if err != nil {
				panic("no connect to postgres:" + err.Error())
			}

			got, err := stroge.Category().GetByID(test.args.ctx, test.args.req)
			if (err != nil) != test.wantErr {
				t.Errorf("categoryRepo.GetByID() error = %v, wantErr %v", err, test.wantErr)
				return
			}

			if !test.wantErr && got != nil {
				if !reflect.DeepEqual(got.Id, test.want.Id) ||
					!reflect.DeepEqual(got.Title, test.want.Title) ||
					!reflect.DeepEqual(got.CreatedAt, test.want.CreatedAt) ||
					!reflect.DeepEqual(got.UpdatedAt, test.want.UpdatedAt) {
					t.Errorf("%s: expected %+v, but got %+v", test.name, test.want, got)
				}
			} else if !test.wantErr {
				t.Errorf("%s: unexpected nil response", test.name)
			}
		})
	}
}

func Test_categoryRepo_GetList(t *testing.T) {
	cfg := config.Load()

	_, err := NewConnectionPostgres(&cfg)
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL: %v", err)
	}

	// type args struct {
	// 	ctx   context.Context
	// 	req   *models.GetListCategoryRequest
	// }

	tests := []struct {
		Name     string
		Input    *models.GetListCategoryRequest
		Expected int
		WantErr  bool
	}{
		{
			Name: "Test with valid parameters",
			Input: &models.GetListCategoryRequest{
				Limit: 12,
			},
			Expected: 12,
			WantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			stroge, err := NewConnectionPostgres(&cfg)
			if err != nil {
				panic("no connect to postgres:" + err.Error())
			}

			got, err := stroge.Category().GetList(context.Background(), test.Input)
			if (err != nil) != test.WantErr {
				t.Errorf("categoryRepo.GetList() error = %v, wantErr %v", err, test.WantErr)
				return
			}

			if !test.WantErr && got != nil {
				if len(got.Categories) != test.Expected {
					t.Errorf("%s: expected %d categories, but got %d", test.Name, test.Expected, len(got.Categories))
				}
			} else if !test.WantErr {
				t.Errorf("%s: unexpected nil response", test.Name)
			}
		})
	}
}

func TestDeleteCategory(t *testing.T) {
	tests := []struct {
		Name    string
		Input   string
		WantErr bool
	}{
		{
			Name:    "Test with valid ID",
			Input:   "7aebba3d-90d6-4037-ad02-2dc86050994c",
			WantErr: false, 
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := config.Load()

			stroge, err := NewConnectionPostgres(&cfg)
			if err != nil {
				panic("no connect to postgres:" + err.Error())
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			err = stroge.Category().Delete(ctx, &models.CategoryPrimaryKey{Id: test.Input})
			if test.WantErr && err == nil {
				t.Errorf("%s: expected an error, but got nil", test.Name)
				return
			}

			if !test.WantErr && err != nil {
				t.Errorf("%s: unexpected error: %v", test.Name, err)
				return
			}

			if !test.WantErr {
				_, err := stroge.Category().GetByID(ctx, &models.CategoryPrimaryKey{Id: test.Input})
				if err == nil {
					t.Errorf("%s: expected category to be deleted, but it still exists", test.Name)
					return
				}
			}
		})
	}
}


func TestUpdateCategory(t *testing.T) {
	tests := []struct {
		Name    string
		Input   *models.UpdateCategory
		WantErr bool
	}{
		{
			Name: "Updated",
			Input: &models.UpdateCategory{
				Id:    "82932a1e-b3fe-4d07-8e40-d118a6dd22e4", 
				Title: "UpdatedTitle",
			},
			WantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			cfg := config.Load()

			stroge, err := NewConnectionPostgres(&cfg)
			if err != nil {
				panic("no connect to postgres:" + err.Error())
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			_, err = stroge.Category().Update(ctx, test.Input)
			if test.WantErr {
				if err == nil {
					t.Errorf("%s: expected an error, but got nil", test.Name)
				} else {
					fmt.Printf("%s: got expected error: %v\n", test.Name, err)
				}
				return
			}

			if err != nil {
				t.Errorf("%s: unexpected error during update: %v", test.Name, err)
				return
			}

			updatedCategory, err := stroge.Category().GetByID(ctx, &models.CategoryPrimaryKey{Id: test.Input.Id})
			if err != nil {
				t.Errorf("%s: error getting updated category: %v", test.Name, err)
				return
			}

			if updatedCategory.Title != test.Input.Title {
				t.Errorf("%s: expected title %s, but got %s", test.Name, test.Input.Title, updatedCategory.Title)
			}
		})
	}
}
