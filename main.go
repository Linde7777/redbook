package main

func main() {
	router := InitWebServer()
	router.Run(":8080")
}
