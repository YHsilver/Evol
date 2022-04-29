package domain

import (
	"evol"
	"evol/repo/memory"
)

var repo evol.ReadRepository = &memory.ModelRepo{}

func QueryUserOrders(uid string) []Order {
	//repo.Find(uid )

	return Orders[uid] //access mock db
}
