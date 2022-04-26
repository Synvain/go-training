package gotraining

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"sync"
	"time"

	mypackage "example.com/go-training/internal/mypackage"
	util "example.com/go-training/pkg/decorators"
	log "github.com/sirupsen/logrus"
)

/*******************************************************************************
	Variable default values : 0, nil, ""
********************************************************************************/
func UsingVarDefautlValues() {
	var a int // default value: 0
	var b string // default value: ""
	var c *struct{} // default value: nil

	fmt.Printf("a=%d, b='%s', c=%v\n", a, b, c)
}

/*******************************************************************************
	Functions parameters and return values
********************************************************************************/
// Public function (accessible from outside the package: `gotraining.Sum()`)
func Sum(a int, b int) int {
	return a + b
}

// Private function
func diff(a int, b int) (int) {
	return a - b
}

// Returning errors (always last return argument)
func div(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

// Using named return values
func div2(a int, b int) (result int, err error) {
	if b == 0 {
		err = fmt.Errorf("cannot divide %d by %d", a, b)
		return
	}
	result = a / b
	return
}

// Handling errors
func div_then_add(a int, b int) (result int, err error) {
	if result, err = div2(a, b); err != nil {
		return
	}
	result = Sum(result, b)
	return
}

// Scope of variables
// This function doesn't work as expected as it never returns errors
// The `err` within the if shadows the one in the function declaration
func div_then_diff2(a int, b int) (result int, err error) {
	if quotient, err := div(a, b); err == nil {
		result = diff(quotient, b)
	}
	return
}

type MyObject struct {
	a int
	b int
}

func MethodUsingACopyOfMyObject(o MyObject) {
	o.a = o.a * 2
	o.b = o.b * 2
}

func MethodUsingMyObject(o *MyObject) {
	o.a = o.a * 2
	o.b = o.b * 2
}

// = vs :=
func UsingFunctions() {
	var result1 int
	// := must have at least one new variable on the left side
	// with := the type of the variable is inferred from the right side

	result1, _ = div_then_add(10, 5) // not interested in the error return value
	fmt.Printf("result: %d\n", result1)

	result2, err := div_then_add(10, 0)
	fmt.Printf("result: %d, err: '%s'\n", result2, err.Error())

	// = can be used to assign a value to an existing variable
	var result3 int
	result3, _ = div_then_add(10, 5)
	fmt.Printf("result: %d\n", result3)
	result3, _ = div_then_diff2(10, 0)
	fmt.Printf("result: %d, err: '%s'\n", result3, err.Error())

	o := MyObject{a: 3, b: 5}
	MethodUsingACopyOfMyObject(o)
	fmt.Printf("After MethodUsingACopyOfMyObject -> o.a: %d, o.b: %d\n", o.a, o.b)
	MethodUsingMyObject(&o)
	fmt.Printf("After MethodUsingMyObject        -> o.a: %d, o.b: %d\n", o.a, o.b)
}

/*******************************************************************************
	Packages
********************************************************************************/
func UsingPackages() {
	mypackage.MyFunc()
	mypackage.MyOtherFunc()
}

/*******************************************************************************
	Custom types
********************************************************************************/
type day int

// Enums
const (
	// Enum
	Monday		day = iota
	Tuesday		day = iota
	Wednesday	day = iota
	Thursday	day = iota
	Friday		day = iota
	Saturday	day = iota
	Sunday		day = iota
)

// Adding methods to any types
func (i day) String() string {
	// using a switch, but we could've used a map
	switch i {
	case Monday:
		return "Monday"
	case Tuesday:
		return "Tuesday"
	case Wednesday:
		return "Wednesday"
	case Thursday:
		return "Thursday"
	case Friday:
		return "Friday"
	case Saturday:
		return "Saturday"
	case Sunday:
		return "Sunday"
	default: // Used if no other case is selected
		return "Unknown"
	}
}

func UsingCusomTypes() {
	// Using custom types
	var d day = Friday
	fmt.Printf("%s, %s\n", d, Saturday)
}

/*********************************************************************************
	Interfaces/Reflection
	- An interface is a `contract` that defines what methods a type must implement
**********************************************************************************/
type MyInterface interface {
	Method1()
	Method2() error
}

type MyClass struct {
	// Context variables
	a int // private variable to the class
	b int
	c int
	D int // public variable to the user of the class
}

func (c MyClass) Method1() {
	fmt.Println("MyClass Method1")
}

func (c MyClass) Method2() error {
	fmt.Println("MyClass Method2")
	c.privateMethod()
	return nil
}

func (c MyClass) Method3() {
	fmt.Println("MyClass Method3")
}

func (c MyClass) privateMethod() {
	fmt.Println("MyClass privateMethod")
}

type MyOtherClass struct {
	d int
	e int
}

func (c MyOtherClass) Method1() {
	fmt.Println("MyOtherClass Method1", c)
}

func (c MyOtherClass) Method2() error {
	fmt.Println("MyOtherClass Method2", c)
	return nil
}

func (c MyOtherClass) Method99() error {
	fmt.Println("MyOtherClass Method99")
	return nil
}


func MethodUsingMyInterface(i MyInterface) {
	i.Method1()
	i.Method2()
}

// Anything can be passed as argument to this function
// Watch out, go will panic if something not supported is attempted on the pas object/pointer
func MethodUsingAnyInterface(i interface{}) {
	fmt.Println("Type: ", reflect.TypeOf(i))

	// using type assertions
	if c, ok := i.(*MyClass); ok {
		fmt.Println("i is a *MyClass")
		i.(*MyClass).Method1()
		c.Method1()
	}

	if c, ok := i.(MyClass); ok {
		fmt.Println("i is a MyClass")
		i.(MyClass).Method1()
		c.Method2()
	}

	if _, ok := i.(*MyOtherClass); ok {
		fmt.Println("i is a *MyOtherClass")
		i.(*MyOtherClass).Method1()
	}

	if _, ok := i.(int); ok {
		fmt.Println("i is a int", i)
	}

	// or using switch statement
	switch c := i.(type) {
	case *MyClass:
		c.Method1()
	default:
		fmt.Println("Unknown: ", reflect.TypeOf(i))
		fmt.Printf("Unknown: '%T'\n", i)
	}
}

// Convention for allocating classes New<myclassname>
func NewMyClass(a int, b int) MyInterface {
	// Allocate memory for a new MyClass
	// As long as the pointer is no longer referenced anywhere, the memory is freed by the GC
	return &MyClass{
		a: a,
		b: b,
		c: a + b,
		D: a * b,
	}
}

func NewMyOtherClass(d int) MyInterface {
	// Allocate memory for a new MyOtherClass
	return &MyOtherClass{
		d: d,
		e: d * 2,
	}
}

func UsingClassesAndInterfaces() {
	// Using interfaces
	i := NewMyClass(1, 2)
	MethodUsingMyInterface(i)
	j := NewMyOtherClass(2)
	MethodUsingMyInterface(j)


	MethodUsingAnyInterface(i)
	MethodUsingAnyInterface(j)

	m := MyClass{}
	MethodUsingAnyInterface(m)

	n := 1
	MethodUsingAnyInterface(n)

}

/*******************************************************************************
	Memory allocation
********************************************************************************/
func UsingMemoryAllocation() {
	// Allocate memory for MyClass
	j := MyClass{}
	j.Method1()

	// Allocate memory for a slice with a capacity of 1
	k := make([]MyClass, 1)
	k[0].Method1()

	// Allocate memory for a map
	l := make(map[int]int) // same as map[int]int{}
	l[1] = 2
	fmt.Println(l)

	// literals can be used to allocate objects
	m := []MyClass{ {a: 1, D: 3} }
	n := map[string]string{ "a": "b" }

	// struct{} is a struct with no fields, somtime used with un-buffered channels for synchronization
	var o *struct{}  // Pointer to a struct with no members
	o = &struct{}{}  // Allocate memory for the struct

	fmt.Println(m, n, o)
}

/*******************************************************************************
	Maps
********************************************************************************/
type MyComplexMap map[string]map[string]string

// Maps used to check if a value is in a list
var Provices = map[string]bool{
	"QC": true,
	"ON": true,
	"BC": true,
	"AB": true,
	"SK": true,
	"MB": true,
	"NL": true,
	"NB": true,
	"NS": true,
	"PE": true,
	"YT": true,
	"NT": true,
	"NU": true,
}

func UsingMaps() {
	m := MyComplexMap{ "QC": { "a": "b" } }
	fmt.Println(m["QC"])

	// e will be the value, ok will be true if found
	e, ok := Provices["EE"]
	fmt.Println(e, ok)

	// Checking if a value is in a map (map's sugar syntax)
	if Provices["QC"] {
		fmt.Println("QC is a province")
	}

	if !Provices["AA"] {
		fmt.Println("AA is not a province")
	}

	m2 := make(map[string]string)
	m2["a"] = "b"
	fmt.Println("m2:['a']: ", m2["a"])
	delete(m2, "a")
	fmt.Println("m2['a']: ", m2["a"])

	m3 := map[string]string{"c": "d", "e": "f"}
	fmt.Println("m3: ", m3)
}

/*******************************************************************************
	Slices
********************************************************************************/
func UsingSlices() {
	// Creating a array
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Println(s)
	// capacity is the the array size
	// length is the number of elements in the array
	fmt.Println(len(s), cap(s))

	// Creating a slice with the array
	s = s[0:4]
	fmt.Println(s)
	fmt.Println(len(s), cap(s))

	// Back to the original array
	s = s[0:cap(s)]
	fmt.Println(s)
	fmt.Println(len(s), cap(s))

	// Appending to slices (will reallocate the array because the capacity is not enough)
	s = append(s, 11)
	fmt.Println(s)
	fmt.Println(len(s), cap(s))
}

/*******************************************************************************
	Ranges
********************************************************************************/
func UsingRanges() {
	// on slices
	for i, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10} {
		fmt.Println(i, v)
	}
	for i, v := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}[0:4] {
		fmt.Println(i, v)
	}
	// on maps: key value pairs
	for k, v := range Provices {
		fmt.Println(k, v)
	}
	// keys only
	for k := range Provices {
		fmt.Println(k)
	}
	// values only
	for _, v := range Provices {
		fmt.Println(v)
	}
}

