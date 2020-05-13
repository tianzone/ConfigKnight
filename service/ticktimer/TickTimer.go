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

	"time"
)

// The name type of timer
type Timer struct{
	// The vaule of countting
	countVal int

	// The channel which as a switch to control start/stop the timer.
	enableSwt chan bool
	enable bool

	// The cannel which transport the new interval setting value
	newInterval chan int
	interval int
}

// The web request handler
func (t Timer) ServeHTTP( w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()

	// Print the request body
	if buf, err := ioutil.ReadAll( r.Body ); err != nil{
		fmt.Printf( "Real request's body met error: %s\n", err )
	}else{
		fmt.Println( buf )
	}

	// Respond the req
	fmt.Fprint( w, "I have received your messages..." )
}

func ( t *Timer ) Run() error{
	// Satr up the countting goroutine
	go t.counttingRoutine()

	// Loop to listen the channels
	for{
		select{
		case enable := <-t.enableSwt:{
			fmt.Printf( "%t timer....\n", enable )
			t.enable = enable
		}
		case interval := <-t.newInterval:{
			fmt.Printf( "Set the interval from %d to %d\n", t.interval, interval )
			t.interval = interval
		}
		}
	}
}

func (t *Timer) counttingRoutine() error{
	fmt.Printf( "Countting routine startup...\n" )
	for{
		if t.enable == true{
			t.countVal += t.interval
		}

		// Sleep some seconds that specify in interval
		time.Sleep( time.Duration( t.interval ) * time.Millisecond )
	}
}