package sqlxt

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"text/template"

	"github.com/go-yaml/yaml"
)

var (
	//ErrKey error key
	ErrKey = errors.New("key error")
	//ErrGroupNotFound group not found error
	ErrGroupNotFound = errors.New("group not found")
	//ErrKeyNotFound key not found error
	ErrKeyNotFound = errors.New("key not found")
)

//Parse parse yaml sql template
type Parse struct {
	path      string
	templates map[string]map[string]*template.Template
	locker    *sync.RWMutex
}

//NewParse new parse
func NewParse(path string) (*Parse, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)
	tmp := make(map[string]map[string]string)
	tmps := make(map[string]map[string]*template.Template)
	err = yaml.Unmarshal(buf.Bytes(), tmp)
	if err != nil {
		return nil, err
	}
	for k, ts := range tmp {
		tm := make(map[string]*template.Template)
		for key, t := range ts {
			temp := template.Must(template.New(k + "-" + "key").Parse(t))
			tm[key] = temp
		}
		tmps[k] = tm
	}
	return &Parse{
		path:      path,
		templates: tmps,
		locker:    new(sync.RWMutex),
	}, nil
}

//Get get a sql
// key的格式为 "group-add" group 为sql的分组，add表示是add sql语句
func (p *Parse) Get(key string) (*template.Template, error) {
	keys := strings.Split(key, "-")
	if len(keys) != 2 {
		return nil, ErrKey
	}
	p.locker.RLock()
	defer p.locker.RUnlock()
	g, ok := p.templates[keys[0]]
	if !ok {
		return nil, ErrGroupNotFound
	}
	k, ok := g[keys[1]]
	if !ok {
		return nil, ErrKeyNotFound
	}
	return k, nil
}
