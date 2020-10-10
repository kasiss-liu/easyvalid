package examples

import "github.com/kasiss-liu/easyvalid/_examples/Animal"

//Dog is type of animal
// easyvalid:valid:d
type Dog struct {
	am   Animal.Animal `evalid:"notEmpty|msg:am属性不能为空"`
	Food [2]string     `evalid:"len:gt(1)|msg:至少1种食物;len:lte(2)|msg:最多2种食物"`
	name string        `evalid:"value:eq(xiaogou)|msg:名称异常;notEmpty|msg:不能为空"`
	Age  uint          `evalid:"value:lte(2)|msg:年龄必须2周岁以上"`
}

func (d *Dog) GetName() string {
	return "d.am.Name"
}

func (d *Dog) GetFood() string {
	return "d.Food"
}

// easyvalid:valid:c
type Cat struct {
	am   *Animal.Animal `evalid:"notEmpty|msg:am属性不能为空"`
	Food []string       `evalid:"len:gt(1)|msg:至少1种食物;len:lte(2)|msg:最多2种食物"`
}

func (c *Cat) GetName() string {
	return "d.am.Name"
}

func (c *Cat) GetFood() string {
	return "d.Food"
}
