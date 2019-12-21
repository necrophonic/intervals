package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/go-yaml/yaml"
	"github.com/schollz/progressbar"
	flag "github.com/spf13/pflag"
)

type (
	configYaml struct {
		Sounds struct {
			Beep string
		}
	}
)

var config configYaml

func init() {
	r, err := os.Open("config.yml")
	if err != nil {
		fmt.Println("Failed to open config file:", err)
		os.Exit(1)
	}

	d := yaml.NewDecoder(r)

	if err := d.Decode(&config); err != nil {
		fmt.Println("Failed to decode config:", err)
		os.Exit(1)
	}
}

func main() {

	var sets int
	flag.IntVarP(&sets, "sets", "s", 1, "the number of times to repeat the interval set")

	var rest time.Duration
	flag.DurationVarP(&rest, "rest", "r", 0, "amount of seconds rest between sets")

	var delay time.Duration
	flag.DurationVarP(&delay, "delay", "d", 0, "amount of seconds delay before the first set starts")

	var intervals []time.Duration
	flag.DurationSliceVarP(&intervals, "intervals", "i", nil, "(optional) number of sets")

	flag.Parse()

	if intervals == nil {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Println("Intervals: ", intervals)

	f, err := os.Open(config.Sounds.Beep)
	if err != nil {
		log.Fatal(err)
	}

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	var wg sync.WaitGroup

	// Initial delay (if any)
	if delay > 0 {
		fmt.Printf("Delaying start for %v ...\n", delay)
		time.Sleep(delay)
	}
	soundBeep(&wg, streamer)

	for set := 1; set <= sets; set++ {
		fmt.Printf("Begin set %d\n", set)
		
		for _, interval := range intervals {
			fmt.Printf("\tStarting interval: %v", interval)

			ticks := int(interval.Seconds())

			bar := progressbar.New(ticks)
			for s := 0 ; s < ticks ; s++ {
				time.Sleep(1 * time.Second)
				bar.Add(1)
			}
			// time.Sleep(interval)
			fmt.Println("BEEP!")
			soundBeep(&wg, streamer)
		}
		// Rest
		if set != sets {
			if rest > 0 {
				fmt.Printf("[REST] %v\n", rest)
				time.Sleep(rest)
				soundBeep(&wg, streamer)
			}
		}
	}

	wg.Wait()
	fmt.Println("All done!")
}

func soundBeep(wg *sync.WaitGroup, streamer beep.StreamSeekCloser) {
	wg.Add(1)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		streamer.Seek(0)
		wg.Done()
	})))
}
