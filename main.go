package main

func main() {
	r := setupRouter()
	r.Run(":3000") // listen and serve on 0.0.0.0:8080
}
