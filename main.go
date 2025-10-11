package main

import (
    "encoding/json"
	"flag"
	"github.com/Zordddd/learning/packages"
	"fmt"
	"log"
	"sort"
	// "os"
)

var sun = flags.SunFlag("sun", true, "солнце")

func main() {
    flag.Parse()
    fmt.Println(*sun)
    n := BackupOne()
    fmt.Println(n)
    var l List
    l.AddNode(1)
    l.AddNode(2)
    l.AddNode(3)
    l.AddNode(4)

    l.Print()

    var a array
    a.Append(1, 2, 3)
    fmt.Println(a)
}

func BackupOne() (value interface{}) {
	defer func() {
		if p := recover(); p != nil {
            value = p
		}
	}()
	panic(1)
}

type Address struct {
    City string
    Street string
}

type Person struct {
    Name string
    Address
}

type array []int

func (a *array) Append(numbers ...int) {
    *a = array(append([]int(*a), numbers...))
}

type Timer interface {
    Time(time int) int
}

func Factorial(x int) int {
    if x <= 1 {
        return 1
    }
    return x * Factorial(x-1)
}

func Crypt() {
    alf := make(map[rune]int)
    srAlf := make([]rune, 0, 33)
    str := "ОЮЭЮКУС УСОШСЕН ДОАФЮЕЮЕАЮ ДОЮШЛКСЖЯИЯС ЛНРНЫ ДЮОЖНЮ ЭАМОС ДЮОЮЛКСЕНЖУА Ж ЛНЖОЮФЮЕЕНФ ФАОЮ. ЕС ЮЮ НЛЕНЖЮ ЖН ЖОЮФИ ЖКНОНЫ ФАОНЖНЫ ЖНЫЕЦ РЦЯ ЛНИШСЕ НШАЕ АЙ ЕСАРНЯЮЮ ЛКНЫУАП ЖНЮЕЕН-ФНОЛУАП ЭАМОНЖ ЖЮЯАУНРОАКСЕАА" 
    newstr := make([]rune, len(str))
    sum := 0.

    for _, c := range str {
        alf[c]++
        sum++
    }
    for k := range alf {
        srAlf = append(srAlf, k)
    }
    sort.Slice(srAlf, func(i int, j int) bool {
        return alf[srAlf[i]] < alf[srAlf[j]]
    })

    for _, k := range srAlf {
        fmt.Printf("%c\t%d\t%g\n", k, alf[k], float64(alf[k]) / sum)
    }

    for i, v := range str {
        if v == 'Н' {
            newstr[i] = 'о'
        }else if v == 'Ю' {
            newstr[i] = 'е'
        } else if v == 'Ж' {
            newstr[i] = 'в'
        } else if v == 'С' {
            newstr[i] = 'а'
        }else if v == 'Е' {
            newstr[i] = 'н'
        }else if v == 'О' {
            newstr[i] = 'р'
        }else if v == 'А' {
            newstr[i] = 'и'
        }else if v == 'Р' {
            newstr[i] = 'б'
        }else if v == 'Я' {
            newstr[i] = 'л'
        }else if v == 'Ц' {
            newstr[i] = 'ы'
        }else if v == 'Ф' {
            newstr[i] = 'м'
        }else if v == 'Л' {
            newstr[i] = 'с'
        }else if v == 'Д' {
            newstr[i] = 'п'
        }else if v == 'Ы' {
            newstr[i] = 'й'
        }else if v == 'К' {
            newstr[i] = 'т'
        }else if v == 'И' {
            newstr[i] = 'я'
        }else if v == 'У' {
            newstr[i] = 'к'
        }else if v == 'Ш' {
            newstr[i] = 'д'
        }else if v == 'П' {
            newstr[i] = 'х'
        }else if v == 'Э' {
            newstr[i] = 'ш'
        }else if v == 'М' {
            newstr[i] = 'Ф'
        }else if v == 'Й' {
            newstr[i] = 'з'
        }else {
            newstr[i] = v
        }
        fmt.Printf("%c", newstr[i])
    }
    
}

func DeepFirstSerch(graph map[string][]string, currentDot string) {
    // if graph[currentDot] == nil {
    //     fmt.Println("Wrong start")
    //     os.Exit(1)
    // }

    visited[currentDot] = true
    if len(visited) != 1 {
        fmt.Printf(" -> %s", currentDot)
    } else {
        fmt.Printf("%s", currentDot)
    }

    for _, v := range graph[currentDot] {
        if ok := visited[v]; !ok {
            DeepFirstSerch(graph, v)
        } else {
            continue
        }
    }
}

func MarshalingTest() {
    fmt.Println("hi")
    names := make(map[string]int)
    names["Dima"]++
    fmt.Printf("%#v\n", names)

    client := Person{Name : "Mama", Address: Address{City: "Likino-Dulevo", Street: "Stepana Morozkina"}}
    var newClient Person
    
    jsonClient, err := json.MarshalIndent(client, "", "\t")

    if err != nil {
        log.Fatalf("%v\n", err)
    }
    fmt.Printf("%s\n", jsonClient)

    err = json.Unmarshal(jsonClient, &newClient)
    if err != nil {
        log.Fatalf("Error Unmarshaling : %v", err)
    }

    jsonNames, err := json.Marshal(names)
    if err != nil {
        log.Fatalf("Error Marshaling : %v", err)
    }

    fmt.Printf("%+s\n%+s\n", jsonNames, newClient)

}
var visited = map[string]bool{}
var Prereqs = map[string][]string{
    "algorithms":            {"data structures"},
    "calculus":              {"linear algebra"},
    "compilers":             {"data structures", "formal languages", "computer organization"},
    "data structures":       {"discrete math"},
    "databases":             {"data structures"},
    "discrete math":         {"intro to programming"},
    "formal languages":      {"discrete math"},
    "networks":              {"operating systems"},
    "operating systems":     {"data structures", "computer organization"},
    "programming languages": {"data structures", "computer organization"},
}

type Node struct {
    value int
    next *Node
}

type List struct {
    head *Node
    tail *Node
}

func (l *List) AddNode(value int) {
    newNode := &Node{value : value}

    if l.head == nil {
        l.head = newNode
        l.tail = newNode
    } else {
        l.tail.next = newNode
        l.tail = newNode
    }
}

func (l *List) Print() {
    currentNode := l.head
    for currentNode != nil {
        fmt.Printf("%d ", currentNode.value)
        currentNode = currentNode.next
    }
    fmt.Printf("\n")
}

