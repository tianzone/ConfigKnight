/*
*
* ticktimer is a timer tool that user can set the counting interval via on web client, and
* also show up the countting time on the web UI. 
*
*/

package ticktimer

import(
	"fmt"

	"net/http"
	"io/ioutil"
	"sync"

	"time"
	"encoding/json"

	"log"

	"bytes"
	"strconv"
)

// TODO: Implement the data struct of respond for web request.
type Response struct{
	Code int `json: "code"`
	Data map[string]int `json: "data"`
}

// TODO: Complete the data struct of request's body
type RequestBody struct{
	Action string
	Params map[string]string
}
// ---------------------------------------------------------------------------

// The name type of timer
type Timer struct{
	// The vaule of countting
	countVal int

	// The channel which as a switch to control start/stop the timer.
	enableSwt chan bool
	enable bool
	enableMutex sync.Mutex

	// The cannel which transport the new interval setting value
	newInterval chan int
	interval int
}

// The web request handler
func (t *Timer) ServeHTTP( w http.ResponseWriter, r *http.Request){
	var action string
	var params = make( map[string]string )
	switch r.Method{
		case "GET":{
			data := r.URL.Query()
			action = data.Get( "action" )

			// Debugging
			fmt.Println( "GET request => action: ", action, " data: ", data )
		}
		case "POST":{
			reqBody := RequestBody{}
			// if n, err := r.Body.Read( buf ); err != nil{
			if buf, err := ioutil.ReadAll( r.Body ); err != nil{
				log.Fatal( "Body.Read: ", err )
			}else if err = json.Unmarshal( buf, &reqBody ); err != nil{
				log.Fatal( "Unmarshal: ", err, ", Body: ", bytes.NewBuffer( buf ).String() )
			}
			action = reqBody.Action
			params = reqBody.Params
			// Debugging
			fmt.Println( "POST request => action: ", action, " params: ", params )
		}
	}
	defer r.Body.Close()
	fmt.Println( "New request => action: ", action, " params: ", params )

	// TODO: Manage the handlers into a map.
	var respString []byte
	switch action{
		case "getInfos":{
			respString = t.ReqHandleGetInfos( params )
		}
		case "applySetting":{
			respString = t.ReqHandleSetting( params )
		}
		default:{
			respString = bytes.NewBufferString( "{ 'code': 2; 'msg': 'No valid handler' }" ).Bytes()
		}
	}

	// Set the content type of respond
	w.Header().Set( "content-type", "application/json" )
	w.Write( respString )
}

// Handling of request the whole info. 
func ( t *Timer ) ReqHandleGetInfos( params map[string]string ) []byte{
	// Respond the req
	data := make( map[string]int )
	data["countVal"] = t.countVal
	data["interval"] = t.interval
	if t.enable{data["enable"] = 1 }else{ data["enable"] = 0 }
	response := Response{
		Code: 0,
		Data: data,
	}
	respondStr, err := json.Marshal( response )
	if err != nil{
		return bytes.NewBufferString( "{ 'code': 1 }" ).Bytes()
	}
	return respondStr
}

// Handle the setting
func ( t *Timer ) ReqHandleSetting( params map[string]string ) []byte{
	// TODO: Handle the errors
	if val, err := strconv.Atoi( params["interval"] ); err != nil{
		log.Fatal( err )
	}else{
		t.newInterval <- val
	}

	if params["enable"] == "0"{
		t.enableSwt <- false
	}else{
		t.enableSwt <- true
	}

	return bytes.NewBufferString( "{ 'code': 0 }" ).Bytes()
}

func ( t *Timer ) Run() error{
	// TODO: Can we just initialize the data of the name struct right here or are there any ways
	// to implement it?
	t.interval = 3
	t.countVal = 0
	t.enable = false

	// Start the processing of the timer
	go func(){
		// Satr up the countting goroutine
		go t.counttingRoutine()
	
		// Loop to listen the channels
		for{
			select{
			case enable := <-t.enableSwt:{
				fmt.Printf( "%t timer....\n", enable )
				t.enableMutex.Unlock()
				t.enable = enable
			}
			case interval := <-t.newInterval:{
				fmt.Printf( "Set the interval from %d to %d\n", t.interval, interval )
				t.interval = interval
			}
			}
		}
	}()

	return nil
}

func (t *Timer) counttingRoutine() error{
	fmt.Printf( "Countting routine startup...\n" )
	for{
		if t.enable == true{
			t.countVal += t.interval
		}else{
			// Block until the timer is enabled.
			fmt.Printf( "Timer was stoped...\n" )
			t.enableMutex.Lock()
			fmt.Printf( "Timer was started...\n" )
		}

		// Sleep some seconds that specify in interval
		time.Sleep( time.Duration( t.interval ) * time.Second )
		fmt.Printf( "Current count: --%d--\n", t.countVal )
	}
}