/*******************************************************************************
	Time
********************************************************************************/
func UsingTime() {
	/***********************************************************************************
	 * Layout to parse/format time based on `Mon Jan 2 15:04:05 -0700 MST 2006`        *
	 * It is human readable!!!                                                         *
	 ***********************************************************************************/

	// Time manipulation
	t := time.Now()
	fmt.Println(t)
	fmt.Println(t.Add(1 * time.Hour + 30 * time.Minute))

	// Parsing time
	//                   layout                 string to parse
	t, err := time.Parse("2006-01-02 15:04:05 MST", "2022-04-19 08:22:38 EDT")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)

	// Format time
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	fmt.Println(t.Format("2006-01-02"))
	fmt.Println(t.Format("15:04:05"))
	fmt.Println(t.Format("15:04:05/2006/MST/Jan-2"))

	// Using time.Duration
	hour := time.Duration(1 * time.Hour)
	second := 1 * time.Second
	millisecond := 1 * time.Millisecond
	nanosecond := 1 * time.Nanosecond
	fmt.Printf("hour: %s, second: %s, millisecond: %s nanosecond: %s\n", hour, second, millisecond, nanosecond)

	// Parsing time.Duration
	d, err := time.ParseDuration("1h30m")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(d)

}

/*******************************************************************************
	Structure Metadata/JSON
********************************************************************************/
type MyStruct struct {
	A int `json:"a"`
	B int `json:"b"`
	c int // Private (lowercase) fields are not serialized to JSON
	D int `json:"c,omitempty"` // omitempty is used to ignore empty fields when serializing to JSON
	E []int `json:"e"`
}
func UsingStructMetadata() {
	out, _ := json.MarshalIndent(MyStruct{A: 1, B: 2, c: 3, E: []int{1, 2, 3}}, "", "  ")
	fmt.Println(string(out))
}

