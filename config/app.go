package config

import "os"

func GetAppEnvironment() string {
	env := os.Getenv("APP_ENV")
	if env != "" {
		return env
	}
	return "development"
}

func GetAppName() string {
	env := os.Getenv("APP_NAME")
	if env != "" {
		return env
	}
	return ""
}

func GetAppPort() string {
	env := os.Getenv("APP_PORT")
	if env != "" {
		return env
	}
	return "8000"
}

func JWTSecretKey() []byte {
	env := os.Getenv("JWT_SECRET_KEY")
	if env != "" {
		return []byte(env)
	}
	return []byte("secret")
}
