package models

import (
	"errors"
	"strings"
	"sync"
)

var (
	ErrNotFound       = errors.New("not found")
	ErrValidation     = errors.New("validation error")
	ErrDuplicateName  = errors.New("duplicate name")
	ErrRelationAbsent = errors.New("related entity does not exist")
)

type Store struct {
	mu sync.RWMutex

	books      map[int]Book
	authors    map[int]Author
	categories map[int]Category

	nextBookID     int
	nextAuthorID   int
	nextCategoryID int
}

func NewStore() *Store {
	return &Store{
		books:          make(map[int]Book),
		authors:        make(map[int]Author),
		categories:     make(map[int]Category),
		nextBookID:     1,
		nextAuthorID:   1,
		nextCategoryID: 1,
	}
}

func (s *Store) ListBooks(authorID, categoryID int) []Book {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Book, 0, len(s.books))
	for _, b := range s.books {
		if authorID > 0 && b.AuthorID != authorID {
			continue
		}
		if categoryID > 0 && b.CategoryID != categoryID {
			continue
		}
		res = append(res, b)
	}
	return res
}

func (s *Store) GetBook(id int) (Book, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	b, ok := s.books[id]
	if !ok {
		return Book{}, ErrNotFound
	}
	return b, nil
}

func (s *Store) CreateBook(in Book) (Book, error) {
	if strings.TrimSpace(in.Title) == "" || in.Price < 0 {
		return Book{}, ErrValidation
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.authors[in.AuthorID]; !ok {
		return Book{}, ErrRelationAbsent
	}
	if _, ok := s.categories[in.CategoryID]; !ok {
		return Book{}, ErrRelationAbsent
	}

	in.ID = s.nextBookID
	s.nextBookID++
	s.books[in.ID] = in
	return in, nil
}

func (s *Store) UpdateBook(id int, in Book) (Book, error) {
	if strings.TrimSpace(in.Title) == "" || in.Price < 0 {
		return Book{}, ErrValidation
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return Book{}, ErrNotFound
	}
	if _, ok := s.authors[in.AuthorID]; !ok {
		return Book{}, ErrRelationAbsent
	}
	if _, ok := s.categories[in.CategoryID]; !ok {
		return Book{}, ErrRelationAbsent
	}

	in.ID = id
	s.books[id] = in
	return in, nil
}

func (s *Store) DeleteBook(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.books[id]; !ok {
		return ErrNotFound
	}
	delete(s.books, id)
	return nil
}

func (s *Store) ListAuthors() []Author {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Author, 0, len(s.authors))
	for _, a := range s.authors {
		res = append(res, a)
	}
	return res
}

func (s *Store) CreateAuthor(name string) (Author, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Author{}, ErrValidation
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, a := range s.authors {
		if strings.EqualFold(a.Name, name) {
			return Author{}, ErrDuplicateName
		}
	}

	a := Author{ID: s.nextAuthorID, Name: name}
	s.nextAuthorID++
	s.authors[a.ID] = a
	return a, nil
}

func (s *Store) ListCategories() []Category {
	s.mu.RLock()
	defer s.mu.RUnlock()

	res := make([]Category, 0, len(s.categories))
	for _, c := range s.categories {
		res = append(res, c)
	}
	return res
}

func (s *Store) CreateCategory(name string) (Category, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return Category{}, ErrValidation
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.categories {
		if strings.EqualFold(c.Name, name) {
			return Category{}, ErrDuplicateName
		}
	}

	c := Category{ID: s.nextCategoryID, Name: name}
	s.nextCategoryID++
	s.categories[c.ID] = c
	return c, nil
}
