package config

type password string

// Configurations exported
type Configurations struct {
	Server   ServerConfigurations
	Database DatabaseConfigurations
	Logging  LoggingConfigurations
}

// ServerConfigurations exported
type ServerConfigurations struct {
	Port  int  `json:"port"`
	Debug bool `json:"debug"`
}

// DatabaseConfigurations exported
type DatabaseConfigurations struct {
	DBHost     string   `json:"db_host"`
	DBPort     int      `json:"db_port"`
	DBName     string   `json:"db_name"`
	DBUser     string   `json:"db_user"`
	DBPassword password `json:"db_password"`
}

type LoggingConfigurations struct {
	Main map[string]interface{} `json:"main"`
}

// Marshaler ignores the field value completely.
func (password) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}
