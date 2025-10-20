package main

import (
	"fmt"
	"log"
	apitelegram "modulo/ApiTelegram"
	"modulo/connect"
	"time"
)

func main() {
	// conexion a la base de datos
	conexion, err := connect.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// cierre de conexion
	defer conexion.Close()
	fmt.Println("Conexion exitosa a PosgreSQL")

	// Conexion y usabilidad bot Telegram
	apitelegram.BotTelegram(conexion)

	time.Sleep(30 * time.Minute)
}
