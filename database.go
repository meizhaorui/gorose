package gorose

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/meizhaorui/gorose/utils"
	"strconv"
	"strings"
)

var (
	regex = []string{"=", ">", "<", "!=", "<>", ">=", "<=", "like", "not like", "in", "not in", "between", "not between"}
	//Dbstruct Database
	tx *sql.Tx
)

//type MapData map[string]interface{}
//type MultiData []map[string]interface{}

//var instance *Database
//var once sync.Once
//func GetInstance() *Database {
//	once.Do(func() {
//		instance = &Database{}
//	})
//	return instance
//}

var beforeParseWhereData [][]interface{}

// Database is data mapper struct
type Database struct {
	connection   *Connection
	table        string          // table
	fields       []string          // fields
	where        [][]interface{} // where
	order        string          // order
	limit        int             // limit
	offset       int             // offset
	join         [][]interface{} // join
	distinct     bool            // distinct
	count        string          // count
	sum          string          // sum
	avg          string          // avg
	max          string          // max
	min          string          // min
	group        string          // group
	having       string          // having
	data         interface{}     // data
	LastInsertId int64           // insert last insert id
	trans        bool
	SqlLogs      []string
	LastSql      string
}

// Fields : select fields
func (dba *Database) Fields(fields ...string) *Database {
	dba.fields = fields
	return dba
}

// AddFields : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddFields method:
func (dba *Database) AddFields(fields ...string) *Database {
	dba.fields = append(dba.fields, fields...)
	return dba
}

// Select : equals Fields()
func (dba *Database) Select(fields ...string) *Database {
	return dba.Fields(fields...)
}

// AddSelect : If you already have a query builder instance and you wish to add a column to its existing select clause, you may use the AddSelect method:
func (dba *Database) AddSelect(fields ...string) *Database {
	dba.fields = append(dba.fields, fields...)
	return dba
}

// Table : select table
func (dba *Database) Table(table string) *Database {
	dba.table = table
	return dba
}

// Data : insert or update data
func (dba *Database) Data(data interface{}) *Database {
	dba.data = data
	return dba
}

// Group : select group by
func (dba *Database) Group(group string) *Database {
	dba.group = group
	return dba
}

// GroupBy : equals Group()
func (dba *Database) GroupBy(group string) *Database {
	return dba.Group(group)
}

// Having : select having
func (dba *Database) Having(having string) *Database {
	dba.having = having
	return dba
}

// Order : select order by
func (dba *Database) Order(order string) *Database {
	dba.order = order
	return dba
}

// OrderBy : equal order
func (dba *Database) OrderBy(order string) *Database {
	return dba.Order(order)
}

// Limit : select limit
func (dba *Database) Limit(limit int) *Database {
	dba.limit = limit
	return dba
}

// Offset : select offset
func (dba *Database) Offset(offset int) *Database {
	dba.offset = offset
	return dba
}

// Take : select limit
func (dba *Database) Take(limit int) *Database {
	return dba.Limit(limit)
}

// Skip : select offset
func (dba *Database) Skip(offset int) *Database {
	return dba.Offset(offset)
}

// Page : select page
func (dba *Database) Page(page int) *Database {
	dba.offset = (page - 1) * dba.limit
	return dba
}

// Where : query or execute where condition, the relation is and
func (dba *Database) Where(args ...interface{}) *Database {
	// 如果只传入一个参数, 则可能是字符串、一维对象、二维数组

	// 重新组合为长度为3的数组, 第一项为关系(and/or), 第二项为具体传入的参数 []interface{}
	w := []interface{}{"and", args}

	dba.where = append(dba.where, w)

	return dba
}

// OrWhere : like where , but the relation is or,
func (dba *Database) OrWhere(args ...interface{}) *Database {
	w := []interface{}{"or", args}
	dba.where = append(dba.where, w)

	return dba
}

// WhereNull : like where , where filed is null,
func (dba *Database) WhereNull(arg string) *Database {
	return dba.Where("arg", " is ", "null")
}

// WhereNotNull : like where , where filed is not null,
func (dba *Database) WhereNotNull(arg string) *Database {
	return dba.Where("arg", " is not ", "null")
}

