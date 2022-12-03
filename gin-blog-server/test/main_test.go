package test

import (
	"fmt"
	"testing"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"github.com/mitchellh/mapstructure"
)

type Common struct {
	ID int `mapstructure:"id"`
}

type User struct {
	Common        `mapstructure:",squash"`
	Name          string `mapstructure:"user_name"`
	Age           int32  `mapstructure:"user_age"`
	RequestMethod string `json:"request_method" mapstructure:"request_method"`
}

func TestMapStruct(t *testing.T) {
	res := map[string]any{}
	user := User{
		Common:        Common{ID: 1},
		Name:          "Zhangsan",
		Age:           11,
		RequestMethod: "GET",
	}
	mapstructure.Decode(user, &res)
	fmt.Println(res)
}

func TestCopier(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18}
	m := map[string]any{}

	copier.Copy(&m, user)

	fmt.Printf("%#v \n", user)
	fmt.Printf("%#v \n", m)
}

type Foo struct {
	A string
	B int64
}

func TestMergo(t *testing.T) {
	user := User{Name: "Jinzhu", Age: 18, RequestMethod: "GET"}
	m := map[string]any{}

	mergo.Map(&m, user)
	fmt.Println(m)
}
