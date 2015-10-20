// internal/github.com/twinj/uuid/saver_test.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package uuid

/****************
 * Date: 21/06/15
 * Time: 6:46 PM
 ***************/

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

const (
	saveDuration time.Duration = 3
)

var (
	config = StateSaverConfig{SaveReport: true, SaveSchedule: saveDuration * time.Second}
)

func init() {
	SetupFileSystemStateSaver(config)
}

// Tests that the schedule is run on the timeDuration
func TestUUID_State_saveSchedule(t *testing.T) {

	if state.saver != nil {
		count := 0

		now := time.Now()
		state.next = timestamp() + Timestamp(config.SaveSchedule/100)

		for i := 0; i < 20000; i++ {
			if timestamp() >= state.next {
				count++
			}
			NewV1()
			time.Sleep(1 * time.Millisecond)
		}
		d := time.Since(now)
		timesSaved := int(d.Seconds()) / int(saveDuration)
		if count != timesSaved {
			t.Errorf("Should be as many saves as %d second increments but got: %d instead of %d", saveDuration, count, timesSaved)
		}
	}
}

// Tests that the schedule saves properly when uuid are called in go routines
func TestUUID_State_saveScheduleGo(t *testing.T) {

	if state.saver != nil {

		size := 5000
		ids := make([]UUID, size)

		var wg sync.WaitGroup
		wg.Add(size)

		var count int32
		mutex := &sync.Mutex{}

		now := time.Now()
		state.next = timestamp() + Timestamp(config.SaveSchedule/100)

		for i := 0; i < size; i++ {
			go func(index int) {
				defer wg.Done()
				if timestamp() >= state.next {
					atomic.AddInt32(&count, 1)
				}
				u := NewV1()
				mutex.Lock()
				ids[index] = u
				mutex.Unlock()
				time.Sleep(100 * time.Nanosecond)
			}(i)
		}
		wg.Wait()
		duration := time.Since(now)

		for j := size - 1; j >= 0; j-- {
			for k := 0; k < size; k++ {
				if k == j {
					continue
				}
				if Equal(ids[j], ids[k]) {
					t.Error("Should not create the same V1 UUID", ids[k], ids[j])
				}
			}
		}

		timesSaved := int(duration.Seconds()) / int(saveDuration)
		if int(count) != timesSaved {
			t.Errorf("Should be as many saves as %d second increments but got: %d instead of %d", saveDuration, count, timesSaved)
		}
	}
}