// OrWhereNull : like WhereNull , the relation is or,
func (dba *Database) OrWhereNull(arg string) *Database {
	return dba.OrWhere("arg", " is ", "null")
}

// OrWhereNotNull : like WhereNotNull , the relation is or,
func (dba *Database) OrWhereNotNull(arg string) *Database {
	return dba.OrWhere("arg", " is not ", "null")
}

// WhereIn : a given column's value is contained within the given array
func (dba *Database) WhereIn(field string, arr []interface{}) *Database {
	return dba.Where(field, " in ", arr)
}

// WhereNotIn : the given column's value is not contained in the given array
func (dba *Database) WhereNotIn(field string, arr []interface{}) *Database {
	return dba.Where(field, " not in ", arr)
}

// OrWhereIn : as WhereIn, the relation is or
func (dba *Database) OrWhereIn(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, " in ", arr)
}

// OrWhereNotIn : as WhereNotIn, the relation is or
func (dba *Database) OrWhereNotIn(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, " not in ", arr)
}

// WhereBetween : a column's value is between two values:
func (dba *Database) WhereBetween(field string, arr []interface{}) *Database {
	return dba.Where(field, " between ", arr)
}

// WhereNotBetween : a column's value lies outside of two values:
func (dba *Database) WhereNotBetween(field string, arr []interface{}) *Database {
	return dba.Where(field, " not between ", arr)
}

// OrWhereBetween : like WhereNull , the relation is or,
func (dba *Database) OrWhereBetween(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, " not between ", arr)
}

// OrWhereNotBetween : like WhereNotNull , the relation is or,
func (dba *Database) OrWhereNotBetween(field string, arr []interface{}) *Database {
	return dba.OrWhere(field, " not in ", arr)
}

// Join : select join query
func (dba *Database) Join(args ...interface{}) *Database {
	//dba.parseJoin(args, "INNER")
	dba.join = append(dba.join, []interface{}{"INNER", args})

	return dba
}

// InnerJoin : equals join
func (dba *Database) InnerJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "INNER")
	dba.join = append(dba.join, []interface{}{"INNER", args})

	return dba
}

// LeftJoin : like join , the relation is left
func (dba *Database) LeftJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "LEFT")
	dba.join = append(dba.join, []interface{}{"LEFT", args})

	return dba
}

// RightJoin : like join , the relation is right
func (dba *Database) RightJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.join = append(dba.join, []interface{}{"RIGHT", args})

	return dba
}

// CrossJoin : like join , the relation is cross
func (dba *Database) CrossJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.join = append(dba.join, []interface{}{"CROSS", args})

	return dba
}

// UnionJoin : like join , the relation is union
func (dba *Database) UnionJoin(args ...interface{}) *Database {
	//dba.parseJoin(args, "RIGHT")
	dba.join = append(dba.join, []interface{}{"UNION", args})

	return dba
}

// Distinct : select distinct
func (dba *Database) Distinct() *Database {
	dba.distinct = true

	return dba
}

// First : Retrieving A Single Row / Column From A Table
func (dba *Database) First(args ...interface{}) (map[string]interface{}, error) {
	//var result map[string]interface{}
	//func (dba *Database) First() interface{} {
	dba.limit = 1
	// 构建sql
	sqls, err := dba.BuildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	res, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}

	// 之所以不在 Query 中统一Reset, 是因为chunk会复用到查询相关条件
	//dba.Reset()

	if len(res) == 0 {
		return nil, nil
	}

	return res[0], nil
}

// Get : select more rows , relation limit set
func (dba *Database) Get() ([]map[string]interface{}, error) {
	//func (dba *Database) Get() interface{} {
	// 构建sql
	sqls, err := dba.BuildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	result, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, nil
	}

	//if JsonEncode == true {
	//	jsons, _ := json.Marshal(result)
	//	return string(jsons)
	//}

	//dba.Reset()

	return result, nil
}

// Pluck : Retrieving A List Of Column Values
func (dba *Database) Pluck(args ...string) (interface{}, error) {
	res,err := dba.Get()
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, nil
	}
	switch len(args) {
	case 1:
		var pluckTmp []interface{}
		for _,item:=range res {
			pluckTmp = append(pluckTmp, item[args[0]])
		}
		return pluckTmp,nil
	case 2:
		var pluckTmp = make(map[interface{}]interface{})
		for _,item:=range res {
			pluckTmp[item[args[1]]] = item[args[0]]
		}
		return pluckTmp,nil
	default:
		return nil,errors.New("params error")
	}
}

