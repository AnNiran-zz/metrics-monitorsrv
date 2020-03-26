package counter

import (
	"time"
	"github.com/twotwotwo/sorts/sortutil"
	"github.com/argcv/stork/log"
)

// refresher checks if the current time later than
// RefreshTimeInterval when called
func (m *Metrics) refresher() {

	// Check is it time to refresh
	currentTime := time.Now().Unix()

	// Allow 1 second deviation in refresh rate because of the 2 seconds sleep
	if (currentTime - m.LastRefreshTime) >= RefreshTimeInterval {
		m.refresh()
		m.LastRefreshTime = currentTime
		log.Infof("Refreshed metrics data ...")
	}
}

// refresh updates metric records to be within the range of 60 minutes
// refreshing is done on each 30 seconds
func (m *Metrics) refresh() {

	// if there are no records we do not need to check anything
	if len(m.Data.List) == 0 {
		return
	}

	// create starting time of the range we need to keep data for
	timeStart := time.Now().Unix() - int64(TimeRange)

	// if the first record is within the last hour - we do not need
	// to compare anything yet
	if m.Data.FirstRecord >= timeStart {
		log.Info("All records are in the specified time range now")
		return
	}

	// Loop over the each metric
	for _, metric := range m.Data.List {
		// if there are no records for the metric we do not need to check anything
		// and move to the next metric
		if len(metric.Data.ValueSet) == 0 {
			continue
		}

		// if the first record of the metric is within the last hour,
		// we do not need to compare anything yet
		if m.Data.FirstRecord >= timeStart {
			continue
		}

		metric.Lock.Lock()
		
		// sort values in increasing order
		// keep the metric values in a separate variable because we do not want to loop over
		// changing length
		sortedValues := sortValues(metric.Data.ValueSet)

		// range is over Unixtime -> value 
		// starting from the least int64 key in the map
		// in this way the loops until finding a value that should not
		// be excluded are theoretically less than simply looping
		for time, value := range sortedValues {
			// if the ith value is equal or greater than the set time range -
			// we do not need to loop anymore because i+n will be greater as well
			if time >= timeStart {
				break
			}
			
			// remove value if it is not within the setup time range
			if time < timeStart {
				metric.Data.Sum -= value
				delete(metric.Data.ValueSet, time)
				continue
			}
		}

		metric.Lock.Unlock()
	}
	
	return
}

// sortValues sort the map saving the times and vlaues of a metric 
// in increasing order
func sortValues(values timeValues) timeValues {
	sortedValues := make(timeValues)
	keys := make([]int64, len(values))

	for k := range values {
		keys = append(keys, k)
	}

	sortutil.Int64s(keys)
	for _, k := range keys {
		sortedValues[k] = values[k]
	}

	return sortedValues
}
