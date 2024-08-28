package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sync"
)

type Data struct {
	gorm.Model
	GoRoutine int
	Step      int
	Content   string
}

var wg sync.WaitGroup

func checkError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func Create(db *gorm.DB, GoRoutine int, Step int, Content string) {
	defer wg.Done() // COMMENT FOR FIRST AND SECOND WAY
	data := new(Data)
	data.GoRoutine = GoRoutine
	data.Step = Step
	data.Content = Content
	db.Create(data)
	fmt.Println("Call (GoRoutine):", GoRoutine, "Step:", Step)
}

func main() {
	db, err := gorm.Open(sqlite.Open("example.db"), &gorm.Config{})
	checkError(err)

	err = db.AutoMigrate(&Data{})
	checkError(err)

	//// FIRST WAY - INSERT MULTIPLE WITHOUT GO ROUTINES
	//for i := 1; i <= 10; i++ {
	//	for j := 1; j <= 10; j++ {
	//		Create(db, j, i, "Text")
	//	}
	//}

	//// SECOND WAY - INSERT MULTIPLE USING GO ROUTINES (WITHOUT WAITING FOR GO ROUTINES FINISHES)
	//// PROBABLY NO ONE INSERTS (OR NOT ALL INSERTS) ARE EXECUTED, BECAUSE MAIN FINISHES BEFORE THE GO ROUTINES FINISHES
	//for i := 1; i <= 10; i++ {
	//	for j := 1; j <= 10; j++ {
	//		go Create(db, j, i, "Text")
	//	}
	//}
	//time.Sleep(100 * time.Millisecond) // WAIT A TIME TO FINISH SOME GO ROUTINES

	// THIRD WAY - INSERT MULTIPLE USING GO ROUTINES (AND WAITING FOR GO ROUTINES FINISHES USING A WAIT GROUP)
	// WE CAN CONTROL THE END OF GO ROUTINES USING A *CHANNEL* (IF YOU NEED A RESPONSE FROM THE GO ROUTINES OR DON'T KNOW
	// THE NUMBER OF EXECUTIONS OR *WAIT GROUP* IF YOU KNOW THE NUMBER OF EXECUTIONS
	// *** CAUTION ***
	// IN THE BELOW EXAMPLE, WE WILL SEND CLOSE TO 100 INSERT IN A FEW PERIOD - WE NEED TO BE CAREFUL ABOUT LOTS OF
	// DB CALLS, ENDPOINT CALLS, FILE OPERATIONS TO AVOID ERRORS.
	// ANOTHER THING WE NEED TO BE CAREFUL IS ABOUT POSSIBLE DEADLOCKS (NOT IN THIS CASE) WHEN MULTIPLE GO ROUTINES ACCESS OR UPDATE
	// A SAME RESOURCE - WE NEED TO USE A "MUTEX" LOCK / UNLOCK TO AVOID IT
	workers := 100
	wg.Add(workers)

	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			go Create(db, j, i, "Text")
		}
	}
	wg.Wait()
}
