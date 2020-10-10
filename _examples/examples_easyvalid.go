// Code generated by easyvalid v0.0.2 . DO NOT EDIT.

package examples

import (
    "errors"
    "github.com/kasiss-liu/easyvalid/_examples/Animal"
)

func (c Cat) EasyValid () (err error) {
    if err = c.valid_bb282d14dcab326bee3ffc600f9c6635f44d68ea_am(); err != nil {
        return
    }
    if err = c.valid_bb282d14dcab326bee3ffc600f9c6635f44d68ea_Food(); err != nil {
        return
    }
    return nil
}
func (c Cat) valid_bb282d14dcab326bee3ffc600f9c6635f44d68ea_am() (err error) {
    if c.am == nil { return errors.New(`am属性不能为空`) }
    return nil
}

func (c Cat) valid_bb282d14dcab326bee3ffc600f9c6635f44d68ea_Food() (err error) {
    if !(len(c.Food) > 1) { return errors.New(`至少1种食物`) }
    if !(len(c.Food) <= 2) { return errors.New(`最多2种食物`) }
    return nil
}

func (d Dog) EasyValid () (err error) {
    if err = d.valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_am(); err != nil {
        return
    }
    if err = d.valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_Food(); err != nil {
        return
    }
    if err = d.valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_name(); err != nil {
        return
    }
    if err = d.valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_Age(); err != nil {
        return
    }
    return nil
}
func (d Dog) valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_am() (err error) {
    if d.am == (Animal.Animal{}) { return errors.New(`am属性不能为空`) }
    return nil
}

func (d Dog) valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_Food() (err error) {
    if !(len(d.Food) > 1) { return errors.New(`至少1种食物`) }
    if !(len(d.Food) <= 2) { return errors.New(`最多2种食物`) }
    return nil
}

func (d Dog) valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_name() (err error) {
    if !(d.name == "xiaogou") { return errors.New(`名称异常`) }
    if d.name == "" { return errors.New(`不能为空`) }
    return nil
}

func (d Dog) valid_f8feab49dc77f3828533615d7950b2f9937f0ec2_Age() (err error) {
    if !(d.Age >= 2) { return errors.New(`年龄必须2周岁以上`) }
    return nil
}


