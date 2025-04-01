package data

import "time"

func (d *Data) UpdateClientLastSeen(clientID string) {

	d.clientsMutex.Lock()
	defer d.clientsMutex.Unlock()

	d.clients[clientID] = time.Now().Unix()
}

func (d *Data) GetClientLastSeen(clientID string) int64 {

	d.clientsMutex.Lock()
	defer d.clientsMutex.Unlock()

	return d.clients[clientID]
}

func (d *Data) GetClients() map[string]int64 {

	d.clientsMutex.Lock()
	defer d.clientsMutex.Unlock()

	// take copy of clients
	clients := make(map[string]int64)
	for clientID, lastSeen := range d.clients {

		clients[clientID] = lastSeen
	}

	return clients
}
