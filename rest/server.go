package rest

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/go-martini/martini"
	"github.com/mgbaozi/asrest/rest/exception"
)

type Server struct {
	*martini.ClassicMartini
	DBUrl  string
	DBName string
}

func NewServer() *Server {
	server := &Server{ClassicMartini: martini.Classic(), DBUrl: "localhost", DBName: "test"}
	server.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
	})
	return server
}

func get_table_name(model interface{}) string {
	model_value := reflect.ValueOf(model)
	method := model_value.MethodByName("TableName")
	var table_name string
	if method.IsValid() {
		values := method.Call([]reflect.Value{})
		table_name = values[0].Interface().(string)
	} else {
		type_name := strings.ToLower(model_value.Type().String())
		splited := strings.Split(type_name, ".")
		table_name = splited[len(splited)-1]
	}
	return table_name
}

func (self *Server) Run() {
	Connect(self.DBUrl)
	self.ClassicMartini.Run()
}

func (self *Server) RunOnAddr(addr string) {
	Connect(self.DBUrl)
	self.ClassicMartini.RunOnAddr(addr)
}

func (self *Server) All(url string, model interface{}) {
	self.Get(url, model)
}

func (self *Server) Get(url string, model interface{}) {
	tags := getTags("param", model)
	model_type := reflect.ValueOf(model).Type()
	table_name := get_table_name(model)
	self.ClassicMartini.Get(url+"/:id", func(params martini.Params) (code int, response string) {
		code = 200
		var result interface{}
		defer func() {
			json_data, _ := json.Marshal(result)
			response = string(json_data)
		}()
		collection := DB(self.DBName).C(table_name)
		data := reflect.New(model_type)
		id, err := toObjectId(params["id"])
		if err != nil {
			code = 404
			result = exception.NotFound()
			return
		}
		query := collection.FindId(id)
		count, err := query.Count()
		if err != nil || count < 1 {
			code = 404
			result = exception.NotFound()
			return
		}
		result = data.Interface()
		query.One(result)
		return
	})
	self.ClassicMartini.Get(url, func(req *http.Request) string {
		collection := DB(self.DBName).C(table_name)
		params := req.URL.Query()
		query := bson.M{}
		for key, value := range params { //params is a map[string][]string
			if index, ok := tags[key]; ok {
				bson_key := getTag(tags[key], "bson", model)
				if len(bson_key) > 0 {
					query[bson_key] = ConverseType(model, index, value[0])
				}
			}
		}
		slice := reflect.MakeSlice(reflect.SliceOf(model_type), 0, 0)
		result := reflect.New(slice.Type())
		result.Elem().Set(slice)
		collection.Find(&query).All(result.Interface())
		data, _ := json.Marshal(result.Interface())
		return string(data)
	})
}
