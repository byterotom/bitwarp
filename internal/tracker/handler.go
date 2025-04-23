package tracker

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	pbtr "github.com/Sp92535/proto/tracker/pb"
	"github.com/redis/go-redis/v9"
)

// ip limit to send
const LIMIT = 50


// implemented function to return resource holders of a particular file
func (tr *TrackerServer) GetResourceHolders(ctx context.Context, req *pbtr.GetResourceHoldersRequest) (*pbtr.GetResourceHoldersResponse, error) {

	var res pbtr.GetResourceHoldersResponse
	var err error
	var currentSize int = 0

	res.Holder = make([]*pbtr.HolderRow, len(req.Status))
	
	// get ips for only the chunks that client doesn't have
	for chunkNo, ok := range req.Status {
		if ok {
			continue
		}

		// construct key
		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)

		res.Holder[chunkNo] = &pbtr.HolderRow{}

		// remove expired members before giving new ips
		RemoveExpired(key, ctx)

		// donot append if limit reached -> note: not exiting due to remaining initialization of ips
		if currentSize < LIMIT {
			// append random number of ips from 1-5
			res.Holder[chunkNo].Ips, err = Rdb.ZRandMember(ctx, key, random(4)).Result()
			currentSize += len(res.Holder[chunkNo].Ips)
		}

		if err != nil {
			log.Printf("error fetching holders from redis: %v", err)
			return nil, err
		}
	}
	return &res, nil
}

// implemented function to register client's node server ip in chunk sets
func (tr *TrackerServer) RegisterResourceHolder(ctx context.Context, req *pbtr.RegisterResourceHolderRequest) (*pbtr.Empty, error) {

	for chunkNo, ok := range req.Status {
		if !ok {
			continue
		}

		// construct key
		key := fmt.Sprintf("%s:%d", req.FileHash, chunkNo)
		// construct score
		score := float64(time.Now().Add(TTL).Unix())

		// add to set along with score(expiration time)
		err := Rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: req.Address}).Err()
		if err != nil {
			log.Printf("error adding holder to redis: %v", err)
			return nil, err
		}
	}
	return &pbtr.Empty{}, nil
}

// function to get random number
func random(x int) int {
	return rand.IntN(x) + 1
}