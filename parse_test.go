package sqlxt

import (
	"sync"
	"testing"
	"text/template"
)

func TestNewParse(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    *Parse
		wantErr bool
	}{
		{
			"错误的path",
			args{
				path: "/12333",
			},
			nil,
			true,
		},
		{
			"错误的YAML格式",
			args{
				path: "./test_data/error_yaml.yaml",
			},
			nil,
			true,
		},
		{
			"正确的New parse",
			args{
				path: "./test_data/test_yaml.yaml",
			},
			&Parse{
				path: "./test_data/test_yaml.yaml",
				templates: map[string]map[string]*template.Template{
					"user": map[string]*template.Template{
						"insert": template.Must(template.New("user-insert").Parse("insert into users ({{.Name}},{{.Mobile}})")),
					},
				},
				locker: new(sync.RWMutex),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewParse(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewParse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("NewParse() = %v, want %v", got, tt.want)
			// }
		})
	}
}

func TestParse_Get(t *testing.T) {
	type fields struct {
		path      string
		templates map[string]map[string]*template.Template
		locker    *sync.RWMutex
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *template.Template
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "测试错误的key",
			fields: fields{
				path: "./test_data/1.yaml",
				templates: map[string]map[string]*template.Template{
					"users": map[string]*template.Template{
						"add": nil,
					},
				},
				locker: new(sync.RWMutex),
			},
			args: args{
				key: "User",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "测试不存在的组",
			fields: fields{
				path: "./test_data/1.yaml",
				templates: map[string]map[string]*template.Template{
					"users": map[string]*template.Template{
						"add": template.Must(template.New("users-add").Parse("insert into users({{.Name}},{{.Mobile}})")),
					},
				},
				locker: new(sync.RWMutex),
			},
			args: args{
				key: "User-add",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "测试不存在的Key",
			fields: fields{
				path: "./test_data/1.yaml",
				templates: map[string]map[string]*template.Template{
					"users": map[string]*template.Template{
						"add": template.Must(template.New("users-add").Parse("insert into users({{.Name}},{{.Mobile}})")),
					},
				},
				locker: new(sync.RWMutex),
			},
			args: args{
				key: "users-update",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "测试正确的key",
			fields: fields{
				path: "./test_data/1.yaml",
				templates: map[string]map[string]*template.Template{
					"users": map[string]*template.Template{
						"add": template.Must(template.New("users-add").Parse("insert into users({{.Name}},{{.Mobile}})")),
					},
				},
				locker: new(sync.RWMutex),
			},
			args: args{
				key: "users-add",
			},
			want:    template.Must(template.New("users-add").Parse("insert into users({{.Name}},{{.Mobile}})")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Parse{
				path:      tt.fields.path,
				templates: tt.fields.templates,
				locker:    tt.fields.locker,
			}
			_, err := p.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if got != tt.want {
			// 	t.Errorf("Parse.Get() = %v, want %v", got, tt.want)
			// }
		})
	}
}