/*******************************************************************************
	Anonymous function
********************************************************************************/
func MyFunctionThatTakesAFunc(f func(int, int) int) int {
	return f(1, 2)
}

func UsingAnonymousFunctions() {
	res := MyFunctionThatTakesAFunc(func(a int, b int) int {
		return a + b
	})
	fmt.Println(res)
}

/*******************************************************************************
	Loops
********************************************************************************/
func UsingLoops() {
	j := 0
top_loop:
	for {
		j++
nested_loop:
		for i := 0; i < 100; i++ {
			if i >= 2 {
				fmt.Println("Breaking nested_loop")
				break nested_loop
			}
			fmt.Println("sleep 1 second")
			time.Sleep(1 * time.Second)
		}
		if j >= 2 {
			fmt.Println("Breaking top_loop")
			break top_loop
		}
	}
}

/*******************************************************************************
	Defer
********************************************************************************/
func OpenSomething() (int, error) {
	s := 1
	fmt.Println("OpenSomething: ", s)
	return s, nil
}
func CloseSomething(s int) error {
	fmt.Println("CloseSomething: ", s)
	return nil
}
func UsingDefer() {
	var s int
	var err error

	if s, err = OpenSomething(); err != nil {
		panic(err)
	}
	defer CloseSomething(s)
	// go on with your function without worrying about closing s
}