// Value : If you don't even need an entire row, you may extract a single value from a record using the  value method. This method will return the value of the column directly:
func (dba *Database) Value(arg string) (interface{}, error) {
	res, err := dba.First()
	if err != nil {
		return nil, err
	}
	if val, ok := res[arg]; ok {
		return val, nil
	}
	return nil, errors.New("the field is not exists")
}

// Count : select count rows
func (dba *Database) Count(args ...interface{}) (int, error) {
	fields := "*"
	if len(args) > 0 {
		fields = utils.ParseStr(args[0])
	}
	res, err := dba.buildUnion("count", fields)
	if err != nil {
		return 0, err
	}
	return int(res.(int64)), nil
}

// Sum : select sum field
func (dba *Database) Sum(sum string) (interface{}, error) {
	return dba.buildUnion("sum", sum)
}

// Avg : select avg field
func (dba *Database) Avg(avg string) (interface{}, error) {
	return dba.buildUnion("avg", avg)
}

// Max : select max field
func (dba *Database) Max(max string) (interface{}, error) {
	return dba.buildUnion("max", max)
}

// Min : select min field
func (dba *Database) Min(min string) (interface{}, error) {
	return dba.buildUnion("min", min)
}

// Chunk : select chunk more data to piceses block
func (dba *Database) Chunk(limit int, callback func([]map[string]interface{})) {
	var step = 0
	var offset = dba.offset
	for {
		dba.limit = limit
		dba.offset = offset + step*limit

		// 查询当前区块的数据
		sqls, _ := dba.BuildQuery()
		data, _ := dba.Query(sqls)

		if len(data) == 0 {
			//dba.Reset()
			break
		}

		callback(data)

		//fmt.Println(res)
		if len(data) < limit {
			//dba.Reset()
			break
		}
		step++
	}
}

// Loop : select more data to piceses block from begin
func (dba *Database) Loop(limit int, callback func([]map[string]interface{})) {
	dba.limit = limit
	for {
		// 查询当前区块的数据
		sqls, _ := dba.BuildQuery()
		data, _ := dba.Query(sqls)
		if len(data) == 0 {
			break
		}

		callback(data)

		if len(data) < limit {
			break
		}
	}
}

// BuildSql : build sql string , but not execute sql really
// operType : select/insert/update/delete
func (dba *Database) BuildSql(operType string) (string,error) {
	if operType=="select"{
		return dba.BuildQuery()
	} else {
		return dba.BuildExecut(operType)
	}
}

// BuildQuery : build query string
func (dba *Database) BuildQuery() (string, error) {
	// 聚合
	unionArr := []string{
		dba.count,
		dba.sum,
		dba.avg,
		dba.max,
		dba.min,
	}
	var union string
	for _, item := range unionArr {
		if item != "" {
			union = item
			break
		}
	}
	// distinct
	distinct := utils.If(dba.distinct, "distinct ", "")
	// fields
	fields := utils.If(len(dba.fields) == 0, "*", strings.Join(dba.fields,",")).(string)
	// table
	table := dba.connection.CurrentConfig["prefix"] + dba.table
	// join
	parseJoin, err := dba.parseJoin()
	if err != nil {
		return "", err
	}
	join := parseJoin
	// where
	beforeParseWhereData = dba.where
	parseWhere, err := dba.parseWhere()
	if err != nil {
		return "", err
	}
	where := utils.If(parseWhere == "", "", " WHERE "+parseWhere).(string)
	// group
	group := utils.If(dba.group == "", "", " GROUP BY "+dba.group).(string)
	// having
	having := utils.If(dba.having == "", "", " HAVING "+dba.having).(string)
	// order
	order := utils.If(dba.order == "", "", " ORDER BY "+dba.order).(string)
	// limit
	limit := utils.If(dba.limit == 0, "", " LIMIT "+strconv.Itoa(dba.limit))
	// offset
	offset := utils.If(dba.offset == 0, "", " OFFSET "+strconv.Itoa(dba.offset))

	//sqlstr := "select " + fields + " from " + table + " " + where + " " + order + " " + limit + " " + offset
	sqlstr := fmt.Sprintf("SELECT %s%s FROM %s%s%s%s%s%s%s%s",
		distinct, utils.If(union != "", union, fields), table, join, where, group, having, order, limit, offset)

	//fmt.Println(sqlstr)
	// Reset Database struct

	return sqlstr, nil
}

