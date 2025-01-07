package components

import (
	"fmt"
	"sync"

	"github.com/mbndr/figlet4go"
)

var once sync.Once
var logo string

func GetLogo() string {
	once.Do(func() {
		ascii := figlet4go.NewAsciiRender()
		var err error
		logo, err = ascii.Render("GophKeeper")
		if err != nil {
			fmt.Println("Ошибка при рендеринге текста:", err)
			panic(err)
		}
	})
	return logo
}