/*******************************************************************************
	Go routines
********************************************************************************/
func SomeRoutine(c chan bool) {
	fmt.Println("SomeRoutine: Hello from some go routine, I'm going to sleep for 3 seconds")
	time.Sleep(3 * time.Second)
	fmt.Println("SomeRoutine: I am done sleeping")
	c <- true
}
func UsingGoRoutines() {
	go func() {
		fmt.Println("Anonymous function: Hello from an anonymous go routine")
	}()

	c := make(chan bool) // Un-buffered channel: Sends and receives block until the other side is ready
						 // Used to synchronize goroutines without explicit locks or condition
	go SomeRoutine(c)
	fmt.Println("UsingGoRoutines: Waiting for some go routine to finish")
	<- c
	fmt.Println("UsingGoRoutines: Done")
}

/*******************************************************************************
	Queues/Channels/Select with timeout
********************************************************************************/
type MyItem struct {
	A string
	B int
}
func RoutineUsingQueues(done chan bool, terminate chan bool, q chan MyItem) {
loop:
	for {
		select {
		case <-terminate:
			fmt.Println("Terminating")
			break loop
		case item, ok := <- q:
			if !ok {
				fmt.Println("Channel closed")
				break loop
			}
			fmt.Println("Received item: ", item.A, item.B)
		case <-time.After(1 * time.Second):
			fmt.Println("Timeout")
			break loop
		}
	}
	done <- true
}
func UsingQueues() {
	done := make(chan bool) // Un-buffered channel: Sends and receives block until the other side is ready
	terminate := make(chan bool)
	q := make(chan MyItem, 3) // Buffered channel of capacity 3, can post 3 items until sender blocks
	q <- MyItem{A: "a", B: 1}
	q <- MyItem{A: "b", B: 2}
	q <- MyItem{A: "c", B: 3}

	go RoutineUsingQueues(done, terminate, q) // Will end because of timeout after emptying the queue
	<- done

	go RoutineUsingQueues(done, terminate, q) // Will end because of terminate
	terminate <- true
	<- done

	go RoutineUsingQueues(done, terminate, q) // Will end because of closed channel
	close(q)
	<- done
}

/*******************************************************************************
	range on channels
********************************************************************************/
func UsingRangeOnChannels() {
	c := make(chan string, 3) // Buffered channel of capacity 3
	c <- "a"
	c <- "b"
	c <- "c"
	close(c) // signal that no more items will be sent

	for i := range c {
		fmt.Println(i)
	}
}

/*******************************************************************************
	Class Embedding/Mutex
********************************************************************************/
type MyClassEmbedding struct {
	sync.Mutex // Embed the Mutex interface
	A int
	B int
}

func (c *MyClassEmbedding) What() string {
	return "MyClassEmbedding"
}

type MyClassEmbedding2 struct {
	MyClassEmbedding
	C int
}

func (c *MyClassEmbedding2) What() string {
	return "MyClassEmbedding2"
}

type MyClassEmbedding3 struct {
	*MyClassEmbedding2
	D int
}

func UsingClassEmbedding() {
	c1 := MyClassEmbedding{}
	c1.Lock()
	defer c1.Unlock()
	c1.A = 1
	c1.B = 2
	fmt.Printf("MyClassEmbedding: %v\n", c1)  // warning: call of fmt.Printf copies lock value
	fmt.Printf("MyClassEmbedding: %v\n", &c1) // Using pointer to c1 to prevent copying object

	// Initialize all fields specifically
	c2 := MyClassEmbedding2{MyClassEmbedding: MyClassEmbedding{A: 1, B: 2}, C: 3}
	c2.Lock()
	defer c2.Unlock()

	// Initialize a single field, others will be initialize their per type default
	c3 := MyClassEmbedding2{C:3}
	c3.A = 4

	// Embed pre-allocated base object
	c4 := MyClassEmbedding3{&c2, 4}

	fmt.Printf("c1: %v c2: %v, c3: %v, c4:%v\n", &c1, &c2, &c3, &c4)
	fmt.Printf("c1: %s c2: %s, c3: %s, c4:%s\n", c1.What(), c2.What(), c3.What(), c4.What())
}

