package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

type table struct {
	grid   [3][3]int
	value  int
	emptyI int
	emptyJ int
}

type queue struct {
	firstNode *node
}

func qInit(n *node) queue {
	q := queue{n}
	return q
}

func (q *queue) add(n *node) {
	//q.show()
	temp := q.firstNode
	for {
		if temp.next == nil || temp.next == temp {
			temp.next = n
			return
		}
		if temp.next.tableInfo.value >= n.tableInfo.value {
			n.next = temp.next
			temp.next = n
			return
		}
		temp = temp.next
	}
}

func (q *queue) show() {
	t := q.firstNode
	fmt.Println("start queue")
	for i := 0; t.next != nil; i++ {
		fmt.Println(i, ":", t.tableInfo.value)
		if t.next != nil {
			t = t.next
			if t == t.next {
				break
			}
		}
	}
	fmt.Println("end queue")
	time.Sleep(time.Second * 1)
}

func (q *queue) count() int {
	t := q.firstNode
	i := 0
	for ; t.next != nil; i++ {
		t = t.next
	}
	return i
}

type node struct {
	tableInfo table
	next      *node
}

func (t *table) getValue() int {
	t.value = 0
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 2 && j == 2 {
				if t.grid[i][j] != -1 {
					t.value++
				}
				return t.value
			}
			if t.grid[i][j] != i*3+j+1 {
				t.value++
			}
		}
	}
	return t.value
}

/*func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
*/

func (t *table) getMovable() [][2]int {
	result := make([][2]int, 0)
	if t.emptyI-1 >= 0 {
		result = append(result, [2]int{t.emptyI - 1, t.emptyJ})
	}
	if t.emptyI+1 <= 2 {
		result = append(result, [2]int{t.emptyI + 1, t.emptyJ})
	}
	if t.emptyJ-1 >= 0 {
		result = append(result, [2]int{t.emptyI, t.emptyJ - 1})
	}
	if t.emptyJ+1 <= 2 {
		result = append(result, [2]int{t.emptyI, t.emptyJ + 1})
	}
	t.getValue()
	return result
}

func (t table) move(i, j int) table {
	t.grid[t.emptyI][t.emptyJ], t.grid[i][j] = t.grid[i][j], t.grid[t.emptyI][t.emptyJ]
	t.emptyI = i
	t.emptyJ = j
	return t
}

func (t *table) generateTable() {
	t.grid = [3][3]int{}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			t.grid[i][j] = i*3 + j + 1
		}
	}
	t.grid[2][2] = -1
	t.emptyI = 2
	t.emptyJ = 2
}

func (t *table) randomizeTable() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			rand.Seed(time.Now().UnixNano())
			i1, j1 := rand.Intn(3-i), rand.Intn(3-j)
			t.grid[2-i][2-j], t.grid[i1][j1] = t.grid[i1][j1], t.grid[2-i][2-j]
		}
	}
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if t.grid[i][j] == -1 {
				t.emptyI = i
				t.emptyJ = j
				return
			}
		}
	}
}

func (t *table) showTable() {
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			fmt.Print(t.grid[i][j], " | ")
		}
		fmt.Println(t.grid[i][2])
	}
}

func sort(head table) {
	q := qInit(&node{head, nil})
	closedList := make([][3][3]int, 0)
	for i := 0; q.firstNode != nil; i++ {
		if i%1000 == 0 {
			PrintMemUsage()
			fmt.Printf("#%7d\n", i)
			fmt.Printf("vals: %d\n", q.count())
			fmt.Printf("minval: %d\n", q.firstNode.tableInfo.value)
			fmt.Println("-----------------")
		}
		if q.firstNode.tableInfo.value == 0 {
			q.firstNode.tableInfo.showTable()
			fmt.Println(">>> Ready!")
			return
		}
		movable := q.firstNode.tableInfo.getMovable()
		for _, v := range movable {
			newTable := q.firstNode.tableInfo.move(v[0], v[1])
			newTable.getValue()
			b := true
			for _, v1 := range closedList {
				if v1 == newTable.grid {
					b = false
					break
				}
			}
			if b {
				q.add(&node{
					tableInfo: newTable,
				})
				closedList = append(closedList, newTable.grid)
			}
		}
		q.firstNode = q.firstNode.next
		///time.Sleep(time.Second*3)
	}
}
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
func main() {
	t := table{}

	t.generateTable()
	t.randomizeTable()
	t.showTable()
	fmt.Println(t.getValue())
	sort(t)
}