// BuildExecut : build execute query string
func (dba *Database) BuildExecut(operType string) (string, error) {
	// insert : {"name":"fizz, "website":"fizzday.net"} or {{"name":"fizz2", "website":"www.fizzday.net"}, {"name":"fizz", "website":"fizzday.net"}}}
	// update : {"name":"fizz", "website":"fizzday.net"}
	// delete : ...
	var update, insertkey, insertval, sqlstr string
	if operType != "delete" {
		update, insertkey, insertval = dba.buildData()
	}

	beforeParseWhereData = dba.where
	res, err := dba.parseWhere()
	if err != nil {
		return res, err
	}
	where := utils.If(res == "", "", " WHERE "+res).(string)

	tableName := dba.connection.CurrentConfig["prefix"] + dba.table
	switch operType {
	case "insert":
		sqlstr = fmt.Sprintf("insert into %s (%s) values %s", tableName, insertkey, insertval)
	case "update":
		sqlstr = fmt.Sprintf("update %s set %s%s", tableName, update, where)
	case "delete":
		sqlstr = fmt.Sprintf("delete from %s%s", tableName, where)
	}
	//fmt.Println(sqlstr)
	//dba.Reset()

	return sqlstr, nil
}

// buildData : build inert or update data
func (dba *Database) buildData() (string, string, string) {
	// insert
	var dataFields []string
	var dataValues []string
	// update or delete
	var dataObj []string

	data := dba.data

	switch data.(type) {
	case string:
		dataObj = append(dataObj, data.(string))
	case []map[string]interface{}: // insert multi datas ([]map[string]interface{})
		datas := data.([]map[string]interface{})
		for key, _ := range datas[0] {
			if utils.InArray(key, dataFields) == false {
				dataFields = append(dataFields, key)
			}
		}
		for _, item := range datas {
			var dataValuesSub []string
			for _, key := range dataFields {
				if item[key] == nil {
					dataValuesSub = append(dataValuesSub, "null")
				} else {
					dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(item[key]))
				}
			}
			dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
		}
		//case "map[string]interface {}":
	default: // update or insert
		datas := make(map[string]string)
		switch data.(type) {
		case map[string]interface{}:
			for key, val := range data.(map[string]interface{}) {
				if val == nil {
					datas[key] = "null"
				} else {
					datas[key] = utils.ParseStr(val)
				}
			}
		case map[string]int:
			for key, val := range data.(map[string]int) {
				datas[key] = utils.ParseStr(val)
			}
		case map[string]string:
			for key, val := range data.(map[string]string) {
				datas[key] = val
			}
		}

		//datas := data.(map[string]interface{})
		var dataValuesSub []string
		for key, val := range datas {
			// insert
			dataFields = append(dataFields, key)
			//dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			if val == "null" {
				dataValuesSub = append(dataValuesSub, "null")
			} else {
				dataValuesSub = append(dataValuesSub, utils.AddSingleQuotes(val))
			}
			// update
			//dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
			if val == "null" {
				dataObj = append(dataObj, key+"=null")
			} else {
				dataObj = append(dataObj, key+"="+utils.AddSingleQuotes(val))
			}
		}
		// insert
		dataValues = append(dataValues, "("+strings.Join(dataValuesSub, ",")+")")
	}

	return strings.Join(dataObj, ","), strings.Join(dataFields, ","), strings.Join(dataValues, ",")
}

