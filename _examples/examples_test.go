package examples

import (
	"testing"

	"github.com/kasiss-liu/easyvalid/_examples/Animal"
)

func TestDog(t *testing.T) {
	dog := &Dog{
		am: Animal.Animal{},
		Food: [2]string{
			"meat",
			"fish",
		},
	}
	t.Log(dog.EasyValid())
}
