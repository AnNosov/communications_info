package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/AnNosov/communications_info/internal/entity"
)

type ResultInfoCache struct {
	*entity.Cache
}

func New(defaultExpiration time.Duration) *ResultInfoCache {

	var data entity.ResultSetT

	cache := entity.Cache{
		Data:              data,
		DefaultExpiration: defaultExpiration,
		Created:           time.Now(),
	}
	log.Println("----- cache created at: ", cache.Created, " -----")

	return &ResultInfoCache{
		Cache: &cache,
	}
}

func (c *ResultInfoCache) Set(cache entity.ResultSetT, duration time.Duration) {

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.DefaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		c.Expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.Data = cache
}

func (c *ResultInfoCache) Get() (entity.ResultSetT, error) {

	c.RLock()

	defer c.RUnlock()

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if c.Expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем ошибку
		if time.Now().UnixNano() > c.Expiration {
			return entity.ResultSetT{}, fmt.Errorf("cache is outdated")
		}

	}

	return c.Cache.Data, nil
}
