package configs


type Configuration struct{
	DB 			DbConfig
	Auth 		AuthConfig
	Server 		ServerConfig
}


type DbConfig struct{
	Host		string
	Port		string
	DbName		string
	Username	string
	Password	string
}


type AuthConfig struct{
	Secret		string
}

type ServerConfig struct{
	Host 	string
	Port 	string
}


// func LoadConfig() (*Configuration, error) {
// 	var config Configuration
// 	var err error



// }