// buildUnion : build union select
func (dba *Database) buildUnion(union, field string) (interface{}, error) {
	unionStr := union + "(" + field + ") as " + union
	switch union {
	case "count":
		dba.count = unionStr
	case "sum":
		dba.sum = unionStr
	case "avg":
		dba.avg = unionStr
	case "max":
		dba.max = unionStr
	case "min":
		dba.min = unionStr
	}

	// 构建sql
	sqls, err := dba.BuildQuery()
	if err != nil {
		return nil, err
	}

	// 执行查询
	result, err := dba.Query(sqls)
	if err != nil {
		return nil, err
	}

	dba.Reset("union")

	//fmt.Println(union, reflect.TypeOf(union), " ", result[0][union])
	if len(result) > 0 {
		return result[0][union], nil
	}

	var tmp int64 = 0
	return tmp, nil
}

/**
 * 将where条件中的参数转换为where条件字符串
 * example: {"id",">",1}, {"age", 18}
 */
// parseParams : 将where条件中的参数转换为where条件字符串
func (dba *Database) parseParams(args []interface{}) (string, error) {
	paramsLength := len(args)
	argsReal := args

	// 存储当前所有数据的数组
	var paramsToArr []string

	switch paramsLength {
	case 3: // 常规3个参数:  {"id",">",1}
		if !utils.InArray(argsReal[1], regex) {
			return "", errors.New("where parameter is wrong")
		}

		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, argsReal[1].(string))

		switch argsReal[1] {
		case "like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(utils.ParseStr(argsReal[2])))
		case "not like":
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(utils.ParseStr(argsReal[2])))
		case "in":
			paramsToArr = append(paramsToArr, "("+utils.Implode(argsReal[2], ",")+")")
		case "not in":
			paramsToArr = append(paramsToArr, "("+utils.Implode(argsReal[2], ",")+")")
		case "between":
			tmpB := argsReal[2].([]string)
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
		case "not between":
			tmpB := argsReal[2].([]string)
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(tmpB[0])+" and "+utils.AddSingleQuotes(tmpB[1]))
		default:
			paramsToArr = append(paramsToArr, utils.AddSingleQuotes(argsReal[2]))
		}
	case 2:
		//if !utils.TypeCheck(args[0], "string") {
		//	panic("where条件参数有误!")
		//}
		//fmt.Println(argsReal)
		paramsToArr = append(paramsToArr, argsReal[0].(string))
		paramsToArr = append(paramsToArr, "=")
		paramsToArr = append(paramsToArr, utils.AddSingleQuotes(argsReal[1]))
	}
	return strings.Join(paramsToArr, " "), nil
}

// parseJoin : parse the join paragraph
func (dba *Database) parseJoin() (string, error) {
	var join []interface{}
	var returnJoinArr []string
	joinArr := dba.join

	for _, join = range joinArr {
		var w string
		var ok bool
		var args []interface{}

		if len(join) != 2 {
			return "", errors.New("join conditions are wrong")
		}

		// 获取真正的where条件
		if args, ok = join[1].([]interface{}); !ok {
			return "", errors.New("join conditions are wrong")
		}

		argsLength := len(args)
		switch argsLength {
		case 1:
			w = args[0].(string)
		case 2:
			w = args[0].(string) + " ON " + args[1].(string)
		case 4:
			w = args[0].(string) + " ON " + args[1].(string) + " " + args[2].(string) + " " + args[3].(string)
		default:
			return "", errors.New("join format error")
		}

		returnJoinArr = append(returnJoinArr, " "+join[0].(string)+" JOIN "+w)
	}

	return strings.Join(returnJoinArr, " "), nil
}

