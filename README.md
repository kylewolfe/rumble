rumble [![GoDoc](http://godoc.org/github.com/kylewolfe/rumble?status.svg)](http://godoc.org/github.com/kylewolfe/rumble) 
======

RumbleDB is an abstraction for boltdb that aims to provide a clean API similiar to that of gopkg.in/mgo.v2 without hiding boltdb away from you completely and without locking you in to a
specefic encoding.

## Why?

The mgo API is awesome, and so is boltdb.

## RumbleDB Out of the Box

```go
var db *rumble.DB
var err error
   
if db, err = rumble.New("test.db"); err != nil {
	panic(err)
}
       
// structs
bucket := db.Bucket("foo")
for i := 0; i < 3; i++ {
	foo := &struct {
		Id    bson.ObjectId `rumble:"id"`
		Fizz  string
		Count int
	}{
		Fizz:  "buzz",
		Count: i,
	}
	if err = bucket.Put(foo); err != nil {
		panic(err)
	}
	fmt.Printf("newly created id: %s\n", foo.Id.Hex()) // ids generated on the fly like mgo
}

// iteration
i := bucket.NewIter()
foo := &struct {
	Id    bson.ObjectId `rumble:"id"`
	Fizz  string
	Count int
}{}
for i.Next(foo) {
	fmt.Printf("created: %s\n", foo.Id.Time())
}

// maps
bucket = db.Bucket("bar")
m := bson.M{"foo": "bar"}
if err = bucket.Put(m); err != nil {
	panic(err)
}
fmt.Println(m)

// newly created id: 56c4ffb89e56a73ced4227d6
// newly created id: 56c4ffb89e56a73ced4227d7
// newly created id: 56c4ffb89e56a73ced4227d8
// created: 2016-02-17 18:18:16 -0500 EST
// created: 2016-02-17 18:18:16 -0500 EST
// created: 2016-02-17 18:18:16 -0500 EST
// map[_id:[86 196 255 185 158 86 167 60 237 66 39 217] foo:bar]
```

## Bring Your Own Encoding

RumbleDB provides encoding functionality from bson out of the box, but you can use whatever you'd like.

```go
db, _ := rumble.New("my.db")
db.Marshal = json.Marshal
db.Unmarshal = json.Unmarshal
```

You can also use your own ID format

```go
var i uint32 = 0
db.NewId = func() []byte {
	return []byte(i := atomic.AddUint32(&rowCounter, 1))
}
```
