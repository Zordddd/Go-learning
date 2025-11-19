// 30
package main

import (
	"fmt"
	"sort"
)

func main() {
	Visioner()

}

func Visioner() {
	message := []byte("ОНЯФЩН ФЪШИУПНРФФ – КЩИ ЪЧ ЩИЭБГИ ЩЧЕЪФЮЧХГНЗ ОНЬНЮН, ЪИ Ф ИУВНЪФОНРФИЪЪНЗ. ЪЧИДЕИЬФПИ УНОУНДИЩНЩБЖИЭФЩФГФ Ф ЖУИРЧЬСУЛ, ГИЩИУЛЧ ДСЬСЩ УЧВЭНПЧЪЩФУИЦНЩБ ЬИХЩСЖ Г ФЪШИУПНРФФ Ф ИДЧХЖЧЮФЦНЩБ ЧЧХИЕУНЪЪИХЩБ. ЦНМЪИ ЩНГМЧ ЖУИЦИЬФЩБ УЧВСЭЗУЪЛЧ НСЬФЩЛ ДЧОИЖНХЪИХЩФ Ф ИРЧЪФЦНЩБ КШШЧГЩФЦЪИХЩБЖУФЪФПНЧПЛЕ ПЧУ.")
	//keyLength := 4
	//groups := make([][]byte, keyLength)
	//
	//for i := 0; i < keyLength; i++ {
	//	groups[i] = make([]byte, 0, len(message)/keyLength+keyLength)
	//}
	//
	//for i, ch := range message {
	//	groups[i%keyLength] = append(groups[i%keyLength], ch)
	//}
	//
	//for i := range groups {
	//	Crypt(string(groups[i]))
	//	fmt.Println()
	//}
	//
	//// Правильное дешифрование
	//keys := []int{13, 9, 12, 11} // Пробуем ваш ключ
	//
	//for groupIdx := 0; groupIdx < keyLength; groupIdx++ {
	//	for i := range groups[groupIdx] {
	//		// Дешифрование: (cipher - key) mod 26
	//		decrypted := (int(groups[groupIdx][i]) - 'A' - keys[groupIdx]) % 26
	//		if decrypted < 0 {
	//			decrypted += 26
	//		}
	//		groups[groupIdx][i] = byte(decrypted + 'A')
	//	}
	//}
	//
	//// Собираем результат
	//result := make([]byte, 0, len(message))
	//for i := 0; i < len(groups[0]); i++ {
	//	for groupIdx := 0; groupIdx < keyLength; groupIdx++ {
	//		if i < len(groups[groupIdx]) {
	//			result = append(result, groups[groupIdx][i])
	//		}
	//	}
	//}
	//
	//fmt.Printf("Расшифрованный текст: %s\n", string(result))
	Crypt(string(message))
}

func Crypt(str string) {
	alf := make(map[rune]int)
	srAlf := make([]rune, 0, 33)
	//str := "ОЛХВЗХЭВЧЯ ЭВЧИТХЗЪЭВВЛПУЖСЕСЦАЦЯЭАЖЭЭМЗОЭВИШЛПМЛЯЛЦЭПЭЖЭВДУФУЖСРЭЯЭМУЖ ЗЪСЛЦЗЖЛОЧВЖЧЙОЗИЭЫЖСЯЗИЛАСЪЗЭЯЛМЪЗФЧЯСИИШСЕСАГЗЭЗИВЛЯЗКЩЗЕЖЗЭМЛМЛИОЛЦЗЖСИЦЭЪСАЮВЛЛЖИВСЪЯЛЦЛЖСХСЪБЖЗШЛПЦЛШУПЭЖВСЪБЖЛЕСФЗШИЗЯЛАСЖЖЛЙЗИВЛЯЗЗШЯЗОВЛМЯСФЗЗЭМЛИЗИВЭПСЖЭТАЪТЭВИТВСЙЖЛОЗИБКАВЛПАЗЦЭАШСШЛПЛЖСЗЕАЭИВЖСИЛАЯЭПЭЖЖЛПУПЗЯУЦЪТЕСИЭШЯЭХЗАСЖЗТИА ЛЭЙЖСЦОЗИЗЛЖЖЭЗИОЛЪБЕЛАСЪЖЗШСШЛМЛОЛЪЖЛЫЭЖЖЛМЛГЗФЯС"
	newstr := make([]rune, len(str))
	sum := 0.
	for i, ch := range str {
		newstr[i] = ch
	}
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
		fmt.Printf("%c\t%d\t%g\n", k, alf[k], float64(alf[k])/sum)
	}

	for i, v := range str {
		if v == 'И' {
			newstr[i] = 'о'
		} else if v == 'Ф' {
			newstr[i] = 'е'
		} else if v == 'К' {
			newstr[i] = 'э'
		} else if v == 'Щ' {
			newstr[i] = 'т'
		}
		fmt.Printf("%c", newstr[i])
	}

}
