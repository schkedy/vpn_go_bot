package cache

import "context"

type Storage interface{
	Get(ctx context.Context,key string) 
}