/*******************************************************************************
	Logging
********************************************************************************/
type MyError struct{
	reason string
}
func (c MyError) Error() string {
	return c.reason
}

func initLogger() {
	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
		FieldMap: log.FieldMap{
			log.FieldKeyTime: "@timestamp",
			log.FieldKeyFunc: "@caller",
		},
	})
}
func UsingLogging() {
	f := func () error {
		return MyError{"Oopsy daisy like Antoine would say"}
	}

	if err := f(); err != nil {
		log.WithError(err).WithField("Antoine", "Imbeau").Error("Call to f() failed")
	}
}

/*******************************************************************************
	Wait for a condition for a maximum of time
********************************************************************************/
func WaitForSomething(something* bool, timeout time.Duration) error {
	for t := time.After(timeout); ; {
		select {
		case <-t:
			return fmt.Errorf("Timeout")
		default: // Run if no other case is ready
			if *something {
				fmt.Println("Something happened")
				return nil
			}
			fmt.Println("sleeping 1 second")
			time.Sleep(1 * time.Second)
		}
	}
}

func UsingSelectToWaitForSomething() {
	something := false
	err := WaitForSomething(&something, 1 * time.Second)
	if err != nil {
		log.WithError(err).Error("Error in UsingSelectToWaitForSomething()")
	}
	go func() {
		fmt.Println("sleeping 2 seconds")
		time.Sleep(2 * time.Second)
		something = true
	}()
	WaitForSomething(&something, 5 * time.Second)
}

/*******************************************************************************
	Context and Cancellation
********************************************************************************/
func WaitForSomethingElse(ctx context.Context, somethingElse* bool) {
	for {
	select {
		case <-ctx.Done():
			log.WithField("error", ctx.Err()).Warn("WaitForSomethingElse was canceled")
			return
		default: // Run if no other case is ready
			if *somethingElse {
				log.Warn("Something else happened")
				return
			}
			fmt.Println("sleeping 1 second")
			time.Sleep(1 * time.Second)
		}
	}
}

func UsingContextToCancel() {
	something := false
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		fmt.Println("sleeping 2 seconds")
		time.Sleep(2 * time.Second)
		cancel()
	}()
	WaitForSomethingElse(ctx, &something)
}

/*******************************************************************************
	Function Decorator
********************************************************************************/
func FunctionThatMayFail() error {
	return fmt.Errorf("oopsy daisy like Antoine would say")
}

func UsingDecorator() {
	err := util.Retry(3, time.Second, 2, func(attempt int) (retry bool, err error) {
		retry = true
		err = FunctionThatMayFail()
		fmt.Printf("Attempt %d: %v\n", attempt, err)
		return
	})
	if err != nil {
		fmt.Printf("It finally failed: %s\n", err.Error())
	}
}

/*******************************************************************************
	Http and Function/Function Decorating/Context with deadline
********************************************************************************/
type MyResponseWriter struct {
	http.ResponseWriter // Embed the ResponseWriter interface
	statusCode int
}

func (w *MyResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func MyHttpWriterDecorator(logger *log.Entry, handlerfunc http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		myw := &MyResponseWriter{w, 0}
		wrappedHandler := http.Handler(handlerfunc)
		wrappedHandler.ServeHTTP(myw, r)

		statusCode := myw.statusCode
		logger.WithFields(log.Fields{"Url": r.URL.Path, "statusCode": statusCode}).Info("MyHttpWriterDecorator")
		switch {
		case statusCode <= 199:
			// metrics.HttpResponseSent1XXCount.Inc(1)
		case statusCode <= 299:
			// metrics.HttpResponseSent2XXCount.Inc(1)
		case statusCode <= 399:
			// metrics.HttpResponseSent3XXCount.Inc(1)
		case statusCode <= 499:
			// metrics.HttpResponseSent4XXCount.Inc(1)
		case statusCode <= 599:
			// metrics.HttpResponseSent5XXCount.Inc(1)
		case statusCode <= 699:
			// metrics.HttpResponseSent6XXCount.Inc(1)
		default:
			// metrics.HttpResponseSentUnknownCount.Inc(1)
		}
	})
}

