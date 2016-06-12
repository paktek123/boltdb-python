// Copyright 2015 The go-python Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// package hi exposes a few Go functions to be wrapped and used from Python.
package hi

import (
	"fmt"
  "log"
//  "os"
	"github.com/go-python/gopy/_examples/cpkg"
  "github.com/boltdb/bolt"
)

const (
	Version  = "0.1" // Version of this package
	Universe = 42    // Universe is the fundamental constant of everything
)

var (
	Debug    = false                            // Debug switches between debug and prod
	Anon     = Person{Age: 1, Name: "<nobody>"} // Anon is a default anonymous person
	IntSlice = []int{1, 2}                      // A slice of ints
	IntArray = [2]int{1, 2}                     // An array of ints
  world = []byte("world")
  key = []byte("hello")
  _sample = "hello"
  _boltholder = make(map[string]*bolt.DB)
  keyholder = make(map[string][]byte)
)

func Stuff() {
  db, err := bolt.Open("/tmp/bolt.db", 0644, nil)
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()
  err = db.View(func(tx *bolt.Tx) error {
    bucket := tx.Bucket(world)
    if bucket == nil {
      return fmt.Errorf("Bucket %q not found!", world)
    }
    val := bucket.Get(key)
    fmt.Println(string(val))
    return nil
  })
}

// Hi prints hi from Go
func Hi() {
	cpkg.Hi()
}

// Hello prints a greeting from Go
func Hello(s string) {
	cpkg.Hello(s)
}

// Concat concatenates two strings together and returns the resulting string.
func Concat(s1, s2 string) string {
	return s1 + s2
}

// LookupQuestion returns question for given answer.
func LookupQuestion(n int) (string, error) {
	if n == 42 {
		return "Life, the Universe and Everything", nil
	} else {
		return "", fmt.Errorf("Wrong answer: %v != 42", n)
	}
}

// Add returns the sum of its arguments.
func Add(i, j int) int {
	return i + j
}



type Boltdb struct {
  Path string
  Mode string
  Connection string
}

func (bo Boltdb) CreateBucket(name string) BoltBucket {
  db := _boltholder["database"]
  db.Update(func(tx *bolt.Tx) error {
    _, err := tx.CreateBucket([]byte(name))
    if err != nil {
        return fmt.Errorf("create bucket: %s", err)
    }
    return nil
  })  
  return BoltBucket{
    Name: name,
  }
}  

type BoltBucket struct {
  Name string
  //Bucket bolt.BucketStats
}

func (bucket BoltBucket) Put(key string, value string) map[string]string {
  db := _boltholder["database"]
  db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket.Name))
    err := b.Put([]byte(key), []byte(value))
    return err
  })
  key_string := fmt.Sprintf("%v",key)
  value_string := fmt.Sprintf("%v", value)
  m := map[string]string{key_string: value_string}
  return m
}

func (bucket BoltBucket) Get(key string) string {
  db := _boltholder["database"]
  byte_key := []byte(key)
  db.View(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket.Name))
    keyholder[key] = b.Get(byte_key)
    return nil
  })
  return string(keyholder[key])
} 

func (bucket BoltBucket) Delete(key string) error {
  db := _boltholder["database"]
  byte_key := []byte(key)
  db.Update(func(tx *bolt.Tx) error {
    b := tx.Bucket([]byte(bucket.Name))
    b.Delete(byte_key)
    return nil
  })
  return nil
}

func NewBoltdb(path string) Boltdb {
  //real_mode := int(mode)
  db, err := bolt.Open(path, 0600, nil)
  if err != nil {
    panic(err)
  }
  _boltholder["database"] = db
  return Boltdb{
    Path: path,
    Mode: "0600",
    Connection: path,
  }
}

// Person is a simple struct
type Person struct {
	Name string
	Age  int
}

// NewPerson creates a new Person value
func NewPerson(name string, age int) Person {
	return Person{
		Name: name,
		Age:  age,
	}
}

// NewPersonWithAge creates a new Person with a specific age
func NewPersonWithAge(age int) Person {
	return Person{
		Name: "stranger",
		Age:  age,
	}
}

// NewActivePerson creates a new Person with a certain amount of work done.
func NewActivePerson(h int) (Person, error) {
	var p Person
	err := p.Work(h)
	return p, err
}

func (p Person) String() string {
	return fmt.Sprintf("hi.Person{Name=%q, Age=%d}", p.Name, p.Age)
}

// Greet sends greetings
func (p *Person) Greet() string {
	return p.greet()
}

// greet sends greetings
func (p *Person) greet() string {
	return fmt.Sprintf("Hello, I am %s", p.Name)
}

// Work makes a Person go to work for h hours
func (p *Person) Work(h int) error {
	cpkg.Printf("working...\n")
	if h > 7 {
		return fmt.Errorf("can't work for %d hours!", h)
	}
	cpkg.Printf("worked for %d hours\n", h)
	return nil
}

// Salary returns the expected gains after h hours of work
func (p *Person) Salary(h int) (int, error) {
	if h > 7 {
		return 0, fmt.Errorf("can't work for %d hours!", h)
	}
	return h * 10, nil
}

// Couple is a pair of persons
type Couple struct {
	P1 Person
	P2 Person
}

// NewCouple returns a new couple made of the p1 and p2 persons.
func NewCouple(p1, p2 Person) Couple {
	return Couple{
		P1: p1,
		P2: p2,
	}
}

func (c *Couple) String() string {
	return fmt.Sprintf("hi.Couple{P1=%v, P2=%v}", c.P1, c.P2)
}

// Float is a kind of float32
type Float float32

// Floats is a slice of floats
type Floats []Float

// Eval evals float64
type Eval func(f float64) float64
