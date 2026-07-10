package conv

import (
    "database/sql/driver"
    "fmt"
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/schema"
)

// Date stores date-only (YYYY-MM-DD) information for DB and JSON.
type Date struct {
    time.Time
}

const dateLayout = "2006-01-02"

func (d Date) MarshalJSON() ([]byte, error) {
    if d.Time.IsZero() {
        return []byte("null"), nil
    }
    s := fmt.Sprintf("\"%s\"", d.Time.Format(dateLayout))
    return []byte(s), nil
}

func (d *Date) UnmarshalJSON(b []byte) error {
    s := string(b)
    if s == "null" || s == `""` {
        d.Time = time.Time{}
        return nil
    }
    if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
        s = s[1 : len(s)-1]
    }
    t, err := time.Parse(dateLayout, s)
    if err != nil {
        t2, err2 := time.Parse(time.RFC3339, s)
        if err2 != nil {
            return err
        }
        t = t2
    }
    d.Time = t
    return nil
}

func (d Date) Value() (driver.Value, error) {
    if d.Time.IsZero() {
        return nil, nil
    }
    return d.Time.Format(dateLayout), nil
}

func (d *Date) Scan(value interface{}) error {
    if value == nil {
        d.Time = time.Time{}
        return nil
    }
    switch v := value.(type) {
    case time.Time:
        d.Time = v
    case []byte:
        t, err := time.Parse(dateLayout, string(v))
        if err != nil {
            t2, err2 := time.Parse(time.RFC3339, string(v))
            if err2 != nil {
                return err
            }
            d.Time = t2
        } else {
            d.Time = t
        }
    case string:
        t, err := time.Parse(dateLayout, v)
        if err != nil {
            t2, err2 := time.Parse(time.RFC3339, v)
            if err2 != nil {
                return err
            }
            d.Time = t2
        } else {
            d.Time = t
        }
    default:
        return fmt.Errorf("cannot scan %T into Date", value)
    }
    return nil
}

func (Date) GormDataType() string {
    return "date"
}

func (Date) GormDBDataType(db *gorm.DB, field *schema.Field) string {
    switch db.Dialector.Name() {
    case "mysql":
        return "date"
    case "sqlite":
        return "date"
    case "postgres":
        return "date"
    default:
        return "date"
    }
}
