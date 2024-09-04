package helper

import "os"

var SECRET_KEY string = os.Getenv("secret")
