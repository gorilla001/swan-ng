package mesos

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/bbklab/swan-ng/mesos/protobuf/mesos"
	"github.com/bbklab/swan-ng/mesos/protobuf/sched"
)

var (
	pub *publisher // global publisher
)

func init() {
	pub = &publisher{
		m:         make(map[subcriber]topicFunc),
		timeout:   time.Second * 5,
		bufferLen: 1024,
	}
}

type publisher struct {
	sync.RWMutex                         // protect m
	m            map[subcriber]topicFunc // hold all of subcribers

	timeout   time.Duration // send topic timeout
	bufferLen int           // buffer length of each subcriber channel
}
type subcriber chan interface{}
type topicFunc func(v interface{}) bool

// numHitTopic return nb of waitting subscribers who cares about the specified value.
func (p *publisher) numHitTopic(v interface{}) (n int) {
	p.RLock()
	defer p.RUnlock()
	for _, tf := range p.m {
		if tf == nil || tf(v) {
			n++
		}
	}
	return
}

// subcribeAll adds a new subscriber that receive all messages.
func (p *publisher) subcribeAll() subcriber {
	return p.subcribe(nil)
}

// subcribe adds a new subscriber that filters messages sent by a topic.
func (p *publisher) subcribe(tf topicFunc) subcriber {
	ch := make(subcriber, p.bufferLen)
	p.Lock()
	p.m[ch] = tf
	p.Unlock()
	return ch
}

// evict removes the specified subscriber from receiving any more messages.
func (p *publisher) evict(sub subcriber) {
	p.Lock()
	delete(p.m, sub)
	close(sub)
	p.Unlock()
}

func (p *publisher) publish(v interface{}) {
	p.RLock()
	defer p.RUnlock()

	var wg sync.WaitGroup
	wg.Add(len(p.m))
	// broadcasting with concurrency
	for sub, tf := range p.m {
		go func(sub subcriber, v interface{}, tf topicFunc) {
			defer wg.Done()
			p.send(sub, v, tf)
		}(sub, v, tf)
	}
	wg.Wait()
}

func (p *publisher) send(sub subcriber, v interface{}, tf topicFunc) {
	// if a subcriber setup topic filter func and not matched by the topic filter
	// skip send message to this subcriber
	if tf != nil && !tf(v) {
		return
	}

	// send with timeout
	if p.timeout > 0 {
		select {
		case sub <- v:
		case <-time.After(p.timeout):
			log.Println("send to subcriber timeout after", p.timeout.String())
		}
		return
	}

	// directely send
	sub <- v
}

//
// utils exported for external usage
//

// SubscribeOffer waitting for the first proper mesos offer
// which satisified the resources request.
func (c *Client) SubscribeOffer(rs []*mesos.Resource) *mesos.Offer {
	var (
		sets = make([]int, 0, 0) // subset idxes of matched offers
		idx  int                 // the final choosen idx
	)

	tf := func(v interface{}) bool {
		ev, ok := v.(*sched.Event)
		if !ok {
			return false
		}

		if ev.GetType() != sched.Event_OFFERS {
			return false
		}

		if n := len(ev.Offers.Offers); n == 0 {
			return false
		}

		// TODO
		// pick up offers those satisified requested resources.
		for i, f := range ev.Offers.Offers {
			if isOfferMatch(f, rs) {
				sets = append(sets, i) // for demo only
			}
		}

		return len(sets) > 0
	}

	sub := pub.subcribe(tf)
	defer pub.evict(sub)

	ev := <-sub
	mesosOffers := ev.(*sched.Event).Offers.Offers

	// by random
	rand.Seed(time.Now().Unix())
	idx = sets[rand.Intn(len(sets))]

	return mesosOffers[idx]
}

// TODO(nmg)
func isOfferMatch(offer *mesos.Offer, res []*mesos.Resource) bool {
	for _, r := range res {
		for _, m := range offer.Resources {
			if r.GetName() == m.GetName() {
			}
		}
	}

	return true
}

// SubscribeUpdate waitting for the specified taskID until finished or error.
func (c *Client) SubscribeTaskUpdate(taskID string) *sched.Event {
	tf := func(v interface{}) bool {
		ev, ok := v.(*sched.Event)
		if !ok {
			return false
		}

		if ev.GetType() != sched.Event_UPDATE {
			return false
		}

		status := ev.GetUpdate().GetStatus()
		return taskID == status.TaskId.GetValue()
	}

	sub := pub.subcribe(tf)
	defer pub.evict(sub)

	ev := <-sub
	return ev.(*sched.Event)
}

// IsTaskDone check that if a task is done or not according by task status.
func (c *Client) IsTaskDone(status *mesos.TaskStatus) bool {
	state := status.GetState()

	switch state {
	case mesos.TaskState_TASK_RUNNING,
		mesos.TaskState_TASK_FINISHED,
		mesos.TaskState_TASK_FAILED,
		mesos.TaskState_TASK_KILLED,
		mesos.TaskState_TASK_ERROR,
		mesos.TaskState_TASK_LOST,
		mesos.TaskState_TASK_DROPPED,
		mesos.TaskState_TASK_GONE:
		return true
	}

	return false
}

func (c *Client) DetectError(status *mesos.TaskStatus) error {
	var (
		state = status.GetState()
		//data  = status.GetData() // docker container inspect result
	)

	switch state {
	case mesos.TaskState_TASK_FAILED,
		mesos.TaskState_TASK_ERROR,
		mesos.TaskState_TASK_LOST,
		mesos.TaskState_TASK_DROPPED,
		mesos.TaskState_TASK_UNREACHABLE,
		mesos.TaskState_TASK_GONE,
		mesos.TaskState_TASK_GONE_BY_OPERATOR,
		mesos.TaskState_TASK_UNKNOWN:
		bs, _ := json.Marshal(map[string]interface{}{
			"state":   state.String(),
			"message": status.GetMessage(),
			"source":  status.GetSource().String(),
			"reason":  status.GetReason().String(),
			"healthy": status.GetHealthy(),
		})
		return errors.New(string(bs))
	}

	return nil
}
