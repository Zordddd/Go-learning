package main

//func Visioner() {
//	message := []byte("KBYGMSBYAELUGNAWGNGHEYHPHLBNJAEXJCMGZBSBGILNUPRWGHVWHIAYSQGGTXGSAJKHUGNTTHMRTVCFWXFEHQAUCHVLLALLSRGJSZSIIFASYYFUHWLSGAIVRKGNUEOFXHHBQFJZBCUAEWOGSAETDRMSPBQYRPTCGECPVMSGZXSVWISSVVREHUAMEMAVVNLBSGAIYFWOYVPAEOIWZGHUGWIGAMRXVVNLBSEBKLGKGALBTDNMXRTWTZMSNBHXUQFVF")
//	keyLength := 4
//	groups := make([][]byte, keyLength)
//
//	for i := 0; i < keyLength; i++ {
//		groups[i] = make([]byte, 0, len(message)/keyLength+keyLength)
//	}
//
//	for i, ch := range message {
//		groups[i%keyLength] = append(groups[i%keyLength], ch)
//	}
//
//	// Правильное дешифрование
//	keys := []int{13, 9, 12, 11} // Пробуем ваш ключ
//
//	for groupIdx := 0; groupIdx < keyLength; groupIdx++ {
//		for i := range groups[groupIdx] {
//			// Дешифрование: (cipher - key) mod 26
//			decrypted := (int(groups[groupIdx][i]) - 'A' - keys[groupIdx]) % 26
//			if decrypted < 0 {
//				decrypted += 26
//			}
//			groups[groupIdx][i] = byte(decrypted + 'A')
//		}
//	}
//
//	// Собираем результат
//	result := make([]byte, 0, len(message))
//	for i := 0; i < len(groups[0]); i++ {
//		for groupIdx := 0; groupIdx < keyLength; groupIdx++ {
//			if i < len(groups[groupIdx]) {
//				result = append(result, groups[groupIdx][i])
//			}
//		}
//	}
//
//	fmt.Printf("Расшифрованный текст: %s\n", string(result))
//}

//func Crypt(str string) {
//	alf := make(map[rune]int)
//	srAlf := make([]rune, 0, 33)
//	//str := "ОЛХВЗХЭВЧЯ ЭВЧИТХЗЪЭВВЛПУЖСЕСЦАЦЯЭАЖЭЭМЗОЭВИШЛПМЛЯЛЦЭПЭЖЭВДУФУЖСРЭЯЭМУЖ ЗЪСЛЦЗЖЛОЧВЖЧЙОЗИЭЫЖСЯЗИЛАСЪЗЭЯЛМЪЗФЧЯСИИШСЕСАГЗЭЗИВЛЯЗКЩЗЕЖЗЭМЛМЛИОЛЦЗЖСИЦЭЪСАЮВЛЛЖИВСЪЯЛЦЛЖСХСЪБЖЗШЛПЦЛШУПЭЖВСЪБЖЛЕСФЗШИЗЯЛАСЖЖЛЙЗИВЛЯЗЗШЯЗОВЛМЯСФЗЗЭМЛИЗИВЭПСЖЭТАЪТЭВИТВСЙЖЛОЗИБКАВЛПАЗЦЭАШСШЛПЛЖСЗЕАЭИВЖСИЛАЯЭПЭЖЖЛПУПЗЯУЦЪТЕСИЭШЯЭХЗАСЖЗТИА ЛЭЙЖСЦОЗИЗЛЖЖЭЗИОЛЪБЕЛАСЪЖЗШСШЛМЛОЛЪЖЛЫЭЖЖЛМЛГЗФЯС"
//	//newstr := make([]rune, len(str))
//	sum := 0.
//
//	for _, c := range str {
//		alf[c]++
//		sum++
//	}
//	for k := range alf {
//		srAlf = append(srAlf, k)
//	}
//	sort.Slice(srAlf, func(i int, j int) bool {
//		return alf[srAlf[i]] < alf[srAlf[j]]
//	})
//
//	for _, k := range srAlf {
//		fmt.Printf("%c\t%d\t%g\n", k, alf[k], float64(alf[k])/sum)
//	}
//
//	//for i, v := range str {
//	//    //if v == 'Л' {
//	//    //    newstr[i] = 'о'
//	//    //}else if v == 'З' {
//	//    //    newstr[i] = 'е'
//	//    //} else if v == 'Ж' {
//	//    //    newstr[i] = 'и'
//	//    //} else if v == 'Э' {
//	//    //    newstr[i] = 'н'
//	//    //}else if v == 'С' {
//	//    //    newstr[i] = 'а'
//	//    //}else if v == 'И' {
//	//    //    newstr[i] = 'т'
//	//    //}else if v == 'В' {
//	//    //    newstr[i] = 'с'
//	//    //}else if v == 'H' {
//	//    //    newstr[i] = 'l'
//	//    //}else if v == 'T' {
//	//    //    newstr[i] = 'x'
//	//    //}else if v == 'B' {
//	//    //    newstr[i] = 'm'
//	//    //}else if v == 'L' {
//	//    //    newstr[i] = 'g'
//	//    //}else if v == 'U' {
//	//    //    newstr[i] = 'y'
//	//    //}else if v == 'D' {
//	//    //    newstr[i] = 'c'
//	//    //}else if v == 'Y' {
//	//    //    newstr[i] = 's'
//	//    //}else if v == 'E' {
//	//    //    newstr[i] = 'i'
//	//    //}else if v == 'M' {
//	//    //    newstr[i] = 'p'
//	//    //}else if v == 'W' {
//	//    //    newstr[i] = 'n'
//	//    //}else if v == 'F' {
//	//    //    newstr[i] = 'w'
//	//    //}else if v == 'A' {
//	//    //    newstr[i] = 'f'
//	//    //}else if v == 'O' {
//	//    //    newstr[i] = 'r'
//	//    //}else if v == 'Z' {
//	//    //    newstr[i] = 'h'
//	//    //}else if v == 'G' {
//	//    //    newstr[i] = 'k'
//	//    //}else if v == 'V' {
//	//    //    newstr[i] = 'v'
//	//    //}else if v == 'R' {
//	//    //    newstr[i] = 'z'
//	//    //}else {
//	//    //    newstr[i] = v
//	//    //}
//	//    fmt.Printf("%c", newstr[i])
//	//}
//
//}
