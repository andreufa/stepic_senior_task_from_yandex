package main

type Book struct {
    ID     int
    Title  string
    Author string
    Genre  string
    Year   int
}

type BookIndex struct {
    ByAuthor map[string][]Book  // Индекс по автору
    ByGenre  map[string][]Book  // Индекс по жанру
    ByYear   map[int][]Book      // Индекс по году
    Books    []Book              // Все книги
}

func BuildIndex(books []Book) *BookIndex {
    index := &BookIndex{
        ByAuthor: make(map[string][]Book),
        ByGenre:  make(map[string][]Book),
        ByYear:   make(map[int][]Book),
        Books:    books,
    }
    
    // Строим все индексы за один проход
    for _, book := range books {
        index.ByAuthor[book.Author] = append(index.ByAuthor[book.Author], book)
        index.ByGenre[book.Genre] = append(index.ByGenre[book.Genre], book)
        index.ByYear[book.Year] = append(index.ByYear[book.Year], book)
    }
    
    return index
}

// Поиск теперь мгновенный!
func (idx *BookIndex) FindByAuthor(author string) []Book {
    return idx.ByAuthor[author] // O(1) - даже с миллионом книг!
}

func (idx *BookIndex) FindByGenre(genre string) []Book {
    return idx.ByGenre[genre] // O(1)
}

func (idx *BookIndex) FindByYear(year int) []Book {
    return idx.ByYear[year] // O(1)
}