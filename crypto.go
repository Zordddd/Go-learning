package main

import (
	"fmt"
	"sort"
)

func Crypt() {
    alf := make(map[rune]int)
    srAlf := make([]rune, 0, 33)
    str := "ОЛХВЗХЭВЧЯ ЭВЧИТХЗЪЭВВЛПУЖСЕСЦАЦЯЭАЖЭЭМЗОЭВИШЛПМЛЯЛЦЭПЭЖЭВДУФУЖСРЭЯЭМУЖ ЗЪСЛЦЗЖЛОЧВЖЧЙОЗИЭЫЖСЯЗИЛАСЪЗЭЯЛМЪЗФЧЯСИИШСЕСАГЗЭЗИВЛЯЗКЩЗЕЖЗЭМЛМЛИОЛЦЗЖСИЦЭЪСАЮВЛЛЖИВСЪЯЛЦЛЖСХСЪБЖЗШЛПЦЛШУПЭЖВСЪБЖЛЕСФЗШИЗЯЛАСЖЖЛЙЗИВЛЯЗЗШЯЗОВЛМЯСФЗЗЭМЛИЗИВЭПСЖЭТАЪТЭВИТВСЙЖЛОЗИБКАВЛПАЗЦЭАШСШЛПЛЖСЗЕАЭИВЖСИЛАЯЭПЭЖЖЛПУПЗЯУЦЪТЕСИЭШЯЭХЗАСЖЗТИА ЛЭЙЖСЦОЗИЗЛЖЖЭЗИОЛЪБЕЛАСЪЖЗШСШЛМЛОЛЪЖЛЫЭЖЖЛМЛГЗФЯС" 
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
        if v == 'Л' {
            newstr[i] = 'о'
        }else if v == 'З' {
            newstr[i] = 'е'
        } else if v == 'Ж' {
            newstr[i] = 'и'
        } else if v == 'Э' {
            newstr[i] = 'н'
        }else if v == 'С' {
            newstr[i] = 'а'
        }else if v == 'И' {
            newstr[i] = 'т'
        }else if v == 'В' {
            newstr[i] = 'с'
        }else if v == 'H' {
            newstr[i] = 'l'
        }else if v == 'T' {
            newstr[i] = 'x'
        }else if v == 'B' {
            newstr[i] = 'm'
        }else if v == 'L' {
            newstr[i] = 'g'
        }else if v == 'U' {
            newstr[i] = 'y'
        }else if v == 'D' {
            newstr[i] = 'c'
        }else if v == 'Y' {
            newstr[i] = 's'
        }else if v == 'E' {
            newstr[i] = 'i'
        }else if v == 'M' {
            newstr[i] = 'p'
        }else if v == 'W' {
            newstr[i] = 'n'
        }else if v == 'F' {
            newstr[i] = 'w'
        }else if v == 'A' {
            newstr[i] = 'f'
        }else if v == 'O' {
            newstr[i] = 'r'
        }else if v == 'Z' {
            newstr[i] = 'h'
        }else if v == 'G' {
            newstr[i] = 'k'
        }else if v == 'V' {
            newstr[i] = 'v'
        }else if v == 'R' {
            newstr[i] = 'z'
        }else {
            newstr[i] = v
        }
        fmt.Printf("%c", newstr[i])
    }
    
}
