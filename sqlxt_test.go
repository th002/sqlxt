package sqlxt

import (
	"reflect"
	"testing"

	"github.com/jmoiron/sqlx"

	"fmt"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

var db *sqlx.DB

func init() {
	var err error
	d, mock, err := sqlmock.New()
	if err != nil {
		fmt.Printf("an error '%s' was not expected when opening a stub database connection", err)
		return
	}
	rows := sqlmock.NewRows([]string{"user_id", "product_id"}).
		AddRow(2, 3)
	mock.ExpectBegin()
	mock.ExpectQuery("select (.+) from product_viewers").WillReturnRows(rows)
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	db = sqlx.NewDb(d, "root@/blog")
}

func TestNew(t *testing.T) {
	parser, err := NewParse("./test_data/test_yaml.yaml")
	if err != nil {
		t.Error(err)
	}
	type args struct {
		db   *sqlx.DB
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *SqlxT
		wantErr bool
	}{
		{
			name: "测试正常的sql语句",
			args: args{
				db:   db,
				path: "./test_data/test_yaml.yaml",
			},
			want: &SqlxT{
				db:     db,
				parser: parser,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.db, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSqlxT_Exec(t *testing.T) {
	parser, err := NewParse("./test_data/test_yaml.yaml")
	if err != nil {
		t.Error(err)
	}
	type fields struct {
		db     *sqlx.DB
		parser *Parse
	}
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "测试正常的sql",
			fields: fields{
				db:     db,
				parser: parser,
			},
			args: args{
				key: "user-insert",
				v:   map[string]int{"UserID": 2, "ProductID": 3},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &SqlxT{
				db:     tt.fields.db,
				parser: tt.fields.parser,
			}
			if err := st.Exec(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SqlxT.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSqlxT_Get(t *testing.T) {
	parser, err := NewParse("./test_data/test_yaml.yaml")
	if err != nil {
		t.Error(err)
	}
	type fields struct {
		db     *sqlx.DB
		parser *Parse
	}
	type args struct {
		key string
		v   interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "测试正常get",
			fields: fields{
				db:     db,
				parser: parser,
			},
			args: args{
				key: "user-get",
				v:   make(map[string]int),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			st := &SqlxT{
				db:     tt.fields.db,
				parser: tt.fields.parser,
			}
			if err := st.Get(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("SqlxT.Get() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestSqlxT_Select(t *testing.T) {
// 	type fields struct {
// 		db     *sqlx.DB
// 		parser *Parse
// 	}
// 	type args struct {
// 		key string
// 		v   interface{}
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 	// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			st := &SqlxT{
// 				db:     tt.fields.db,
// 				parser: tt.fields.parser,
// 			}
// 			if err := st.Select(tt.args.key, tt.args.v); (err != nil) != tt.wantErr {
// 				t.Errorf("SqlxT.Select() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
