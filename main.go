package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)


var machine = []string{"Drone", "Wheeled", "HeavyLifter"}

type Product struct {
	id     int
	name   string
	weight float64
}
type Robot struct {
	id          int
	model       string
	battery     float64
	isAvailable bool
}

type Mover interface {
	Move(distance float64) error
}

func Recharge(robot *Robot) {
	robot.battery = 100.0
}

func WeightLimit(r *Robot, p *Product) bool {

	maxweight := 0
	switch r.model {
	case machine[0]:
		maxweight = 2
	case machine[1]:
		maxweight = 345
	case machine[2]:
		maxweight = 750

	}

	if p.weight > float64(maxweight) {
		fmt.Println("Ошибка вес слишком большой")
		return false
	} else {
		return true
	}

}

func (d *Robot) Mover(distance float64) bool {

	switch d.model {
	case machine[0]:
		energycost := distance * 1.5
		if d.battery < energycost {
			fmt.Println("У робота нет заряда")
			return false
		} else {
			d.battery -= energycost
			return true
		}
	case machine[1]:
		energycost := distance * 0.8
		if d.battery < energycost {
			fmt.Println("У робота нет заряда")
			return false
		} else {
			d.battery -= energycost
			return true
		}
	case machine[2]:
		energycost := distance * 2.5
		if d.battery < energycost {
			fmt.Println("У робота нет заряда")
			return false
		} else {
			d.battery -= energycost
			return true
		}
	default:
		energycost := distance
		if d.battery < energycost {
			fmt.Println("У робота нет заряда")
			return false
		} else {
			d.battery -= energycost
			return true
		}
	}

}

func DeliverTask(r *Robot, p Product, distance int, ch chan string, wg *sync.WaitGroup) {
	defer wg.Done()
	if WeightLimit(r, &p) {
		result := r.Mover(float64(distance))
		if result {
			fmt.Println("Робот", r.id, "Выехал")

			time.Sleep(3 * time.Second)

			fmt.Println(strings.ToUpper("Робот доставил заказ"))

			ch <- strings.ToUpper("Заказ " + p.name + " доставлен")
		} else {
			ch <- "Заказ не доставлен. Не смогли доставить: " + p.name

		}
	} else {
		ch <- "Заказ слишком тяжелый. Не смогли доставить : " + p.name
	}
	r.isAvailable = true

}

func main() {
	const BaseDeliveryRate = 10
	const MaxBatteryLevel = 100.0
	const EmergencyThreshold = 0.20

	inventory := map[string]int{"Box1": 2, "Box2": 4, "Box3": 6}
	availableMachine := []Robot{
		{id: 1, model: "Drone", battery: 100, isAvailable: true},
		{id: 2, model: "Wheeled", battery: 100, isAvailable: true},
		{id: 3, model: "HeavyLifter", battery: 100, isAvailable: true},
	}
	bufferCapacity := [5]Product{}

	orderCount := 0

	for {

		var text3 string
		fmt.Println("Запустить процесс доставки?\n Напишите 1 если хотите добавить товар в заказ\n Напишите 2 если хотите увидеть состаяние склада\n Напишите 3 если хотите запустить процесс доставки\n Если хотите выйдет напиши 'exit'")
		fmt.Scan(&text3)
		activeJobs := 0
		if text3 == "1" {

			if orderCount >= 5 {
				fmt.Println("Слишком много заказов на данный момент повторите позже")
			} else {
				fmt.Println("Введите имя заказа: ")
				name := ""
				fmt.Scan(&name)

				if v, ok := inventory[name]; ok && v > 0 {
					inventory[name]--
					fmt.Println("Введите вес заказа: ")
					weight := 0
					fmt.Scan(&weight)
					newProduct := Product{orderCount, name, float64(weight)}
					bufferCapacity[orderCount] = newProduct
					orderCount++
				} else {
					fmt.Println("Такого товара нет на складе")
				}

			}

		}

		if text3 == "2" {
			for i, v := range inventory {
				fmt.Printf("Товар: %-10s Количество: %d\n", i, v)
			}
			for _, v := range availableMachine {
				fmt.Printf("id: %d Робот: %s Заряд батерии: %f Доступен:%t\n", v.id, v.model, v.battery, v.isAvailable)
			}

			fmt.Println(bufferCapacity)
		}
		if text3 == "3" {
			var wg sync.WaitGroup

			ch := make(chan string)

			for _, v := range bufferCapacity {
				if v.name == "" {
					continue
				}
				flag := true
				for i, r := range availableMachine {
					if availableMachine[i].isAvailable && WeightLimit(&r, &v) {
						orderCount -= 1
						flag = false
						wg.Add(1)
						activeJobs += 1
						availableMachine[i].isAvailable = false
						go DeliverTask(&availableMachine[i], v, 10, ch, &wg)
						break
					}
				}
				if flag {
					fmt.Println("Нет доступных исполнителей для заказа: " + v.name)
				}

			}
			go func() {
				wg.Wait()
				close(ch)
			}()
			for v := range ch {
				fmt.Println(v)
			}

			bufferCapacity = [5]Product{}
			orderCount = 0
		} else if text3 == "exit" {
			break
		}

	}

}