// parseWhere : parse where condition
func (dba *Database) parseWhere() (string, error) {
	// 取出所有where
	wheres := dba.where
	// where解析后存放每一项的容器
	var where []string

	for _, args := range wheres {
		// and或者or条件
		var condition string = args[0].(string)
		// 统计当前数组中有多少个参数
		params := args[1].([]interface{})
		paramsLength := len(params)

		switch paramsLength {
		case 3: // 常规3个参数:  {"id",">",1}
			res, err := dba.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)

		case 2: // 常规2个参数:  {"id",1}
			res, err := dba.parseParams(params)
			if err != nil {
				return res, err
			}
			where = append(where, condition+" "+res)
		case 1: // 二维数组或字符串
			switch paramReal := params[0].(type) {
			case string:
				where = append(where, condition+" ("+paramReal+")")
			case map[string]interface{}: // 一维数组
				var whereArr []string
				for key, val := range paramReal {
					whereArr = append(whereArr, key+"="+utils.AddSingleQuotes(val))
				}
				where = append(where, condition+" ("+strings.Join(whereArr, " and ")+")")
			case [][]interface{}: // 二维数组
				var whereMore []string
				for _, arr := range paramReal { // {{"a", 1}, {"id", ">", 1}}
					whereMoreLength := len(arr)
					switch whereMoreLength {
					case 3:
						res, err := dba.parseParams(arr)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					case 2:
						res, err := dba.parseParams(arr)
						if err != nil {
							return res, err
						}
						whereMore = append(whereMore, res)
					default:
						return "", errors.New("where data format is wrong")
					}
				}
				where = append(where, condition+" ("+strings.Join(whereMore, " and ")+")")
			case func():
				// 清空where,给嵌套的where让路,复用这个节点
				dba.where = [][]interface{}{}

				// 执行嵌套where放入Database struct
				paramReal()
				// 再解析一遍后来嵌套进去的where
				wherenested, err := dba.parseWhere()
				if err != nil {
					return "", err
				}
				// 嵌套的where放入一个括号内
				where = append(where, condition+" ("+wherenested+")")
			default:
				return "", errors.New("where data format is wrong")
			}
		}
	}

	// 还原初始where, 以便后边调用
	dba.where = beforeParseWhereData

	return strings.TrimLeft(
		strings.TrimLeft(strings.TrimLeft(
			strings.Trim(strings.Join(where, " "), " "),
			"and"), "or"),
		" "), nil
}

// parseExecute : parse execute condition
func (dba *Database) parseExecute(stmt *sql.Stmt, operType string, vals []interface{}) (int64, error) {
	var rowsAffected int64
	var err error
	defer stmt.Close()
	result, errs := stmt.Exec(vals...)
	if errs != nil {
		return 0, errs
	}

	if operType == "insert" {
		// get last insert id
		lastInsertId, err := result.LastInsertId()
		if err == nil {
			dba.LastInsertId = lastInsertId
		}
	}
	// get rows affected
	rowsAffected, err = result.RowsAffected()

	// 如果是事务, 则重置所有参数
	if dba.trans == true {
		dba.Reset("transaction")
	}

	return rowsAffected, err
}

// Insert : insert data and get affected rows
func (dba *Database) Insert() (int, error) {
	sqlstr, err := dba.BuildExecut("insert")
	if err != nil {
		return 0, err
	}

	res, err := dba.Execute(sqlstr)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

// insertGetId : insert data and get id
func (dba *Database) InsertGetId() (int, error) {
	_, err := dba.Insert()
	if err != nil {
		return 0, err
	}
	return int(dba.LastInsertId), nil
}

// Update : update data
func (dba *Database) Update() (int, error) {
	sqlstr, err := dba.BuildExecut("update")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, err
	}
	return int(res), nil
}

// Delete : delete data
func (dba *Database) Delete() (int, error) {
	sqlstr, err := dba.BuildExecut("delete")
	if err != nil {
		return 0, err
	}

	res, errs := dba.Execute(sqlstr)
	if errs != nil {
		return 0, err
	}
	return int(res), nil
}

// Increment : auto Increment +1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
// we can use this method as decrement with the third param as "-"
func (dba *Database) Increment(args ...interface{}) (int, error) {
	argLen := len(args)
	var field string
	var value string = "1"
	var mode string = "+"
	switch argLen {
	case 1:
		field = args[0].(string)
	case 2:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = utils.ParseStr(args[1])
		case int64:
			value = utils.ParseStr(args[1])
		case float32:
			value = utils.ParseStr(args[1])
		case float64:
			value = utils.ParseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
	case 3:
		field = args[0].(string)
		switch args[1].(type) {
		case int:
			value = utils.ParseStr(args[1])
		case int64:
			value = utils.ParseStr(args[1])
		case float32:
			value = utils.ParseStr(args[1])
		case float64:
			value = utils.ParseStr(args[1])
		case string:
			value = args[1].(string)
		default:
			return 0, errors.New("第二个参数类型错误")
		}
		mode = args[2].(string)
	default:
		return 0, errors.New("参数数量只允许1个,2个或3个")
	}
	dba.Data(field + "=" + field + mode + value)
	return dba.Update()
}

