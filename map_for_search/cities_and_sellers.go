package main

// У нас есть объект [Продавец ID] -> [Список городов, где он осуществляет услуги]
// Необходимо по запрошенным городам вернуть такой же объект только с продавцами, у которых есть желаемые населенные пункты,
// лишнее надо откинуть

// Пример

// Пример1
// ### in
// sellers = {
// 1: ['Москва', 'Самара', 'Ростов'],
// 2: ['Москва', 'Самара', 'Ростов', 'Казань', 'Курган', 'Пенза'],
// 3: ['Самара', 'Ростов', 'Курган', 'Пенза'],
// 4: ['Москва', 'Казань', 'Тула'],
// }

// cities = ['Москва', 'Казань', 'Тула']
// ### out
// {
// 1: [Москва],
// 2: [Москва, Казань],
// 4: [Москва, Казань, Тула],
// }


func MatchSellers(sellers map[int][]string, cities [] string) map[int][]string{
	result := make(map[int][]string)

	mSities := make(map[string]bool)

	for _, city := range cities{
		mSities[city] = true
	}

	for id, sellerSities := range sellers{
		var mathces []string
		for _, s := range sellerSities{
			if mSities[s]{
				mathces = append(mathces, s)
			}
		}
		if len(mathces) > 0{
			result[id] = mathces
		}
	}

	return  result
}