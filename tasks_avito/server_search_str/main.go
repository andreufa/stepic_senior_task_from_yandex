package main

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Есть метод, который осуществляет поиск на сервер заданной строки

// Реализовать метод опроса нескольких серверов.
// Условия:
// 1. Сервера должны опрашиваться одновременно.
// 2. Мы ждем первого ответа от любого сервера.
//    Как только ответ получен метод должен его вернуть.
//    Остальные ответы можно игнорировать.
// 3. Метод возвращает ошибку только в том сулчае, если все серверы ответили ошибкой.

func AskAllServers(servers []string, query string) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var wg sync.WaitGroup

	res := make(chan string)

	for _, srv := range servers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			select {
			case <-ctx.Done():
				return
			default:
				resp, err := ServerSearch(srv, query)
				if err != nil {
					return
				}
				select {
				case res <- resp[0]:
					cancel() // отменяем остальные
				case <-ctx.Done():
					// кто-то уже отправил результат
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	v, ok := <-res
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil

}

func ServerSearch(server string, query string) ([]string, error) {
	if server == query {
		return []string{server}, nil
	}
	return nil, errors.New("not found in this server")
}

func main() {
	servers := []string{"ffff1", "ssss", "ddddd"}
	res, err := AskAllServers(servers, "ffff")
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
