package main

import (
	"zestream/service"
)

func dev() {
	// utils.Fetch("https://file-examples.com/storage/fe332cf53a63a4bd5991eb4/2017/04/file_example_MP4_480_1_5MG.mp4", "earth.mp4")
	service.GenerateDash("www.mp4")
}

func main() {

	dev()

	// e := godotenv.Load()

	// if e != nil {
	// 	fmt.Print(e)
	// }

	// r := routes.Init()

	// port := os.Getenv(constants.PORT)

	// if port == "" {
	// 	port = constants.DEFAULT_PORT
	// }

	// err := http.ListenAndServe(port, nil)

	// if err != nil {
	// 	fmt.Println(err)
	// }

	// r.Run(":" + port)

}
