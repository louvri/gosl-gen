
package test_table
import (
	"fmt"
	"strings"
)

type Model struct {
    Id int64 `db:"id"`
    FieldValue string `db:"field_value"`
    Defaults map[string]interface{}
}

func (m *Model) SetDefaults(defaults map[string]interface{}){
    m.Defaults = defaults
}

func (m *Model) Get(field string) interface{}{
    switch field {
                case "id","Id": return m.Id
                case "field_value","FieldValue": return m.FieldValue
    }
    return nil
}

func (m *Model) Patch(data Model) {
            if data.Id != 0 {
                m.Id = data.Id
            }     
            if data.FieldValue != "" {
                m.FieldValue = data.FieldValue
            }
}

func (m *Model) ToString(builder *strings.Builder) {
            if builder.Len() > 0 {
                builder.WriteString(",")
            }
                builder.WriteString(fmt.Sprintf("%d",m.Id))
            if builder.Len() > 0 {
                builder.WriteString(",")
            }
                builder.WriteString("\"")    
                builder.WriteString(m.FieldValue)
                builder.WriteString("\"")
}

func (m *Model) ToMap(builder map[string]interface{},dontIgnoreEmtpy ...bool) map[string]interface{} {
    if builder == nil {
        builder = make(map[string]interface{})
    }
    var dontIgnore bool
    if len(dontIgnoreEmtpy) > 0 {
         dontIgnore = dontIgnoreEmtpy[0]
    }
        if dontIgnore {
                value := int64(0)
                if  m.Defaults["id"] != nil {
                    value = m.Defaults["id"].(int64)
                }
                if m.Id != value {
                    builder["id"] = m.Id
                } else {
                    builder["id"] = value
                } 
        } else {
                 value := int64(0)
                if  m.Defaults["id"] != nil {
                    value = m.Defaults["id"].(int64)
                }
                if m.Id != value {
                    builder["id"] = m.Id
                } 
        }
        if dontIgnore {
                value := ""
                if  m.Defaults["fieldValue"] != nil {
                    value = m.Defaults["fieldValue"].(string)
                }
                if m.FieldValue != value {
                    builder["fieldValue"] = m.FieldValue
                } else {
                     builder["fieldValue"] = value
                } 
        } else {
                value := ""
                if  m.Defaults["fieldValue"] != nil {
                    value = m.Defaults["fieldValue"].(string)
                }
                if m.FieldValue != value {
                    builder["fieldValue"] = m.FieldValue
                } 
        }
    return builder
}

func (m *Model) ToMapWithFilter(builder map[string]interface{},filter []string,dontIgnoreEmtpy ...bool) map[string]interface{} {
    if builder == nil {
        builder = make(map[string]interface{})
    }
    var dontIgnore bool
    if len(dontIgnoreEmtpy) > 0 {
         dontIgnore = dontIgnoreEmtpy[0]
    }
    for _,item := range filter {
            switch strings.ToLower(item) {
                case "id":
                     if dontIgnore {
                            value := int64(0)
                            if  m.Defaults["id"] != nil {
                                value = m.Defaults["id"].(int64)
                            }
                            if m.Id != value {
                                builder["id"] = m.Id
                            } else {
                                builder["id"] = value
                            } 
                    } else {
                            value := int64(0)
                            if  m.Defaults["id"] != nil {
                                value = m.Defaults["id"].(int64)
                            }
                            if m.Id != value {
                                builder["id"] = m.Id
                            } 
                    }
                case "fieldvalue":
                     if dontIgnore {
                            value := ""
                            if  m.Defaults["fieldvalue"] != nil {
                                value = m.Defaults["fieldvalue"].(string)
                            }
                            if m.FieldValue != value {
                                builder["fieldvalue"] = m.FieldValue
                            } else {
                                builder["fieldvalue"] = value
                            } 
                    } else {
                            value := ""
                            if  m.Defaults["fieldvalue"] != nil {
                                value = m.Defaults["fieldvalue"].(string)
                            }
                            if m.FieldValue != value {
                                builder["fieldvalue"] = m.FieldValue
                            } 
                    }
            }
    }
    return builder
}

//note: use the actual column name registered at the database
func (m *Model) ToStringWithFilter(builder *strings.Builder, filter []string, injector func(field string)) {
    for _,item := range filter {
        switch strings.ToLower(item) {
            case "id":
                if builder.Len() > 0 {
                    builder.WriteString(",")
                }
                    builder.WriteString(fmt.Sprintf("%d",m.Id))
            case "fieldvalue":
                if builder.Len() > 0 {
                    builder.WriteString(",")
                }
                    builder.WriteString("\"")    
                    builder.WriteString(m.FieldValue)
                    builder.WriteString("\"")
            default:
                injector(item)
        }
    }
}