func StartMyHttpServer() {
	logger := log.WithField("logger", "http/server")

	// Undecorated handler
	http.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Fatal(err) // Would exit program with status 1
		}

		logger.WithField("body", string(body)).Info("Request received")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello http client"))
	})

	// Decorated to log the status code
	http.Handle("/foo", MyHttpWriterDecorator(logger, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logger.Fatal(err) // Would exit program with status 1
		}

		logger.WithField("body", string(body)).Info("Request received")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello http client"))
	}))

	// Decorated to log the status code
	http.Handle("/wait", MyHttpWriterDecorator(logger, func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		logger.Warn("sleeping 3 seconds")
		time.Sleep(3 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))

	go func() {
		http.ListenAndServe(":8080", nil)
	}()
}

func MySimpleHttpClientPost(logger *log.Entry) {
	resp, err := http.Post("http://localhost:8080/foo", "text", bytes.NewBuffer([]byte("Hello http server")))
	if err != nil {
		logger.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Fatal(err)
	}

	logger.WithFields(log.Fields{"body":string(body), "status": resp.Status}).Info("Response")
}

func MySimpleHttpClientPostWithDeadline(logger *log.Entry) {
	tr := &http.Transport{}
    client := &http.Client{Transport: tr}

	req, err := http.NewRequest("POST", "http://localhost:8080/wait", nil)
	if err != nil {
		logger.Fatal(err)
	}

	/* Note: http.Client has a `Timeout` parameter.
		But if we want multiple request to be completed within a certain time, we use context.WithDeadline
		Otherwise we would have to manage each request's timeout individually vs the the time budget for all the requests.
	*/
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(2 * time.Second))
	req = req.WithContext(ctx)
	defer cancel()

	_, err = client.Do(req)
	if err != nil {
		logger.Error(err)
	}
	// Wait for http request to complete on the server side (for demo purposes)
	time.Sleep(2 * time.Second)
}

func UsingHttpServer() {
	StartMyHttpServer()

	logger := log.WithField("logger", "http/client")
	MySimpleHttpClientPost(logger)
	MySimpleHttpClientPostWithDeadline(logger)
}

/*******************************************************************************
	RunCourses
********************************************************************************/
type Example struct {
	Name string
	Function func()
}

var Examples = []Example{
	{ Name: "UsingVarDefautlValues", Function: UsingVarDefautlValues },
	{ Name: "UsingFunctions", Function: UsingFunctions },
	{ Name: "UsingPackages", Function: UsingPackages },
	{ Name: "UsingCusomTypes", Function: UsingCusomTypes },
	{ Name: "UsingClassesAndInterfaces", Function: UsingClassesAndInterfaces },
	{ Name: "UsingMemoryAllocation", Function: UsingMemoryAllocation },
	{ Name: "UsingMaps", Function: UsingMaps },
	{ Name: "UsingSlices", Function: UsingSlices },
	{ Name: "UsingRanges", Function: UsingRanges },
	{ Name: "UsingTime", Function: UsingTime},
	{ Name: "UsingStructMetadata", Function: UsingStructMetadata },
	{ Name: "UsingAnonymousFunctions", Function: UsingAnonymousFunctions },
	{ Name: "UsingLoops", Function: UsingLoops },
	{ Name: "UsingDefer", Function: UsingDefer },
	{ Name: "UsingGoRoutines", Function: UsingGoRoutines },
	{ Name: "UsingQueues", Function: UsingQueues },
	{ Name: "UsingRangeOnChannels", Function: UsingRangeOnChannels },
	{ Name: "UsingClassEmbedding", Function: UsingClassEmbedding },
	{ Name: "UsingSelectToWaitForSomething", Function: UsingSelectToWaitForSomething },
	{ Name: "UsingContextToCancel", Function: UsingContextToCancel },
	{ Name: "UsingLogging", Function: UsingLogging },
	{ Name: "UsingDecorator", Function: UsingDecorator },
	{ Name: "UsingHttpServer", Function: UsingHttpServer },
}

func RunCourses() {
	initLogger()

	for _, example := range Examples {
		fmt.Printf("\n****************** %s ******************\n", example.Name)
		example.Function()
	}
}