// Decrement : auto Decrement -1 default
// we can define step (such as 2, 3, 6 ...) if give the second params
func (dba *Database) Decrement(args ...interface{}) (int, error) {
	arglen := len(args)
	switch arglen {
	case 1:
		args = append(args, 1)
		args = append(args, "-")
	case 2:
		args = append(args, "-")
	default:
		return 0, errors.New("Decrement参数个数有误")
	}
	return dba.Increment(args...)
}

func (dba *Database) Begin() {
	tx, _ = dba.connection.DB.Begin()
	dba.trans = true
}
func (dba *Database) Commit() {
	tx.Commit()
	dba.trans = false
}
func (dba *Database) Rollback() {
	tx.Rollback()
	dba.trans = false
}

// Reset : reset union select
func (dba *Database) Reset(source string) {
	if source == "transaction" {
		//this = new(Database)
		dba.table = ""
		dba.fields = []string{}
		dba.where = [][]interface{}{}
		dba.order = ""
		dba.limit = 0
		dba.offset = 0
		dba.join = [][]interface{}{}
		dba.distinct = false
		dba.group = ""
		dba.having = ""
		var tmp interface{}
		dba.data = tmp
	}
	dba.count = ""
	dba.sum = ""
	dba.avg = ""
	dba.max = ""
	dba.min = ""
}

// JsonEncode : parse json
func (dba *Database) JsonEncode(data interface{}) string {
	res, _ := utils.JsonEncode(data)
	return res
}

// Query : query instance of sql.DB.Query
func (dba *Database) Query(args ...interface{}) ([]map[string]interface{}, error) {
	//defer DB.Close()
	tableData := make([]map[string]interface{}, 0)

	lenArgs := len(args)
	var sqlstring string
	var vals []interface{}

	sqlstring = args[0].(string)

	if lenArgs > 1 {
		for k, v := range args {
			if k > 0 {
				vals = append(vals, v)
			}
		}
	}
	// 记录sql log
	dba.LastSql = fmt.Sprintf(strings.Replace(sqlstring, "%", "%%", -1), vals...)
	dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)

	stmt, err := dba.connection.DB.Prepare(sqlstring)
	if err != nil {
		return tableData, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(vals...)
	if err != nil {
		return tableData, err
	}

	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return tableData, err
	}

	count := len(columns)

	values := make([]interface{}, count)
	scanArgs := make([]interface{}, count)

	for rows.Next() {
		for i := 0; i < count; i++ {
			scanArgs[i] = &values[i]
		}
		rows.Scan(scanArgs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			//fmt.Println(val, reflect.TypeOf(val))
			if b, ok := val.([]byte); ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	return tableData, nil
}

// Execute : query instance of sql.DB.Execute
func (dba *Database) Execute(args ...interface{}) (int64, error) {
	//defer DB.Close()
	lenArgs := len(args)
	var sqlstring string
	var vals []interface{}

	sqlstring = args[0].(string)

	if lenArgs > 1 {
		for k, v := range args {
			if k > 0 {
				vals = append(vals, v)
			}
		}
	}
	// 记录sqlLog
	dba.LastSql = fmt.Sprintf(sqlstring, vals...)
	dba.SqlLogs = append(dba.SqlLogs, dba.LastSql)

	var operType string = strings.ToLower(sqlstring[0:6])
	if operType == "select" {
		return 0, errors.New("this method does not allow select operations, use Query")
	}

	var stmt *sql.Stmt
	var err error
	if dba.trans == true {
		stmt, err = tx.Prepare(sqlstring)
	} else {
		stmt, err = dba.connection.DB.Prepare(sqlstring)
	}

	if err != nil {
		return 0, err
	}
	return dba.parseExecute(stmt, operType, vals)
}

// Transaction : is a simple usage of trans
func (dba *Database) Transaction(closure func() (error)) (bool, error) {
	//defer func() {
	//	if err := recover(); err != nil {
	//		dba.Rollback()
	//		panic(err)
	//	}
	//}()

	dba.Begin()
	err := closure()
	if err != nil {
		dba.Rollback()
		return false, err
	}
	dba.Commit()

	return true, nil
}
