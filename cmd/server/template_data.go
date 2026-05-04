package main

import "botai/botaigolang/internal/service" // [THEM_MOI]

type templateData struct { // [THEM_MOI]
	Form     any             // [THEM_MOI]
	Account  service.Account // [THEM_MOI]
	UserID   string          // [THEM_MOI]
	Orders   any             // [THEM_MOI]
	IsAuthen bool            // [SUA]
}
