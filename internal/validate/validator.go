package validate

// DBValidator отвечает за проверку данных по базе.
// Можно добавить на случай спама одними и теми же хэшами, на случай если нашего livetime не хватит.
//type DBValidator interface {
//	Validate(data []byte) error
//}

type Validator interface {
	Validate(ipAddress string, data []byte) error
}

type Connection struct {
	validators []Validator
}

func NewConnection(validators []Validator) *Connection {
	return &Connection{validators: validators}
}

func (c *Connection) Validate(ipAddress string, data []byte) (err error) {
	for _, v := range c.validators {
		if err = v.Validate(ipAddress, data); err != nil {
			break
		}
	}
	return err
}
