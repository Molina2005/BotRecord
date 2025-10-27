package main

import (
	"fmt"
	"log"
	apitelegram "modulo/ApiTelegram"
	"modulo/connect"
	"time"
)

func main() {
	// Conexion a la base de datos
	conexion, err := connect.Connect()
	if err != nil {
		log.Fatal(err)
	}
	// Cierre de conexion
	defer conexion.Close()
	fmt.Println("Conexion exitosa a PosgreSQL")

	// Conexion y usabilidad bot Telegram
	apitelegram.BotTelegram(conexion)
	// Tiempo espera
	time.Sleep(30 * time.Minute)
}
